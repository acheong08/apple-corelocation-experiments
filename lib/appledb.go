package lib

import (
	"bytes"
	"geostalk/proto"
	"io"
	"log"
	"net/http"

	protobuf "google.golang.org/protobuf/proto"
)

func QueryBssid(bssids []string, moreResults bool) proto.AppleWLoc {
	block := proto.AppleWLoc{}
	block.WifiDevices = make([]*proto.WifiDevice, len(bssids))
	for i, bssid := range bssids {
		block.WifiDevices[i] = &proto.WifiDevice{Bssid: &bssid}
	}
	zero32 := int32(0)
	one32 := int32(1)
	block.UnknownValue1 = &zero32
	if moreResults {
		block.ReturnSingleResult = &zero32
	} else {
		block.ReturnSingleResult = &one32
	}
	// Serialize to bytes
	serializedBlock, err := protobuf.Marshal(&block)
	if err != nil {
		panic(err)
	}
	data := []byte{0x00, 0x01, 0x00, 0x05}
	data = append(data, []byte("en_US")...)
	data = append(data, 0x00, 0x13)
	data = append(data, []byte("com.apple.locationd")...)
	data = append(data, 0x00, 0x0a)
	data = append(data, []byte("8.1.12B411")...)
	data = append(data, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00)
	data = append(data, byte(len(serializedBlock)))
	data = append(data, serializedBlock...)
	// Make HTTP request
	req, _ := http.NewRequest(http.MethodPost, "https://gs-loc.apple.com/clls/wloc", bytes.NewReader(data))
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
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatalf("Request failed with status code %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	block = proto.AppleWLoc{}
	err = protobuf.Unmarshal(body[10:], &block)
	if err != nil {
		panic(err)
	}
	return block
}
