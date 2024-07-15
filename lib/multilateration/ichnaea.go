package multilateration

import "wloc/lib"

type Request struct {
	APs []AccessPoint `json:"wifiAccessPoints"`
}

type AccessPoint struct {
	Mac            string `json:"macAddress"`
	Age            uint   `json:"age"`
	SignalStrength int    `json:"signalStrength"`
	Location       lib.Location
	Score          float64
}
