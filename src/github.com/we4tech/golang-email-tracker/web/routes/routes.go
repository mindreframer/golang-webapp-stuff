package routes

import (
	// Core modules
	"net/http"
	"path/filepath"

	// Controllers
	"controllers"
)

func GetServerPath(prefix string) http.Dir {
	v, _ := filepath.Abs(prefix)
	return http.Dir(v)
}

func init() {
	// Map all controllers
	app := controllers.New()
	http.HandleFunc("/users/signup", app.Signup)
	http.HandleFunc("/users/create", app.CreateUser)
	http.HandleFunc("/users/session/create", app.CreateSession)
	http.HandleFunc("/users/logout", app.Authenticated(app.DestroySession))

	http.HandleFunc("/codes", app.Authenticated(app.CodesRoot))
	http.HandleFunc("/codes/create", app.Authenticated(app.CreateCode))
	http.HandleFunc("/codes/destroy", app.Authenticated(app.DestroyCode))
	http.HandleFunc("/codes/start-tracking", app.Authenticated(app.StartTracking))
	http.HandleFunc("/api/codes", app.Authenticated(app.ApiCode))
	http.HandleFunc("/codes/track", app.TrackCode)

	http.HandleFunc("/", app.Root)

	// Map static files
	http.Handle("/stylesheets/", http.StripPrefix("/stylesheets/", http.FileServer(http.Dir("./public/stylesheets/"))))
	http.Handle("/javascripts/", http.StripPrefix("/javascripts/", http.FileServer(http.Dir("./public/javascripts/"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./public/images/"))))
	http.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir("./public/fonts/"))))
}
