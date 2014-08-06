package analytics

import (
	"fmt"
	"math"
	"sort"
	"strconv"
)

func FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}

type Bin struct {
	Max, Min float64
	Weight   float64
}

type Histo struct {
	Bins []Bin
}

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

func HistoCumulative(a []float64, n_bins int) Histo {
	normalize := true
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

	var h Histo
	h.Bins = make([]Bin, n_bins, n_bins)

	var keys []float64
	for k := range cdf {
		keys = append(keys, k)
	}

	sort.Float64s(keys)
	i := 0
	for _, k := range keys {
		h.Bins[i].Max = k
		h.Bins[i].Min = k - binStep
		if normalize {
			h.Bins[i].Weight = float64(cdf[k]) / float64(l)

		} else {
			h.Bins[i].Weight = float64(cdf[k])
		}
		//fmt.Println("Key: ", k, " Bins[i].Max: ", h.Bins[i].Max)
		i++
	}

	/*
		// Print Out Map
		for k := range cdf {
			keys = append(keys, k)
		}
		sort.Float64s(keys)
		for _, k := range keys {
			//fmt.Println("Key: ", k, " Value: ", cdf[k])
		}
		// END Print Out Map
	*/

	fmt.Println("HistoCumulative() analytics:")
	fmt.Printf("Max: %f\n", max)
	fmt.Printf("Min: %f\n", min)
	fmt.Printf("Len: %d\n", l)

	return h
}
