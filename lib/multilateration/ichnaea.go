package multilateration

import "wloc/lib"

type Request struct {
	CellTowers Cell
}

type Cell struct {
	lib.TowerInfo
	Age            uint `json:"age"`
	SignalStrength int  `json:"signalStrength"`
	Location       lib.Location
	Score          float64
}
