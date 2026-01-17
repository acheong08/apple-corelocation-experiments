package spiral_test

import (
	"github.com/acheong08/apple-corelocation-experiments/lib/spiral"
	"testing"
)

func TestSpiral(t *testing.T) {
	s := spiral.NewSpiral(0, 0)
	points := make([]struct{ x, y int }, 25)

	// Fill the array with points
	for i := range points {
		points[i].x, points[i].y = s.Next()
	}

	// Check if the 5x5 area around (0, 0) is filled
	isFilled := func(x, y int) bool {
		for _, p := range points {
			if p.x == x && p.y == y {
				return true
			}
		}
		return false
	}

	for x := -2; x <= 2; x++ {
		for y := -2; y <= 2; y++ {
			if !isFilled(x, y) {
				t.Errorf("Point (%d, %d) is not filled", x, y)
			}
		}
	}
}
