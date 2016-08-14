package main

import (
	_ "expvar"
	"log"
	"net/http"
	"time"
)

//
const MAJOR = "0.1"

// DO NOT UPDATE. This is populated by the build. See the Makefile.
var VERSION = "0"

//
func main() {
	func() {
		srv := &http.Server{
			Addr:         ":3001",
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		}
		log.Println(srv.ListenAndServe())
	}()
}
