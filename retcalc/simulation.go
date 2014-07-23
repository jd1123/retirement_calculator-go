package retcalc

import (
	"math/rand"
	"time"
)

// Simulation Object and some Methods
type Sim []float64

// Methods
func (s Sim) GrowthFactor(startYear int) float64 {
	fac := 1.0
	for i := startYear; i < len(s)-1; i++ {
		fac *= (1 + s[i])
	}
	return fac
}

func (s Sim) GrowthFactorWithTaxes(startYear int, eff_tax float64) float64 {
	fac := 1.0
	for i := startYear; i < len(s); i++ {
		fac *= (1 + s[i]*(1-eff_tax))
	}
	return fac
}

// Just a function
// Gives a float64 slice of lognormal returns with
// mean mean and stdev stdev of sample_size length
// Also seeds the random number generator
func Simulation(mean, stdev float64, sample_size int) []float64 {
	rand.Seed(time.Now().UTC().UnixNano())
	sim := make([]float64, sample_size, sample_size)
	for i := 0; i < sample_size; i++ {
		sim[i] = rand.NormFloat64()*stdev + mean
	}
	return sim
}
