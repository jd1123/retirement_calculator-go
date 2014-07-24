package analytics

import (
	"fmt"
	"math"
	"sort"
)

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

func HistoCumulative(a []float64, n_bins int) []byte {
	max := MaxF64(a)
	min := MinF64(a)
	l := len(a)
	cdf := make(map[float64]int)
	bins := make([]float64, n_bins, n_bins)
	binStep := max / float64(n_bins)

	for i := range bins {
		bins[i] = binStep * float64(i+1)
		cdf[bins[i]] = 0
	}

	for k, _ := range cdf {
		for j := range a {
			if a[j] > k-binStep {
				cdf[k]++
			}
		}
	}

	// Print Out Map
	var keys []float64
	for k := range cdf {
		keys = append(keys, k)
	}
	sort.Float64s(keys)
	for _, k := range keys {
		fmt.Println("Key: ", k, " Value: ", cdf[k])
	}
	// END Print Out Map

	fmt.Println("HistoCumulative() analytics:")
	fmt.Printf("Max: %f\n", max)
	fmt.Printf("Min: %f\n", min)
	fmt.Printf("Len: %d\n", l)

	return []byte{0}
}
