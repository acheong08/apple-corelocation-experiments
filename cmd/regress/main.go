package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

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
	fmt.Printf("Latitude formula: \n%v\nAccuracy: %f\n", latReg.Formula, latReg.R2)
	fmt.Printf("Longitude formula: \n%v\nAccuracy: %f\n", longReg.Formula, longReg.R2)
	predictMorton(latReg, longReg, 33.744890516666665, -116.1865365)
}

func predictMorton(latReg, longReg *regression.Regression, lat, long float64) {
	latMor, _ := latReg.Predict([]float64{lat})
	longMor, _ := longReg.Predict([]float64{long})
	log.Println(latMor, longMor)
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
