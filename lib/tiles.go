package lib

import (
	"fmt"
	"io"
	"net/http"
	"wloc/pb"

	"google.golang.org/protobuf/proto"
)

func GetTile(tileKey int64) (*pb.WifiTile, error) {
	req, err := http.NewRequest("GET", "https://gspe85-ssl.ls.apple.com/wifi_request_tile", nil)
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
	return wifuTile, err
}
