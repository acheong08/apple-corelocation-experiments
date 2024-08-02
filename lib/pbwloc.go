package lib

import (
	"bytes"
	"net/http"
	"wloc/pb"
)

func RequestPbcWloc(p *pb.PbcWlocRequest) error {
	b, err := serializeProto(p, initialPbcWlocBytes)
	if err != nil {
		return err
	}
	req, _ := http.NewRequest(http.MethodPost, "https://gsp10-ssl.apple.com/hcy/pbcwloc", bytes.NewReader(b))
	for key, val := range headers {
		req.Header.Set(key, val)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
