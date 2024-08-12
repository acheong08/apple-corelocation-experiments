package main

import (
	"encoding/binary"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"wloc/lib"
	"wloc/pb"

	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	initialBytes []byte
	lat, lon     int64 = lib.IntFromCoord(34.074403, 8), lib.IntFromCoord(143.573894, 8)
)

func init() {
	initialBytes, _ = hex.DecodeString("0001000000010000")
}

func p64(i int) *int64 {
	i64 := int64(i)
	return &i64
}

func main() {
	e := echo.New()
	e.POST("/clls/wloc", func(c echo.Context) error {
		log.Println("Recieved request")
		body := c.Request().Body
		defer body.Close()
		b, err := io.ReadAll(body)
		if err != nil {
			return c.String(400, "failed to read body")
		}
		var p pb.AppleWLoc
		if err := proto.Unmarshal(b[57:], &p); err != nil {
			return c.String(400, "failed to parse protobuf")
		}
		for i := range len(p.GetWifiDevices()) {
			log.Println(p.WifiDevices[i].Bssid)
			p.WifiDevices[i].Location = &pb.Location{
				Latitude:                 &lat,
				Longitude:                &lon,
				HorizontalAccuracy:       p64(39),
				VerticalAccuracy:         p64(1000),
				Altitude:                 p64(530),
				UnknownValue4:            p64(3),
				MotionActivityType:       p64(63),
				MotionActivityConfidence: p64(467),
			}
		}
		p.NumCellResults = nil
		p.NumWifiResults = nil
		p.DeviceType = nil
		b, err = SerializeProto(&p, initialBytes)
		if err != nil {
			return c.String(500, "failed to encode protobuf")
		}
		log.Printf("%x\n", b)
		return c.Blob(200, "", b)
	})
	log.Fatal(http.ListenAndServe(":9090", e))
}

func SerializeProto(p protoreflect.ProtoMessage, initial []byte) ([]byte, error) {
	if p == nil {
		panic("protobuf is nil")
	}
	b, err := proto.Marshal(p)
	if err != nil {
		return nil, err
	}
	int16Len := make([]byte, 2)
	binary.BigEndian.PutUint16(int16Len, uint16(len(b)))
	if initial != nil {
		b = append(initial, append(int16Len, b...)...)
	}
	return b, nil
}
