-- Migration script to fix datetime format in location_changes table
-- The timestamps are in Go's time.Time format with monotonic clock data

-- First, create a backup of the original table
CREATE TABLE location_changes_backup AS SELECT * FROM location_changes;

-- Create a new table with proper datetime column
CREATE TABLE location_changes_new (
    id INTEGER PRIMARY KEY,
    bssid TEXT,
    old_lat REAL,
    old_lon REAL,
    new_lat REAL,
    new_lon REAL,
    change_time DATETIME,
    old_tile_key INTEGER,
    new_tile_key INTEGER
);

-- Extract the datetime part (before " m=") and convert to proper SQLite datetime
-- The format is "2025-09-04 04:59:51.897008864 +0000 UTC m=+66001.593305669"
-- We want just "2025-09-04 04:59:51.897008864"
INSERT INTO location_changes_new (
    id, bssid, old_lat, old_lon, new_lat, new_lon, change_time, old_tile_key, new_tile_key
)
SELECT 
    id,
    bssid,
    old_lat,
    old_lon,
    new_lat,
    new_lon,
    DATETIME(SUBSTR(change_time, 1, INSTR(change_time, ' +') - 1)) as change_time,
    old_tile_key,
    new_tile_key
FROM location_changes;

-- Drop the original table and rename the new one
DROP TABLE location_changes;
ALTER TABLE location_changes_new RENAME TO location_changes;

-- Verify the migration worked
SELECT 'Migration complete. Sample records:' as status;
SELECT id, bssid, change_time FROM location_changes LIMIT 5;

-- Show statistics
SELECT 
    'Total records migrated:' as info,
    COUNT(*) as count
FROM location_changes;

SELECT 
    'Date range:' as info,
    MIN(change_time) as earliest,
    MAX(change_time) as latest
FROM location_changes;