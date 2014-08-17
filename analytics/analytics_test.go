package analytics

import (
	"retirement_calculator-go/retcalc"
	"testing"
)

func TestHistoCumulative(t *testing.T) {
	rc := retcalc.NewRetCalc()
	a := rc.RunIncomes()
}

func TestIncomeTaxLiability(t *testing.T) {
	if IncomeTaxLiability(600000)-194068 > 500 {
		t.Errorf("TestIncomeTaxLiability(600000) does not match known value")
	}
}
