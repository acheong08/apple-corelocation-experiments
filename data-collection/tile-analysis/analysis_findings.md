# WiFi BSSID Location Tracking Analysis

## Overview

Analysis of location changes for WiFi BSSIDs collected from Apple's location services data, stored in SQLite database with 342,356 location change records across 9,327 unique BSSIDs.

## Summary Statistics

- **Total location changes**: 342,356
- **Number of batches**: 343
- **Overall average distance**: 0.7 m
- **Min distance**: 0.0 m
- **Max distance**: 24.9 m
- **Standard deviation**: 1.0 m

### 1. Update Pattern Analysis

#### Overall Statistics

```sql
SELECT COUNT(*) as total_records,
       COUNT(DISTINCT bssid) as unique_bssids,
       COUNT(DISTINCT old_tile_key) as unique_tiles
FROM location_changes;
```

- **Total Records**: 342,356 location changes
- **Unique BSSIDs**: 9,327
- **Unique Tiles**: 51

#### Time Interval Analysis

```sql
WITH time_intervals AS (
    SELECT
        bssid,
        change_time,
        LAG(change_time) OVER (PARTITION BY bssid ORDER BY change_time) as prev_change,
        (julianday(change_time) - julianday(LAG(change_time) OVER (PARTITION BY bssid ORDER BY change_time))) * 24 * 60 as minutes_between
    FROM location_changes
)
SELECT minutes_between
FROM time_intervals
WHERE minutes_between IS NOT NULL AND minutes_between < 1440;
```

**Two Distinct Groups Identified:**

1. **Short Intervals (≤200 minutes): 95.2% of data** during update zone
   - Count: 310,548 intervals
   - Mean: 11.3 minutes
   - Median: 5.0 minutes
   - Range: 4.3 to 100.1 minutes

2. **Long Intervals (≥1100 minutes): 4.8% of data**
   - Count: 15,725 intervals
   - Mean: 1,315.1 minutes (~22 hours)
   - Median: 1,350.0 minutes (~22.5 hours)
   - Range: 1,130 to 1,435 minutes (~18.8 hours to ~23.9 hours)

Essentially, AP locations are guaranteed to have updates at least once a day, and there is a 5 hour range within which updates are allowed to occur.

**Gap**: 100.1 to 1,130.0 minutes with no data points

### 2. Bulk Update Event Detection

#### Update Window Analysis

```sql
WITH update_windows AS (
    SELECT
        strftime('%Y-%m-%d %H:%M:00', change_time) as minute_window,
        COUNT(DISTINCT old_tile_key) as tiles_involved,
        COUNT(*) as total_updates
    FROM location_changes
    GROUP BY strftime('%Y-%m-%d %H:%M:00', change_time)
    HAVING COUNT(*) >= 100
)
SELECT minute_window, tiles_involved, total_updates
FROM update_windows
ORDER BY total_updates DESC;
```

**Top Update Windows:**

- 2025-09-07 06:31:00: 5,277 updates across 23 tiles
- 2025-09-08 07:41:00: 5,209 updates across 23 tiles
- 2025-09-07 06:26:00: 5,055 updates across 22 tiles
- 2025-09-07 06:36:00: 4,119 updates across 20 tiles

**Pattern**: Daily updates with peak activity during morning hours (6-8 AM).

### 3. Tile Change Analysis

#### Cross-Tile Movement Check

```sql
SELECT old_tile_key, new_tile_key, COUNT(*) as changes
FROM location_changes
WHERE old_tile_key != new_tile_key
GROUP BY old_tile_key, new_tile_key;
```

**Findings**:

- Only 26 actual tile changes across 4 tile pairs
- 99.99% of BSSIDs remain in the same tile
- Minimal cross-tile movement detected

#### Tile Update Distribution

```sql
SELECT
    old_tile_key,
    COUNT(*) as updates_in_tile,
    COUNT(DISTINCT bssid) as unique_bssids_in_tile
FROM location_changes
GROUP BY old_tile_key
ORDER BY updates_in_tile DESC
LIMIT 15;
```

**Top Tiles by Update Volume**:

- Tile 78946923: 77,088 updates (1,857 BSSIDs)
- Tile 78946921: 33,857 updates (966 BSSIDs)
- Tile 98908354: 30,863 updates (1,088 BSSIDs)

### 4. Spatial-Temporal Correlation

#### Tile Overlap During Updates

```sql
WITH update_windows AS (
    SELECT
        strftime('%Y-%m-%d %H:%M:00', change_time) as minute_window,
        COUNT(DISTINCT old_tile_key) as tiles_involved,
        COUNT(*) as total_updates
    FROM location_changes
    GROUP BY strftime('%Y-%m-%d %H:%M:00', change_time)
)
SELECT
    CASE
        WHEN tiles_involved = 1 THEN 'Single tile'
        WHEN tiles_involved <= 5 THEN 'Few tiles'
        WHEN tiles_involved <= 15 THEN 'Many tiles'
        ELSE 'Massive update'
    END as update_pattern,
    COUNT(*) as occurrences
FROM update_windows
GROUP BY update_pattern;
```

**Result**: Mix of update patterns - single tile, few tiles, many tiles, and some larger coordinated updates affecting 15+ tiles simultaneously.

### 5. Location Oscillation Analysis

```sql
-- Analyze the number of unique locations per BSSID
WITH location_counts AS (
    SELECT
        bssid,
        COUNT(DISTINCT ROUND(old_lat, 6) || ',' || ROUND(old_long, 6)) as unique_old_locations,
        COUNT(DISTINCT ROUND(new_lat, 6) || ',' || ROUND(new_long, 6)) as unique_new_locations,
        COUNT(*) as total_changes
    FROM location_changes
    GROUP BY bssid
),
all_locations AS (
    SELECT
        bssid,
        ROUND(old_lat, 6) as lat,
        ROUND(old_long, 6) as long
    FROM location_changes
    UNION
    SELECT
        bssid,
        ROUND(new_lat, 6) as lat,
        ROUND(new_long, 6) as long
    FROM location_changes
),
unique_location_counts AS (
    SELECT
        bssid,
        COUNT(DISTINCT lat || ',' || long) as total_unique_locations
    FROM all_locations
    GROUP BY bssid
)
SELECT
    u.total_unique_locations,
    l.total_changes,
    COUNT(u.bssid) as bssid_count,
    ROUND(CAST(l.total_changes AS FLOAT) / u.total_unique_locations, 2) as changes_per_location
FROM unique_location_counts u
JOIN location_counts l ON u.bssid = l.bssid
GROUP BY u.total_unique_locations, l.total_changes
ORDER BY u.total_unique_locations DESC, l.total_changes DESC;
```

**Key Findings:**

- Most BSSIDs oscillate between 1-5 unique locations
- BSSIDs with 5 unique locations show up to 80 changes per batch, with 10.4 changes per location being most common (1,276 BSSIDs)
- Some BSSIDs have up to 78 unique locations with 47 changes (15.6 changes per location)
- Single-location BSSIDs can have up to 43 changes, indicating precision refinement rather than movement
- Updates occur daily with peak activity during morning hours (6-8 AM)
- Multiple tiles often update simultaneously, meaning that updates are not sharded by tiles.
