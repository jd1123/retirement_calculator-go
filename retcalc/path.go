// The Path Object and methods

package retcalc

import (
	"fmt"
	"math"
	"time"
)

type YearlyEntry struct {
	Age                                              int
	Year                                             time.Time
	SOY_taxable_balance, EOY_taxable_balance         float64
	SOY_non_taxable_balance, EOY_non_taxable_balance float64
	Taxable_returns, Non_taxable_returns             float64
	Rate_of_return                                   float64
	Taxable_contribution, Non_taxable_contribution   float64
	Yearly_expenses                                  float64
	Deficit                                          float64
	Retired                                          bool
}

// Path type and Path methods
// NOTE: The income methods work but do not take into
// account taxes
// NOTE: taxes are important
type Path struct {
	Yearly_entries []YearlyEntry
	Sim            []float64
	Inflation_rate float64
}

// Prints info about the path
func (p Path) Print_path() {
	sum := 0.0
	for i := 0; i < len(p.Yearly_entries); i++ {
		//fmt.Println(p[i])
		sum += p.Yearly_entries[i].Rate_of_return
		//fmt.Printf("Age: %d EOY Balance: %f\n", p[i].Age, p[i].EOY_taxable_balance+p[i].EOY_non_taxable_balance)
		fmt.Printf("Age: %d EOY_B: %f cont: %f r: %f\n", p.Yearly_entries[i].Age,
			p.Yearly_entries[i].EOY_taxable_balance+p.Yearly_entries[i].EOY_non_taxable_balance,
			p.Yearly_entries[i].Non_taxable_contribution+p.Yearly_entries[i].Taxable_contribution, p.Sim[i])
	}
	fmt.Printf("Average RoR: %f", sum/float64(len(p.Yearly_entries)))
}

// Returns the final balance
func (p Path) Final_balance() float64 {
	l := len(p.Yearly_entries) - 1
	return p.Yearly_entries[l].EOY_taxable_balance + p.Yearly_entries[l].EOY_non_taxable_balance
}

// Returns the factors and the sum of factors
// to compute income from a final balance
// NOTE: I think this should be renamed IncomeFactors, and
// also think it probably should be elsewhere. It is not reusing
// code effectively due to its location in path.go
// Going to move it to a method of retcalc
func (p Path) Factors() (float64, []float64) {
	factors := make([]float64, len(p.Sim), len(p.Sim))
	s_factors := 0.0

	for i := range p.Sim {
		sum := 1.0
		for j := i + 1; j < len(p.Sim); j++ {
			sum *= (1 + p.Sim[j])
		}
		if p.Yearly_entries[i].Retired {
			factors[i] = sum * math.Pow(1+p.Inflation_rate, float64(i))
			s_factors += factors[i]
		} else {
			factors[i] = 0
		}
	}

	s := 0.0
	for i := range factors {
		if i != 0 {
			if factors[i-1] != 0 && factors[i] != 0 {
				s += factors[i]
			}
		}
	}

	return s, factors
}

// Implement the sort interface on a group of Paths
type PathGroup []Path

func (p PathGroup) Len() int {
	return len(p)
}

func (p PathGroup) Less(i, j int) bool {
	li, lj := len(p[i].Yearly_entries), len(p[j].Yearly_entries)
	//vi := p[i].End_Balances()
	//vj := p[j].End_Balances()
	vi := p[i].Yearly_entries[li-1].EOY_taxable_balance + p[i].Yearly_entries[li-1].EOY_non_taxable_balance
	vj := p[j].Yearly_entries[lj-1].EOY_taxable_balance + p[j].Yearly_entries[lj-1].EOY_non_taxable_balance
	return vi < vj
}

func (p PathGroup) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p PathGroup) End_balances() []float64 {
	eb := make([]float64, len(p), len(p))
	for i := range p {
		eb[i] = p[i].Final_balance()
	}
	return eb
}
