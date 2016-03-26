package retcalc

import (
	"encoding/json"
	"fmt"
	"math"
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
	Sims                                           []Sim
	Non_Taxable_contribution, Taxable_contribution float64
	Non_Taxable_balance, Taxable_balance           float64
	Yearly_retirement_expenses                     float64
	PortfolioSelection                             Portfolio
	PortfolioString                                string
	Inflation_rate                                 float64
	all_paths                                      PathGroup
	SessionId                                      string
}

// METHODS
func (r RetCalc) ShowRetCalc() {
	fmt.Println("----Showing you the RetCalc----")
	fmt.Printf("Age: %d\n", r.Age)
	fmt.Printf("Retirement Age: %d\n", r.Retirement_age)
	fmt.Printf("Terminal Age: %d\n", r.Terminal_age)
	fmt.Printf("Years: %d\n", r.Years)
	fmt.Printf("Non_Taxable_contribution %f\n", r.Non_Taxable_contribution)
	fmt.Printf("Taxable_contribution %f\n", r.Taxable_contribution)
	fmt.Printf("Inflation_rate %f\n", r.Inflation_rate)
	fmt.Println("Portfolio Selection:", r.PortfolioSelection)
	fmt.Println()
}

// Runs all the paths - dont do this in the constructor
// because the data will be lost when passing a smaller,
// client side object that is portable
func (r RetCalc) RunAllPaths() PathGroup {
	all_paths := make(PathGroup, len(r.Sims), len(r.Sims))
	pathChan := make(chan Path)
	for i := range r.Sims {
		go func() {
			pathChan <- RunPath(r, r.Sims[i])
			//all_paths[i] = RunPath(r, r.Sims[i])
		}()
	}
	//sort.Sort(all_paths)
	i := 0
	for j := range pathChan {
		all_paths[i] = j
		i++
		if i == len(r.Sims) {
			close(pathChan)
		}
	}
	return all_paths
}

// The path struct should implement this logic, it is misplaced
func (r RetCalc) RunIncomes() []float64 {
	incomes := make([]float64, len(r.Sims), len(r.Sims))
	for i := range r.Sims {
		untaxed_total_wealth := r.Non_Taxable_balance * r.Sims[i].GrowthFactor(0)
		taxed_total_wealth := r.Taxable_balance * r.Sims[i].GrowthFactorWithTaxes(0, r.Effective_tax_rate)
		//ft, _ := r.All_paths[i].Factors()
		//fmt.Printf("Income from taxed accts starting: %f nt accts: %f\n", taxed_total_wealth/ft,
		//	(untaxed_total_wealth/(1+r.Effective_tax_rate))/ft)
		sum_t, sum_ut := 0.0, 0.0
		for j := range r.Sims[i] {
			if j+r.Age < r.Retirement_age {
				untaxed_total_wealth += r.Non_Taxable_contribution * r.Sims[i].GrowthFactor(j)
				taxed_total_wealth += r.Taxable_contribution * r.Sims[i].GrowthFactorWithTaxes(j, r.Effective_tax_rate)
				sum_ut += r.Sims[i].GrowthFactor(j)
				sum_t += r.Sims[i].GrowthFactorWithTaxes(j, r.Effective_tax_rate)
			}
		}
		f, _ := r.IncomeFactors(i)
		ft, _ := r.IncomeFactorsWithTaxes(i)
		incomes[i] = (taxed_total_wealth/ft + untaxed_total_wealth*(1-r.Effective_tax_rate)/f)
	}
	sort.Float64s(incomes)
	return incomes
}

// returns the income for a specific path in all_paths
func (r RetCalc) IncomeOnPath(pathIndex int) float64 {
	untaxed_total_wealth := r.Non_Taxable_balance * r.Sims[pathIndex].GrowthFactor(0)
	taxed_total_wealth := r.Taxable_balance * r.Sims[pathIndex].GrowthFactorWithTaxes(0, r.Effective_tax_rate)
	sum_t, sum_ut := 0.0, 0.0
	for j := range r.Sims[pathIndex] {
		if j+r.Age < r.Retirement_age {
			untaxed_total_wealth += r.Non_Taxable_contribution * r.Sims[pathIndex].GrowthFactor(j)
			taxed_total_wealth += r.Taxable_contribution * r.Sims[pathIndex].GrowthFactorWithTaxes(j, r.Effective_tax_rate)
			sum_ut += r.Sims[pathIndex].GrowthFactor(j)
			sum_t += r.Sims[pathIndex].GrowthFactorWithTaxes(j, r.Effective_tax_rate)
		}
	}
	f, _ := r.IncomeFactors(pathIndex)
	ft, _ := r.IncomeFactorsWithTaxes(pathIndex)
	income := (taxed_total_wealth/ft + untaxed_total_wealth*(1-r.Effective_tax_rate)/f)
	return income
}

// returns the minimum income that you will have with percentile probability
func (r RetCalc) PercentileIncome(percentile float64) float64 {
	incomes := r.RunIncomes()
	ix := int(percentile * float64(r.N))
	return incomes[ix]
}

// returns the path of the minumum income with percentile probability
func (r RetCalc) PercentilePath(percentile float64) Path {
	ix := int(float64(r.N) * percentile)
	return r.all_paths[ix]
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

// Sets a simulation - used for testing
func (r RetCalc) SetSim(ix int, newSim []float64) {
	for i := range r.Sims[ix] {
		r.Sims[ix][i] = newSim[i]
	}
	r.all_paths[ix] = RunPath(r, r.Sims[ix])
}

// returns the inflation factors for the retcalc object
func (r RetCalc) InflationFactors() []float64 {
	inflationFactors := make([]float64, r.Years, r.Years)
	for i := 0; i < r.Years; i++ {
		inflationFactors[i] = math.Pow(1+r.Inflation_rate, float64(i+1))
	}
	return inflationFactors
}

// returns the income factors for the retcalc objects for sim simIdx
func (r RetCalc) IncomeFactors(simIdx int) (float64, []float64) {
	sim := r.Sims[0]
	l := len(sim)
	incomeFactors := make([]float64, l, l)
	inflationFactors := r.InflationFactors()
	sumFactors := 0.0
	for i := range incomeFactors {
		if r.Age+i > r.Retirement_age {
			incomeFactors[i] = r.Sims[simIdx].GrowthFactor(i) * inflationFactors[i]
			sumFactors += incomeFactors[i]
		} else {
			incomeFactors[i] = 0.0
		}
	}

	return sumFactors, incomeFactors
}

// returns the income factors with taxes for sim simIdx
func (r RetCalc) IncomeFactorsWithTaxes(simIdx int) (float64, []float64) {
	sim := r.Sims[0]
	l := len(sim)
	incomeFactors := make([]float64, l, l)
	inflationFactors := r.InflationFactors()
	sumFactors := 0.0
	for i := range incomeFactors {
		if r.Age+i > r.Retirement_age {
			incomeFactors[i] = r.Sims[simIdx].GrowthFactorWithTaxes(i, r.Effective_tax_rate) * inflationFactors[i]
			sumFactors += incomeFactors[i]
		} else {
			incomeFactors[i] = 0.0
		}
	}

	return sumFactors, incomeFactors
}

// returns the growth factor for starting year startYear and sim simIdx
func (r RetCalc) GrowthFactor(startYear, simIdx int) float64 {
	return r.Sims[simIdx].GrowthFactor(startYear)
}

// Constructors

// This constructor will populate a RetCalc from
// JSON input from the web ----- NEEDS work
// FIXME: This should return an error value
func NewRetCalcFromJSON(json_obj []byte) RetCalc {
	var r RetCalc
	err := json.Unmarshal(json_obj, &r)
	if err != nil {
		fmt.Println("ERROR in NewRetCalcFromJSON()")
		fmt.Println(err)
	}
	r.Years = r.Terminal_age - r.Age + 1
	if r.N == 0 {
		r.N = 10000
	}
	r.PortfolioSelection = PortfolioStrings[r.PortfolioString]
	/*
		if r.PortfolioSelection == BLANKPORTFOLIO {
			//r.PortfolioSelection = HIGHRISKPORTFOLIO
			r.PortfolioSelection = LOWRISKPORTFOLIO
		}
	*/
	if r.Inflation_rate == 0 {
		r.Inflation_rate = 0.035
	}

	r.Sims = make([]Sim, r.N, r.N)
	for i := range r.Sims {
		r.Sims[i] = Simulation(r.PortfolioSelection, r.Years)
	}
	r.all_paths = r.RunAllPaths()

	return r
}

// A very basic RetCalc with standard info
func NewRetCalc() RetCalc {
	r := RetCalc{}

	r.N = 20000
	r.Age = 22
	r.Retirement_age = 65
	r.Terminal_age = 90
	r.Years = r.Terminal_age - r.Age + 1
	r.Effective_tax_rate = 0.30
	r.Returns_tax_rate = 0.30
	r.Non_Taxable_contribution = 17500.0
	r.Taxable_contribution = 0
	r.Non_Taxable_balance = 0
	r.Yearly_retirement_expenses = float64(100000)
	r.Taxable_balance = 0.0
	r.PortfolioSelection = HIGHRISKPORTFOLIO
	r.Inflation_rate = 0.035
	r.Sims = make([]Sim, r.N, r.N)
	r.all_paths = make([]Path, r.N, r.N)
	for i := range r.Sims {
		r.Sims[i] = Simulation(r.PortfolioSelection, r.Years)
		r.all_paths[i] = RunPath(r, r.Sims[i])
	}
	sort.Sort(r.all_paths)

	return r
}
