/*

client.go

My retirement calculator in Go

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.

Copyright 2014 Johnnydiabetic
*/

package main

import (
	"encoding/json"
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
	GetResponse()
}
