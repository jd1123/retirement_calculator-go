package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"retirement_calculator-go/analytics"
	"retirement_calculator-go/retcalc"
	"strconv"

	"github.com/gorilla/mux"
)

// consts to define the url paths
const AllDataPrefix = "/alldata/"
const IncomesPrefix = "/incomes/"
const PathPrefix = "/paths/"
const InputPrefix = "/input/"
const InPathPrefix = "/inpath/"

func RegisterHandlers() {
	r := mux.NewRouter()

	// Actual used functions
	r.HandleFunc(InputPrefix, error_handler(RecalcFromWebInput)).Methods("POST")
	r.HandleFunc(PathPrefix, error_handler(SinglePath)).Methods("GET")

	// Functions for testing
	//r.HandleFunc(IncomesPrefix, error_handler(IncomesJSON)).Methods("GET")
	r.HandleFunc(AllDataPrefix, error_handler(Retcalc_basic)).Methods("GET")
	r.HandleFunc(InPathPrefix, error_handler(PathInfo)).Methods("GET")

	//Set Up the Handlers
	http.Handle(InputPrefix, r)
	http.Handle(PathPrefix, r)

	http.Handle(AllDataPrefix, r)
	//http.Handle(IncomesPrefix, r)
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

// Not used yet
func PathInfo(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Creates a RetCalc object from user input and returns
// a histogram of the incomes from this input
func RecalcFromWebInput(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	myRetCalc := retcalc.NewRetCalcFromJSON(body)

	fmt.Println()
	fmt.Println("POST request recieved - RecalcFromWebInput()")
	fmt.Printf("Recived: %s\n", string(body))

	fmt.Println("SessionId: " + myRetCalc.SessionId)

	// Save a file with the simulations for future reference
	go func() {
		jcalc, err := json.Marshal(myRetCalc)
		filename := myRetCalc.SessionId
		pth := "tmp/" + filename
		err = ioutil.WriteFile(pth, jcalc, 0644)
		if err != nil {
			panic(err)
		}
	}()

	return json.NewEncoder(w).Encode(analytics.HistoCumulative(myRetCalc.RunIncomes(), 250))
}

// Creates and returns a basic RetCalc object
func Retcalc_basic(w http.ResponseWriter, r *http.Request) error {
	rc := retcalc.NewRetCalc()
	return json.NewEncoder(w).Encode(rc)
}

/*
// this is for testing - returns income for a default RetCalc
func IncomesJSON(w http.ResponseWriter, r *http.Request) error {
	rc := retcalc.NewRetCalc()
	return json.NewEncoder(w).Encode(retcalc.HistoFromSlice(rc.RunIncomes()))
}
*/

// This function looks for two HTTP headers:
// X-Session-Id to get the SessionId and
// X-Percentile-Req to check which path the user clicked on
// Returns a json encoded path to display to the user
func SinglePath(w http.ResponseWriter, r *http.Request) error {
	// Process HTTP headers
	sessId := r.Header["X-Session-Id"][0]
	percentile, _ := strconv.ParseFloat(r.Header["X-Percentile-Req"][0], 64)
	filename := "tmp/" + string(sessId)

	// Error check HTTP Headers
	if percentile > 1.0 || percentile < 0.0 {
		fmt.Println("ERROR: invalid percentile requested - setting to 0.5")
		percentile = 0.5
	}

	if _, err := os.Stat(filename); err != nil {
		fmt.Printf("File does not exist: panicking")
		panic(err)
	}

	// Try opening file
	savedSim, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	// Now do the retcalc processing
	rc := retcalc.NewRetCalcFromJSON(savedSim)
	fmt.Println("Percentile Requested:", percentile, "SessionID", sessId)

	// Debug stuff - leave it in here
	// j, _ := json.Marshal(rc.PercentilePath(percentile))
	// fmt.Println(string(j))

	return json.NewEncoder(w).Encode(rc.PercentilePath(percentile))
}
