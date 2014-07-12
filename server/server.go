/*

   server.go

   This file is part of retirement_calculator.

   Retirement_calculator is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   Retirement_calculator is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with retirement_calculator.  If not, see <http://www.gnu.org/licenses/>.

   Copyright 2014 Johnnydiabetic
*/

package server

import (
	"encoding/json"
	"log"
	"net/http"
	"retirement_calculator-go/retcalc"

	"github.com/gorilla/mux"
)

const PathPrefix = "/recalc/"

func RegisterHandlers() {
	r := mux.NewRouter()
	r.HandleFunc(PathPrefix, error_handler(Retcalc_basic)).Methods("GET")
	r.HandleFunc(PathPrefix, error_handler(HelloWorld)).Methods("POST")
	http.Handle(PathPrefix, r)
}

// badRequest is handled by setting the status code in the reply to StatusBadRequest.
type badRequest struct{ error }

// notFound is handled by setting the status code in the reply to StatusNotFound.
type notFound struct{ error }

func error_handler(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err == nil {
			return
		}
		switch err.(type) {
		case badRequest:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case notFound:
			http.Error(w, "task not found", http.StatusNotFound)
		default:
			log.Println(err)
			http.Error(w, "oops", http.StatusInternalServerError)
		}
	}
}

func Retcalc_basic(w http.ResponseWriter, r *http.Request) error {
	rc := retcalc.NewRetCalc()
	return json.NewEncoder(w).Encode(rc)
}

func HelloWorld(w http.ResponseWriter, r *http.Request) error {
	type c struct {
		This int
	}
	s := c{2}
	return json.NewEncoder(w).Encode(s)
}
