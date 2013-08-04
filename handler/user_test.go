package handler

import (
	"github.com/bmizerany/assert"
	"net/url"
	"testing"
)

func TestGitHubRequestAccessURL(t *testing.T) {
	clientId := "client_id"
	redirectURI := "http://example.com"
	url, err := url.Parse(gitHubRequestAccessURL(clientId, redirectURI))
	query := url.Query()

	assert.Equal(t, nil, err)
	assert.Equal(t, clientId, query.Get("client_id"))
	assert.Equal(t, redirectURI, query.Get("redirect_uri"))
	assert.Equal(t, "repo,user", query.Get("scope"))
	assert.T(t, len(query.Get("scope")) > 0)
}
