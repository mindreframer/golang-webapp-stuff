// Package dummy implements dummy storage for API entities
package dummy

import (
	. "github.com/ku-ovdp/api/entities"
	. "github.com/ku-ovdp/api/repository"
	"time"
)

func (d dummyBackend) NewSessionRepository(repositories RepositoryGroup) SessionRepository {
	return sessionRepository{
		dummySessionData,
		repositories["projects"].(ProjectRepository),
		repositories,
	}
}

type sessionRepo map[int]Session

type sessionRepository struct {
	sessionRepo
	projects ProjectRepository
	group    RepositoryGroup
}

var dummySessionData = map[int]Session{
	1: {Id: 1, ProjectId: 1,
		Created: time.Now().Add(time.Hour * -24 * 14),
		FormValues: []FormFieldValue{
			{FieldSlug: "age", Value: 42},
			{FieldSlug: "gender", Value: "Male"},
			{FieldSlug: "parkinsons", Value: true},
		},
	},
}

func (sr sessionRepository) Get(id int) (Session, error) {
	if obj, ok := sr.sessionRepo[id]; ok {
		return obj, nil
	} else {
		return Session{}, NewErrNotFound(Session{}, id)
	}
}

func (sr sessionRepository) Put(session Session) (Session, error) {
	session.Id = len(sr.sessionRepo) + 1
	sr.sessionRepo[session.Id] = session
	return session, nil
}

func (sr sessionRepository) Remove(id int) error {
	delete(sr.sessionRepo, id)
	return nil
}

func (sr sessionRepository) Scan(projectId int, from, to int) ([]Session, error) {
	results := []Session{}
	if _, err := sr.projects.Get(projectId); err != nil {
		return results, err
	}
	for id, value := range sr.sessionRepo {
		if id < from {
			continue
		}
		if id > to && to != 0 {
			continue
		}
		if value.ProjectId != projectId {
			continue
		}
		results = append(results, value)
	}
	return results, nil
}

func (sr sessionRepository) Group() RepositoryGroup {
	return sr.group
}
