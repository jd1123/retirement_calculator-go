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
	"time"

	"code.google.com/p/plotinum/plot"
	"code.google.com/p/plotinum/plotter"
)

func RunPath(r RetCalc, s []float64) Path {
	ye := make([]YearlyEntry, r.Years, r.Years)

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
			retired := false
			ye[i] = YearlyEntry{age, st_date, SOY_taxable_balance, EOY_taxable_balance, SOY_non_taxable_balance,
				EOY_non_taxable_balance, Taxable_returns, Non_taxable_returns, Rate_of_return, Taxable_contribution,
				Non_taxable_contribution, Yearly_expenses, Deficit, retired}

		} else {
			//Not First year calculations
			retired := false
			st_date := time.Date(ye[i-1].Year.Year()+1, 1, 1, 0, 0, 0, 0, time.UTC)
			age := r.Age + i
			SOY_taxable_balance := ye[i-1].EOY_taxable_balance
			SOY_non_taxable_balance := ye[i-1].EOY_non_taxable_balance
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

			EOY_taxable_balance := SOY_taxable_balance + Taxable_returns + Taxable_contribution
			EOY_non_taxable_balance := SOY_non_taxable_balance + Non_taxable_returns +
				Non_taxable_contribution

			// Deduce Expenses
			if age > r.Retirement_age {
				retired = true
			}
			ye[i] = YearlyEntry{age, st_date, SOY_taxable_balance, EOY_taxable_balance, SOY_non_taxable_balance,
				EOY_non_taxable_balance, Taxable_returns, Non_taxable_returns, Rate_of_return,
				Taxable_contribution, Non_taxable_contribution, 0, 0, retired}
		}

	}

	return Path{ye, s, r.Inflation_rate}
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
	All_paths                                      PathGroup
}

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

func NewRetCalc_from_json(json_obj []byte) RetCalc {
	var r RetCalc
	err := json.Unmarshal(json_obj, &r)
	if err != nil {
		fmt.Println("Error")
	}

	r.sims = make([][]float64, r.N, r.N)
	for i := range r.sims {
		r.sims[i] = Simulation(r.Expected_rate_of_return, r.Asset_volatility, r.Years)
		r.All_paths[i] = RunPath(r, r.sims[i])
	}

	return r
}

func NewRetCalc() RetCalc {
	r := RetCalc{}

	r.N = 2500
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
	r.sims = make([][]float64, r.N, r.N)
	r.All_paths = make([]Path, r.N, r.N)
	for i := range r.sims {
		r.sims[i] = Simulation(r.Expected_rate_of_return, r.Asset_volatility, r.Years)
		r.All_paths[i] = RunPath(r, r.sims[i])
	}
	sort.Sort(r.All_paths)

	return r
}

func histogram(all_paths PathGroup) {
	//eb := all_paths.End_balances()
	eb := make([]float64, len(all_paths), len(all_paths))
	for i := range all_paths {
		eb[i] = all_paths[i].Income_from_path()
	}
	v := make(plotter.Values, len(eb))
	for i := range v {
		v[i] = eb[i]
	}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Histogram"
	h, err := plotter.NewHist(v, 100)
	if err != nil {
		panic(err)
	}
	h.Normalize(1)
	p.Add(h)

	if err := p.Save(4, 4, "hist.png"); err != nil {
		panic(err)
	}
	fmt.Println(h)
}

// main, mainly for testing
func main() {
	r := NewRetCalc()
	/*
		var all_paths path.PathGroup
		all_paths = make([]path.Path, r.N, r.N)

		for i := 0; i < r.N; i++ {
			all_paths[i] = RunPath(r, r.sims[i])
		}*/

	histogram(r.All_paths)
	fmt.Println(r)
	/*
		p := all_paths[0]
		_, f := p.Factors()
		fmt.Printf("Return\tFactor\n")
		for i := 0; i < len(p.Yearly_entries); i++ {
			fmt.Printf("%f\t%f\t\n", p.Sim[i], f[i])
		}*/
	/*
		for i := range all_paths {
			fmt.Printf("Income from path: %f\n", all_paths[i].Income_from_path())
		}
		all_paths[2500].Print_path()
		fmt.Println()
		fmt.Println(all_paths[2500].Final_balance())
		s, _ := all_paths[2500].Factors()
		fmt.Println(s)
		fmt.Println(all_paths[2500].Final_balance() / s)
		fmt.Println(all_paths[2500].Income_from_path())
	*/
}
