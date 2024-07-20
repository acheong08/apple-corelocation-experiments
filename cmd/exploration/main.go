package main

import (
	"fmt"
	"log"

	shp "github.com/jonas-p/go-shp"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/planar"
	"github.com/paulmach/orb/project"
)

func main() {
	// Open the shapefile
	shape, err := shp.Open("assets/land_polygons.shp")
	if err != nil {
		log.Fatal(err)
	}
	defer shape.Close()
	var polygons []orb.Polygon
	for shape.Next() {
		_, p := shape.Shape()
		var orbPoly orb.Polygon
		switch poly := p.(type) {
		case *shp.Polygon:
			for i := range poly.NumParts {
				var points []shp.Point
				if poly.NumParts == 1 {
					points = poly.Points
				} else if i == poly.NumParts-1 {
					points = poly.Points[poly.Parts[i]:]
				} else {
					points = poly.Points[poly.Parts[i]:poly.Parts[i+1]]
				}
				ring := make(orb.Ring, len(points))
				for i := range len(points) {
					ring[i] = orb.Point{points[i].X, points[i].Y}
				}
				orbPoly = append(orbPoly, ring)
			}

		default:
			panic("unexpected shape type")
		}
		if len(orbPoly) == 0 {
			continue
		}
		polygons = append(polygons, orbPoly)
	}
	fmt.Println(len(polygons))
	points := map[string]struct {
		X float64
		Y float64
	}{
		"New York City":  {X: -74.0060, Y: 40.7128},
		"London":         {X: -0.1276, Y: 51.5074},
		"Tokyo":          {X: 139.6917, Y: 35.6895},
		"Sydney":         {X: 151.2093, Y: -33.8688},
		"Rio de Janeiro": {X: -43.1729, Y: -22.9068},
		"Cairo":          {X: 31.2357, Y: 30.0444},
		"Moscow":         {X: 37.6173, Y: 55.7558},
		"Cape Town":      {X: 18.4241, Y: -33.9249},
		"Beijing":        {X: 116.4074, Y: 39.9042},
		"Los Angeles":    {X: -118.2437, Y: 34.0522},
		"Mid-Atlantic":   {X: -30.0000, Y: 30.0000},
		"Mid-Pacific":    {X: -150.0000, Y: 0.0000},
		"North Pole":     {X: 0.0000, Y: 90.0000},
		"South Pole":     {X: 0.0000, Y: -90.0000},
		"Prime Meridian": {X: 0.0000, Y: 0.0000},
	}
	for name, p := range points {
		point := project.Point(orb.Point{p.X, p.Y}, project.WGS84.ToMercator)
		found := false
		for i, poly := range polygons {
			if planar.PolygonContains(poly, point) {
				fmt.Printf("%s is in polygon %d\n", name, i)
				found = true
				break
			}
		}
		if !found {
			fmt.Println(name, "is not in any polygon")
		}
	}
}
