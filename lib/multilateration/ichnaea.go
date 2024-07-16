package multilateration

import "wloc/lib"

type Request struct {
	APs []AccessPoint `json:"wifiAccessPoints"`
}

type AccessPoint struct {
	Mac            string `json:"macAddress"`
	SignalStrength int    `json:"signalStrength"`
	Location       lib.Location
}
