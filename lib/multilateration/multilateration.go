package multilateration

import (
	"math"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/optimize"
)

func leastSquareFunc(point [2]float64, points []AccessPoint) *mat.VecDense {
	result := make([]float64, len(points))
	for i, p := range points {
		dist := Distance(p.Location.Lat, p.Location.Long, point[0], point[1])
		signalFactor := math.Pow(float64(p.SignalStrength), 2)
		result[i] = dist / signalFactor
	}
	return mat.NewVecDense(len(result), result)
}

type ResultType struct {
	Lat      float64
	Lon      float64
	Accuracy float64
}

func CalculatePosition(networks []AccessPoint) (lat, lon, accuracy float64) {
	// Initial position as the mean over all networks
	points := make([]float64, len(networks)*2)
	for i, net := range networks {
		points[i*2] = net.Location.Lat
		points[i*2+1] = net.Location.Long
	}

	initial := average(points)

	// Perform least squares optimization
	problem := optimize.Problem{
		Func: func(x []float64) float64 {
			result := leastSquareFunc([2]float64{x[0], x[1]}, networks)
			sum := 0.0
			for i := 0; i < result.Len(); i++ {
				sum += math.Pow(result.AtVec(i), 2)
			}
			return sum
		},
	}

	result, err := optimize.Minimize(problem, initial, nil, nil)
	if err != nil || result.Status != optimize.GradientThreshold {
		// No solution found, use initial estimate
		lat, lon = initial[0], initial[1]
	} else {
		lat, lon = result.X[0], result.X[1]
	}

	// Guess the accuracy as the 95th percentile of the distances
	distances := make([]float64, len(networks))
	for i, net := range networks {
		distances[i] = Distance(lat, lon, net.Location.Lat, net.Location.Long)
	}
	accuracy = percentile(distances, 95)

	return lat, lon, accuracy
}
