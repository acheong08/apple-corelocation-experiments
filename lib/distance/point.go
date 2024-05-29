package distance

import "math"

type Point struct {
	Id   string
	x, y float64
}

func (p *Point) Distance(d Point) float64 {
	return math.Abs(p.x-d.x) + math.Abs(p.y-d.y)
}

func Closest(point Point, points []Point) Point {
	var closest Point
	var closestDistance float64
	for _, p := range points {
		distance := point.Distance(p)
		if distance < closestDistance {
			closest = p
			closestDistance = distance
		}
	}
	return closest
}
