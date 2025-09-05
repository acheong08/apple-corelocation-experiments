package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"wloc/lib/morton"

	_ "modernc.org/sqlite"
)

func main() {
	var (
		dbPath = flag.String("db", "bssid_tracking.db", "Path to SQLite database file")
		level  = flag.Int("level", 13, "Tile level for morton encoding")
		dryRun = flag.Bool("dry-run", false, "Show what would be done without making changes")
	)
	flag.Parse()

	if _, err := os.Stat(*dbPath); os.IsNotExist(err) {
		log.Fatalf("Database file does not exist: %s", *dbPath)
	}

	db, err := sql.Open("sqlite", *dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	if err := addTileKeyColumns(db, *dryRun); err != nil {
		log.Fatalf("Failed to add tile_key columns: %v", err)
	}

	if err := updateCurrentBSSIDTileKeys(db, *level, *dryRun); err != nil {
		log.Fatalf("Failed to update current_bssids tile keys: %v", err)
	}

	if err := updateLocationChangesTileKeys(db, *level, *dryRun); err != nil {
		log.Fatalf("Failed to update location_changes tile keys: %v", err)
	}

	log.Println("Successfully completed tile key calculation and insertion for both current_bssids and location_changes tables")
}

func addTileKeyColumns(db *sql.DB, dryRun bool) error {
	// Add tile_key column to current_bssids table
	alterCurrentSQL := "ALTER TABLE current_bssids ADD COLUMN tile_key INTEGER"

	if dryRun {
		log.Printf("DRY RUN: Would execute: %s", alterCurrentSQL)
	} else {
		var columnExists bool
		checkSQL := `SELECT COUNT(*) FROM pragma_table_info('current_bssids') WHERE name='tile_key'`
		if err := db.QueryRow(checkSQL).Scan(&columnExists); err != nil {
			return fmt.Errorf("failed to check if current_bssids.tile_key column exists: %w", err)
		}

		if !columnExists {
			if _, err := db.Exec(alterCurrentSQL); err != nil {
				return fmt.Errorf("failed to add tile_key column to current_bssids: %w", err)
			}
			log.Println("Added tile_key column to current_bssids table")
		} else {
			log.Println("tile_key column already exists in current_bssids, skipping creation")
		}
	}

	// Add tile_key columns to location_changes table (for old and new locations)
	alterOldTileSQL := "ALTER TABLE location_changes ADD COLUMN old_tile_key INTEGER"
	alterNewTileSQL := "ALTER TABLE location_changes ADD COLUMN new_tile_key INTEGER"

	if dryRun {
		log.Printf("DRY RUN: Would execute: %s", alterOldTileSQL)
		log.Printf("DRY RUN: Would execute: %s", alterNewTileSQL)
		return nil
	}

	// Check and add old_tile_key column
	var oldColumnExists bool
	checkOldSQL := `SELECT COUNT(*) FROM pragma_table_info('location_changes') WHERE name='old_tile_key'`
	if err := db.QueryRow(checkOldSQL).Scan(&oldColumnExists); err != nil {
		return fmt.Errorf("failed to check if location_changes.old_tile_key column exists: %w", err)
	}

	if !oldColumnExists {
		if _, err := db.Exec(alterOldTileSQL); err != nil {
			return fmt.Errorf("failed to add old_tile_key column to location_changes: %w", err)
		}
		log.Println("Added old_tile_key column to location_changes table")
	} else {
		log.Println("old_tile_key column already exists in location_changes, skipping creation")
	}

	// Check and add new_tile_key column
	var newColumnExists bool
	checkNewSQL := `SELECT COUNT(*) FROM pragma_table_info('location_changes') WHERE name='new_tile_key'`
	if err := db.QueryRow(checkNewSQL).Scan(&newColumnExists); err != nil {
		return fmt.Errorf("failed to check if location_changes.new_tile_key column exists: %w", err)
	}

	if !newColumnExists {
		if _, err := db.Exec(alterNewTileSQL); err != nil {
			return fmt.Errorf("failed to add new_tile_key column to location_changes: %w", err)
		}
		log.Println("Added new_tile_key column to location_changes table")
	} else {
		log.Println("new_tile_key column already exists in location_changes, skipping creation")
	}

	return nil
}

func updateCurrentBSSIDTileKeys(db *sql.DB, level int, dryRun bool) error {
	selectSQL := "SELECT bssid, lat, long FROM current_bssids"

	if dryRun {
		log.Printf("DRY RUN: Would query current_bssids: %s", selectSQL)
	}

	rows, err := db.Query(selectSQL)
	if err != nil {
		return fmt.Errorf("failed to query current_bssids: %w", err)
	}
	defer rows.Close()

	var updates []struct {
		bssid   string
		tileKey int64
	}

	for rows.Next() {
		var bssid string
		var lat, long float64

		if err := rows.Scan(&bssid, &lat, &long); err != nil {
			return fmt.Errorf("failed to scan current_bssids row: %w", err)
		}

		tileKey := morton.Encode(lat, long, level)

		if dryRun {
			log.Printf("DRY RUN: current_bssids BSSID %s at (%.6f, %.6f) -> tile_key %d", bssid, lat, long, tileKey)
		} else {
			updates = append(updates, struct {
				bssid   string
				tileKey int64
			}{bssid, tileKey})
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error reading current_bssids rows: %w", err)
	}

	if dryRun {
		log.Printf("DRY RUN: Would update %d current_bssids records", len(updates))
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction for current_bssids: %w", err)
	}
	defer tx.Rollback()

	updateSQL := "UPDATE current_bssids SET tile_key = ? WHERE bssid = ?"
	stmt, err := tx.Prepare(updateSQL)
	if err != nil {
		return fmt.Errorf("failed to prepare current_bssids update statement: %w", err)
	}
	defer stmt.Close()

	updatedCount := 0
	for _, update := range updates {
		if _, err := stmt.Exec(update.tileKey, update.bssid); err != nil {
			return fmt.Errorf("failed to update current_bssids BSSID %s: %w", update.bssid, err)
		}
		updatedCount++

		if updatedCount%1000 == 0 {
			log.Printf("Updated %d/%d current_bssids records...", updatedCount, len(updates))
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit current_bssids transaction: %w", err)
	}

	log.Printf("Successfully updated %d current_bssids records with tile keys (level %d)", updatedCount, level)
	return nil
}

func updateLocationChangesTileKeys(db *sql.DB, level int, dryRun bool) error {
	selectSQL := "SELECT id, old_lat, old_long, new_lat, new_long FROM location_changes"

	if dryRun {
		log.Printf("DRY RUN: Would query location_changes: %s", selectSQL)
	}

	rows, err := db.Query(selectSQL)
	if err != nil {
		return fmt.Errorf("failed to query location_changes: %w", err)
	}
	defer rows.Close()

	var updates []struct {
		id         int64
		oldTileKey int64
		newTileKey int64
	}

	for rows.Next() {
		var id int64
		var oldLat, oldLong, newLat, newLong float64

		if err := rows.Scan(&id, &oldLat, &oldLong, &newLat, &newLong); err != nil {
			return fmt.Errorf("failed to scan location_changes row: %w", err)
		}

		oldTileKey := morton.Encode(oldLat, oldLong, level)
		newTileKey := morton.Encode(newLat, newLong, level)

		if dryRun {
			log.Printf("DRY RUN: location_changes ID %d: old (%.6f, %.6f) -> tile_key %d, new (%.6f, %.6f) -> tile_key %d",
				id, oldLat, oldLong, oldTileKey, newLat, newLong, newTileKey)
		} else {
			updates = append(updates, struct {
				id         int64
				oldTileKey int64
				newTileKey int64
			}{id, oldTileKey, newTileKey})
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error reading location_changes rows: %w", err)
	}

	if dryRun {
		log.Printf("DRY RUN: Would update %d location_changes records", len(updates))
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction for location_changes: %w", err)
	}
	defer tx.Rollback()

	updateSQL := "UPDATE location_changes SET old_tile_key = ?, new_tile_key = ? WHERE id = ?"
	stmt, err := tx.Prepare(updateSQL)
	if err != nil {
		return fmt.Errorf("failed to prepare location_changes update statement: %w", err)
	}
	defer stmt.Close()

	updatedCount := 0
	for _, update := range updates {
		if _, err := stmt.Exec(update.oldTileKey, update.newTileKey, update.id); err != nil {
			return fmt.Errorf("failed to update location_changes ID %d: %w", update.id, err)
		}
		updatedCount++

		if updatedCount%1000 == 0 {
			log.Printf("Updated %d/%d location_changes records...", updatedCount, len(updates))
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit location_changes transaction: %w", err)
	}

	log.Printf("Successfully updated %d location_changes records with tile keys (level %d)", updatedCount, level)
	return nil
}
