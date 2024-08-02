package main

import (
	"fmt"
	"math"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
)

func main() {
	p := orb.Point{79.66789068052458, 5.4002188606820907}
	fmt.Println(geo.PointAtBearingAndDistance(p, 0, 10))
	// Create polygon around point with radius x
	poly := createCircularPolygon(p, 10)
	if !poly.Bound().Center().Equal(p) {
		fmt.Printf("Expected %v, got %v", p, poly.Bound().Center())
	}
}

func createCircularPolygon(center orb.Point, radius float64) orb.Polygon {
	const numSides = 36
	var polygon orb.Polygon
	var ring orb.Ring

	for i := 0; i < numSides; i++ {
		angle := float64(i) * (2 * math.Pi / numSides)
		point := geo.PointAtBearingAndDistance(center, angle, radius)
		ring = append(ring, orb.Point{point.Lon(), point.Lat()})
	}

	// Close the ring
	ring = append(ring, ring[0])
	polygon = append(polygon, ring)

	return polygon
}
