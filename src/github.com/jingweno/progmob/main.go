package main

import (
	"github.com/bmizerany/pat"
	"github.com/jingweno/progmob/database"
	"github.com/jingweno/progmob/handler"
	"github.com/jingweno/progmob/repository"
	//"github.com/kr/secureheader"
	"html/template"
	"log"
	"net/http"
	"os"
)

var (
	db        *database.Database
	templates *template.Template
)

func main() {
	var err error
	if templates, err = template.New("").Funcs(funcs).
		ParseFiles(
		"templates/index.html",
	); err != nil {
		log.Fatal(err)
	}

	dbURL := os.Getenv("DATABASE_URL")
	db, err = database.New(dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = repository.Setup(db)
	if err != nil {
		log.Fatal(err)
	}

	handler.Templates = templates

	m := pat.New()
	m.Get("/", handler.HandleError(handler.Home))
	m.Get("/login", handler.HandleError(handler.Login))
	m.Get("/login/callback", handler.HandleError(handler.LoginCallback))
	m.Get("/logout", handler.HandleError(handler.Logout))
	m.Get("/user/repos.json", handler.Authenticate(handler.Repos))
	m.Get("/user/mobs.json", handler.Authenticate(handler.Mobs))
	m.Get("/mobs/:owner/:name", handler.Authenticate(handler.Mob))
	m.Get("/assets/bootstrap/css/:css", http.StripPrefix("/assets", http.FileServer(http.Dir("assets"))))
	m.Get("/assets/fonts/:fonts", http.StripPrefix("/assets", http.FileServer(http.Dir("assets"))))
	m.Get("/assets/css/:css", http.StripPrefix("/assets", http.FileServer(http.Dir("assets"))))
	m.Get("/assets/js/:js", http.StripPrefix("/assets", http.FileServer(http.Dir("assets"))))
	http.Handle("/", m)

	port := os.Getenv("PORT")
	//secureheader.DefaultConfig.PermitClearLoopback = true
	log.Printf("Server starting at port %s", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf(`{"func":"ListenAndServe", "error":%q}`, err)
	}
}
