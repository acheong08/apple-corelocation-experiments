package multilateration

import "github.com/acheong08/apple-corelocation-experiments/lib"

type Request struct {
	APs []AccessPoint `json:"wifiAccessPoints"`
}

type AccessPoint struct {
	Mac            string `json:"macAddress"`
	SignalStrength int    `json:"signalStrength"`
	Location       lib.Location
}
