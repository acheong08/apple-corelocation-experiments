package main

import (
	"encoding/json"
	"fmt"
	"wloc/lib"
	"wloc/pb"
)

var (
	fakeMac      = "CC:D0:83:A1:6A:A3"
	fakeLat      = int64(5158125200)
	fakeLong     = int64(-298583100)
	fakeHA       = int64(23)
	fakeUnknown4 = int64(3)
	fakeAlt      = int64(17)
	fakeVA       = int64(10)
	fakeMAT      = int64(63)
	fakeMAC      = int64(207)
	one64        = int64(1)
	// one32        = int32(1)
	zero32 = int32(0)
)

// Push a fake MAC address onto Apple's WPS database
func main() {
	applWloc := pb.AppleWLoc{
		UnknownValue0: &one64,
		UnknownValue1: &zero32,
		// NumResults:    &zero32,
	}
	applWloc.WifiDevices = make([]*pb.WifiDevice, 1)
	fakeDevice := &pb.WifiDevice{Bssid: fakeMac}
	fakeDevice.Location = &pb.WifiDevice_Location{
		Latitude:                 &fakeLat,
		Longitude:                &fakeLong,
		HorizontalAccuracy:       &fakeHA,
		UnknownValue4:            &fakeUnknown4,
		Altitude:                 &fakeAlt,
		VerticalAccuracy:         &fakeVA,
		MotionActivityType:       &fakeMAT,
		MotionActivityConfidence: &fakeMAC,
	}
	applWloc.WifiDevices[0] = fakeDevice
	respWloc, err := lib.RequestWloc(&applWloc)
	if err != nil {
		panic(err)
	}
	jsonBytes, _ := json.MarshalIndent(respWloc, "", " ")
	fmt.Println(string(jsonBytes))
}
