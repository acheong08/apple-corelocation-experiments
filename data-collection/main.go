package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var (
		tileFile = flag.String("tiles", "", "Path to text file containing comma-separated tile keys")
		dbPath   = flag.String("db", "bssid_tracking.db", "Path to SQLite database file")
		interval = flag.Duration("interval", 5*time.Minute, "Collection interval (e.g., 5m, 30s, 1h)")
	)
	flag.Parse()

	if *tileFile == "" {
		fmt.Fprintf(os.Stderr, "Usage: %s -tiles <file> [-db <path>] [-interval <duration>]\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	tileContent, err := os.ReadFile(*tileFile)
	if err != nil {
		log.Fatalf("Failed to read tile file %s: %v", *tileFile, err)
	}

	tileKeys, err := parseTileKeys(string(tileContent))
	if err != nil {
		log.Fatalf("Failed to parse tile keys: %v", err)
	}

	log.Printf("Loaded %d tile keys from %s", len(tileKeys), *tileFile)

	db, err := NewDatabase(*dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	collector := NewCollector(db, tileKeys, *interval)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Received shutdown signal, stopping collection...")
		cancel()
	}()

	log.Printf("Starting BSSID location tracking with %d tile keys, checking every %v", len(tileKeys), *interval)
	log.Printf("Database: %s", *dbPath)
	log.Println("Press Ctrl+C to stop")

	if err := collector.Start(ctx); err != nil && err != context.Canceled {
		log.Fatalf("Collection failed: %v", err)
	}

	log.Println("Collection stopped")
}
