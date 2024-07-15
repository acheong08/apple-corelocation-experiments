package multilateration

import (
	"math"
	"sort"

	"github.com/jftuga/geodist"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/optimize"
)

func Distance(lat1, lon1, lat2, lon2 float64) float64 {
	a := geodist.Coord{Lat: lat1, Lon: lon1}
	b := geodist.Coord{Lat: lat2, Lon: lon2}
	_, km, err := geodist.VincentyDistance(a, b)
	if err == nil {
		return km
	}
	_, km = geodist.HaversineDistance(a, b)
	return km
}

func LeastSquareFunc(point [2]float64, points []Cell) *mat.VecDense {
	result := make([]float64, len(points))
	for i, p := range points {
		dist := Distance(p.Location.Lat, p.Location.Long, point[0], point[1])
		ageFactor := math.Min(math.Sqrt(2000.0/float64(p.Age)), 1.0)
		signalFactor := math.Pow(float64(p.SignalStrength), 2)
		result[i] = dist * ageFactor / signalFactor
	}
	return mat.NewVecDense(len(result), result)
}

func AggregateMacPosition(networks []Cell, minimumAccuracy float64) (lat, lon, accuracy float64) {
	// Guess initial position as the weighted mean over all networks
	points := make([]float64, len(networks)*2)
	weights := make([]float64, len(networks))

	for i, net := range networks {
		points[i*2] = net.Location.Lat
		points[i*2+1] = net.Location.Long
		ageFactor := math.Min(math.Sqrt(2000.0/float64(net.Age)), 1.0)
		signalFactor := math.Pow(float64(net.SignalStrength), 2)
		weights[i] = net.Score * ageFactor / signalFactor
	}

	initial := weightedAverage(points, weights)

	// Perform least squares optimization
	problem := optimize.Problem{
		Func: func(x []float64) float64 {
			result := LeastSquareFunc([2]float64{x[0], x[1]}, networks)
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
	if accuracy < minimumAccuracy {
		accuracy = minimumAccuracy
	}

	return lat, lon, accuracy
}

func weightedAverage(points []float64, weights []float64) []float64 {
	sum := make([]float64, 2)
	totalWeight := 0.0
	for i := 0; i < len(weights); i++ {
		sum[0] += points[i*2] * weights[i]
		sum[1] += points[i*2+1] * weights[i]
		totalWeight += weights[i]
	}
	return []float64{sum[0] / totalWeight, sum[1] / totalWeight}
}

func percentile(data []float64, percentile float64) float64 {
	sort.Float64s(data)
	index := (percentile / 100) * float64(len(data)-1)
	if index == float64(int(index)) {
		return data[int(index)]
	}
	lower := data[int(index)]
	upper := data[int(index)+1]
	return lower + (upper-lower)*(index-float64(int(index)))
}
