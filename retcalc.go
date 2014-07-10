/*

My retirement calculator in Go

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.

Copyright 2014 Johnnydiabetic
*/

package main

import (
	"encoding/json"
	"fmt"
	"math"
	"retirement_calculator-go/path"
	"retirement_calculator-go/simulation"
	"sort"
	"time"
)

func run_path(r RetCalc, s []float64) path.Path {
	p := make([]path.YearlyEntry, r.Years, r.Years)

	for i := 0; i < r.Years; i++ {
		if i == 0 {
			st_date := time.Date(time.Now().Year()+1, 1, 1, 0, 0, 0, 0, time.UTC)
			age := r.Age
			SOY_taxable_balance := r.Taxable_balance
			SOY_non_taxable_balance := r.Non_Taxable_balance
			Rate_of_return := s[i]
			Taxable_contribution := r.Taxable_contribution
			Non_taxable_contribution := r.Non_Taxable_contribution
			Taxable_returns := Rate_of_return * SOY_taxable_balance
			Non_taxable_returns := Rate_of_return * SOY_non_taxable_balance
			Yearly_expenses := float64(0)

			EOY_taxable_balance := SOY_taxable_balance + Taxable_returns + Taxable_contribution
			EOY_non_taxable_balance := SOY_non_taxable_balance + Non_taxable_returns + Non_taxable_contribution
			Deficit := 0.0

			p[i] = path.YearlyEntry{age, st_date, SOY_taxable_balance, EOY_taxable_balance, SOY_non_taxable_balance,
				EOY_non_taxable_balance, Taxable_returns, Non_taxable_returns, Rate_of_return, Taxable_contribution,
				Non_taxable_contribution, Yearly_expenses, Deficit}

		} else {
			//Not First year calculations

			st_date := time.Date(p[i-1].Year.Year()+1, 1, 1, 0, 0, 0, 0, time.UTC)
			age := r.Age + i
			SOY_taxable_balance := p[i-1].EOY_taxable_balance
			SOY_non_taxable_balance := p[i-1].EOY_non_taxable_balance
			Rate_of_return := s[i]

			Taxable_contribution := 0.0
			Non_taxable_contribution := 0.0
			if age <= r.Retirement_age {
				Taxable_contribution = r.Taxable_contribution
				Non_taxable_contribution = r.Non_Taxable_contribution
			}

			Taxable_returns := 0.0
			if SOY_taxable_balance > 0 {
				Taxable_returns = Rate_of_return * SOY_taxable_balance * (1 - r.Returns_tax_rate)
			}

			Non_taxable_returns := 0.0
			if SOY_non_taxable_balance > 0 {
				Non_taxable_returns = Rate_of_return * SOY_non_taxable_balance
			}

			Yearly_expenses := float64(0)
			if age > r.Retirement_age {
				Yearly_expenses = r.Yearly_retirement_expenses * math.Pow((1+r.Inflation_rate), float64(i))
			}

			EOY_taxable_balance := SOY_taxable_balance + Taxable_returns + Taxable_contribution
			EOY_non_taxable_balance := SOY_non_taxable_balance + Non_taxable_returns +
				Non_taxable_contribution

			// Deduce Expenses
			Deficit := 0.0
			if (SOY_non_taxable_balance + Non_taxable_returns) < Yearly_expenses {
				Yearly_expenses -= EOY_non_taxable_balance + Non_taxable_returns
				EOY_non_taxable_balance = 0
				EOY_taxable_balance -= Yearly_expenses
			} else {
				EOY_non_taxable_balance -= Yearly_expenses * (1 + r.Effective_tax_rate)
			}

			p[i] = path.YearlyEntry{age, st_date, SOY_taxable_balance, EOY_taxable_balance, SOY_non_taxable_balance,
				EOY_non_taxable_balance, Taxable_returns, Non_taxable_returns, Rate_of_return,
				Taxable_contribution, Non_taxable_contribution, Yearly_expenses, Deficit}
		}

	}

	return p
}

// RetCalc Object, the meat
// also the RetCalc "constructors"
type RetCalc struct {
	Age, Retirement_age, Terminal_age              int
	Effective_tax_rate, Returns_tax_rate           float64
	Years                                          int
	N                                              int
	sims                                           [][]float64
	Non_Taxable_contribution, Taxable_contribution float64
	Non_Taxable_balance, Taxable_balance           float64
	Yearly_retirement_expenses                     float64
	Yearly_social_security_income                  float64
	Asset_volatility, Expected_rate_of_return      float64
	Inflation_rate                                 float64
}

func NewRetCalc_from_json(json_obj []byte) RetCalc {
	var r RetCalc
	err := json.Unmarshal(json_obj, &r)
	if err != nil {
		fmt.Println("Error")
	}

	r.sims = make([][]float64, r.N, r.N)
	for i := range r.sims {
		r.sims[i] = simulation.Simulation(r.Expected_rate_of_return, r.Asset_volatility, r.Years)
	}

	return r
}

func NewRetCalc() RetCalc {
	r := RetCalc{}

	r.N = 100000
	r.Age = 22
	r.Retirement_age = 65
	r.Terminal_age = 90
	r.Years = r.Terminal_age - r.Age
	r.Effective_tax_rate = 0.30
	r.Returns_tax_rate = 0.30
	r.Non_Taxable_contribution = 17500
	r.Taxable_contribution = 2500
	r.Non_Taxable_balance = 0
	r.Yearly_retirement_expenses = float64(60000)
	r.Yearly_social_security_income = 0.0
	r.Taxable_balance = 0.0
	r.Asset_volatility = 0.15
	r.Expected_rate_of_return = 0.07
	r.Inflation_rate = 0.035

	r.sims = make([][]float64, r.N, r.N)
	for i := range r.sims {
		r.sims[i] = simulation.Simulation(r.Expected_rate_of_return, r.Asset_volatility, r.Years)
	}

	return r
}

// main, mainly for testing
func main() {
	r := NewRetCalc()
	var all_paths path.PathGroup
	all_paths = make([]path.Path, r.N, r.N)

	for i := 0; i < r.N; i++ {
		all_paths[i] = run_path(r, r.sims[i])
	}

	sort.Sort(all_paths)
	ix := int(float64(r.N) * 0.20)
	all_paths[ix].Print_path()

	fmt.Println(len(all_paths))
	//p := run_path(r, r.sims[0])
	//p.print_path()
	b, _ := json.Marshal(r)
	fmt.Println(string(b))
}
