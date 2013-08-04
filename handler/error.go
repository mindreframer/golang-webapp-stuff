package handler

import (
	"log"
	"net/http"
)

func HandleError(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			serveError(w, err)
			log.Printf("handling %q: %v", r.RequestURI, err)
		}
	}
}
