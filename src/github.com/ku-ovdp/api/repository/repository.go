// Package repository implements the storage interface for API entities
package repository

import (
	. "github.com/ku-ovdp/api/entities"
)

type ProjectRepository interface {
	Get(id int) (Project, error)
	Put(project Project) (Project, error)
	Remove(id int) error
	Scan(from, to int) ([]Project, error)
	Group() RepositoryGroup
}

type SessionRepository interface {
	Get(id int) (Session, error)
	Put(session Session) (Session, error)
	Remove(id int) error
	Scan(projectId int, from, to int) ([]Session, error)
	Group() RepositoryGroup
}

type VoiceSampleRepository interface {
	Get(sessionId, id int) (VoiceSample, error)
	Put(sample VoiceSample) (VoiceSample, error)
	Remove(sessionId, id int) error
	Scan(sessionId int, from, to int) ([]VoiceSample, error)
	Group() RepositoryGroup
}

type RepositoryGroup map[string]interface{}

func NewRepositoryGroup() RepositoryGroup {
	return make(RepositoryGroup, 0)
}

func (rg RepositoryGroup) Projects() ProjectRepository {
	return rg["projects"].(ProjectRepository)
}

func (rg RepositoryGroup) Sessions() SessionRepository {
	return rg["sessions"].(SessionRepository)
}

func (rg RepositoryGroup) VoiceSamples() VoiceSampleRepository {
	return rg["samples"].(VoiceSampleRepository)
}
