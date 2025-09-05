package main

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

type Database struct {
	db *sql.DB
}

type BSSID struct {
	BSSID    string
	Lat      float64
	Long     float64
	LastSeen time.Time
}

type LocationChange struct {
	ID         int64
	BSSID      string
	OldLat     float64
	OldLong    float64
	NewLat     float64
	NewLong    float64
	ChangeTime time.Time
}

func NewDatabase(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	d := &Database{db: db}
	if err := d.createTables(); err != nil {
		return nil, err
	}

	return d, nil
}

func (d *Database) createTables() error {
	createCurrentBSSIDs := `
	CREATE TABLE IF NOT EXISTS current_bssids (
		bssid TEXT PRIMARY KEY,
		lat REAL NOT NULL,
		long REAL NOT NULL,
		last_seen TIMESTAMP NOT NULL
	);`

	createLocationChanges := `
	CREATE TABLE IF NOT EXISTS location_changes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		bssid TEXT NOT NULL,
		old_lat REAL NOT NULL,
		old_long REAL NOT NULL,
		new_lat REAL NOT NULL,
		new_long REAL NOT NULL,
		change_time TIMESTAMP NOT NULL
	);`

	if _, err := d.db.Exec(createCurrentBSSIDs); err != nil {
		return err
	}

	if _, err := d.db.Exec(createLocationChanges); err != nil {
		return err
	}

	return nil
}

func (d *Database) GetBSSID(bssid string) (*BSSID, error) {
	query := "SELECT bssid, lat, long, last_seen FROM current_bssids WHERE bssid = ?"
	row := d.db.QueryRow(query, bssid)

	var b BSSID
	err := row.Scan(&b.BSSID, &b.Lat, &b.Long, &b.LastSeen)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &b, nil
}

func (d *Database) InsertBSSID(bssid string, lat, long float64) error {
	query := "INSERT OR REPLACE INTO current_bssids (bssid, lat, long, last_seen) VALUES (?, ?, ?, ?)"
	_, err := d.db.Exec(query, bssid, lat, long, time.Now())
	return err
}

func (d *Database) UpdateBSSIDLocation(bssid string, oldLat, oldLong, newLat, newLong float64) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	updateQuery := "UPDATE current_bssids SET lat = ?, long = ?, last_seen = ? WHERE bssid = ?"
	_, err = tx.Exec(updateQuery, newLat, newLong, time.Now(), bssid)
	if err != nil {
		return err
	}

	logQuery := "INSERT INTO location_changes (bssid, old_lat, old_long, new_lat, new_long, change_time) VALUES (?, ?, ?, ?, ?, ?)"
	_, err = tx.Exec(logQuery, bssid, oldLat, oldLong, newLat, newLong, time.Now())
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (d *Database) UpdateLastSeen(bssid string) error {
	query := "UPDATE current_bssids SET last_seen = ? WHERE bssid = ?"
	_, err := d.db.Exec(query, time.Now(), bssid)
	return err
}

func (d *Database) Close() error {
	return d.db.Close()
}
