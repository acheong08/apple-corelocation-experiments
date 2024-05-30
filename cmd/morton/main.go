package main

import (
	"fmt"
	"log"
	"wloc/lib/morton"

	"github.com/buckhx/tiles"
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
		log.Println(morton.Unpack(tileKey))
		lat, long := morton.Decode(tileKey)
		fmt.Println(lat, long)
		return nil
	})
	experiment := cli.NewSubCommand("experiment", "Experimental encoding")
	experiment.Float64Flag("lat", "latitude", &lat)
	experiment.Float64Flag("long", "longitude", &long)
	experiment.Action(func() error {
		t := tiles.FromCoordinate(lat, long, 13)
		p := t.ToPixel()
		t2, _ := p.ToTile()
		tileKey := morton.Pack(t2.Y, t2.X)
		fmt.Println(tileKey)
		return nil
	})
	err := cli.Run()
	if err != nil {
		panic(err)
	}
}
