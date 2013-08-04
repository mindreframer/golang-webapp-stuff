package handler

import (
	"github.com/jingweno/progmob/cookie"
	"github.com/jingweno/progmob/repository"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) error {
	c, err := cookie.Read(r, "user")
	var userLogin, userAvatarURL string

	if err == nil {
		objectId := c.Values["id"]
		if objectId != "" {
			user := repository.User().FirstByObjectId(objectId)
			if user != nil {
				userLogin = user.Login
				userAvatarURL = user.AvatarURL
			}
		}
	}

	return Templates.ExecuteTemplate(w, "index.html", struct {
		UserLogin     string
		UserAvatarURL string
	}{
		UserLogin:     userLogin,
		UserAvatarURL: userAvatarURL,
	})
}
