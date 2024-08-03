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
	"google.golang.org/protobuf/reflect/protoreflect"
)

func serializeProto(p protoreflect.ProtoMessage, initial []byte) ([]byte, error) {
	if p == nil {
		panic("protobuf is nil")
	}
	b, err := proto.Marshal(p)
	if err != nil {
		return nil, err
	}
	if initial != nil {
		b = append(initial, append([]byte{byte(len(b))}, b...)...)
	}
	return b, nil
}

func RequestWloc(block *pb.AppleWLoc, options ...Modifier) (*pb.AppleWLoc, error) {
	args := newWlocArgs()
	if len(options) != 0 {
		for _, option := range options {
			if option != nil {
				option(&args)
			}
		}
	}
	// Serialize to bytes
	serializedBlock, err := serializeProto(block, initialWlocBytes)
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
	for key, val := range headers {
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

var zero int32

func QueryBssid(bssids []string, maxResults int32, options ...Modifier) ([]AP, error) {
	block := &pb.AppleWLoc{
		NumCellResults: &zero,
		DeviceType: &pb.DeviceType{
			OperatingSystem: "iPhone OS17.5/21F79",
			Model:           "iPhone12,1",
		},
	}
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
	i := 0
	for _, d := range block.GetWifiDevices() {
		long := CoordFromInt(d.GetLocation().GetLongitude(), -8)
		lat := CoordFromInt(d.GetLocation().GetLatitude(), -8)
		alt := CoordFromInt(d.GetLocation().GetAltitude(), -8)
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
		i++
	}
	resp = resp[:i]
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
				Long: CoordFromInt(c.GetLocation().GetLongitude(), -8),
				Lat:  CoordFromInt(c.GetLocation().GetLatitude(), -8),
				Alt:  CoordFromInt(c.GetLocation().GetAltitude(), -8),
			},
		}
	}
	return cells, nil
}

func CoordFromInt(n int64, pow int) float64 {
	return float64(n) * math.Pow10(pow)
}

func IntFromCoord(coord float64, pow int) int64 {
	return int64(coord * math.Pow10(pow))
}
