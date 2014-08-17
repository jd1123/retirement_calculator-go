package mortgage

import (
	"fmt"
	"testing"
)

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

func TestMortgageCalc(t *testing.T) {
	m := NewMortgageCalc(360, 559000, 900, 0.2, 0.0395, 0.065, 0.015, 0.02, 0.035)
	m.PrintMortgageCalc()
}

func TestAmmortizationTable(t *testing.T) {
	m := NewMortgageCalc(360, 559000, 900, 0.2, 0.0395, 0.065, 0.015, 0.02, 0.035)
	m.PrintAmmortizationTable()
}

func TestYearlyInterest(t *testing.T) {
	m := NewMortgageCalc(360, 559000, 900, 0.2, 0.0395, 0.065, 0.015, 0.02, 0.035)
	yr := m.yearlyInterest()
	for i := 1; i <= len(yr); i++ {
		fmt.Println("yr:", i, " int:", yr[i])
	}
}
