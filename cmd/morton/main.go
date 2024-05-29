package main

import (
	"fmt"
	"log"
	"wloc/lib/morton"

	"github.com/leaanthony/clir"
)

func main() {
	cli := clir.NewCli("morton", "Convert between GPS and morton encoded tile key coordinates", "0.0.1")
	var lat float64
	var long float64
	encode := cli.NewSubCommand("encode", "Encode GPS coordinates to morton tile key")
	encode.Float64Flag("lat", "latitude", &lat)
	encode.Float64Flag("long", "longitude", &long)
	encode.Action(func() error {
		tileKey := morton.Encode(lat, long)
		fmt.Println("Tile Key: ", tileKey)
		return nil
	})
	var tileKey int64
	decode := cli.NewSubCommand("decode", "Decode morton tile key to GPS coordinates")
	decode.Int64Flag("tile", "tile key", &tileKey)
	decode.Action(func() error {
		log.Println(morton.UnpackAppleCoord(tileKey))
		lat, long := morton.Decode(tileKey)
		fmt.Println(lat, long)
		return nil
	})
	experiment := cli.NewSubCommand("experiment", "Experimental decoding")
	experiment.Int64Flag("tile", "tile key", &tileKey)
	experiment.Action(func() error {
		level := 0
		row := 0
		column := 0
		quadKey := tileKey
		for quadKey > 1 {
			mask := 1 << level
			if quadKey&0x1 != 0 {
				column |= mask
			}
			if quadKey&0x2 != 0 {
				row |= mask
			}
			level++
			quadKey = (quadKey - (quadKey & 0x3)) / 4
		}
		fmt.Println(row, column, level)
		return nil
	})
	err := cli.Run()
	if err != nil {
		panic(err)
	}
}
