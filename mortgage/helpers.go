package mortgage

import "math"

func MortgageMonthlyPayment(loanAmount, r float64, termInMonths int) float64 {
	c := r / 12.0
	return loanAmount * (c * math.Pow((1+c), float64(termInMonths))) / (math.Pow((1+c), float64(termInMonths)) - 1)

}

func MonthlyMoneyGrowthVector(amount, r float64, termInMonths int) []float64 {
	out := make([]float64, termInMonths, termInMonths)
	for i := 0; i < termInMonths; i++ {
		out[i] = amount * GF_m(r, i) //math.Pow((1+c), float64(i))
	}
	return out
}

func MonthlyMoneyGrowthTerminal(amount, r float64, termInMonths int) float64 {
	return amount * GF_m(r, termInMonths) // math.Pow((1+c), float64(termInMonths))
}

func SumProduct(a, b []float64) float64 {
	if len(a) != len(b) {
		panic("lengths of slices are not equal!")
	}
	sum := 0.0
	for i := 0; i < len(a); i++ {
		sum += a[i] * b[i]
	}
	return sum / float64(len(a))
}

// Growth factor and discount factor functions
func DF_d(r float64, termInDays int) float64 {
	c := r / 365.25
	return 1 / math.Pow((1+c), float64(termInDays))
}

func GF_d(r float64, termInDays int) float64 {
	c := r / 365.25
	return 1 / math.Pow((1+c), float64(termInDays))
}

func DF_m(r float64, termInMonths int) float64 {
	c := r / 12.0
	return 1 / math.Pow((1+c), float64(termInMonths))
}

func GF_m(r float64, termInMonths int) float64 {
	c := r / 12.0
	return math.Pow((1 + c), float64(termInMonths))
}

func DF(r float64, termInYears int) float64 {
	return 1 / math.Pow((1+r), float64(termInYears))
}

func GF(r float64, termInYears int) float64 {
	return math.Pow((1 + r), float64(termInYears))
}
