package controllers

import (
	"net/http"
	"appengine"
	"models"
	"github.com/gorilla/sessions"
)

func (c *AppController) CreateSession(w http.ResponseWriter, r *http.Request) {
	session := c.UserSession(r)
	ctx := appengine.NewContext(r)
	user, err := models.Authenticate(ctx, r)

	if err == nil {
		c.ClearUserSession(session)
		c.SetCurrentUser(session, user)
		c.AddNotice(session, "You have successfully Logged in.")
		session.Save(r, w)

		ar := &ActionResponse{RedirectTo: "/codes"}
		ar.Perform(w, r)
	} else {
		ar := &ActionResponse{RedirectTo: "/users/signup"}
		c.AddNotice(session, err.Error())
		session.Save(r, w)
		ar.Perform(w, r)
	}
}

func (c *AppController) DestroySession(w http.ResponseWriter, r *http.Request, session *sessions.Session) *ActionResponse {
	c.ClearUserSession(session)
	c.AddNotice(session, "You have successfully logged out!")

	return &ActionResponse{ RedirectTo: "/" }
}
