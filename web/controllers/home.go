package controllers

import (
	"net/http"
)

func (c *AppController) Root(w http.ResponseWriter, r *http.Request) {
	session := c.UserSession(r)

	context := &DataContext{
		LoggedIn: c.LoggedIn(session),
		CurrentUser: c.GetCurrentUser(session),
		Notice: c.GetNotice(session),
	}

	ar := &ActionResponse{
		Render: "views/home/index.html",
		Context: context,
	}

	session.Save(r, w)

	ar.Perform(w, r)
}
