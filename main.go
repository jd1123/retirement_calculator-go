package main

import (
	"net/http"
	"retirement_calculator-go/server"
)

func main() {
	server.RegisterHandlers()
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.ListenAndServe(":8080", nil)
}
