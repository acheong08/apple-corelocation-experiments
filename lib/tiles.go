package lib

import (
	"fmt"
	"io"
	"net/http"
	"wloc/lib/mac"
	"wloc/pb"

	"google.golang.org/protobuf/proto"
)

func GetTile(tileKey int64, options ...Modifier) ([]AP, error) {
	var tileURL string
	args := newWlocArgs()
	for _, option := range options {
		option(&args)
	}
	switch args.region {
	case Options.China:
		tileURL = "https://gspe85-cn-ssl.ls.apple.com"
	case Options.International:
		tileURL = "https://gspe85-ssl.ls.apple.com"
	}
	tileURL = tileURL + "/wifi_request_tile"
	req, err := http.NewRequest("GET", tileURL, nil)
	if err != nil {
		return nil, err
	}
	for key, val := range map[string]string{
		"Accept":          "*/*",
		"Connection":      "keep-alive",
		"X-tilekey":       fmt.Sprintf("%d", tileKey),
		"User-Agent":      "geod/1 CFNetwork/1496.0.7 Darwin/23.5.0",
		"Accept-Language": "en-US,en-GB;q=0.9,en;q=0.8",
		"X-os-version":    "17.5.21F79",
	} {
		req.Header.Set(key, val)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	wifuTile := &pb.WifiTile{}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = proto.Unmarshal(b, wifuTile)
	if err != nil {
		return nil, err
	}
	aps := make([]AP, 1000)
	max := 0
	for _, region := range wifuTile.GetRegion() {
		for _, device := range region.GetDevices() {
			if device == nil || device.Bssid == 0 {
				continue
			}
			aps[max] = AP{
				BSSID: mac.Decode(device.GetBssid()),
				Location: Location{
					Lat:  coordFromInt(int64(device.GetEntry().GetLat()), -7),
					Long: coordFromInt(int64(device.GetEntry().GetLong()), -7),
				},
			}
			max++
		}
	}
	return aps[:max], nil
}
