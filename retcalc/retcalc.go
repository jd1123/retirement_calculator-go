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

package retcalc

import (
	"encoding/json"
	"fmt"
	"sort"
)

// Struct definitions

// RetCalc Object, the meat
// also the RetCalc "constructors"
// This needs to be reduced - we will have it only hold all the
// Necessary info and run the calcs on API call
type RetCalc struct {
	Age, Retirement_age, Terminal_age              int
	Effective_tax_rate, Returns_tax_rate           float64
	Years                                          int
	N                                              int
	sims                                           []Sim
	Non_Taxable_contribution, Taxable_contribution float64
	Non_Taxable_balance, Taxable_balance           float64
	Yearly_retirement_expenses                     float64
	Yearly_social_security_income                  float64
	Asset_volatility, Expected_rate_of_return      float64
	Inflation_rate                                 float64
	All_paths                                      PathGroup
}

// This is the IDEAL RetCalc Object but
// required further refacoring that I cannot do
// while on the train and having to take a massive
// dump with the guy next to me sharting his pants
type RetCalc_web_input struct {
	Age, Retirement_age, Terminal_age              int
	Effective_tax_rate, Returns_tax_rate           float64
	Years                                          int
	N                                              int
	Non_Taxable_contribution, Taxable_contribution float64
	Non_Taxable_balance, Taxable_balance           float64
	Yearly_social_security_income                  float64
	Asset_volatility, Expected_rate_of_return      float64
	Inflation_rate                                 float64
}

// METHODS

// Runs all the paths - dont do this in the constructor
// because the data will be lost when passing a smaller,
// client side object that is portable
func (r RetCalc) AllPaths() PathGroup {
	all_paths := make(PathGroup, len(r.sims), len(r.sims))
	for i := range r.sims {
		all_paths[i] = RunPath(r, r.sims[i])
	}
	sort.Sort(all_paths)
	return all_paths
}

// The path struct should implement this logic, it is misplaced
func (r RetCalc) RunIncomes() []float64 {
	incomes := make([]float64, len(r.sims), len(r.sims))
	for i := range r.sims {
		taxed_total_wealth := 0.0
		untaxed_total_wealth := 0.0
		sum_t, sum_ut := 0.0, 0.0
		for j := range r.sims[i] {
			if j+r.Age < r.Retirement_age {
				untaxed_total_wealth += r.Non_Taxable_contribution * r.sims[i].GrowthFactor(j)
				taxed_total_wealth += r.Taxable_contribution * r.sims[i].GrowthFactorWithTaxes(j, r.Effective_tax_rate)
				sum_ut += r.sims[i].GrowthFactor(j)
				sum_t += r.sims[i].GrowthFactorWithTaxes(j, r.Effective_tax_rate)
			}
		}
		f, _ := r.All_paths[i].Factors()
		incomes[i] = (taxed_total_wealth + untaxed_total_wealth*(1-r.Effective_tax_rate)) / f
	}
	return incomes
}

func (r RetCalc) PercentilePath(percentile float64) Path {
	ix := int(float64(r.N) * percentile)
	return r.All_paths[ix]
}

func (r RetCalc) IncomeProbability() float64 {
	incomes := r.RunIncomes()
	counter := 0
	for i := range incomes {
		if incomes[i] >= r.Yearly_retirement_expenses {
			counter++
		}
	}
	fmt.Printf("Counter: %i Income: %f N: %i\n", counter, r.Yearly_retirement_expenses, r.N)
	return float64(counter) / float64(r.N)
}

// Constructors

// This constructor will populate a RetCalc from
// JSON input from the web ----- NEEDS work
func NewRetCalc_from_json(json_obj []byte) RetCalc {
	var r RetCalc
	err := json.Unmarshal(json_obj, &r)
	if err != nil {
		fmt.Println("Error")
	}

	r.sims = make([]Sim, r.N, r.N)
	for i := range r.sims {
		r.sims[i] = Simulation(r.Expected_rate_of_return, r.Asset_volatility, r.Years)
		r.All_paths[i] = RunPath(r, r.sims[i])
	}

	return r
}

// A default RetCalc
func NewRetCalc() RetCalc {
	r := RetCalc{}

	r.N = 10000
	r.Age = 22
	r.Retirement_age = 65
	r.Terminal_age = 90
	r.Years = r.Terminal_age - r.Age
	r.Effective_tax_rate = 0.30
	r.Returns_tax_rate = 0.30
	r.Non_Taxable_contribution = 17500.0
	r.Taxable_contribution = 0
	r.Non_Taxable_balance = 0
	r.Yearly_retirement_expenses = float64(100000)
	r.Yearly_social_security_income = 0.0
	r.Taxable_balance = 0.0
	r.Asset_volatility = 0.15
	r.Expected_rate_of_return = 0.07
	r.Inflation_rate = 0.035
	r.sims = make([]Sim, r.N, r.N)
	r.All_paths = make([]Path, r.N, r.N)
	for i := range r.sims {
		r.sims[i] = Simulation(r.Expected_rate_of_return, r.Asset_volatility, r.Years)
		r.All_paths[i] = RunPath(r, r.sims[i])
	}
	sort.Sort(r.All_paths)

	return r
}

// Dumbed Down version of constructor
func NewRetCalc_b() RetCalc {
	r := RetCalc{}

	r.N = 10000
	r.Age = 22
	r.Retirement_age = 65
	r.Terminal_age = 90
	r.Years = r.Terminal_age - r.Age
	r.Effective_tax_rate = 0.30
	r.Returns_tax_rate = 0.30
	r.Non_Taxable_contribution = 17500
	r.Taxable_contribution = 0
	r.Non_Taxable_balance = 0
	r.Yearly_retirement_expenses = float64(60000)
	r.Yearly_social_security_income = 0.0
	r.Taxable_balance = 0.0
	r.Asset_volatility = 0.15
	r.Expected_rate_of_return = 0.07
	r.Inflation_rate = 0.035
	r.sims = make([]Sim, r.N, r.N)
	for i := range r.sims {
		r.sims[i] = Simulation(r.Expected_rate_of_return, r.Asset_volatility, r.Years)
	}
	return r
}
