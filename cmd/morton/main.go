package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"

	"github.com/gojuno/go.morton"
	"github.com/leaanthony/clir"
	"github.com/sajari/regression"
)

var (
	lats  regression.DataPoints
	longs regression.DataPoints
)

func main() {
	// Load jsonl
	f, err := os.Open("./tileCoordsToGPS.jsonl")
	if err != nil {
		panic(err)
	}
	// For each line, load json
	dec := json.NewDecoder(f)

	// fmt.Printf("Latitude formula: \n%v\nAccuracy: %f\n", latReg.Formula, latReg.R2)
	// fmt.Printf("Longitude formula: \n%v\nAccuracy: %f\n", longReg.Formula, longReg.R2)
	cli := clir.NewCli("morton", "Convert between GPS and morton encoded tile key coordinates", "0.0.1")
	var lat float64
	var long float64
	encode := cli.NewSubCommand("encode", "Encode GPS coordinates to morton tile key")
	encode.Float64Flag("lat", "latitude", &lat)
	encode.Float64Flag("long", "longitude", &long)
	encode.Action(func() error {
		for {
			var entry tileCoordEntry
			if err := dec.Decode(&entry); err != nil {
				if err == io.EOF {
					break
				}
				panic(err)
			}
			// Add to regression
			lats = append(lats, regression.DataPoint(float64(entry.Morton[0]), []float64{entry.GPS[0]}))
			longs = append(longs, regression.DataPoint(float64(entry.Morton[1]), []float64{entry.GPS[1]}))
		}
		latReg := trainRegression(lats)
		longReg := trainRegression(longs)
		latMor, _ := latReg.Predict([]float64{lat})
		longMor, _ := longReg.Predict([]float64{long})
		fmt.Println(latMor, longMor)
		m := morton.Make64(2, 32)
		fmt.Println("Tile Key: ", m.Pack(uint64(math.Round(longMor)), uint64(math.Round(latMor))))
		return nil
	})
	var tileKey int64
	decode := cli.NewSubCommand("decode", "Decode morton tile key to GPS coordinates")
	decode.Int64Flag("tile", "tile key", &tileKey)
	decode.Action(func() error {
		for {
			var entry tileCoordEntry
			if err := dec.Decode(&entry); err != nil {
				if err == io.EOF {
					break
				}
				panic(err)
			}
			// Add to regression
			lats = append(lats, regression.DataPoint(float64(entry.GPS[0]), []float64{float64(entry.Morton[0])}))
			longs = append(longs, regression.DataPoint(float64(entry.GPS[1]), []float64{float64(entry.Morton[1])}))
		}
		latReg := trainRegression(lats)
		longReg := trainRegression(longs)
		m := morton.Make64(2, 32)
		longLat := m.Unpack(tileKey)
		long, _ := longReg.Predict([]float64{float64(longLat[0])})
		lat, _ := latReg.Predict([]float64{float64(longLat[1])})
		fmt.Println(lat, long)
		return nil
	})
	err = cli.Run()
	if err != nil {
		panic(err)
	}
}

func trainRegression(data regression.DataPoints) *regression.Regression {
	r := new(regression.Regression)
	r.SetObserved("GPS")
	r.SetVar(0, "Morton")
	for _, point := range data {
		r.Train(point)
	}
	err := r.Run()
	if err != nil {
		panic(err)
	}
	return r
}

type tileCoordEntry struct {
	GPS    []float64 `json:"coord"`
	Morton []int     `json:"morton"`
}
