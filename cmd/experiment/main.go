package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"wloc/lib"
	"wloc/lib/morton"
	"wloc/lib/shapefiles"
)

func main() {
	if !shapefiles.IsInWater(82.940327, -180.000000) {
		panic("something went wrong with shapefiles")
	}
	gen := NewGenerator()
	err := gen.LoadState("state.gob")
	if err != nil {
		log.Println("Failed to load state. This is expected for new runs")
	}
	gen.Start()
	defer func() {
		err := recover()
		if err != nil {
			log.Println("Panic caught ", err)
		}
		err = gen.SaveState("state.gob")
		if err != nil {
			b, _ := json.Marshal(gen)
			fmt.Println(string(b))
			return
		}
		log.Println("State saved")
	}()
	database := InitDatabase()
	ctx, cancel := context.WithCancel(context.Background())
	wait := sync.WaitGroup{}
	for i := range 8 {
		wait.Add(1)
		go func() {
			Datafetcher(ctx, &database, gen.Channel())
			wait.Done()
		}()
		log.Println("Started thread ", i)
	}
	// Catch ctrl+c for graceful exit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	log.Println("Exiting gracefully...")
	cancel()
	// Wait for fetchers to finish
	wait.Wait()
	log.Println("Tasks completed")
}

func Datafetcher(ctx context.Context, database *db, c <-chan Coordinate) {
	for coord := range c {
		lat, lon := morton.FromTile(coord.X, coord.Y, 13)
		if shapefiles.IsInWater(lat, lon) {
			log.Println("Coordinate is in water, skipping")
			continue
		}
		select {
		case <-ctx.Done():
			return
		default:
			code := morton.Pack(coord.X, coord.Y, 13)
			aps, err := lib.GetTile(code)
			if err != nil {
				if err.Error() == "unexpected status code: 404" {
					log.Printf("Nothing found at %f %f", lat, lon)
					continue
				}
				log.Println("Something went from in tile call: ", err)
			}
			log.Printf("Found %d access points\n", len(aps))
			database.Add(aps)
		}
	}
}
