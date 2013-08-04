package handler

import (
	"encoding/json"
	"net/http"
)

func redirectTo(w http.ResponseWriter, r *http.Request, path string) {
	http.Redirect(w, r, path, http.StatusTemporaryRedirect)
}

func serveError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func serveUnauthorized(w http.ResponseWriter) {
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}

func serveJSON(w http.ResponseWriter, object interface{}) error {
	return json.NewEncoder(w).Encode(object)
}
