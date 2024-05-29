package main

import (
	"log"
	"math"
	"strconv"
	"wloc/lib"
	"wloc/lib/distance"
	"wloc/lib/mac"
	"wloc/lib/morton"
	"wloc/lib/spiral"
	"wloc/pb"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.POST("/gps", func(c echo.Context) error {
		sLat := c.FormValue("lat")
		sLong := c.FormValue("long")
		if sLat == "" || sLong == "" {
			return c.String(400, "Bad Request")
		}
		// Parse as float64
		lat, err := strconv.ParseFloat(sLat, 64)
		if err != nil {
			return c.String(400, "Bad Request")
		}
		long, err := strconv.ParseFloat(sLong, 64)
		if err != nil {
			return c.String(400, "Bad Request")
		}
		mLat, mLong := morton.PredictAppleCoord(lat, long)
		sp := spiral.NewSpiral(mLat, mLong)
		var tile *pb.WifiTile
		for i := 0; i < 20; i++ {
			mLat, mLong = sp.Next()
			tile, err = lib.GetTile(morton.PackAppleCoord(mLat, mLong))
			if err != nil {
				tile = nil
				continue
			}
			break
		}
		if tile == nil {
			return c.String(500, "Internal Server Error")
		}
		var points []distance.Point
		for _, r := range tile.GetRegion() {
			for _, d := range r.GetDevices() {
				points = append(points, distance.Point{
					Id: mac.Decode(d.GetBssid()),
					Y:  float64(d.GetEntry().GetLat()) * math.Pow10(-7),
					X:  float64(d.GetEntry().GetLong()) * math.Pow10(-7),
				})
			}
		}
		closest := distance.Closest(distance.Point{
			Id: "click",
			Y:  lat,
			X:  long,
		}, points)
		// Try to get closer via the wloc API
		for {
			devices, err := lib.QueryBssid([]string{closest.Id}, true)
			if err != nil {
				log.Println(err)
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
				Id: "click",
				Y:  lat,
				X:  long,
			}, points)
			if newClosest.Id == closest.Id {
				break
			}
			closest = newClosest
		}
		return c.JSON(200, map[string]any{
			"closest": closest,
			"points":  points,
		})
	})
	e.Logger.Fatal(e.Start(":1974"))
}
