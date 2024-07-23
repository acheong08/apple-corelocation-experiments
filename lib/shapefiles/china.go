package shapefiles

import (
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/planar"
)

func IsInChina(lat, lon float64) bool {
	return planar.PolygonContains(China, orb.Point{lon, lat})
}
