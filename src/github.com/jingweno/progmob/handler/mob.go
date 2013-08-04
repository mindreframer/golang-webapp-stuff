package handler

import (
	"fmt"
	"github.com/jingweno/octokat"
	"github.com/jingweno/progmob/model"
	"github.com/jingweno/progmob/repository"
	"net/http"
)

func Mobs(w http.ResponseWriter, r *http.Request, currentUser *model.User) error {
	mobs, err := repository.Mob().ForUser(currentUser)
	if err != nil {
		return err
	}

	return serveJSON(w, mobs)
}

func Mob(w http.ResponseWriter, r *http.Request, currentUser *model.User) error {
	q := r.URL.Query()
	owner := q.Get(":owner")
	if owner == "" {
		return fmt.Errorf("Missing parameter owner")
	}
	name := q.Get(":name")
	if name == "" {
		return fmt.Errorf("Missing parameter name")
	}

	gitHubRepo, err := fetchGitHubRepo(currentUser, owner, name)
	if err != nil {
		serveUnauthorized(w)
		return nil
	}

	mob, err := upsertMob(currentUser, gitHubRepo)
	if err != nil {
		return err
	}

	err = addMobToUser(currentUser, mob)

	return serveJSON(w, mob)
}

func fetchGitHubRepo(user *model.User, owner, name string) (*octokat.Repository, error) {
	return repository.GitHub(user.AccessToken).Repo(owner, name)
}

func upsertMob(user *model.User, repo *octokat.Repository) (*model.Mob, error) {
	mob := model.Mob{
		Owner:        repo.Owner.Login,
		RepoName:     repo.Name,
		RepoFullName: repo.FullName,
	}

	return repository.Mob().Upsert(&mob)
}

func addMobToUser(user *model.User, mob *model.Mob) error {
	return repository.User().AddMob(user, mob)
}
