package shapefiles

import (
	"github.com/acheong08/apple-corelocation-experiments/lib/morton"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/planar"
	"github.com/paulmach/orb/project"
)

const level = 7

func IsInWater(lat, lon float64) bool {
	if lon == -180 || lon == 180 {
		return true
	}
	merc := project.WGS84.ToMercator(orb.Point{lon, lat})
	polies, ok := Waters[morton.Encode(lat, lon, level)]
	if ok && planar.MultiPolygonContains(polies, merc) {
		return true
	}
	return false
}
