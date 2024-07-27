package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"slices"
	"wloc/lib/morton"

	shp "github.com/jonas-p/go-shp"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/planar"
	"github.com/paulmach/orb/project"
)

const level = 7

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: shp-to-orb path/to/shp path/to/output")
	}
	path := os.Args[1]
	// Open the shapefile
	shape, err := shp.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer shape.Close()
	var polygons orb.MultiPolygon
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
	if slices.Contains(os.Args, "-map") {
		// Build a map of zoom level 9 morton blocks
		watery := make(map[int64][]orb.Polygon)
		count := 0
		for _, poly := range polygons {
			center := project.Mercator.ToWGS84(poly.Bound().Bound().Center())
			code := morton.Encode(center.Lat(), center.Lon(), level)
			watery[code] = append(watery[code], poly)
			count++
		}
		log.Println(count)
		save(watery)
	} else {
		if len(polygons) == 1 {
			save(polygons[0])
		} else {
			save(polygons)
		}
	}
	test(polygons)
}

func save(a any) {
	// Save polygons as gob file
	f, err := os.Create(os.Args[2])
	if err != nil {
		log.Fatalf("Failed to create destination file %s: %s", os.Args[2], err.Error())
	}
	enc := gob.NewEncoder(f)
	if err := enc.Encode(a); err != nil {
		log.Fatalf("Failed to encode polygons: %s", err.Error())
	}
}

func test(poly orb.MultiPolygon) {
	fmt.Println(len(poly))
	points := map[string]struct {
		X float64
		Y float64
	}{
		"North Pacific": {32.890398, 146.864834},
		"China":         {45.964474, 119.773672},
		"Penang coast":  {5.419154, 100.343326},
		"Penang bridge": {5.304548, 100.359499},
	}
	for name, p := range points {
		point := project.Point(orb.Point{p.Y, p.X}, project.WGS84.ToMercator)
		found := false
		for i, poly := range poly {
			if planar.PolygonContains(poly, point) {
				center := project.Mercator.ToWGS84(poly.Bound().Center())
				code := morton.Encode(center.Lat(), center.Lon(), level)

				fmt.Printf("%s is in polygon %d with code %d\n", name, i, code)
				found = true
				break
			}
		}
		if !found {
			fmt.Println(name, "is not in any polygon")
		}
	}
}
