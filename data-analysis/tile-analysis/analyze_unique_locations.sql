-- Analyze the number of unique locations per BSSID
-- This helps identify BSSIDs that oscillate between multiple locations

-- Main query: Count unique locations per BSSID
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
    -- Get all unique locations (both old and new) per BSSID
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

-- Main results
SELECT 
    u.total_unique_locations,
    l.total_changes,
    COUNT(u.bssid) as bssid_count,
    ROUND(CAST(l.total_changes AS FLOAT) / u.total_unique_locations, 2) as changes_per_location
FROM unique_location_counts u
JOIN location_counts l ON u.bssid = l.bssid
GROUP BY u.total_unique_locations, l.total_changes
ORDER BY u.total_unique_locations DESC, l.total_changes DESC;
