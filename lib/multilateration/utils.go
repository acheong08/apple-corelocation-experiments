package multilateration

import "sort"

func average(points []float64) []float64 {
	result := make([]float64, len(points)/2)
	for i := 0; i < len(points); i += 2 {
		result[i/2] = (points[i] + points[i+1]) / 2
	}
	return result
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
