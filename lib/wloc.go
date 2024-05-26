package lib

import (
	"bytes"
	"io"
	"log"
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
	copyMultiByte(data, []byte{0x00, 0x01, 0x00, 0x05}, []byte("en_US"), []byte{0x00, 0x13}, []byte("com.apple.locationd"), []byte{0x00, 0x0a}, []byte("14.5.23F79"), []byte{0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00}, []byte{byte(len(serializedWloc))})
	data = append(data, serializedWloc...)

	return data, nil
}

func RequestWloc(block *pb.AppleWLoc) (*pb.AppleWLoc, error) {
	// Serialize to bytes
	serializedBlock, err := serializeWlocRequest(block)
	if err != nil {
		return nil, err
	}
	// Make HTTP request
	req, _ := http.NewRequest(http.MethodPost, "https://gs-loc.apple.com/clls/wloc", bytes.NewReader(serializedBlock))
	for key, val := range map[string]string{
		"Content-Type":   "application/x-www-form-urlencoded",
		"Accept":         "*/*",
		"Accept-Charset": "utf-8",
		// "Accept-Encoding": "gzip, deflate",
		"Accept-Language": "en-us",
		"User-Agent":      "locationd/1753.17 CFNetwork/711.1.12 Darwin/14.0.0",
	} {
		req.Header.Set(key, val)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatalf("Request failed with status code %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	respBlock := pb.AppleWLoc{}
	err = proto.Unmarshal(body[10:], &respBlock)
	if err != nil {
		return nil, err
	}
	return &respBlock, nil
}

func QueryBssid(bssids []string, maxResults bool) (*pb.AppleWLoc, error) {
	zero32 := int32(0)
	one32 := int32(1)
	block := pb.AppleWLoc{}
	block.WifiDevices = make([]*pb.WifiDevice, len(bssids))
	for i, bssid := range bssids {
		block.WifiDevices[i] = &pb.WifiDevice{Bssid: bssid}
	}
	if maxResults {
		block.NumResults = &zero32
	} else {
		block.NumResults = &one32
	}
	return RequestWloc(&block)
}

func copyMultiByte(dst []byte, srcs ...[]byte) {
	n := 0
	for _, src := range srcs {
		copy(dst[n:], src)
		n += len(src)
	}
}