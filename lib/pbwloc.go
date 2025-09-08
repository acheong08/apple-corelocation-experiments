package lib

import (
	"bytes"
	"errors"
	"net/http"
	"wloc/pb"
)

func RequestPbcWloc(p *pb.PbcWlocRequest) error {
	b, err := SerializeProto(p, pbcWlocArpcRequest)
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
	if resp.StatusCode != 200 {
		return errors.New("server replied with non-200 status code")
	}
	return nil
}
