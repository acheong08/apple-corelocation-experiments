package shapefiles

import (
	"log"
	"wloc/lib/morton"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/planar"
	"github.com/paulmach/orb/project"
)

const level = 7

func IsInWater(lat, lon float64) bool {
	merc := project.WGS84.ToMercator(orb.Point{lon, lat})
	log.Printf("%f, %f", merc.X(), merc.Y())
	log.Println("Expected: ", morton.Encode(lat, lon, level))
	polies, ok := Waters[morton.Encode(lat, lon, level)]
	log.Println("Polylen ", len(polies))
	if ok && planar.MultiPolygonContains(polies, merc) {
		return true
	}
	return false
}
