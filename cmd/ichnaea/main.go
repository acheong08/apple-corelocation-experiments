package main

import (
	"log"
	"net/http"
	"slices"
	"github.com/acheong08/apple-corelocation-experiments/lib"
	"github.com/acheong08/apple-corelocation-experiments/lib/multilateration"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.POST("/v1/geolocate", func(c echo.Context) error {
		var req multilateration.Request
		if err := c.Bind(&req); err != nil {
			return c.String(400, "ur request bad")
		}
		log.Println("Request made from ", c.RealIP())
		macs := make([]string, len(req.APs))
		for i, ap := range req.APs {
			macs[i] = ap.Mac
		}
		results, err := lib.QueryBssid(macs, int32(len(macs)), nil)
		if err != nil {
			return c.String(500, "uh oh. apple not working")
		}
		log.Printf("Results: %d, Requested: %d\n", len(results), len(req.APs))
		merger := make(map[string]multilateration.AccessPoint)
		for _, result := range results {
			if result.Location.Long == -180 {
				continue
			}
			if i := slices.Index(macs, result.BSSID); i != -1 {
				merger[result.BSSID] = multilateration.AccessPoint{
					Mac:            result.BSSID,
					Location:       result.Location,
					SignalStrength: req.APs[i].SignalStrength,
				}
			}
		}
		req.APs = make([]multilateration.AccessPoint, len(merger))
		cunt := 0
		for _, ap := range merger {
			req.APs[cunt] = ap
			cunt++
		}
		log.Println("Final length ", len(req.APs))
		for _, ap := range req.APs {
			log.Printf("%s: %4f, %4f", ap.Mac, ap.Location.Lat, ap.Location.Long)
		}
		lat, lon, accuracy := multilateration.CalculatePosition(req.APs)
		return c.JSON(200, map[string]any{
			"location": map[string]float64{
				"lat": lat,
				"lng": lon,
			},
			"accuracy": accuracy,
		})
	})
	log.Println("Starting server")
	if err := http.ListenAndServe("127.0.0.1:1975", e); err != nil {
		panic(err)
	}
}
