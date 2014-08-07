/*

main.go

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
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"retirement_calculator-go/server"
	"runtime/pprof"
	"strconv"
)

const listenPort = 8081

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Clean up all tmp files on exit
	go func() {
		for sig := range c {
			log.Printf("captured %v cleaning up and exiting..", sig)

			files, err := ioutil.ReadDir("tmp")
			if err != nil {
				panic(err)
			}

			for _, f := range files {
				fn := "tmp/" + f.Name()
				os.Remove(fn)
			}

			pprof.StopCPUProfile()
			os.Exit(1)
		}
	}()

	listenStr := ":" + strconv.Itoa(listenPort)
	server.RegisterHandlers()
	http.Handle("/", http.FileServer(http.Dir("./static")))
	fmt.Printf("RetCalc server on: Listening on port %d\n", listenPort)
	http.ListenAndServe(listenStr, nil)
}
