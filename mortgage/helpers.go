package mortgage

import "math"

func MortgageMonthlyPayment(loanAmount, r float64, termInMonths int) float64 {
	c := r / 12.0
	return loanAmount * (c * math.Pow((1+c), float64(termInMonths))) / (math.Pow((1+c), float64(termInMonths)) - 1)

}

func MonthlyMoneyGrowthVector(amount, r float64, termInMonths int) []float64 {
	c := r / 12.0
	out := make([]float64, termInMonths, termInMonths)
	for i := 0; i < termInMonths; i++ {
		out[i] = amount * math.Pow((1+c), float64(i))
	}
	return out
}

func MonthlyMoneyGrowthTerminal(amount, r float64, termInMonths int) float64 {
	c := r / 12.0
	return amount * math.Pow((1+c), float64(termInMonths))
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
