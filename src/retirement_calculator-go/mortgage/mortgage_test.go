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

func TestAmmortizationTable(t *testing.T) {
	m := NewMortgageCalc(60, 500000, 900, 0.2, 0.0395, 0.065, 0.015, 0.02, 0.035)
	//	m.PrintAmmortizationTable()
	_, principal, interest := m.AmmortizationTable()
	knownInterest := []float64{1316.66666666667, 1296.78197136419, 1276.831822273, 1256.81600394106, 1236.7343002071, 1216.58649419836,
		1196.37236832817, 1176.09170429366, 1155.74428307337, 1135.32988492489, 1114.84828938251, 1094.2992752548, 1073.68262062226,
		1052.99810283488, 1032.24549850979, 1011.42458352879, 990.535133035978, 969.576921435296,
		948.549722388095, 927.453308810697, 906.28745287194, 885.051925990718, 863.746498833512,
		842.370941311913, 820.925022580139, 799.40851103254, 777.821174301097, 756.162779252912,
		734.433091987694, 712.631877835228, 690.758901352844, 668.813926322871, 646.796715750092,
		624.707031859177, 602.544636092121, 580.309289105666, 558.000750768713, 535.618780159735,
		513.163135564168, 490.633574471808, 468.029853574186, 445.351728761942, 422.598955122191,
		399.771286935876, 376.868477675114, 353.890280000536, 330.836445758612, 307.706725978976,
		284.500870871731, 261.218629824758, 237.859751401006, 214.423983335775, 190.911072533997,
		167.320765067496, 143.652806172251, 119.906940245642, 96.0829108436919, 72.1804606782935,
		48.199331614434, 24.139264667406}
	knownPrincipal := []float64{6040.92009189298, 6060.80478719546, 6080.75493628665, 6100.77075461859, 6120.85245835255,
		6141.00026436129, 6161.21439023148, 6181.49505426599, 6201.84247548628, 6222.25687363476, 6242.73846917714,
		6263.28748330485, 6283.90413793739, 6304.58865572477, 6325.34126004986, 6346.16217503086, 6367.05162552367,
		6388.00983712435, 6409.03703617156, 6430.13344974895, 6451.29930568771, 6472.53483256893, 6493.84025972614,
		6515.21581724774, 6536.66173597951, 6558.17824752711, 6579.76558425855, 6601.42397930674, 6623.15366657196,
		6644.95488072442, 6666.82785720681, 6688.77283223678, 6710.79004280956, 6732.87972670047, 6755.04212246753,
		6777.27746945398, 6799.58600779094, 6821.96797839992, 6844.42362299548, 6866.95318408784, 6889.55690498546,
		6912.23502979771, 6934.98780343746, 6957.81547162377, 6980.71828088454, 7003.69647855911, 7026.75031280104,
		7049.88003258067, 7073.08588768792, 7096.36812873489, 7119.72700715864, 7143.16277522387, 7166.67568602565,
		7190.26599349215, 7213.9339523874, 7237.67981831401, 7261.50384771596, 7285.40629788136, 7309.38742694522,
		7333.44749389224}
	for i := 0; i < m.TermInMonths; i++ {
		if interest[i]-knownInterest[i] > 15 {
			t.Errorf("AmmortizationTable(): Interest should match known value")
		}
		if principal[i]-knownPrincipal[i] > 15 {
			t.Errorf("AmmortizationTable(): Principal should match known value")
		}
	}
	m = NewMortgageCalc(60, 0, 900, 0.2, 0.0395, 0.065, 0.015, 0.02, 0.035)
	_, principal, interest = m.AmmortizationTable()
	for i := 0; i < m.TermInMonths; i++ {
		if interest[i] != 0 || principal[i] != 0 {
			fmt.Println("Interest[", i, "]: ", interest[i], " Principal[", i, "]: ", principal[i])
			t.Errorf("0 case should be equal to zero in AmmortizationTable()")
		}
	}
}

func TestYearlyInterest(t *testing.T) {
	//	m := NewMortgageCalc(360, 559000, 900, 0.2, 0.0395, 0.065, 0.015, 0.02, 0.035)
	//	yr := m.yearlyInterest()
	//	for i := 1; i <= len(yr); i++ {
	//		fmt.Println("yr:", i, " int:", yr[i])
	//	}
}

func TestNominalIncomeTaxBenefit(t *testing.T) {
	m := NewMortgageCalc(360, 580000, 900, 0.25, 0.0395, 0.065, 0.015, 0.02, 0.035)
	//mp, tax, ins := m.TotalMonthlyPayments()
	m.Income = 105000.0
	interest := m.yearlyInterest()[0]
	fmt.Println("Yearly Interest:\t", interest)
	fmt.Println("Taxes:\t", m.SalePrice*m.PropertyTax)

	fmt.Println("Income Tax Benefit:\t", m.nominalMonthlyIncomeTaxBenefit()[0])
	m.Income = 205000
	fmt.Println("Income Tax Benefit:\t", m.nominalMonthlyIncomeTaxBenefit()[0])
	m.Income = 305000
	fmt.Println("Income Tax Benefit:\t", m.nominalMonthlyIncomeTaxBenefit()[0])
	m.Income = 705000
	fmt.Println("Income Tax Benefit:\t", m.nominalMonthlyIncomeTaxBenefit()[0])
	m.Income = 905000
	fmt.Println("Income Tax Benefit:\t", m.nominalMonthlyIncomeTaxBenefit()[0])
}

func TestTotalOwnershipCost(t *testing.T) {
	m := NewMortgageCalc(360, 575000, 900, 0.25, 0.0395, 0.065, 0.015, 0.02, 0.03)
	knownOwn := 450442.0
	totOwn := m.TotalOwnershipCost(m.IRR)
	if totOwn-knownOwn > 100.0 {
		t.Errorf("TotalOwnershipCost() should be equal to known value")
		fmt.Println("totOwn: ", totOwn, " knownOwn: ", knownOwn)
	}
}

func TestSumProduct(t *testing.T) {
	r := SumProduct([]float64{1.0, 2.0, 3.0}, []float64{1.0, 2.0, 3.0})
	if r != 14.0 {
		t.Errorf("SumProduct() does not match known value")
		fmt.Println("SumProduct(): ", r)
	}
}

func TestCopy(t *testing.T) {
	m := NewMortgageCalc(360, 575000, 900, 0.25, 0.0395, 0.065, 0.015, 0.02, 0.03)
	nm := m.copyCalc()
	if &nm == &m {
		t.Errorf("MortgageCalc.copyCalc() should return new instance")
	}
	if nm != m {
		t.Errorf("MortgageCalc.copyCalc() should return an exact replica")
	}
}

// Test is not done - check more stuff - mostly computed fields
func TestNewMortgagecalFromJson(t *testing.T) {
	JsonObj := []byte(`{"SalePrice":570000.0,"DownPaymentPercentage":0.35,"LoanRate":0.039,"IRR":0.07,"UpkeepPctPerYear":0.015,
	"PropertyTax":0.02,"InsurancePerYear":1800.0,"TermInMonths":360,"ExpectedHousingReturn":0.025,"Comission":0.06,"Income":100000.0}`)
	m := NewMortgageCalcFromJSON(JsonObj)
	if m.SalePrice != 570000.0 {
		t.Errorf("m.SalePrice should equal known value")
	}
	if m.monthlyPayment != MortgageMonthlyPayment((1-m.DownPaymentPercentage)*m.SalePrice, m.LoanRate, m.TermInMonths) {
		t.Errorf("m.monthlyPayment should equal output from MortgageMonthlyPayment()")
	}
}
