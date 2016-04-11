package main

import "retirement_calculator-go/mortgage"

func notmain() {
	sp := 580000.0
	yearlyInsurance := 150.0 * 12.0
	propertyTaxRate := 870.0 * 12.0 / sp
	m := mortgage.NewMortgageCalc(360, sp, yearlyInsurance, float64(200000.0/sp), 0.03875, 0.065, 0.01, propertyTaxRate, 0.025)
	m.Income = 105000
	m.PrintMortgageCalc()
	//m.PrintAmmortizationTable()
}
