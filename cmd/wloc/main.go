package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"wloc/lib"
	"wloc/pb"

	"github.com/gptlang/oui/ouidb"
	"github.com/leaanthony/clir"
)

func main() {
	var china bool
	cli := clir.NewCli("wloc", "Retrieve BSSID geolocation using Apple's API", "v0.0.1")
	cli.BoolFlag("china", "Use the China region for the request", &china)
	var displayVendor bool
	getCmd := cli.NewSubCommandInheritFlags("get", "Gets and displays adjacent BSSID locations given an existing BSSID")
	var bssids []string
	var less bool
	getCmd.StringsFlag("bssid", "One or more known bssid strings", &bssids)
	getCmd.BoolFlag("less", "Only return requested BSSID location", &less)
	getCmd.BoolFlag("vendor", "Tells the CLI to append the vendor of the MAC address to outpus", &displayVendor)
	getCmd.Action(func() error {
		if len(bssids) == 0 {
			log.Fatalln("BSSIDs cannot be empty")
		}
		var options []lib.Modifier
		if china {
			options = append(options, lib.Options.WithRegion(lib.Options.China))
		}

		blocks, err := lib.QueryBssid(bssids, !less, options...)
		if err != nil {
			panic(err)
		}
		for _, wifi := range blocks.GetWifiDevices() {
			if displayVendor {
				man, err := ouidb.Lookup(wifi.GetBssid())
				if err != nil {
					man = "Unknown"
				}
				fmt.Printf("BSSID: %s (%s) found at Lat: %f Long: %f\n", wifi.GetBssid(), man, float64(*wifi.GetLocation().Latitude)*math.Pow10(-8), float64(*wifi.GetLocation().Longitude)*math.Pow10(-8))
			} else {
				fmt.Printf("BSSID: %s found at Lat: %f Long: %f\n", wifi.GetBssid(), float64(*wifi.GetLocation().Latitude)*math.Pow10(-8), float64(*wifi.GetLocation().Longitude)*math.Pow10(-8))
			}
		}
		fmt.Println(len(blocks.GetWifiDevices()), "number of devices found in area")
		return nil
	})
	tileKey := int64(81644853)
	tileCmd := cli.NewSubCommandInheritFlags("tile", "Returns a list of BSSIDs and their associated GPS locations")
	tileCmd.Int64Flag("key", "The tile key used to determine region", &tileKey)
	tileCmd.BoolFlag("vendor", "Tells the CLI to append the vendor of the MAC address to outpus", &displayVendor)
	tileCmd.Action(func() error {
		var options []lib.Modifier
		if china {
			options = append(options, lib.Options.WithRegion(lib.Options.China))
		}
		tiles, err := lib.GetTile(tileKey, options...)
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
				if displayVendor {
					manufacturer, err := ouidb.Lookup(mac)
					if err != nil {
						continue
					}
					fmt.Printf("MAC: %s (%s) - %f %f\n", mac, manufacturer, lat, long)
				} else {
					fmt.Printf("MAC: %s - %f %f\n", mac, lat, long)
				}
			}
		}
		return nil
	})
	experiment := cli.NewSubCommandInheritFlags("exp", "Experimental command for WLOC requests")
	experiment.Action(func() error {
		zero := int32(0)
		negative1 := int32(-1)
		block := pb.AppleWLoc{
			NumCellResults: &zero,
			NumWifiResults: &negative1,
			CellTowerRequest: &pb.CellTower{
				Mmc:    502,
				Mnc:    16,
				CellId: 9603074,
				TacId:  43337,
			},
			DeviceType: &pb.DeviceType{
				OperatingSystem: "iPhone OS17.5/21F79",
				Model:           "iPhone12,1",
			},
		}
		resp, err := lib.RequestWloc(&block, nil)
		if err != nil {
			panic(err)
		}
		b, _ := json.MarshalIndent(resp, "", " ")
		fmt.Println(string(b))
		return nil
	})
	err := cli.Run()
	if err != nil {
		log.Fatal(err)
	}
}
