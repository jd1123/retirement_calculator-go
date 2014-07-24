package analytics

import "math"

func MaxF64(a []float64) float64 {
	if len(a) == 0 {
		panic("Empty slice passed!")
	}
	max := math.MaxFloat64 * -1
	for i := range a {
		if a[i] > max {
			max = a[i]
		}
	}
	return max
}

func MinF64(a []float64) float64 {
	if len(a) == 0 {
		panic("Empty slice passed!")
	}
	min := math.MaxFloat64
	for i := range a {
		if a[i] < min {
			min = a[i]
		}
	}
	return min
}

func AvgF64(a []float64) float64 {
	if len(a) == 0 {
		panic("Empty slice passed!")
	}
	avg := 0.0
	for i := range a {
		avg += a[i]
	}
	return avg / float64(len(a))
}
