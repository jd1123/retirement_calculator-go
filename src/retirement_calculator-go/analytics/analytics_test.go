package analytics

import (
	"fmt"
	"math"
	"testing"
)

func TestHistoCumulative(t *testing.T) {
}

func TestIncomeTaxLiability(t *testing.T) {
	if math.Abs(IncomeTaxLiability(600000)-194068.0) > 1.0 {
		t.Errorf("TestIncomeTaxLiability(600000) does not match known value")
	}
	if math.Abs(IncomeTaxLiability(9075)-907.5) > 1.0 {
		fmt.Println(IncomeTaxLiability(9075))
		t.Errorf("IncomeTaxLiability(9075) should equal known value")
	}
	if math.Abs(IncomeTaxLiability(5000)-500.0) > 1 {
		fmt.Println(IncomeTaxLiability(5000))
		t.Errorf("IncomeTaxLiability(5000) should be equal to known value")
	}
}

func TestIncomeTaxBenefit(t *testing.T) {

}
