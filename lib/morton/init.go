package morton

import (
	_ "embed"
	"encoding/json"
	"io"
	"strings"

	"github.com/sajari/regression"
)

//go:embed tileCoordsToGPS.jsonl
var tileCoordsJsonl string
var tileCoords []tileCoordEntry

var (
	mortonToGpsLat  *regression.Regression = new(regression.Regression)
	mortonToGpsLong *regression.Regression = new(regression.Regression)
	gpsToMortonLat  *regression.Regression = new(regression.Regression)
	gpsToMortonLong *regression.Regression = new(regression.Regression)
)

func init() {
	dec := json.NewDecoder(strings.NewReader(tileCoordsJsonl))
	for {
		var entry tileCoordEntry
		if err := dec.Decode(&entry); err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		tileCoords = append(tileCoords, entry)
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
