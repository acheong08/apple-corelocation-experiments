//go:build !js && !wasm

package lib

var (
	wlocArpcRequest    ArpcRequest
	pbcWlocArpcRequest ArpcRequest
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
	wlocArpcRequest = ArpcRequest{
		Version:       "1",
		Locale:        "en-001_001",
		AppIdentifier: "com.apple.locationd",
		OsVersion:     "18.6.2.22G100",
		FunctionId:    1,
		Payload:       []byte{},
	}

	pbcWlocArpcRequest = ArpcRequest{
		Version:       "1",
		Locale:        "en-001_001",
		AppIdentifier: "com.apple.locationd",
		OsVersion:     "17.4.1.21E236",
		FunctionId:    100,
		Payload:       []byte{},
	}
}
