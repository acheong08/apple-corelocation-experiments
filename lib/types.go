package lib

type AP struct {
	BSSID    string
	Location Location
}

type Cell struct {
	Tower    TowerInfo
	Location Location
}

type TowerInfo struct {
	Mmc    uint32 `json:"mobileCountryCode"`
	Mnc    uint32 `json:"mobileNetworkCode"`
	CellId uint32 `json:"cellId"`
	TacId  uint32 `json:"locationAreaCode"`
}

type Location struct {
	Long, Lat, Alt float64
}
