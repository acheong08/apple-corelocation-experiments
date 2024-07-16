package main

import (
	"net/http"
	"wloc/lib"
	"wloc/lib/multilateration"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.POST("/v1/geolocate", func(c echo.Context) error {
		var req multilateration.Request
		if err := c.Bind(&req); err != nil {
			return c.String(400, "ur request bad")
		}
		macs := make([]string, len(req.APs))
		for i, ap := range req.APs {
			macs[i] = ap.Mac
		}
		results, err := lib.QueryBssid(macs, true, nil)
		if err != nil {
			return c.String(500, "uh oh. apple not working")
		}
		for i, result := range results {
			req.APs[i].Location = result.Location
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
	if err := http.ListenAndServe("127.0.0.1:1975", e); err != nil {
		panic(err)
	}
}
