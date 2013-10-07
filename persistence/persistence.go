// package persistence provides an interface to persistence backends
package persistence

import (
	. "github.com/ku-ovdp/api/repository"
)

type Backend interface {
	NewProjectRepository(RepositoryGroup) ProjectRepository
	NewSessionRepository(RepositoryGroup) SessionRepository
	NewSampleRepository(RepositoryGroup) VoiceSampleRepository
	Init(args ...interface{})
}

func Get(identifier string) Backend {
	return backends[identifier]
}

var backends map[string]Backend

func Register(identifier string, backend Backend) {
	backends[identifier] = backend
}

func init() {
	backends = make(map[string]Backend)
}
