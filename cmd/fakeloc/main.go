package main

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
	"wloc/lib"
	"wloc/pb"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
)

func main() {
	bssid := strings.ToLower("DE:AD:FE:D5:BA:BE")
	lat, lon := 9.036947, 80.289058
	center := orb.Point{lon, lat}
	// Create polygon around point with radius x
	poly := createCircularPolygon(center, 6)
	wifiEntries := make([]*pb.PbcWifiEntry, len(poly[0]))
	if len(poly) != 1 {
		panic("multiple rings formed")
	}
	for i, point := range poly[0] {
		t := CocoaTimestamp(time.Now().Add(-time.Duration(i) * time.Hour))
		wifiEntries[i] = &pb.PbcWifiEntry{
			Bssid:   bssid,
			Channel: 161,
			Rssi:    -30,
			Location: &pb.PbcWlocLocation{
				Latitude:                           point.Lat(),
				Longitude:                          point.Lon(),
				HorizontalAccuracy:                 float32(rand.NormFloat64()*5 + 3),
				Altitude:                           30,
				VerticalAccuracy:                   float32(rand.NormFloat64()*6 + 5),
				Timestamp:                          t,
				Provider:                           1,
				MotionVehicleConnected:             0,
				MotionVehicleConnectedStateChanged: 0,
				RawMotionActivity: &pb.MotionActivity{
					Confidence: 2,
					Activity:   pb.MotionActivity_walking,
				},

				MotionActivity: &pb.MotionActivity{
					Confidence: 2,
					Activity:   pb.MotionActivity_unknown,
				},

				DominantMotionActivity: &pb.MotionActivity{
					Confidence: 3,
					Activity:   pb.MotionActivity_walking,
				},
			},
			Hidden:    1,
			Timestamp: t, // Todo
			ScanType:  5,
		}
	}
	for _, entry := range wifiEntries {
		if err := lib.RequestPbcWloc(&pb.PbcWlocRequest{
			WifiEntries: []*pb.PbcWifiEntry{entry},
			DeviceInfo: &pb.DeviceType{
				OperatingSystem: "N104AP",
				Model:           "iPhone OS17.5.1/21F90",
			},
		}); err != nil {
			panic(err)
		}
		fmt.Println("Requested...")
	}
}

func createCircularPolygon(center orb.Point, radius float64) orb.Polygon {
	const numSides = 10
	var polygon orb.Polygon
	var ring orb.Ring

	for i := 0; i < numSides; i++ {
		angle := float64(i) * (2 * math.Pi / numSides)
		point := geo.PointAtBearingAndDistance(center, angle, radius)
		ring = append(ring, orb.Point{point.Lon(), point.Lat()})
	}

	// Close the ring
	ring = append(ring, ring[0])
	polygon = append(polygon, ring)

	return polygon
}

// https://github.com/cgerro/ios-location-trace-study/blob/90f60ac797c2fc541a6b6dcf8ef1c43f669e05d9/utils.py#L237
func CocoaTimestamp(t time.Time) float64 {
	coreDataStartDate := time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
	return t.Sub(coreDataStartDate).Seconds()
}
