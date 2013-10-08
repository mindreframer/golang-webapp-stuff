package controllers

import (
	"models"
	"github.com/gorilla/sessions"
)

type Form interface {
	IsValid() (bool, error)
}

type AppController struct {
	Store *sessions.CookieStore
}

type DataContext struct {
	CurrentUser *models.User
	LoggedIn bool
	Notice   string

	User     *models.User
	Codes    []models.Code
	Code 	 *models.Code
	Err      string
}

type ActionResponse struct {
	RedirectTo  string
	Render      string
	Context    *DataContext
	RenderText  string
	Code        int
	Notice	    string
	NoLayout      bool
	ContentType string
}

type ApiErrorResponse struct {
	Error bool
	Message string
	Type string
}
