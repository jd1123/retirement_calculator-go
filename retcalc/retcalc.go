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
func (r RetCalc) RunAllPaths() PathGroup {
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
		untaxed_total_wealth := r.Non_Taxable_balance * r.sims[i].GrowthFactor(0)
		taxed_total_wealth := r.Taxable_balance * r.sims[i].GrowthFactorWithTaxes(0, r.Effective_tax_rate)
		ft, _ := r.All_paths[i].Factors()
		fmt.Printf("Income from taxed accts starting: %f nt accts: %f\n", taxed_total_wealth/ft, (untaxed_total_wealth/(1+r.Effective_tax_rate))/ft)
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
	sort.Float64s(incomes)
	return incomes
}

func (r RetCalc) PercentileIncome(percentile float64) float64 {
	incomes := r.RunIncomes()
	ix := int(percentile * float64(r.N))
	return incomes[ix]
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
	fmt.Printf("Counter: %d Income: %f N: %d\n", counter, r.Yearly_retirement_expenses, r.N)
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
	}

	return r
}

// A default RetCalc
func NewRetCalc() RetCalc {
	r := RetCalc{}

	r.N = 20000
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

	r.N = 20000
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
