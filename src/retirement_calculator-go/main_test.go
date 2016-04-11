package main

import (
	"fmt"
	"retirement_calculator-go/retcalc"
	"testing"
)

func TestConstructors(t *testing.T) {
	JsonObj := []byte(`{"Age":22, "Retirement_age":65, "Terminal_age":90, "Effective_tax_rate":0.3, "Returns_tax_rate":0.3, "N": 20000, 
					"Non_Taxable_contribution":17500, "Taxable_contribution": 0, "Non_Taxable_balance":0, "Taxable_balance": 0, 
					"Yearly_social_security_income":0, "Asset_volatility": 0.15, "Expected_rate_of_return": 0.07, "Inflation_rate":0.035}`)
	rc, err := retcalc.NewRetCalcFromJSON(JsonObj)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(rc.IncomeOnPath(0))
}
