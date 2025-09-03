package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
	"wloc/lib"
)

type Collector struct {
	db       *Database
	tileKeys []int64
	interval time.Duration
}

func NewCollector(db *Database, tileKeys []int64, interval time.Duration) *Collector {
	return &Collector{
		db:       db,
		tileKeys: tileKeys,
		interval: interval,
	}
}

func (c *Collector) Start(ctx context.Context) error {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	log.Printf("Starting data collection with %d tile keys, interval: %v", len(c.tileKeys), c.interval)

	for {
		select {
		case <-ctx.Done():
			log.Println("Collection stopped")
			return ctx.Err()
		case <-ticker.C:
			if err := c.collectData(); err != nil {
				log.Printf("Error collecting data: %v", err)
			}
		}
	}
}

func (c *Collector) collectData() error {
	for _, tileKey := range c.tileKeys {
		if err := c.processTile(tileKey); err != nil {
			log.Printf("Error processing tile %d: %v", tileKey, err)
			continue
		}
	}
	return nil
}

func (c *Collector) processTile(tileKey int64) error {
	aps, err := lib.GetTile(tileKey)
	if err != nil {
		return fmt.Errorf("failed to get tile %d: %w", tileKey, err)
	}

	log.Printf("Processing tile %d with %d APs", tileKey, len(aps))

	for _, ap := range aps {
		if err := c.processBSSID(ap.BSSID, ap.Location.Lat, ap.Location.Long); err != nil {
			log.Printf("Error processing BSSID %s: %v", ap.BSSID, err)
		}
	}

	return nil
}

func (c *Collector) processBSSID(bssid string, lat, long float64) error {
	existing, err := c.db.GetBSSID(bssid)
	if err != nil {
		return fmt.Errorf("failed to get BSSID %s: %w", bssid, err)
	}

	if existing == nil {
		log.Printf("New BSSID discovered: %s at (%.6f, %.6f)", bssid, lat, long)
		return c.db.InsertBSSID(bssid, lat, long)
	}

	const epsilon = 1e-6
	if hasLocationChanged(existing.Lat, existing.Long, lat, long, epsilon) {
		distance := calculateDistance(existing.Lat, existing.Long, lat, long)
		log.Printf("Location changed for BSSID %s: (%.6f, %.6f) -> (%.6f, %.6f), distance: %.2fm",
			bssid, existing.Lat, existing.Long, lat, long, distance)
		return c.db.UpdateBSSIDLocation(bssid, existing.Lat, existing.Long, lat, long)
	}

	return c.db.UpdateLastSeen(bssid)
}

func hasLocationChanged(oldLat, oldLong, newLat, newLong, epsilon float64) bool {
	return math.Abs(oldLat-newLat) > epsilon || math.Abs(oldLong-newLong) > epsilon
}

func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371000

	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

func parseTileKeys(content string) ([]int64, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, fmt.Errorf("empty file content")
	}

	parts := strings.Split(content, ",")
	tileKeys := make([]int64, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		tileKey, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid tile key '%s': %w", part, err)
		}

		tileKeys = append(tileKeys, tileKey)
	}

	if len(tileKeys) == 0 {
		return nil, fmt.Errorf("no valid tile keys found")
	}

	return tileKeys, nil
}
