package retcalc

import (
	"math/rand"
	"time"
)

type Sim []float64

func (s Sim) GrowthFactor(startYear int) float64 {
	fac := 1.0
	for i := startYear; i < len(s); i++ {
		fac *= (1 + s[i])
	}
	return fac
}

func (s Sim) GrowthFactorWithTaxes(startYear int, eff_tax float64) float64 {
	fac := 1.0
	for i := startYear; i < len(s); i++ {
		fac *= (1 + s[i]*eff_tax)
	}
	return fac
}

func Simulation(mean, stdev float64, sample_size int) []float64 {
	rand.Seed(time.Now().UTC().UnixNano())
	sim := make([]float64, sample_size, sample_size)
	for i := 0; i < sample_size; i++ {
		sim[i] = rand.NormFloat64()*stdev + mean
	}
	return sim
}
