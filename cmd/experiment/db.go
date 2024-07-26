package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"
	"wloc/lib"
	"wloc/lib/mac"
	"wloc/lib/morton"

	_ "modernc.org/sqlite"
)

type db struct {
	db   *sql.DB
	lock sync.Mutex
}

func InitDatabase() db {
	d, err := sql.Open("sqlite", "seeds.db")
	if err != nil {
		panic(fmt.Errorf("Failed to open database: %w", err))
	}
	if _, err := d.Exec(`CREATE TABLE IF NOT EXISTS seeds (
			bssid INTEGER PRIMARY KEY,
			tilekey INTEGER NOT NULL,
			lat REAL NOT NULL,
			lon REAL NOT NULL,
			created INTEGER NOT NULL
		)
		`); err != nil {
		panic(fmt.Errorf("Failed to create table: %w", err))
	}
	return db{db: d}
}

func (d *db) Add(s []lib.AP) {
	d.lock.Lock()
	defer d.lock.Unlock()
	tx, err := d.db.Begin()
	if err != nil {
		panic("transaction failed")
	}

	for _, ap := range s {
		bssid, err := mac.Encode(ap.BSSID)
		if err != nil {
			continue
		}
		_, err = tx.Exec("INSERT OR IGNORE INTO seeds (bssid, tilekey, lat, lon, created) VALUES (?,?,?,?,?)", bssid, morton.Encode(ap.Location.Lat, ap.Location.Long, 13), ap.Location.Lat, ap.Location.Long, time.Now().Unix())
		if err != nil {
			log.Println("Failed to insert into seeds ", bssid)
			continue
		}
	}
	if err = tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			panic("Can't roll back transaction")
		}
		log.Println("Commit failed but was rolled back")
	}
}