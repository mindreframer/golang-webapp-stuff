package controllers

import (
	"fmt"
	"net/http"
	"html/template"
	"models"
	"encoding/gob"
	"github.com/gorilla/sessions"
)

func (ar *ActionResponse) Perform(w http.ResponseWriter, r *http.Request) {
	if ar.RedirectTo != "" {
		http.Redirect(w, r, ar.RedirectTo, http.StatusFound)

	} else if ar.RenderText != "" {

		if ar.Code > 0 {
			w.WriteHeader(ar.Code)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		if ar.ContentType != "" {
			w.Header().Add("Content-Type", ar.ContentType)
		} else {
			w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		}

		fmt.Fprintf(w, ar.RenderText)

	} else if ar.Render != "" {
		var (
			templ *template.Template
			err error
		)

		if ar.NoLayout {
			t, e := template.ParseFiles(ar.Render)
			templ = t
			err = e

		} else {
			t, e := template.ParseFiles("views/layouts/base.html", ar.Render)
			templ = t
			err = e
		}

		if err != nil {
			fmt.Fprintln(w, err.Error())
		} else {
			templ.Execute(w, ar.Context)
		}

	}
}

func (c *AppController) Init() {
	c.Store = sessions.NewCookieStore([]byte("something-very-secret"))
	gob.Register(&models.User{})
}

func (c *AppController) Authenticated(handler func (w http.ResponseWriter, r *http.Request, s *sessions.Session) *ActionResponse) func (w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := c.UserSession(r)

		if c.LoggedIn(session) {
			ar := handler(w, r, session)
			session.Save(r, w)
			ar.Perform(w, r)
		} else {
			c.AddNotice(session, "This page requires authentication, please create an account or sign in with your existing account.")
			session.Save(r, w)

			http.Redirect(w, r, "/users/signup", http.StatusFound)
		}
	}
}


