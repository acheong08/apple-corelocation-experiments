package main

import (
	"fmt"
	"log"
	"math"
	"wloc/lib"

	"github.com/gptlang/oui/ouidb"
	"github.com/leaanthony/clir"
)

func main() {
	cli := clir.NewCli("wloc", "Retrieve BSSID geolocation using Apple's API", "v0.0.1")
	getCmd := cli.NewSubCommand("get", "Gets and displays adjacent BSSID locations given an existing BSSID")
	var bssids []string
	var less bool
	getCmd.StringsFlag("bssid", "One or more known bssid strings", &bssids)
	getCmd.BoolFlag("less", "Only return requested BSSID location", &less)
	getCmd.Action(func() error {
		if len(bssids) == 0 {
			log.Fatalln("BSSIDs cannot be empty")
		}
		blocks, err := lib.QueryBssid(bssids, !less)
		if err != nil {
			panic(err)
		}
		for _, wifi := range blocks.GetWifiDevices() {
			man, err := ouidb.Lookup(wifi.GetBssid())
			if err != nil {
				man = "Unknown"
			}
			fmt.Printf("BSSID: %s (%s) found at Lat: %f Long: %f\n", wifi.GetBssid(), man, float64(*wifi.GetLocation().Latitude)*math.Pow10(-8), float64(*wifi.GetLocation().Longitude)*math.Pow10(-8))
		}
		fmt.Println(len(blocks.GetWifiDevices()), "number of devices found in area")
		return nil
	})
	tileKey := int64(81644853)
	tileCmd := cli.NewSubCommand("tile", "Returns a list of BSSIDs and their associated GPS locations")
	tileCmd.Int64Flag("key", "The tile key used to determine region", &tileKey)
	tileCmd.Action(func() error {
		tiles, err := lib.GetTile(tileKey)
		if err != nil {
			panic(err)
		}
		for _, region := range tiles.GetRegion() {
			for _, device := range region.GetDevices() {
				lat := float64(device.GetEntry().GetLat()) * math.Pow10(-7)
				long := float64(device.GetEntry().GetLong()) * math.Pow10(-7)
				macHex := fmt.Sprintf("%x", device.GetBssid())
				if len(macHex) != 12 {
					// Fill it up to 12 with 0s in front
					for i := 0; i < 13-len(macHex); i++ {
						macHex = "0" + macHex
					}
				}
				// Insert : between every 2 hex values
				mac := ""
				for i := 0; i < len(macHex); i += 2 {
					if i+2 < len(macHex) {
						mac += macHex[i:i+2] + ":"
					} else {
						mac += macHex[i:]
					}
				}
				// manufacturer, err := ouidb.Lookup(mac)
				// if err != nil {
				// 	continue
				// }
				fmt.Printf("MAC: %s - %f %f\n", mac, lat, long)
			}
		}
		return nil
	})
	err := cli.Run()
	if err != nil {
		log.Fatal(err)
	}
}
