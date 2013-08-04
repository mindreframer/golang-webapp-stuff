package handler

import (
	"encoding/json"
	"github.com/jingweno/octokat"
	"github.com/jingweno/progmob/model"
	"github.com/jingweno/progmob/repository"
	"labix.org/v2/mgo/bson"
	"net/http"
)

type repoWithOwner struct {
	Owner string
	Repos []model.Repo
}

func Repos(w http.ResponseWriter, r *http.Request, currentUser *model.User) error {
	gh := repository.GitHub(currentUser.AccessToken)

	reposWithOwners := []repoWithOwner{}

  // user repos
	githubUserRepos, err := gh.Repos()
	if err != nil {
		return err
	}
	userRepos, err := upsertRepos(githubUserRepos)
	if err != nil {
		return err
	}
	userReposWithOwner := repoWithOwner{
		Owner: currentUser.Login,
		Repos: userRepos,
	}
	reposWithOwners = append(reposWithOwners, userReposWithOwner)

  // org repos
	orgs, err := gh.Orgs()
	if err != nil {
		return err
	}

	for _, org := range orgs {
		r, err := gh.OrgRepos(org.Login)
		if err != nil {
			return err
		}

		repos, err := upsertRepos(r)
		if err != nil {
			return err
		}

		repoWithOwner := repoWithOwner{
			Owner: org.Login,
			Repos: repos,
		}

		reposWithOwners = append(reposWithOwners, repoWithOwner)
	}

	repoObjectIds := []bson.ObjectId{}
	for _, ro := range reposWithOwners {
		for _, repo := range ro.Repos {
			repoObjectIds = append(repoObjectIds, repo.ObjectId)
		}
	}

	//currentUser.Repos = repoObjectIds
	//currentUser, err = repository.User().Upsert(currentUser)
	//if err != nil {
		//return err
	//}

	return json.NewEncoder(w).Encode(reposWithOwners)
}

func upsertRepos(githubRepos []octokat.Repository) ([]model.Repo, error) {
	repos := []model.Repo{}
	for _, githubRepo := range githubRepos {
		repo := model.Repo{
			ID:        githubRepo.ID,
			Owner:     githubRepo.Owner.Login,
			Name:      githubRepo.Name,
			FullName:  githubRepo.FullName,
			Private:   githubRepo.Private,
			HTMLURL:   githubRepo.HTMLURL,
			CreatedAt: githubRepo.CreatedAt,
			UpdatedAt: githubRepo.UpdatedAt,
		}
		repos = append(repos, repo)
	}

	repos, err := repository.Repo().UpsertAll(repos...)
	if err != nil {
		return nil, err
	}

	return repos, nil
}
