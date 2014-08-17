package incomeanalytics

import "testing"

func TestIncomeTaxLiability(t *testing.T) {
	if IncomeTaxLiability(600000)-194068 > 500 {
		t.Errorf("TestIncomeTaxLiability(600000) does not match known value")
	}
}
