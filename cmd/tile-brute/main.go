// Brute force the tile API to see how tile keys are allocated
package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"wloc/lib"
)

func main() {
	max := 78768401
	min := 78767287
	for i := min; i < max; i++ {
		// randKey := randInt(min, max)
		randKey := int64(i)
		tiles, err := lib.GetTile(randKey)
		if err != nil {
			log.Printf("%d not found\n", randKey)
			continue
		}
		log.Printf("%d found records\n", randKey)
		f, err := os.Create(fmt.Sprintf("points/%d.txt", randKey))
		if err != nil {
			panic(err)
		}
		defer f.Close()
		for _, region := range tiles.GetRegion() {
			for _, device := range region.GetDevices() {
				_, err := f.WriteString(fmt.Sprintf("%f %f\n", intCoordToFloat(device.Entry.Lat), intCoordToFloat(device.Entry.Long)))
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

func intCoordToFloat(n int32) float64 {
	return float64(n) * math.Pow10(-7)
}

func randInt(min, max int) int64 {
	return int64(rand.Intn(max-min) + min)
}
