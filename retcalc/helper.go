// Helper functions for the RetCalc Package that
// are not methods
// here to unclog retcalc.go

package retcalc

import (
	"fmt"
	"time"

	"code.google.com/p/plotinum/plot"
	"code.google.com/p/plotinum/plotter"
)

// Builds a histogram from a []float64 slice :)
func HistoFromSlice(slice []float64) *plotter.Histogram {
	v := make(plotter.Values, len(slice))
	for i := range v {
		v[i] = slice[i]
	}
	h, err := plotter.NewHist(v, 150)
	if err != nil {
		panic(err)
	}
	return h
}

// Plots the historgram using plotinum
func Histogram(r RetCalc) {
	//eb := all_paths.End_balances()
	eb := make([]float64, len(r.All_paths), len(r.All_paths))
	incs := r.RunIncomes()
	for i := range incs {
		eb[i] = incs[i]
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
	//h.Normalize(1)
	p.Add(h)

	if err := p.Save(4, 4, "hist.png"); err != nil {
		panic(err)
	}
	fmt.Println(h)
}

// Runs a path
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
