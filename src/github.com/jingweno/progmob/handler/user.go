package handler

import (
	"fmt"
	"github.com/jingweno/progmob/cookie"
	"github.com/jingweno/progmob/model"
	"github.com/jingweno/progmob/repository"
	"github.com/nu7hatch/gouuid"
	"net/http"
	"net/url"
	"os"
)

const (
	GitHubAuthorizeURL = "https://github.com/login/oauth/authorize"
)

var gitHubClientId, gitHubClientSecret string

func init() {
	gitHubClientId = os.Getenv("GIT_HUB_CLIENT_ID")
	gitHubClientSecret = os.Getenv("GIT_HUB_CLIENT_SECRET")
}

func Login(w http.ResponseWriter, r *http.Request) error {
	redirectURI := loginRedirectURI(r)
	url := gitHubRequestAccessURL(gitHubClientId, redirectURI)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)

	return nil
}

func LoginCallback(w http.ResponseWriter, r *http.Request) error {
	q := r.URL.Query()
	// TODO: verify state is the same
	state := q.Get("state")
	if state == "" {
		return fmt.Errorf("State can't be blank")
	}
	code := q.Get("code")
	if code == "" {
		return fmt.Errorf("Code can't be blank")
	}

	redirectURI := loginRedirectURI(r)
	accessToken, err := repository.GitHub("").CreateAccessToken(gitHubClientId, gitHubClientSecret, code, redirectURI)
	if err != nil {
		return err
	}

	githubUser, err := repository.GitHub(accessToken.Token).AuthenticatedUser()
	if err != nil {
		return err
	}

	userRepo := repository.User()
	user := &model.User{
		ID:          githubUser.ID,
		Login:       githubUser.Login,
		Email:       githubUser.Email,
		HTMLURL:     githubUser.HTMLURL,
		AvatarURL:   githubUser.AvatarURL,
		AccessToken: accessToken.Token,
		CreatedAt:   githubUser.CreatedAt,
	}
	user, err = userRepo.Upsert(user)
	if err != nil {
		return err
	}

	c := cookie.Cookie{Name: "user", Path: "/", Values: cookie.CookieValues{"id": user.ObjectIdString()}}
	cookie.Write(w, &c)

	redirectTo(w, r, "/")

	return nil
}

func Logout(w http.ResponseWriter, r *http.Request) error {
	cookie.Delete(w, "user")
	redirectTo(w, r, "/")

	return nil
}

func gitHubRequestAccessURL(clientId, redirectURI string) string {
	uuid, _ := uuid.NewV4()
	state := uuid.String()
	scope := "repo,user"
	url, _ := url.Parse(GitHubAuthorizeURL)

	q := url.Query()
	q.Add("client_id", clientId)
	q.Add("state", state)
	q.Add("redirect_uri", redirectURI)
	q.Add("scope", scope)
	url.RawQuery = q.Encode()

	return url.String()
}

func loginRedirectURI(r *http.Request) string {
	host, _ := url.Parse(r.Referer())
	return fmt.Sprintf("%s://%s/login/callback", host.Scheme, r.Host)
}
