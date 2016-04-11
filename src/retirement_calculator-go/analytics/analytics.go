package analytics

import (
	"math"
	"sort"
	"strconv"
)

// Wrapper to make my life easier
func FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}

// Structs to make a histogram
type Bin struct {
	Max, Min float64
	Weight   float64
}

type Histo struct {
	Bins []Bin
}

// Max for float64 slices
func MaxF64(a []float64) float64 {
	if len(a) == 0 {
		panic("Empty slice passed to MaxF64()!")
	}
	max := math.MaxFloat64 * -1
	for i := range a {
		if a[i] > max {
			max = a[i]
		}
	}
	return max
}

// Min for float64 slices
func MinF64(a []float64) float64 {
	if len(a) == 0 {
		panic("Empty slice passed to MinF64()!")
	}
	min := math.MaxFloat64
	for i := range a {
		if a[i] < min {
			min = a[i]
		}
	}
	return min
}

// Average for float64 slices
func AvgF64(a []float64) float64 {
	if len(a) == 0 {
		panic("Empty slice passed to AvgF64()!")
	}
	avg := 0.0
	for i := range a {
		avg += a[i]
	}
	return avg / float64(len(a))
}

// Returns a cumulative histogram for a []float64 slice
func HistoCumulative(a []float64, n_bins int) Histo {
	normalize := true
	max := MaxF64(a)
	//min := MinF64(a)
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
		i++
	}

	// commented out code to look at the histo
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
	/*
		fmt.Println("HistoCumulative() analytics:")
		fmt.Printf("Max: %f\n", max)
		fmt.Printf("Min: %f\n", min)
		fmt.Printf("Len: %d\n", l)
	*/

	return h
}
