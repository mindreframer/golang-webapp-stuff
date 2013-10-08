package controllers

import (
	"fmt"
	"net/http"
	"models"
	"net/url"

	"github.com/gorilla/sessions"
)

func ConvertHeader2Array(httpHeaders http.Header) []models.Header {
	headers := make([]models.Header, len(httpHeaders))

	i := 0
	for key, values := range httpHeaders {
		if key != "Cookie" {
			if len(values) > 0 {
				headers[i] = models.Header{Key: key, Value: values[0]}
				i++
			}
		}
	}

	return headers
}

func (c *AppController) UserSession(r *http.Request) *sessions.Session {
	session, _ := c.Store.Get(r, "user_session")
	return session
}

func (c *AppController) LoggedIn(s *sessions.Session) bool {
	return s.Values["user"] != nil
}

func (c *AppController) SetCurrentUser(session *sessions.Session, user *models.User) {
	session.Values["user"] = user
}

func (c *AppController) ClearUserSession(session *sessions.Session) {
	session.Values = make(map[interface {}]interface {})
}

func (c *AppController) GetCurrentUser(s *sessions.Session) *models.User {
	u := s.Values["user"]
	if u != nil {
		return s.Values["user"].(*models.User)
	} else {
		return nil
	}
}

func (c *AppController) AddNotice(s *sessions.Session, msg string) {
	s.Values["notice"] = msg
}

func (c *AppController) GetNotice(session *sessions.Session) string {
	msg := session.Values["notice"]

	newVal := make(map[interface {}]interface {})

	for k, v := range session.Values {
		if k != "notice" {
			newVal[k] = v
		}
	}
	session.Values = newVal

	if msg != nil {
		return msg.(string)
	} else {
		return ""
	}
}

func (c *AppController) SaveSession(w http.ResponseWriter, r *http.Request) {
	err := c.UserSession(r).Save(r, w)
	if err != nil {
		fmt.Fprintf(w, "Session error: %s", err.Error())
	}
}

func (c *AppController) redirectWhenLoggedIn(session *sessions.Session) *ActionResponse {
	c.AddNotice(session, "You have already logged in.")
	return &ActionResponse{
		RedirectTo: "/",
	}
}

func serveValidationError(user *models.User, validationErr error) *ActionResponse {
	return &ActionResponse{
		RedirectTo: fmt.Sprintf("/users/signup?err=%v&name=%s&email=%s",
			url.QueryEscape(validationErr.Error()), url.QueryEscape(user.Name), url.QueryEscape(user.Email)),
	}
}

func serve404(w http.ResponseWriter) *ActionResponse {
	return &ActionResponse{
		Code: http.StatusNotFound,
		RenderText: "Page not found",
	}
}

func serve500(e error) *ActionResponse {
	return &ActionResponse{ RenderText: e.Error() }
}
