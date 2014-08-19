package analytics

import "sort"

var incomeBracket = map[float64]float64{
	0.0:      0.1,
	9075.0:   0.15,
	36900.0:  0.25,
	89350.0:  0.28,
	186350.0: 0.33,
	405100.0: 0.35,
	406750.0: 0.396,
}

// Given the tax bracket above, compute the income tax liability
func IncomeTaxLiability(income float64) float64 {
	keys := make([]float64, len(incomeBracket))
	i := 0
	for k, _ := range incomeBracket {
		keys[i] = k
		i++
	}
	sort.Float64s(keys)
	incomeLiability := 0.0
	//incomeLeft := income
	for i = 1; i < len(keys); i++ {
		if income > keys[i] {
			if i != len(keys)-1 {
				br := keys[i] - keys[i-1]
				incomeLiability += br * incomeBracket[keys[i-1]]
			} else if i == len(keys)-1 {
				br := income - keys[i]
				incomeLiability += br * incomeBracket[keys[i]]
			}
		} else if income > keys[i-1] && income < keys[i] {
			br := income - keys[i-1]
			incomeLiability += br * incomeBracket[keys[i]]
		}
	}
	return incomeLiability
}
