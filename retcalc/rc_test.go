// These tests are shit - write good ones

package retcalc

import "testing"

func TestIncomeProbability(t *testing.T) {
	rc := NewRetCalc()
	i := rc.IncomeProbability()
	if i == 0 {
		t.Errorf("IncomeProbability() not working on NewRetCalc")
	}
}

func TestNewRetCalcJSon(t *testing.T) {
	var JsonObj []byte = "{'Age':22, 'Retirement_age':65, 'Terminal_age':90, 'Effective_tax_rate':0.3, 'Returns_tax_rate':0.3, 'Years':68, 'N': 20000, 'Non_Taxable_contribution':17500, 'Taxable_contribution': 0, 'Non_Taxable_balance':0, 'Taxable_balance': 0, 'Yearly_social_security_income':0, 'Asset_volatility': 0.15, 'Expected_rate_of_return': 0.07, \'Inflation_rate\':0.035}"
	rc := NewRetCalc_from_json(JsonObj)
	rc.All_paths = rc.RunAllPaths()
}
