package path

import (
	"fmt"
	"time"
)

type YearlyEntry struct {
	Age                                              int
	Year                                             time.Time
	SOY_taxable_balance, EOY_taxable_balance         float64
	SOY_non_taxable_balance, EOY_non_taxable_balance float64
	Taxable_returns, Non_taxable_returns             float64
	Rate_of_return                                   float64
	Taxable_contribution, Non_taxable_contribution   float64
	Yearly_expenses                                  float64
	Deficit                                          float64
}

// Path type and Path methods
type Path []YearlyEntry

func (p Path) Print_path() {
	sum := 0.0
	for i := 0; i < len(p); i++ {
		//fmt.Println(p[i])
		sum += p[i].Rate_of_return
		//fmt.Printf("Age: %d EOY Balance: %f\n", p[i].Age, p[i].EOY_taxable_balance+p[i].EOY_non_taxable_balance)
		fmt.Printf("Age: %d EOY t_Balance: %f EOY nt_balance: %f Expenses: %f\n", p[i].Age,
			p[i].EOY_taxable_balance, p[i].EOY_non_taxable_balance, p[i].Yearly_expenses)
	}
	fmt.Printf("Average RoR: %f", sum/float64(len(p)))
}

// Implement the sort interface on a group of Paths
type PathGroup []Path

func (p PathGroup) Len() int {
	return len(p)
}

func (p PathGroup) Less(i, j int) bool {
	li, lj := len(p[i]), len(p[j])
	vi := p[i][li-1].EOY_taxable_balance + p[i][li-1].EOY_non_taxable_balance
	vj := p[j][lj-1].EOY_taxable_balance + p[j][lj-1].EOY_non_taxable_balance
	return vi < vj
}

func (p PathGroup) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
