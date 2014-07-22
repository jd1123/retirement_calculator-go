// Testing suite for retcalc.go
// going to test the individual functions to see if they
// work with known values against an excel spreadsheet
// NOTE: I don't know if there is a setUp, tearDown functionality like
// in python unittest, so there is a lot of repeated code and the tests
// take longer than they should
package retcalc

import (
	"math"
	"testing"
)

// Testing the json constructor
func TestNewRetCalcJSon(t *testing.T) {
	// Set up json byte array
	JsonObj := []byte(`{"Age":22, "Retirement_age":65, "Terminal_age":90, "Effective_tax_rate":0.3, "Returns_tax_rate":0.3, "N": 20000, "Non_Taxable_contribution":17500, "Taxable_contribution": 0, "Non_Taxable_balance":0, "Taxable_balance": 0, "Yearly_social_security_income":0, "Asset_volatility": 0.15, "Expected_rate_of_return": 0.07, "Inflation_rate":0.035}`)

	rc := NewRetCalc_from_json(JsonObj)
	rc.All_paths = rc.RunAllPaths()
	if rc.Age != 22 && rc.Retirement_age != 65 && rc.Terminal_age != 90 && rc.Effective_tax_rate != 0.3 && rc.N != 20000 {
		t.Errorf("json did not initialize the retcalc object correctly")
	}
	if rc.Years != rc.Terminal_age-rc.Age {
		t.Errorf("rc.Years did not intitialize correctly")
	}

	if len(rc.sims) != rc.N {
		t.Errorf("rc.sims does not have the correct length")
	}
	if len(rc.sims[0]) != rc.Years {
		t.Errorf("rc.sims[i] does not have the correct length")
	}

}

// Test the components of income calculation

func TestGrowthFactors(t *testing.T) {
	JsonObj := []byte(`{"Age":22, "Retirement_age":65, "Terminal_age":90, "Effective_tax_rate":0.3, "Returns_tax_rate":0.3, "N": 20000, "Non_Taxable_contribution":17500, "Taxable_contribution": 0, "Non_Taxable_balance":0, "Taxable_balance": 0, "Yearly_social_security_income":0, "Asset_volatility": 0.15, "Expected_rate_of_return": 0.07, "Inflation_rate":0.035}`)
	r := NewRetCalc_from_json(JsonObj)
	r.All_paths = r.RunAllPaths()
	r.sims[0] = make([]float64, r.Years, r.Years)
	for i := range r.sims[0] {
		r.sims[0][i] = 0.07
	}

	gfactors := make([]float64, len(r.sims[0]), len(r.sims[0]))
	for i := range r.sims[0] {
		gfactors[i] = r.sims[0].GrowthFactor(i)
	}

	known_gfactors := []float64{99.5627497577374, 93.0492988390069, 86.9619615317821, 81.2728612446562, 75.9559450884637, 70.9868645686577, 66.3428640828576, 62.0026767129511, 57.9464268345337, 54.1555390976951, 50.6126533623318, 47.301545198441, 44.207051587328, 41.3150014834841, 38.612150919144, 36.0861223543402, 33.7253479947105, 31.5190168174864, 29.4570250630714, 27.5299299654873, 25.7289065098012, 24.0457070185058, 22.4726233817811, 21.002451758674, 19.6284595875457, 18.3443547547157, 17.1442567801081, 16.0226698879515, 14.974457839207, 13.9948204104738, 13.0792714116578, 12.2236181417362, 11.4239421885385, 10.6765814846155, 9.9781135370238, 9.32533975422784, 8.71527079834378, 8.14511289564839, 7.61225504266205, 7.11425704921687, 6.64883836375408, 6.21386762967671, 5.8073529249315, 5.4274326401229, 5.07236695338589, 4.74052986297746, 4.43040174110043, 4.14056237486022, 3.86968446248618, 3.61652753503382, 3.37993227573254, 3.15881521096499, 2.95216374856541, 2.75903154071534, 2.57853415020125, 2.40984500018808, 2.25219158896082, 2.10485195229984, 1.96715135728957, 1.83845921242016, 1.71818617983192, 1.60578147647843, 1.500730351849, 1.4025517307, 1.31079601, 1.225043, 1.1449, 1.07, 1}

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

// Not yet implemented - test will fail
func TestInflationFactors(t *testing.T) {
	JsonObj := []byte(`{"Age":22, "Retirement_age":65, "Terminal_age":90, "Effective_tax_rate":0.3, "Returns_tax_rate":0.3, "N": 20000, "Non_Taxable_contribution":17500, "Taxable_contribution": 0, "Non_Taxable_balance":0, "Taxable_balance": 0, "Yearly_social_security_income":0, "Asset_volatility": 0.15, "Expected_rate_of_return": 0.07, "Inflation_rate":0.035}`)
	r := NewRetCalc_from_json(JsonObj)
	r.All_paths = r.RunAllPaths()
	r.sims[0] = make([]float64, r.Years, r.Years)
	for i := range r.sims[0] {
		r.sims[0][i] = 0.07
	}
	knownInflationFactors := make([]float64, len(r.sims[0]), len(r.sims[0]))
	infationFactors := r.InflationFactors()

	factorsGood := true
	for i := range inflationFactors {
		if math.Abs(inflationFactors[i]-knownInflationFactors[i]) > 0.1 {
			factorsGood = false
		}
	}

	if !factorsGood {
		t.Errorf("Inflation Factors did not compute correctly")
	}

}

func TestIncomeFactors(t *testing.T) {
}

func TestRunIncome(t *testing.T) {
	JsonObj := []byte(`{"Age":22, "Retirement_age":65, "Terminal_age":90, "Effective_tax_rate":0.3, "Returns_tax_rate":0.3, "N": 20000, "Non_Taxable_contribution":17500, "Taxable_contribution": 0, "Non_Taxable_balance":0, "Taxable_balance": 0, "Yearly_social_security_income":0, "Asset_volatility": 0.15, "Expected_rate_of_return": 0.07, "Inflation_rate":0.035}`)
	rc := NewRetCalc_from_json(JsonObj)
	rc.All_paths = rc.RunAllPaths()

	rc.sims[0] = make([]float64, rc.Years, rc.Years)
	for i := range rc.sims[0] {
		rc.sims[0][i] = 0.07
	}
	if int(rc.IncomeOnPath(0)) != 42978 {
		t.Errorf("IncomeOnPath() not calculating correct income")
	}
}
