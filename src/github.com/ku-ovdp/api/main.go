// Package api implements the REST API for api.openvoicedata.org
package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// Constants
const API_VERSION = 1

// Flags
var persistenceBackend = flag.String("persistenceBackend", "mgo", "Persistence backend (dummy, mgo)")

func main() {
	flag.Parse()

	constructApplication()

	listen()
}

func listen() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	log.Println("attempting to listen on port", port)
	log.Fatalln("ListenAndServe:", http.ListenAndServe(":"+port, nil))
}
