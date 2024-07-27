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
	"time"
	"wloc/lib"
	"wloc/lib/morton"
	"wloc/lib/shapefiles"

	"github.com/schollz/progressbar/v3"
)

var progress = progressbar.Default(MaxTile * MaxTile)

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
	if err := progress.Add(gen.Current.Y*(MaxTile+1) + gen.Current.X); err != nil {
		panic(err)
	}
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
	for i := range 500 {
		wait.Add(1)
		go func() {
			Datafetcher(ctx, &database, gen.Channel())
			wait.Done()
			log.Println("Thread completed")
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
		if err := progress.Add(1); err != nil {
			panic(err)
		}
		lat, lon := morton.FromTile(coord.X, coord.Y, 13)
		if shapefiles.IsInWater(lat, lon) {
			continue
		}
		select {
		case <-ctx.Done():
			return
		default:
			for {
				code := morton.Pack(coord.X, coord.Y, 13)
				aps, err := lib.GetTile(code)
				if err != nil {
					if err.Error() == "unexpected status code: 404" {
						break
					}
					log.Println("Something went from in tile call: ", err)
					time.Sleep(100 * time.Millisecond)
					continue
				}
				log.Printf("\nFound %d access points at %f, %f\n", len(aps), lat, lon)
				database.Add(aps)
				break
			}
		}
	}
}
