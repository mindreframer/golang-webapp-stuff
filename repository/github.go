package repository

import (
	"fmt"
	"github.com/jingweno/octokat"
)

type gitHub struct {
	client *octokat.Client
}

func (gh *gitHub) CreateAccessToken(clientId, clientSecret, code, redirectURI string) (*octokat.AccessToken, error) {
	params := octokat.Params{"client_id": clientId, "client_secret": clientSecret, "code": code, "redirect_uri": redirectURI}
	accessToken, err := octokat.CreateAccessToken(&params)
	if err != nil {
		return nil, fmt.Errorf("Error occurred when creating access token: %s", err)
	}

	return accessToken, nil
}

func (gh *gitHub) AuthenticatedUser() (*octokat.User, error) {
	user, err := gh.client.User("")
	if err != nil {
		return nil, fmt.Errorf("Error occurred when gettting authenticated user: %s", err)
	}

	return user, nil
}

func (gh *gitHub) Repo(owner, name string) (*octokat.Repository, error) {
	repo := octokat.Repo{UserName: owner, Name: name}
	return gh.client.Repository(repo)
}

func (gh *gitHub) Repos() ([]octokat.Repository, error) {
	repos, err := gh.client.Repositories("", nil)
	if err != nil {
		return nil, fmt.Errorf("Error occurred when gettting repositories: %s", err)
	}

	return repos, nil
}

func (gh *gitHub) OrgRepos(org string) ([]octokat.Repository, error) {
	repos, err := gh.client.OrganizationRepositories(org, nil)
	if err != nil {
		return nil, fmt.Errorf("Error occurred when gettting organization repositories for %s: %s", org, err)
	}

	return repos, nil
}

func (gh *gitHub) Orgs() ([]octokat.Organization, error) {
	orgs, err := gh.client.Organizations("", nil)
	if err != nil {
		return nil, fmt.Errorf("Error occurred when gettting organizations: %s", err)
	}

	return orgs, nil
}
