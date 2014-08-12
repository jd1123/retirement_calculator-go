package mortgage

import "testing"

func TestMortgageMonthlyPayment(t *testing.T) {
	loan := 408070.0
	r := 0.0395
	term := 360
	payment := MortgageMonthlyPayment(loan, r, term)
	if payment-1936 > 1 {
		t.Errorf("mortgage.MonthlyMortgagePayment should equal known value")
	}
}

func TestMonthlyMoneyGrowthTerminal(t *testing.T) {
	amount := 100000.0
	r := 0.06
	term := 120
	g := MonthlyMoneyGrowthTerminal(amount, r, term)
	if g-181939.0 > 1 {
		t.Errorf("mortgage.MonthlyMoneyGrowthTerminal() should equal known value")
	}
}
