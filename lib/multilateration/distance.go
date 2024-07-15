package multilateration

import (
	"github.com/jftuga/geodist"
)

func Distance(lat1, lon1, lat2, lon2 float64) float64 {
	a := geodist.Coord{Lat: lat1, Lon: lon1}
	b := geodist.Coord{Lat: lat2, Lon: lon2}
	_, km, err := geodist.VincentyDistance(a, b)
	if err == nil {
		return km
	}
	_, km = geodist.HaversineDistance(a, b)
	return km
}
