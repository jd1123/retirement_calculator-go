// Testing suite for retcalc.go
// going to test the individual functions to see if they
// work with known values against an excel spreadsheet
// NOTE: I don't know if there is a setUp, tearDown functionality like
// in python unittest, so there is a lot of repeated code and the tests
// take longer than they should
package retcalc

import (
	"fmt"
	"math"
	"sort"
	"testing"
)

// Thresholds for testing floating point and integer results
// against known values - ints resulting from a cast
const FP_THRESHOLD = 0.1
const INT_THRESHOLD = 250

// Testing the json constructor
func TestNewRetCalcJSon(t *testing.T) {
	// Set up json byte array
	JsonObj := []byte(`{"Age":22, "Retirement_age":65, "Terminal_age":90, "Effective_tax_rate":0.3, "Returns_tax_rate":0.3, "N": 20000, 
					"Non_Taxable_contribution":17500, "Taxable_contribution": 0, "Non_Taxable_balance":0, "Taxable_balance": 0, 
					"Yearly_social_security_income":0, "Asset_volatility": 0.15, "Expected_rate_of_return": 0.07, "Inflation_rate":0.035}`)

	rc := NewRetCalcFromJSON(JsonObj)
	if rc.Age != 22 && rc.Retirement_age != 65 && rc.Terminal_age != 90 && rc.Effective_tax_rate != 0.3 && rc.N != 20000 {
		t.Errorf("json did not initialize the retcalc object correctly: values do not match known values")
	}
	if rc.Years != rc.Terminal_age-rc.Age+1 {
		t.Errorf("RetCalc.Years did not intitialize correctly: Years did not computer correctly")
	}

	if len(rc.Sims) != rc.N {
		t.Errorf("RetCalc.sims does not have the correct length")
	}
	if len(rc.Sims[0]) != rc.Years {
		t.Errorf("RetCalc.sims[i] does not have the correct length")
	}
	if rc.Sims == nil {
		t.Errorf("Retcalc.sims not initialized")
	}
	if rc.IncomeOnPath(1) < 0.0 {
		t.Errorf("Retcalc.IncomeOnPath() does not work for newly init NewRetCalcFromJSON()")
	}

}

func TestNewRetCalcJSon_withIncompleteInput(t *testing.T) {
	// Set up json byte array
	JsonObj := []byte(`{"Age":22, "Retirement_age":65, "Terminal_age":90, "Effective_tax_rate":0.3, "Returns_tax_rate":0.3, "N": 20000, 
					"Non_Taxable_contribution":17500, "Taxable_contribution": 0, "Non_Taxable_balance":0, "Taxable_balance": 0}`)

	rc := NewRetCalcFromJSON(JsonObj)
	if rc.Age != 22 && rc.Retirement_age != 65 && rc.Terminal_age != 90 && rc.Effective_tax_rate != 0.3 && rc.N != 20000 {
		t.Errorf("json did not initialize the retcalc object correctly: values do not match known values")
	}
	if rc.Years != rc.Terminal_age-rc.Age+1 {
		t.Errorf("RetCalc.Years did not intitialize correctly: Years did not computer correctly")
	}

	if len(rc.Sims) != rc.N {
		t.Errorf("RetCalc.sims does not have the correct length")
	}
	if len(rc.Sims[0]) != rc.Years {
		t.Errorf("RetCalc.sims[i] does not have the correct length")
	}
	if rc.Sims == nil {
		t.Errorf("Retcalc.sims not initialized")
	}
	if rc.IncomeOnPath(1) < 0.0 {
		t.Errorf("Retcalc.IncomeOnPath() does not work for newly init NewRetCalcFromJSON()")
	}

}

// Test the components of income calculation
func TestGrowthFactors(t *testing.T) {
	JsonObj := []byte(`{"Age":22, "Retirement_age":65, "Terminal_age":90, "Effective_tax_rate":0.3, "Returns_tax_rate":0.3, "N": 20000, 
						"Non_Taxable_contribution":17500, "Taxable_contribution": 0, "Non_Taxable_balance":0, "Taxable_balance": 0, 
						"Yearly_social_security_income":0, "Asset_volatility": 0.15, "Expected_rate_of_return": 0.07, "Inflation_rate":0.035}`)
	r := NewRetCalcFromJSON(JsonObj)
	r.Sims[0] = make([]float64, r.Years, r.Years)
	for i := range r.Sims[0] {
		r.Sims[0][i] = 0.07
	}

	gfactors := make([]float64, len(r.Sims[0]), len(r.Sims[0]))
	for i := range r.Sims[0] {
		gfactors[i] = r.Sims[0].GrowthFactor(i)
	}

	known_gfactors := []float64{99.5627497577374, 93.0492988390069, 86.9619615317821, 81.2728612446562, 75.9559450884637, 70.9868645686577, 66.3428640828576,
		62.0026767129511, 57.9464268345337, 54.1555390976951, 50.6126533623318, 47.301545198441, 44.207051587328, 41.3150014834841, 38.612150919144,
		36.0861223543402, 33.7253479947105, 31.5190168174864, 29.4570250630714, 27.5299299654873, 25.7289065098012, 24.0457070185058, 22.4726233817811,
		21.002451758674, 19.6284595875457, 18.3443547547157, 17.1442567801081, 16.0226698879515, 14.974457839207, 13.9948204104738, 13.0792714116578,
		12.2236181417362, 11.4239421885385, 10.6765814846155, 9.9781135370238, 9.32533975422784, 8.71527079834378, 8.14511289564839, 7.61225504266205,
		7.11425704921687, 6.64883836375408, 6.21386762967671, 5.8073529249315, 5.4274326401229, 5.07236695338589, 4.74052986297746,
		4.43040174110043, 4.14056237486022, 3.86968446248618, 3.61652753503382, 3.37993227573254, 3.15881521096499, 2.95216374856541,
		2.75903154071534, 2.57853415020125, 2.40984500018808, 2.25219158896082, 2.10485195229984, 1.96715135728957, 1.83845921242016,
		1.71818617983192, 1.60578147647843, 1.500730351849, 1.4025517307, 1.31079601, 1.225043, 1.1449, 1.07, 1}

	factorsGood := true
	for i := range gfactors {
		if math.Abs(gfactors[i]-known_gfactors[i]) > 0.1 {
			factorsGood = false
		}
	}
	if !factorsGood {
		t.Errorf("Growth Factor calculations incorrect")
	}
}

func TestGrowthFactorsWithTaxes(t *testing.T) {
	knownGrowthFactors := []float64{25.8662337047586, 24.6579920922389, 23.506188839122, 22.4081876445395, 21.3614753522779, 20.363656198549, 19.4124463284547,
		18.5056685685936, 17.6412474438452, 16.8172044269259, 16.0316534098435, 15.2827963868861, 14.5689193392623, 13.8883883119755,
		13.2396456739518, 12.6212065528616, 12.0316554364744, 11.4696429327687, 10.933882681381, 10.4231484093242, 9.93627112423656,
		9.47213643873838, 9.02968201976967, 8.60789515707309, 8.20581044525556, 7.82250757412351, 7.45710922223404, 7.10877904884084,
		6.77671977963855, 6.46017138192427, 6.1584093249993, 5.87074292182965, 5.59651374816935, 5.33509413552846, 5.08588573453618,
		4.84831814541104, 4.62184761240328, 4.40595577922143, 4.20014850259431, 4.00395472125292, 3.81692537774349, 3.6386323906039,
		3.4686676745509, 3.30664220643556, 3.15218513482894, 3.00494293120014, 2.8645785807437, 2.73077081100448, 2.6032133565343,
		2.48161425789733, 2.36569519341976, 2.2551908421542, 2.14984827660077, 2.04942638379482, 1.95369531343643, 1.86243595179832,
		1.77543942020812, 1.69250659695722, 1.61344766154167, 1.53808166019225, 1.466236091699, 1.39774651258246, 1.33245616070778,
		1.27021559648025, 1.210882360801, 1.154320649, 1.100401, 1.049, 1}
	JsonObj := []byte(`{"Age":22, "Retirement_age":65, "Terminal_age":90, "Effective_tax_rate":0.3, "Returns_tax_rate":0.3, "N": 20000, 
						"Non_Taxable_contribution":17500, "Taxable_contribution": 0, "Non_Taxable_balance":0, "Taxable_balance": 0, 
						"Yearly_social_security_income":0, "Asset_volatility": 0.15, "Expected_rate_of_return": 0.07, "Inflation_rate":0.035}`)
	r := NewRetCalcFromJSON(JsonObj)
	r.Sims[0] = make([]float64, r.Years, r.Years)
	for i := range r.Sims[0] {
		r.Sims[0][i] = 0.07
	}

	gfactors := make([]float64, len(r.Sims[0]), len(r.Sims[0]))
	for i := range r.Sims[0] {
		gfactors[i] = r.Sims[0].GrowthFactorWithTaxes(i, r.Effective_tax_rate)
	}
	factorsGood := true
	for i := range gfactors {
		if math.Abs(gfactors[i]-knownGrowthFactors[i]) > FP_THRESHOLD {
			factorsGood = false
			fmt.Printf("gfactor: %f    knownFactor: %f\n", gfactors[i], knownGrowthFactors[i])
		}
	}
	if !factorsGood {
		t.Errorf("Growth Factor calculations incorrect")
	}
}

func TestRetCalcGrowthFactors(t *testing.T) {
	JsonObj := []byte(`{"Age":22, "Retirement_age":65, "Terminal_age":90, "Effective_tax_rate":0.3, "Returns_tax_rate":0.3, "N": 20000, 
						"Non_Taxable_contribution":17500, "Taxable_contribution": 0, "Non_Taxable_balance":0, "Taxable_balance": 0, 
						"Yearly_social_security_income":0, "Asset_volatility": 0.15, "Expected_rate_of_return": 0.07, "Inflation_rate":0.035}`)
	r := NewRetCalcFromJSON(JsonObj)
	retCalcFactor := r.GrowthFactor(0, 0)
	simFactor := r.Sims[0].GrowthFactor(0)
	if retCalcFactor != simFactor {
		t.Errorf("RetCalc.GrowthFactors() does not return the same result as RetCalc.sim[0].GrowthFactor()")
	}
}

// Testing RetCalc.InflationFactors() against known values
func TestInflationFactors(t *testing.T) {
	knownInflationFactors := []float64{1.035, 1.071225, 1.108717875, 1.147523000625, 1.18768630564687, 1.22925532634452, 1.27227926276657, 1.3168090369634,
		1.36289735325712, 1.41059876062112, 1.45996971724286, 1.51106865734636, 1.56395606035348, 1.61869452246585, 1.67534883075216,
		1.73398603982848, 1.79467555122248, 1.85748919551527, 1.9225013173583, 1.98978886346584, 2.05943147368715, 2.1315115752662,
		2.20611448040051, 2.28332848721453, 2.36324498426704, 2.44595855871639, 2.53156710827146, 2.62017195706096, 2.71187797555809,
		2.80679370470263, 2.90503148436722, 3.00670758632007, 3.11194235184127, 3.22086033415572, 3.33359044585117, 3.45026611145596,
		3.57102542535692, 3.69601131524441, 3.82537171127796, 3.95925972117269, 4.09783381141374, 4.24125799481322, 4.38970202463168,
		4.54334159549379, 4.70235855133607, 4.86694110063283, 5.03728403915498, 5.2135889805254, 5.39606459484379, 5.58492685566332,
		5.78039929561154, 5.98271327095794, 6.19210823544147, 6.40883202368192, 6.63314114451079, 6.86530108456866, 7.10558662252857,
		7.35428215431707, 7.61168202971816, 7.8780909007583, 8.15382408228484, 8.43920792516481, 8.73458020254557, 9.04029050963467,
		9.35670067747188, 9.6841852011834, 10.0231316832248, 10.3739412921377, 10.7370292373625}

	JsonObj := []byte(`{"Age":22, "Retirement_age":65, "Terminal_age":90, "Effective_tax_rate":0.3, "Returns_tax_rate":0.3, "N": 20000, 
						"Non_Taxable_contribution":17500, "Taxable_contribution": 0, "Non_Taxable_balance":0, "Taxable_balance": 0, 
						"Yearly_social_security_income":0, "Asset_volatility": 0.15, "Expected_rate_of_return": 0.07, "Inflation_rate":0.035}`)

	r := NewRetCalcFromJSON(JsonObj)
	r.Sims[0] = make([]float64, r.Years, r.Years)
	for i := range r.Sims[0] {
		r.Sims[0][i] = 0.07
	}
	inflationFactors := r.InflationFactors()
	for i := range inflationFactors {
		if math.Abs(inflationFactors[i]-knownInflationFactors[i]) > FP_THRESHOLD {
			t.Errorf("InflationFactors() does not produce known factors")
		}
	}
}

// Good test, I think - caught a bug
func TestIncomeFactors(t *testing.T) {
	JsonObj := []byte(`{"Age":22, "Retirement_age":65, "Terminal_age":90, "Effective_tax_rate":0.3, "Returns_tax_rate":0.3, "N": 20000, 
						"Non_Taxable_contribution":17500, "Taxable_contribution": 0, "Non_Taxable_balance":0, "Taxable_balance": 0, 
						"Yearly_social_security_income":0, "Asset_volatility": 0.15, "Expected_rate_of_return": 0.07, "Inflation_rate":0.035}`)

	r := NewRetCalcFromJSON(JsonObj)
	r.Sims[0] = make([]float64, r.Years, r.Years)
	for i := range r.Sims[0] {
		r.Sims[0][i] = 0.07
	}

	pth := RunPath(r, r.Sims[0])
	sum, factors := r.IncomeFactors(0)
	knownSum := 411.683115
	knownFactors := []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 23.8520881187686, 23.0718796289023, 22.3171919774896, 21.5871903707493, 20.8810673212388, 20.1980417546562, 19.5373581458591,
		18.8982856831441, 18.2801174598636, 17.6821696924849, 17.1037809642261, 16.5443114934336, 16.0031424258914, 15.4796751502781,
		14.9733306360167, 14.4835487927825, 14.0097878509625, 13.5515237623796, 13.1082496206195, 12.6794751003189, 12.2647259147944,
		11.8635432914133, 11.4754834641241, 11.1001171825873, 10.7370292373625}

	if math.Abs(knownSum-sum) > FP_THRESHOLD {
		fmt.Printf("sum: %f  knownSum: %f\n", sum, knownSum) //Comment this out after testing
		t.Errorf("Income factor sum from path.Factors() not computing correctly for known values")
	}

	factorsGood := true
	for i := range factors {
		if math.Abs(factors[i]-knownFactors[i]) > FP_THRESHOLD {
			factorsGood = false
			fmt.Printf("factors: %f  knownFactors: %f sim: %f\n", factors[i], knownFactors[i], pth.Sim[i]) //Comment this after testing
		}
	}
	if !factorsGood {
		t.Errorf("Factors returned from path.Factors() do not match known factors")
	}
}

func TestIncomeFactorsWithTaxes(t *testing.T) {
	JsonObj := []byte(`{"Age":22, "Retirement_age":65, "Terminal_age":90, "Effective_tax_rate":0.3, "Returns_tax_rate":0.3, "N": 20000, 
						"Non_Taxable_contribution":17500, "Taxable_contribution": 0, "Non_Taxable_balance":0, "Taxable_balance": 0, 
						"Yearly_social_security_income":0, "Asset_volatility": 0.15, "Expected_rate_of_return": 0.07, "Inflation_rate":0.035}`)

	r := NewRetCalcFromJSON(JsonObj)
	r.Sims[0] = make([]float64, r.Years, r.Years)
	for i := range r.Sims[0] {
		r.Sims[0][i] = 0.07
	}

	pth := RunPath(r, r.Sims[0])
	sum, factors := r.IncomeFactorsWithTaxes(0)
	knownSum := 316.87
	knownFactors := []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 14.8227047241573, 14.624880256914, 14.4296959636854, 14.2371166085934, 14.0471074260192,
		13.8596341143278, 13.6746628296752, 13.4921601798988, 13.3120932184893, 13.1344294386429, 12.9591367673931,
		12.7861835598207, 12.6155385933407, 12.4471710620664, 12.2810505712476, 12.1171471317838, 11.9554311548105,
		11.7958734463574, 11.6384452020781, 11.4831180020503, 11.3298638056455, 11.1786549464662, 11.0294641273523,
		10.8822644154524, 10.7370292373625}

	if math.Abs(knownSum-sum) > FP_THRESHOLD {
		fmt.Printf("sum: %f  knownSum: %f\n", sum, knownSum) //Comment this out after testing
		t.Errorf("Income factor sum from path.Factors() not computing correctly for known values")
	}

	factorsGood := true
	for i := range factors {
		if math.Abs(factors[i]-knownFactors[i]) > FP_THRESHOLD {
			factorsGood = false
			fmt.Printf("factors: %f  knownFactors: %f sim: %f\n", factors[i], knownFactors[i], pth.Sim[i]) //Comment this after testing
		}
	}
	if !factorsGood {
		t.Errorf("Factors returned from path.Factors() do not match known factors")
	}
}

// After testing the components that go into the calculations, test that incomes
// are computed correctly
func TestIncomeOnPath(t *testing.T) {
	JsonObj := []byte(`{"Age":22, "Retirement_age":65, "Terminal_age":90, "Effective_tax_rate":0.3, "Returns_tax_rate":0.3, "N": 200, 
						"Non_Taxable_contribution":17500, "Taxable_contribution": 0, "Non_Taxable_balance":0, "Taxable_balance": 0, 
						"Yearly_social_security_income":0, "Asset_volatility": 0.15, "Expected_rate_of_return": 0.07, "Inflation_rate":0.035}`)
	rc := NewRetCalcFromJSON(JsonObj)

	knownSim := make([]float64, rc.Years, rc.Years)
	for i := range knownSim {
		knownSim[i] = 0.07
	}

	rc.SetSim(0, knownSim)

	knownIncome := 42978.0
	if math.Abs(rc.IncomeOnPath(0)-knownIncome) > INT_THRESHOLD {
		t.Errorf("IncomeOnPath() does not match known value")
		fmt.Printf("IncomeOnPath(): %f --- knownIncome: %f\n", rc.IncomeOnPath(0), knownIncome)
	}
	knownIncomeTaxed := 26885.0
	JsonObj = []byte(`{"Age":22, "Retirement_age":65, "Terminal_age":90, "Effective_tax_rate":0.3, "Returns_tax_rate":0.3, "N": 200, 
						"Non_Taxable_contribution":0, "Taxable_contribution": 17500, "Non_Taxable_balance":0, "Taxable_balance": 0, 
						"Yearly_social_security_income":0, "Asset_volatility": 0.15, "Expected_rate_of_return": 0.07, "Inflation_rate":0.035}`)
	rc = NewRetCalcFromJSON(JsonObj)
	rc.SetSim(0, knownSim)
	if math.Abs(rc.IncomeOnPath(0)-knownIncomeTaxed) > INT_THRESHOLD {
		t.Errorf("IncomeOnPath() does not match known value using Taxable Accounts")
		fmt.Printf("IncomeOnPath(): %f --- knownIncomeTaxed: %f\n", rc.IncomeOnPath(0), knownIncomeTaxed)
	}

	knownIncome = 16929.0
	JsonObj = []byte(`{"Age":22, "Retirement_age":65, "Terminal_age":90, "Effective_tax_rate":0.3, "Returns_tax_rate":0.3, "N": 200, 
						"Non_Taxable_contribution":0, "Taxable_contribution": 0, "Non_Taxable_balance":100000, "Taxable_balance": 0, 
						"Yearly_social_security_income":0, "Asset_volatility": 0.15, "Expected_rate_of_return": 0.07, "Inflation_rate":0.035}`)
	rc = NewRetCalcFromJSON(JsonObj)
	rc.SetSim(0, knownSim)
	if math.Abs(rc.IncomeOnPath(0)-knownIncome) > INT_THRESHOLD {
		t.Errorf("IncomeOnPath() does not match known value using Non-Taxable starting balance")
		fmt.Printf("IncomeOnPath(): %f --- knownIncome: %f\n", rc.IncomeOnPath(0), knownIncome)
	}

	knownIncome = 8363.0
	JsonObj = []byte(`{"Age":22, "Retirement_age":65, "Terminal_age":90, "Effective_tax_rate":0.3, "Returns_tax_rate":0.3, "N": 200, 
						"Non_Taxable_contribution":0, "Taxable_contribution": 0, "Non_Taxable_balance":0, "Taxable_balance": 100000, 
						"Yearly_social_security_income":0, "Asset_volatility": 0.15, "Expected_rate_of_return": 0.07, "Inflation_rate":0.035}`)
	rc = NewRetCalcFromJSON(JsonObj)
	rc.SetSim(0, knownSim)
	if math.Abs(rc.IncomeOnPath(0)-knownIncome) > INT_THRESHOLD {
		t.Errorf("IncomeOnPath() does not match known value using Non-Taxable starting balance")
		fmt.Printf("IncomeOnPath(): %f --- knownIncome: %f\n", rc.IncomeOnPath(0), knownIncome)
	}
}

// Test SetSim Methods
func TestSetSim(t *testing.T) {
	JsonObj := []byte(`{"Age":22, "Retirement_age":65, "Terminal_age":90, "Effective_tax_rate":0.3, "Returns_tax_rate":0.3, "N": 20000, 
						"Non_Taxable_contribution":17500, "Taxable_contribution": 0, "Non_Taxable_balance":0, "Taxable_balance": 0, 
						"Yearly_social_security_income":0, "Asset_volatility": 0.15, "Expected_rate_of_return": 0.07, "Inflation_rate":0.035}`)
	rc := NewRetCalcFromJSON(JsonObj)
	s := make([]float64, len(rc.Sims[0]), len(rc.Sims[0]))
	for i := range s {
		s[i] = 0.07
	}
	rc.SetSim(0, s)
	for i := range rc.Sims[0] {
		if rc.Sims[0][i] != s[i] {
			t.Errorf("RetCalc.sim does not set correctly with RetCalc.SetSim()")
		}
	}
}

func TestRunIncomes(t *testing.T) {
	JsonObj := []byte(`{"Age":22, "Retirement_age":65, "Terminal_age":90, "Effective_tax_rate":0.3, "Returns_tax_rate":0.3, "N": 20000, 
						"Non_Taxable_contribution":17500, "Taxable_contribution": 0, "Non_Taxable_balance":0, "Taxable_balance": 0, 
						"Yearly_social_security_income":0, "Asset_volatility": 0.15, "Expected_rate_of_return": 0.07, "Inflation_rate":0.035}`)
	rc := NewRetCalcFromJSON(JsonObj)
	runIncomes := rc.RunIncomes()
	sort.Float64s(runIncomes)
	incomePerRun := make([]float64, rc.N, rc.N)
	for i := range incomePerRun {
		incomePerRun[i] = rc.IncomeOnPath(i)
	}
	sort.Float64s(incomePerRun)
	incomesOk := true
	for i := range incomePerRun {
		if incomePerRun[i] != runIncomes[i] {
			incomesOk = false
			fmt.Printf("RunIncomes: %f, IncomeOnPath: %f\n", runIncomes[i], incomePerRun[i])
		}
	}
	if !incomesOk {
		t.Errorf("Incomes do not calculate correctly for RunIncomes()")
	}
	if !sort.Float64sAreSorted(runIncomes) {
		t.Errorf("Incomes from RetCalc.RunIncomes() should be sorted on return")
	}
}

func TestPercentileIncome(t *testing.T) {
	knownLow := 2942.0
	knownHigh := 42978.0
	JsonObj := []byte(`{"Age":22, "Retirement_age":65, "Terminal_age":90, "Effective_tax_rate":0.3, "Returns_tax_rate":0.3, "N": 20, 
						"Non_Taxable_contribution":17500, "Taxable_contribution": 0, "Non_Taxable_balance":0, "Taxable_balance": 0, 
						"Yearly_social_security_income":0, "Asset_volatility": 0.15, "Expected_rate_of_return": 0.07, "Inflation_rate":0.035}`)
	rc := NewRetCalcFromJSON(JsonObj)
	dudSim := make([]float64, rc.Years, rc.Years)
	detSim := make([]float64, rc.Years, rc.Years)
	for i := range dudSim {
		detSim[i] = 0.07
		dudSim[i] = 0.0
	}
	for i := 0; i < rc.N; i++ {
		if i%2 == 0 {
			rc.SetSim(i, dudSim)
		} else {
			rc.SetSim(i, detSim)
		}
	}
	if math.Abs(rc.PercentileIncome(0.25)-knownLow) > INT_THRESHOLD {
		t.Errorf("PercentileIncome should match known value")
	}

	if math.Abs(rc.PercentileIncome(0.75)-knownHigh) > INT_THRESHOLD {
		t.Errorf("PercentileIncome should match known value")
	}

}

/*
func TestRunIncomesAnalytics(t *testing.T) {
	JsonObj := []byte(`{"Age":22, "Retirement_age":65, "Terminal_age":90, "Effective_tax_rate":0.3, "Returns_tax_rate":0.3, "N": 20000,
						"Non_Taxable_contribution":17500, "Taxable_contribution": 0, "Non_Taxable_balance":0, "Taxable_balance": 0,
						"Yearly_social_security_income":0, "Asset_volatility": 0.15, "Expected_rate_of_return": 0.07, "Inflation_rate":0.035}`)
	rc := NewRetCalcFromJSON(JsonObj)
	runIncomes := rc.RunIncomes()
	sort.Float64s(runIncomes)
	min := 100000.0
	max := 0.0
	avg := 0.0
	for i := range runIncomes {
		avg += runIncomes[i]
		if runIncomes[i] > max {
			max = runIncomes[i]
		}
		if runIncomes[i] < min {
			min = runIncomes[i]
		}
	}
	avg = avg / float64(len(runIncomes))
	fmt.Printf("Max: %f\n", max)
	fmt.Printf("Min %f\n", min)
	fmt.Printf("Avg: %f\n", avg)
}
*/
func TestPathFinalBalance(t *testing.T) {
	knownFinalBalance := 25276177.0
	JsonObj := []byte(`{"Age":22, "Retirement_age":65, "Terminal_age":90, "Effective_tax_rate":0.3, "Returns_tax_rate":0.3, "N": 20, 
						"Non_Taxable_contribution":17500, "Taxable_contribution": 0, "Non_Taxable_balance":0, "Taxable_balance": 0, 
						"Yearly_social_security_income":0, "Asset_volatility": 0.15, "Expected_rate_of_return": 0.07, "Inflation_rate":0.035}`)
	rc := NewRetCalcFromJSON(JsonObj)
	detSim := make([]float64, rc.Years, rc.Years)
	for i := range detSim {
		detSim[i] = 0.07
	}
	rc.SetSim(0, detSim)
	rc.all_paths[0] = RunPath(rc, rc.Sims[0])
	fb := rc.all_paths[0].Final_balance()
	if math.Abs(knownFinalBalance-fb) > INT_THRESHOLD {
		t.Errorf("Path.FinalBalance() does not match known value")
		fmt.Printf("Path.FinalBalance(): %f -- Known Final Balance: %f", fb, knownFinalBalance)
		rc.all_paths[0].Print_path()
	}
}

func TestPathEndBalances(t *testing.T) {
	JsonObj := []byte(`{"Age":22, "Retirement_age":65, "Terminal_age":90, "Effective_tax_rate":0.3, "Returns_tax_rate":0.3, "N": 20, 
						"Non_Taxable_contribution":17500, "Taxable_contribution": 0, "Non_Taxable_balance":0, "Taxable_balance": 0, 
						"Yearly_social_security_income":0, "Asset_volatility": 0.15, "Expected_rate_of_return": 0.07, "Inflation_rate":0.035}`)
	rc := NewRetCalcFromJSON(JsonObj)
	eb := rc.all_paths.End_balances()
	for i := 0; i < rc.N; i++ {
		if math.Abs(rc.all_paths[i].Final_balance()-eb[i]) > INT_THRESHOLD {
			t.Errorf("Path.FinalBalance() does not match known value")
			fmt.Printf("Path.FinalBalance(): %f -- Known Final Balance: %f", rc.all_paths[i].Final_balance(), eb[i])
			rc.all_paths[0].Print_path()
		}
	}
}
