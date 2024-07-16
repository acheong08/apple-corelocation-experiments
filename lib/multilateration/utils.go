package multilateration

import "sort"

func average(points []float64) []float64 {
	sum := make([]float64, 2)
	for i := 0; i < len(points); i += 2 {
		sum[0] += points[i]
		sum[1] += points[i+1]
	}
	return []float64{sum[0] / float64(len(points)/2), sum[1] / float64(len(points)/2)}
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
