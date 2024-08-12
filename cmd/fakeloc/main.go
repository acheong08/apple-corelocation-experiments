package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"strings"
	"sync"
	"time"
	"wloc/lib"
	"wloc/pb"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
)

var center = orb.Point{80.289058, 9.036947}

var bssid = strings.ToLower("DE:AD:BA:BE:CA:FE")

func main() {
	maxDistance := 400.0
	minDistance := 0.0
	for i := minDistance; i < maxDistance; i += 50 {
		rssi := int32(((i / maxDistance) * 100.0) * -1)
		log.Printf("Requesting with distance %f and rssi %d", i, rssi)
		requestWithRadiusAndRssi(i, rssi)
	}
}

// Create polygon around point with radius x
func requestWithRadiusAndRssi(radius float64, rssi int32) {
	poly := createCircularPolygon(center, radius, 10)
	wifiEntries := make([]*pb.PbcWifiEntry, len(poly[0]))
	if len(poly) != 1 {
		panic("multiple rings formed")
	}
	for i, point := range poly[0] {
		t := CocoaTimestamp(time.Now().Add(time.Duration(i)))
		wifiEntries[i] = &pb.PbcWifiEntry{
			Bssid:   bssid,
			Channel: 100,
			Rssi:    rssi,
			Location: &pb.PbcWlocLocation{
				Latitude:                           point.Lat(),
				Longitude:                          point.Lon(),
				HorizontalAccuracy:                 float32(rand.NormFloat64()*5 + 3),
				Altitude:                           float32(rand.NormFloat64()*2 + 4),
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
					Activity:   pb.MotionActivity_walking,
				},

				DominantMotionActivity: &pb.MotionActivity{
					Confidence: 3,
					Activity:   pb.MotionActivity_walking,
				},
			},
			Timestamp: t, // Todo
			ScanType:  2,
		}
	}
	wait := sync.WaitGroup{}
	for _, entry := range wifiEntries {
		wait.Add(1)
		go func() {
			defer wait.Done()
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
		}()
	}
	wait.Wait()
}

func createCircularPolygon(center orb.Point, radius float64, numSides int) orb.Polygon {
	var polygon orb.Polygon
	var ring orb.Ring

	for i := 0; i < numSides; i++ {
		angle := float64(i) * (2 * math.Pi / float64(numSides))
		point := geo.PointAtBearingAndDistance(center, angle, rand.NormFloat64()*3+radius)
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
