package controllers

import (
	"fmt"
	"net/http"
	"appengine"
	"models"
	"strconv"

	"github.com/gorilla/sessions"
)

func (c *AppController) CodesRoot(w http.ResponseWriter, r *http.Request, session *sessions.Session) *ActionResponse {
	ctx := appengine.NewContext(r)
	user := c.GetCurrentUser(session)
	codes, err := models.FindCodesByUserId(ctx, user.Id)
	code := &models.Code{}

	if err != nil {
		return &ActionResponse{
			RenderText: fmt.Sprintf("Code loading error: %s", err.Error()),
		}
	} else {
		return &ActionResponse{
			NoLayout: true,
			Render: "views/codes/angular_index.html",
			Context: &DataContext{
				LoggedIn: true,
				CurrentUser: c.GetCurrentUser(session),
				Notice: c.GetNotice(session),
				Code: code,
				Codes: codes,
			},
		}
	}
}

func (c *AppController) CreateCode(w http.ResponseWriter, r *http.Request, session *sessions.Session) *ActionResponse {
	ctx := appengine.NewContext(r)
	user := c.GetCurrentUser(session)
	code := models.PopulateCode(r, true)

	code.UserId = user.Id
	saved, validationErr, err := code.Save(ctx)

	if saved {
		c.AddNotice(session, "Successfully created new mail tracking code.")
		return &ActionResponse{
			RedirectTo: "/codes",
		}
	} else {
		if validationErr != nil {
			c.AddNotice(session, fmt.Sprintf("%v", validationErr))
		} else if err != nil {
			c.AddNotice(session, fmt.Sprintf("%v", err))
		} else {
			c.AddNotice(session, "Failed to create new code")
		}

		return &ActionResponse{
			RedirectTo: "/codes",
		}
	}
}

func (c *AppController) DestroyCode(w http.ResponseWriter, r *http.Request, session *sessions.Session) *ActionResponse {
	ctx := appengine.NewContext(r)
	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)

	if id == 0 {
		c.AddNotice(session, "Request does not include any object id")
	} else {
		removed, err := models.DestroyCodeBy(ctx, id)
		if removed {
			c.AddNotice(session, "Successfully removed!")
		} else {
			fmt.Fprintf(w, "Failed to delete: %s", err.Error())
		}
	}

	return &ActionResponse{RedirectTo: "/codes"}
}

func (c *AppController) TrackCode(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)

	if id == 0 {
		ctx.Infof("Not available")
	} else {
		updated, err := models.UpdateCode(ctx, id, ConvertHeader2Array(r.Header))

		if updated == false {
			ctx.Infof("Error: %s", err.Error())
		} else {
			ctx.Infof("Saved!")
		}
	}

	http.ServeFile(w, r, "./public/images/code.png")
}

func (c *AppController) StartTracking(w http.ResponseWriter, r *http.Request, session *sessions.Session) *ActionResponse {
	ctx := appengine.NewContext(r)
	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)

	if id == 0 {
		c.AddNotice(session, "Not available")
	} else {
		updated, err := models.StartTrackingCode(ctx, id)

		if updated == false {
			c.AddNotice(session, fmt.Sprintf("Error: %s", err.Error()))
		} else {
			c.AddNotice(session, "Saved!")
		}
	}

	return &ActionResponse{
		RedirectTo: "/codes",
	}
}
