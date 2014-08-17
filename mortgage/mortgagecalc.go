package mortgage

import "fmt"

import "retirement_calculator-go/analytics"

type MortgageCalc struct {
	// Inputs
	SalePrice             float64
	DownPaymentPercentage float64
	LoanAmount            float64
	LoanRate              float64
	IRR                   float64
	UpkeepPctPerYear      float64
	PropertyTax           float64
	InsurancePerYear      float64
	TermInMonths          int
	ExpectedHousingReturn float64
	Commission            float64
	Income                float64
	// Computed fields
	monthlyPayment         float64
	monthlyTaxes           float64
	monthlyInsurance       float64
	monthlyUpkeep          float64
	TerminalHousePrice     float64
	FixedMonthlyPayment    float64
	FloatingMonthlyPayment float64
}

type Payment struct {
	fixedPayment    float64
	floatingPayment float64
}

type RentCalc struct {
	MonthlyRent   float64
	RentInflation float64
}

func NewMortgageCalc(termInMonths int, salePrice, insurancePerYear, downPaymentPercentage, loanRate, irr, maintencePerYear, propertyTax,
	expectedHousingReturn float64) MortgageCalc {

	m := MortgageCalc{}
	m.TermInMonths = termInMonths
	m.SalePrice = salePrice
	m.InsurancePerYear = insurancePerYear
	m.DownPaymentPercentage = downPaymentPercentage
	m.LoanRate = loanRate
	m.UpkeepPctPerYear = maintencePerYear
	m.PropertyTax = propertyTax
	m.ExpectedHousingReturn = expectedHousingReturn
	m.IRR = irr
	m.Income = 100000
	m.computeParameters()

	return m
}

func (m *MortgageCalc) computeParameters() {
	m.LoanAmount = (1 - m.DownPaymentPercentage) * m.SalePrice
	m.monthlyPayment = MortgageMonthlyPayment(m.LoanAmount, m.LoanRate, m.TermInMonths)
	m.monthlyTaxes = m.PropertyTax * m.SalePrice / 12.0
	m.monthlyInsurance = m.InsurancePerYear / 12.0
	m.monthlyUpkeep = m.SalePrice * m.UpkeepPctPerYear / 12.0
	m.FixedMonthlyPayment = m.monthlyPayment
	m.FloatingMonthlyPayment = m.monthlyTaxes + m.monthlyUpkeep + m.monthlyInsurance
	m.TerminalHousePrice = m.SalePrice * GF(m.ExpectedHousingReturn, m.TermInMonths/12)
}

func (m MortgageCalc) PrintMortgageCalc() {
	fmt.Println("Mortgage Calc")
	fmt.Println("============")
	fmt.Println("SalePrice:\t\t", m.SalePrice)
	fmt.Println("Down Payment Pct:\t", m.DownPaymentPercentage)
	fmt.Println("DownPayment:\t\t", (1-m.DownPaymentPercentage)*m.SalePrice)
	fmt.Println("Loan Rate:\t\t", m.LoanRate)
	fmt.Println("Mortgage Term:\t\t", m.TermInMonths)
	fmt.Println("Mortgage Payment:\t", m.monthlyPayment)
}

func (m MortgageCalc) DFPaymentsVector(irr float64) []float64 {
	payments := make([]float64, m.TermInMonths+1, m.TermInMonths+1)
	for i := 0; i < len(payments); i++ {
		if i == 0 {
			payments[i] = m.DownPaymentPercentage * m.SalePrice
		} else {
			fixed := m.FixedMonthlyPayment * DF_m(irr, i)
			floating := m.FloatingMonthlyPayment * DF(irr, i) * GF_m(m.ExpectedHousingReturn, i)
			payments[i] = fixed + floating
		}
	}
	return payments
}

// Discounted House Price
func (m MortgageCalc) DFTerminalHousePrice(irr float64) float64 {
	return m.TerminalHousePrice * (1 - m.Commission) * DF_m(irr, m.TermInMonths)
}

// DF of all payments including down payment
func (m MortgageCalc) TotalOwnershipCost(irr float64) float64 {
	return Sum(m.DFPaymentsVector(irr)) + m.DFTerminalHousePrice(irr)
}

func (m MortgageCalc) TotalInterest() float64 {
	totalPayments := MortgageMonthlyPayment(m.LoanAmount, m.LoanRate, m.TermInMonths) * float64(m.TermInMonths)
	return totalPayments - m.LoanAmount
}

func (m MortgageCalc) AmmortizationTable() ([]int, []float64, []float64) {
	mp := MortgageMonthlyPayment(m.LoanAmount, m.LoanRate, m.TermInMonths)
	month := make([]int, m.TermInMonths, m.TermInMonths)
	principal := make([]float64, m.TermInMonths, m.TermInMonths)
	interest := make([]float64, m.TermInMonths, m.TermInMonths)
	remainingBalance := m.LoanAmount
	for i := 0; i < m.TermInMonths; i++ {
		month[i] = i + 1
		interest[i] = remainingBalance * m.LoanRate / 12.0
		principal[i] = mp - interest[i]
		remainingBalance -= principal[i]
	}
	return month, principal, interest
}

func (m MortgageCalc) PrintAmmortizationTable() {
	month, principal, interest := m.AmmortizationTable()
	for i := 0; i < m.TermInMonths; i++ {
		fmt.Println("Month: ", month[i], " Principal: ", principal[i], " Interest: ", interest[i])
	}
}

func (m MortgageCalc) yearlyInterest() map[int]float64 {
	yrlyInterest := make(map[int]float64)
	_, _, interest := m.AmmortizationTable()
	for i := 0; i < m.TermInMonths; i += 12 {
		totalInt := 0.0
		for j := 0; j < 12; j++ {
			//fmt.Println("mnth:", i, " interest:", interest[i+j])
			totalInt += interest[i+j]
		}
		yrlyInterest[(i / 12)] = totalInt
	}
	return yrlyInterest
}

func (m MortgageCalc) nominalMonthlyIncomeTaxBenefit() []float64 {
	yrlyInterest := m.yearlyInterest()
	termInYears := m.TermInMonths / 12
	yearlyTaxes := make([]float64, termInYears, termInYears)
	yearlyIncomeTaxBenefit := make([]float64, termInYears, termInYears)
	monthlyIncomeTaxBenefit := make([]float64, m.TermInMonths, m.TermInMonths)
	for i := 0; i < termInYears; i++ {
		yearlyTaxes[i] = m.SalePrice * GF(m.ExpectedHousingReturn, i) * m.PropertyTax
		yearlyIncomeTaxBenefit[i] = analytics.IncomeTaxLiability(m.Income) - analytics.IncomeTaxLiability(m.Income-yearlyTaxes[i]-yrlyInterest[i])
		for j := 0; j < 12; j++ {
			monthlyIncomeTaxBenefit[i*12+j] = yearlyIncomeTaxBenefit[i] / 12.0
		}
	}
	return monthlyIncomeTaxBenefit
}

func (m MortgageCalc) TotalMonthlyPayments() (float64, float64, float64) {
	mp := MortgageMonthlyPayment(m.LoanAmount, m.LoanRate, m.TermInMonths)
	taxes := m.SalePrice * m.PropertyTax / 12.0
	insurance := m.InsurancePerYear / 12.0
	return mp, taxes, insurance
}
