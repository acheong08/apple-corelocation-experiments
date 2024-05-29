package morton

import (
	_ "embed"
	"encoding/json"
	"log"
	"strings"

	"github.com/sajari/regression"
)

//go:embed tileCoordsToGPS.json
var tileCoordsJson string
var tileCoords []tileCoordEntry

var (
	mortonToGpsLat  *regression.Regression = new(regression.Regression)
	mortonToGpsLong *regression.Regression = new(regression.Regression)
	gpsToMortonLat  *regression.Regression = new(regression.Regression)
	gpsToMortonLong *regression.Regression = new(regression.Regression)
)

func init() {
	dec := json.NewDecoder(strings.NewReader(tileCoordsJson))
	// Decode JSON array
	err := dec.Decode(&tileCoords)
	if err != nil {
		panic(err)
	}
	train(mortonToGpsLat, func(i int) (float64, float64) {
		return tileCoords[i].GPS[0], float64(tileCoords[i].Morton[0])
	})
	train(mortonToGpsLong, func(i int) (float64, float64) {
		return tileCoords[i].GPS[1], float64(tileCoords[i].Morton[1])
	})
	train(gpsToMortonLat, func(i int) (float64, float64) {
		return float64(tileCoords[i].Morton[0]), tileCoords[i].GPS[0]
	})
	train(gpsToMortonLong, func(i int) (float64, float64) {
		return float64(tileCoords[i].Morton[1]), tileCoords[i].GPS[1]
	})
	// Log the accuracy
	log.Println("Morton to GPS Latitude R^2:", mortonToGpsLat.R2)
	log.Println("Morton to GPS Longitude R^2:", mortonToGpsLong.R2)
	log.Println("GPS to Morton Latitude R^2:", gpsToMortonLat.R2)
	log.Println("GPS to Morton Longitude R^2:", gpsToMortonLong.R2)
}

func train(reg *regression.Regression, dataPoint func(i int) (float64, float64)) {
	for i := 0; i < len(tileCoords); i++ {
		x, y := dataPoint(i)
		reg.Train(regression.DataPoint(x, []float64{y}))
	}
	err := reg.Run()
	if err != nil {
		panic(err)
	}
}

type tileCoordEntry struct {
	GPS    []float64 `json:"coord"`
	Morton []int     `json:"morton"`
}
