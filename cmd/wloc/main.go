package main

import (
	"encoding/json"
	"fmt"
	"log"
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

		blocks, err := lib.QueryBssid(bssids, boolToInt[int32](less), options...)
		if err != nil {
			panic(err)
		}
		for _, ap := range blocks {
			if displayVendor {
				man, err := ouidb.Lookup(ap.BSSID)
				if err != nil {
					man = "Unknown"
				}
				fmt.Printf("BSSID: %s (%s) found at Lat: %f Long: %f\n", ap.BSSID, man, ap.Location.Lat, ap.Location.Long)
			} else {
				fmt.Printf("BSSID: %s found at Lat: %f Long: %f\n", ap.BSSID, ap.Location.Lat, ap.Location.Long)
			}
		}
		fmt.Println(len(blocks), "number of devices found in area")
		return nil
	})
	tileKey := int64(81644853)
	tileCmd := cli.NewSubCommandInheritFlags("tile", "Returns a list of BSSIDs and their associated GPS locations")
	tileCmd.Int64Flag("key", "The tile key used to determine region", &tileKey)
	tileCmd.BoolFlag("vendor", "Tells the CLI to append the vendor of the MAC address to outpus", &displayVendor)
	tileCmd.Action(func() error {
		tiles, err := lib.GetTile(tileKey)
		if err != nil {
			panic(err)
		}
		for _, d := range tiles {
			if displayVendor {
				manufacturer, err := ouidb.Lookup(d.BSSID)
				if err != nil {
					continue
				}
				fmt.Printf("MAC: %s (%s) - %f %f\n", d.BSSID, manufacturer, d.Location.Lat, d.Location.Long)
			} else {
				fmt.Printf("MAC: %s - %f %f\n", d.BSSID, d.Location.Lat, d.Location.Long)
			}
		}
		return nil
	})
	experiment := cli.NewSubCommandInheritFlags("exp", "Experimental command for WLOC requests")
	var mcc, mnc, cellid, tacid uint32
	experiment.Uint32Flag("mcc", "Mobile Country Code", &mcc)
	experiment.Uint32Flag("mnc", "Mobile Network Code", &mnc)
	experiment.Uint32Flag("cellid", "Cell ID", &cellid)
	experiment.Uint32Flag("tacid", "Tracking Area Code", &tacid)
	experiment.Action(func() error {
		zero := int32(0)
		block := pb.AppleWLoc{
			NumCellResults: &zero,
			CellTowerRequest: &pb.CellTower{
				Mcc:    mcc,
				Mnc:    mnc,
				CellId: cellid,
				TacId:  tacid,
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

type anyNum interface {
	int | int32 | int64 | uint | uint32 | uint64 | float32 | float64
}

func boolToInt[T anyNum](b bool) T {
	if b {
		return 1
	}
	return 0
}
