//go:build !js && !wasm

package lib

import "encoding/hex"

var (
	initialWlocBytes    []byte
	initialPbcWlocBytes []byte
)

var headers = map[string]string{
	"Content-Type":   "application/x-www-form-urlencoded",
	"Accept":         "*/*",
	"Accept-Charset": "utf-8",
	// "Accept-Encoding": "gzip, deflate",
	"Accept-Language": "en-us",
	"User-Agent":      "locationd/2890.16.16 CFNetwork/1496.0.7 Darwin/23.5.0",
}

func init() {
	var err error
	if initialWlocBytes, err = hex.DecodeString("0001000a656e2d3030315f3030310013636f6d2e6170706c652e6c6f636174696f6e64000c31372e352e312e323146393000000001000000"); err != nil {
		panic(err)
	}
	if initialPbcWlocBytes, err = hex.DecodeString("0001000a656e2d3030315f3030310013636f6d2e6170706c652e6c6f636174696f6e64000d31372e342e312e323145323336000000640000"); err != nil {
		panic(err)
	}
}
