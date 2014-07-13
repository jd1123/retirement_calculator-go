/*
   This file is part of retirement_calculator.

   Retirement_calculator is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   Retirement_calculator is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with retirement_calculator.  If not, see <http://www.gnu.org/licenses/>.
*/

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
		fmt.Printf("Age: %d EOY t_Balance: %f EOY nt_balance: %f Return: %f\n", p.Yearly_entries[i].Age,
			p.Yearly_entries[i].EOY_taxable_balance, p.Yearly_entries[i].EOY_non_taxable_balance, p.Sim[i])
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

	return s_factors, factors
}

// Calculates the yearly income from a particular path
func (p Path) Income_from_path() float64 {
	s_factors, _ := p.Factors()
	return p.Final_balance() / s_factors
}

// Implement the sort interface on a group of Paths
type PathGroup []Path

func (p PathGroup) Len() int {
	return len(p)
}

func (p PathGroup) Less(i, j int) bool {
	//li, lj := len(p[i].Yearly_entries), len(p[j].Yearly_entries)
	vi := p[i].Income_from_path()
	vj := p[j].Income_from_path()
	//vi := p[i].Yearly_entries[li-1].EOY_taxable_balance + p[i].Yearly_entries[li-1].EOY_non_taxable_balance
	//vj := p[j].Yearly_entries[lj-1].EOY_taxable_balance + p[j].Yearly_entries[lj-1].EOY_non_taxable_balance
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
