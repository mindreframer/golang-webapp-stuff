package controllers

import (
	"net/http"
	"appengine"
	"models"
)

func (c *AppController) Signup(w http.ResponseWriter, r *http.Request) {

	session := c.UserSession(r)

	if c.LoggedIn(session) {
		ar := c.redirectWhenLoggedIn(session)
		session.Save(r, w)
		ar.Perform(w, r)
		return
	}

	user := models.PopulateUser(r, false)

	ar := &ActionResponse{
		Render: "views/users/signup.html",
		Context: &DataContext {Err: r.FormValue("err"), User: user, Notice: c.GetNotice(session)},
	}
	session.Save(r, w)
	ar.Perform(w, r)
}

func (c *AppController) CreateUser(w http.ResponseWriter, r *http.Request) {
	session := c.UserSession(r)

	if c.LoggedIn(session) {
		ar := c.redirectWhenLoggedIn(session)
		session.Save(r, w)
		ar.Perform(w, r)
		return
	}

	if r.Method != "POST" {
		ar := serve404(w)
		ar.Perform(w, r)
		return
	}

	ctx := appengine.NewContext(r)
	user := models.PopulateUser(r, true)
	saved, validationErr, err := user.Save(ctx)

	if saved {
		c.SetCurrentUser(session, user)
		c.AddNotice(session, "You have successfully created your account.")
		session.Save(r, w)

		ar := &ActionResponse{RedirectTo: "/codes"}
		ar.Perform(w, r)
	} else {
		if validationErr != nil {
			ar := serveValidationError(user, validationErr)
			ar.Perform(w, r)
		} else {
			ar := serve500(err)
			ar.Perform(w, r)
		}
	}
}
