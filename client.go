package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"retirement_calculator-go/retcalc"
)

func perror(err error) {
	if err != nil {
		panic(err)
	}
}

func GetResponse() {
	url := "http://127.0.0.1:8080/recalc/"
	res, err := http.Get(url)
	perror(err)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	perror(err)
	var r retcalc.RetCalc

	err = json.Unmarshal(body, &r)

}

func main() {
	//	GetResponse()
	rc := retcalc.NewRetCalc_b()
	rc.Yearly_retirement_expenses = 60000.0
	rc.Non_Taxable_contribution = 42000.0
	rc.Non_Taxable_balance = 100000.0
	rc.Taxable_contribution = 18000.0
	rc.Taxable_balance = 100000.0
	rc.All_paths = rc.RunAllPaths()
	//fmt.Println(rc.RunIncomes())
	fmt.Println("\n")
	fmt.Println(rc.IncomeProbability())
	fmt.Println(rc.PercentileIncome(0.2))
}
