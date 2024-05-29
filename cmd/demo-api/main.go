package main

import (
	"fmt"
	"log"
	"math"
	"wloc/lib"
	"wloc/lib/distance"
	"wloc/lib/mac"
	"wloc/lib/morton"
	"wloc/lib/spiral"
	"wloc/pb"

	_ "embed"

	"github.com/labstack/echo/v4"
)

//go:embed static/index.html
var index string

func init() {
	log.SetFlags(log.Lshortfile)
}

type gps struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

type tileCoords struct {
	Coord  []float64 `json:"coord"`
	Morton []int     `json:"morton"`
}

var tileCache = make([]tileCoords, 0)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.HTML(200, index)
	})
	e.GET("/cache", func(c echo.Context) error {
		return c.JSON(200, tileCache)
	})
	e.POST("/gps", func(c echo.Context) error {
		var g gps
		if err := c.Bind(&g); err != nil {
			return c.String(400, "Bad Request")
		}
		if g.Lat < -90 || g.Lat > 90 || g.Long < -180 || g.Long > 180 || g.Lat == 0 || g.Long == 0 {
			return c.String(400, "Bad Request")
		}
		respPoints := make([]distance.Point, 0)
		mLat, mLong := morton.PredictAppleCoord(g.Lat, g.Long)
		sp := spiral.NewSpiral(mLat, mLong)
		var tile *pb.WifiTile
		var err error
		var closest distance.Point
		var points []distance.Point
		for i := 0; i < 20; i++ {
			mLat, mLong = sp.Next()
			lat, long := morton.Decode(morton.PackAppleCoord(mLat, mLong))
			respPoints = append(respPoints, distance.Point{
				Id: fmt.Sprintf("Spiral %d", i),
				Y:  lat,
				X:  long,
			})
			log.Println(lat, long)
			tile, err = lib.GetTile(morton.PackAppleCoord(mLat, mLong))
			if err != nil {
				tile = nil
				log.Println(err)
				continue
			}
			if tile == nil {
				return c.String(500, "Internal Server Error")
			}
			points = make([]distance.Point, 0)
			var avgLat, avgLong int32
			var avgCount int
			for _, r := range tile.GetRegion() {
				for _, d := range r.GetDevices() {
					if d == nil || d.GetBssid() == 0 {
						continue
					}
					points = append(points, distance.Point{
						Id: mac.Decode(d.GetBssid()),
						Y:  float64(d.GetEntry().GetLat()) * math.Pow10(-7),
						X:  float64(d.GetEntry().GetLong()) * math.Pow10(-7),
					})
					avgLat += d.GetEntry().GetLat()
					avgLong += d.GetEntry().GetLong()
					avgCount++
				}
			}
			avgLat /= int32(avgCount)
			avgLong /= int32(avgCount)
			tileCache = append(tileCache, tileCoords{
				Coord:  []float64{float64(avgLat) * math.Pow10(-7), float64(avgLong) * math.Pow10(-7)},
				Morton: []int{int(mLat), int(mLong)},
			})

			closest = distance.Closest(distance.Point{
				Id: "click",
				Y:  g.Lat,
				X:  g.Long,
			}, []distance.Point{closest, distance.Closest(distance.Point{
				Id: "click",
				Y:  g.Lat,
				X:  g.Long,
			}, points)})
			if i >= 9 {
				break
			}
		}
		if tile == nil {
			return c.String(500, "Internal Server Error")
		}

		// Try to get closer via the wloc API
		for {
			devices, err := lib.QueryBssid([]string{closest.Id}, true)
			if err != nil {
				log.Println(closest)
				log.Println("Failed to find BSSID", closest.Id)
				return c.String(500, "Internal Server Error")
			}
			if len(devices.GetWifiDevices()) == 0 {
				log.Println("Could not find given BSSID")
				return c.String(500, "Internal Server Error")
			}
			points = make([]distance.Point, len(devices.GetWifiDevices()))
			for i, device := range devices.GetWifiDevices() {
				points[i] = distance.Point{
					Id: device.GetBssid(),
					Y:  float64(*device.GetLocation().Latitude) * math.Pow10(-8),
					X:  float64(*device.GetLocation().Longitude) * math.Pow10(-8),
				}
			}
			newClosest := distance.Closest(distance.Point{
				Id: "TILE",
				Y:  g.Lat,
				X:  g.Long,
			}, points)
			if newClosest.Id == closest.Id {
				break
			}
			closest = newClosest
		}

		respPoints = append(respPoints, points...)
		return c.JSON(200, map[string]any{
			"closest": closest,
			"points":  respPoints,
		})
	})
	e.Logger.Fatal(e.Start("127.0.0.1:1974"))
}
