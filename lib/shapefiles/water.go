package shapefiles

import (
	"wloc/lib/morton"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/planar"
	"github.com/paulmach/orb/project"
)

const level = 7

func IsInWater(lat, lon float64) bool {
	merc := project.WGS84.ToMercator(orb.Point{lon, lat})
	polies, ok := Waters[morton.Encode(lat, lon, level)]
	if ok && planar.MultiPolygonContains(polies, merc) {
		return true
	}
	return false
}
