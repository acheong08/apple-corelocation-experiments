package lib

type AP struct {
	BSSID string
	Location Location
}

type Cell struct {
	Tower TowerInfo
	Location Location
}

type TowerInfo struct {
	Mmc, Mnc, CellId, TacId uint32
}

type Location struct {
	Long, Lat, Alt float64
}
