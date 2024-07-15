package multilateration

import "sort"

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
