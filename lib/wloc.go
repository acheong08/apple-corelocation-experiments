package lib

import (
	"bytes"
	"errors"
	"io"
	"log"
	"math"
	"net/http"
	"wloc/pb"

	"google.golang.org/protobuf/proto"
)

func serializeWlocRequest(applWloc *pb.AppleWLoc) ([]byte, error) {
	if applWloc == nil {
		panic("nil pointer error")
	}
	serializedWloc, err := proto.Marshal(applWloc)
	if err != nil {
		return nil, err
	}
	data := make([]byte, 50)
	copyMultiByte(data, []byte{0x00, 0x01, 0x00, 0x05}, []byte("en_US"), []byte{0x00, 0x13}, []byte("com.apple.locationd"), []byte{0x00, 0x0a}, []byte("17.5.21F79"), []byte{0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00}, []byte{byte(len(serializedWloc))})
	data = append(data, serializedWloc...)

	return data, nil
}

func RequestWloc(block *pb.AppleWLoc, options ...Modifier) (*pb.AppleWLoc, error) {
	args := newWlocArgs()
	if len(options) != 0 {
		for _, option := range options {
			option(&args)
		}
	}
	// Serialize to bytes
	serializedBlock, err := serializeWlocRequest(block)
	if err != nil {
		return nil, errors.New("failed to serialize protobuf")
	}
	var wlocURL string = "https://gs-loc.apple.com"
	switch args.region {
	case Options.China:
		log.Println("Using China API")
		wlocURL = "https://gs-loc-cn.apple.com"
	}
	wlocURL = wlocURL + "/clls/wloc"
	// Make HTTP request
	req, _ := http.NewRequest(http.MethodPost, wlocURL, bytes.NewReader(serializedBlock))
	for key, val := range map[string]string{
		"Content-Type":   "application/x-www-form-urlencoded",
		"Accept":         "*/*",
		"Accept-Charset": "utf-8",
		// "Accept-Encoding": "gzip, deflate",
		"Accept-Language": "en-us",
		"User-Agent":      "locationd/1753.17 CFNetwork/711.1.12 Darwin/14.0.0",
		jsHeader:          jsHeaderValue,
	} {
		req.Header.Set(key, val)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.New("failed to make request")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		if resp.StatusCode == 0 {
			return nil, errors.New("cors issue probably")
		}
		return nil, errors.New(http.StatusText(resp.StatusCode))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("failed to read response body")
	}
	respBlock := pb.AppleWLoc{}
	err = proto.Unmarshal(body[10:], &respBlock)
	if err != nil {
		return nil, errors.New("failed to unmarshal response protobuf")
	}
	return &respBlock, nil
}

func QueryBssid(bssids []string, maxResults int32, options ...Modifier) ([]AP, error) {
	block := &pb.AppleWLoc{}
	block.WifiDevices = make([]*pb.WifiDevice, len(bssids))
	for i, bssid := range bssids {
		block.WifiDevices[i] = &pb.WifiDevice{Bssid: bssid}
	}
	block.NumWifiResults = &maxResults
	block, err := RequestWloc(block, options...)
	if err != nil {
		return nil, err
	}
	resp := make([]AP, len(block.GetWifiDevices()))
	for i, d := range block.GetWifiDevices() {
		long := coordFromInt(d.GetLocation().GetLongitude(), -8)
		lat := coordFromInt(d.GetLocation().GetLatitude(), -8)
		alt := coordFromInt(d.GetLocation().GetAltitude(), -8)
		if long == -180 && lat == -180 {
			continue
		}
		resp[i] = AP{
			BSSID: d.GetBssid(),
			Location: Location{
				Long: long,
				Lat:  lat,
				Alt:  alt,
			},
		}
	}
	return resp, nil
}

func QueryCell(mmc, mnc, cellid, tacid uint32, numResults int32, options ...Modifier) ([]Cell, error) {
	block := &pb.AppleWLoc{
		NumCellResults: &numResults,
		CellTowerRequest: &pb.CellTower{
			Mmc:    mmc,
			Mnc:    mnc,
			CellId: cellid,
			TacId:  tacid,
		},
		DeviceType: &pb.DeviceType{
			OperatingSystem: "iPhone OS17.5/21F79",
			Model:           "iPhone12,1",
		},
	}
	block, err := RequestWloc(block, options...)
	if err != nil {
		return nil, err
	}
	cells := make([]Cell, len(block.GetCellTowerResponse()))
	for i, c := range block.GetCellTowerResponse() {
		cells[i] = Cell{
			Tower: TowerInfo{
				Mmc:    c.GetMmc(),
				Mnc:    c.GetMnc(),
				CellId: c.GetCellId(),
				TacId:  c.GetTacId(),
			},
			Location: Location{
				Long: coordFromInt(c.GetLocation().GetLongitude(), -8),
				Lat:  coordFromInt(c.GetLocation().GetLatitude(), -8),
				Alt:  coordFromInt(c.GetLocation().GetAltitude(), -8),
			},
		}
	}
	return cells, nil
}

func copyMultiByte(dst []byte, srcs ...[]byte) {
	n := 0
	for _, src := range srcs {
		copy(dst[n:], src)
		n += len(src)
	}
}

func coordFromInt(n int64, pow int) float64 {
	return float64(n) * math.Pow10(pow)
}
