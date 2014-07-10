package simulation

import (
	"math/rand"
	"time"
)

func Simulation(mean, stdev float64, sample_size int) []float64 {
	rand.Seed(time.Now().UTC().UnixNano())
	sim := make([]float64, sample_size, sample_size)
	for i := 0; i < sample_size; i++ {
		sim[i] = rand.NormFloat64()*stdev + mean
	}
	return sim
}
