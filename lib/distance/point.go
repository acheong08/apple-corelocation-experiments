package distance

import "math"

type Point struct {
	Id   string
	X, Y float64
}

func (p *Point) Distance(d Point) float64 {
	return math.Abs(p.X-d.X) + math.Abs(p.Y-d.Y)
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
