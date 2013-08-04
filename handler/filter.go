package handler

import (
	"github.com/jingweno/progmob/cookie"
	"github.com/jingweno/progmob/model"
	"github.com/jingweno/progmob/repository"
	"log"
	"net/http"
)

func Authenticate(f func(http.ResponseWriter, *http.Request, *model.User) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			user *model.User
			err  error
		)

		c, err := cookie.Read(r, "user")
		if err != nil {
			serveUnauthorized(w)
			return
		}

		userId := c.Values["id"]
		if userId == "" {
			serveUnauthorized(w)
			return
		}

		user = repository.User().FirstByObjectId(userId)
		if user == nil {
			serveUnauthorized(w)
			return
		}

		err = f(w, r, user)

		if err != nil {
			serveError(w, err)
			log.Printf("handling %q: %v", r.RequestURI, err)
		}
	}
}
