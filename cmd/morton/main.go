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
	var level int = 13
	cli.IntFlag("level", "Openstreetmap level", &level)
	encode := cli.NewSubCommandInheritFlags("encode", "Encode GPS coordinates to morton tile key")
	encode.Float64Flag("lat", "latitude", &lat)
	encode.Float64Flag("long", "longitude", &long)
	encode.Action(func() error {
		tileKey := morton.Encode(lat, long, level)
		fmt.Println("Tile Key: ", tileKey)
		return nil
	})
	var tileKey int64
	decode := cli.NewSubCommand("decode", "Decode morton tile key to GPS coordinates")
	decode.Int64Flag("tile", "tile key", &tileKey)
	decode.Action(func() error {
		log.Println(morton.Unpack(tileKey))
		lat, long, _ := morton.Decode(tileKey)
		fmt.Println(lat, long)
		return nil
	})

	var mLat int
	var mLong int
	pack := cli.NewSubCommand("pack", "Pack quadkey encoded coordinates to tilekey")
	pack.IntFlag("lat", "latitude", &mLat)
	pack.IntFlag("long", "longitude", &mLong)
	pack.Action(func() error {
		tileKey := morton.Pack(mLat, mLong, 13)
		fmt.Println("Tile Key: ", tileKey)
		return nil
	})

	unpack := cli.NewSubCommand("unpack", "Unpack tilekey to quadkey encoded coordinates")
	unpack.Int64Flag("tile", "tile key", &tileKey)
	unpack.Action(func() error {
		mLat, mLong, _ := morton.Unpack(tileKey)
		fmt.Println(mLat, mLong)
		return nil
	})

	err := cli.Run()
	if err != nil {
		panic(err)
	}
}
