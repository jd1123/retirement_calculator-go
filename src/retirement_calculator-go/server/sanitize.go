// Santize user input

package server

import "fmt"

// Sanitze JSON input for retcalc
func sanitizeJSON(JSONbody string) error {
	fmt.Println("JSON Body Sanitizer")
	fmt.Println(JSONbody)
	return nil
}
