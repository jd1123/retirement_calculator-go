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
	"runtime/pprof"
	"strconv"

	"github.com/jd1123/retirement_calculator-go/server"
)

const listenPort = 8081

func serverMsg(listenPort int) {
	fmt.Println("RetCalc server listening on port " + strconv.Itoa(listenPort))
}

func usage() {
	fmt.Println("Usage: main -p <port number>")
}

func processCmdLnArgs(args []string) int {
	lp := listenPort
	fail := false
	if len(args) > 1 {
		i := 1
		for i < len(args) {
			switch args[i] {

			// Port setting
			case "-p":
				{
					p, err := strconv.Atoi(args[i+1])
					if err != nil {
						usage()
					} else {
						lp = p
					}
					i++
				}
			// Unrecognized argument - fail
			default:
				{
					fail = true
				}
			}
			i++
		}
	}
	if fail {
		usage()
	}
	return lp
}

func main() {
	// Get command line args
	lp := processCmdLnArgs(os.Args)

	// Clean up all tmp files on exit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
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

	// Start listening
	listenStr := ":" + strconv.Itoa(lp)
	server.RegisterHandlers()
	http.Handle("/", http.FileServer(http.Dir("./static")))
	serverMsg(lp)
	http.ListenAndServe(listenStr, nil)
}
