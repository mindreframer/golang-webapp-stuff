package mgo

import (
	"fmt"
	. "github.com/ku-ovdp/api/entities"
	. "github.com/ku-ovdp/api/repository"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

func (m *mgoBackend) NewSessionRepository(repositories RepositoryGroup) SessionRepository {
	if m.db == nil {
		panic("Not initialized! Missing call to Init()?")
	}
	sr := &sessionRepository{
		m,
		m.db.C("sessions"),
		repositories["projects"].(ProjectRepository),
		repositories,
	}

	return sr
}

type sessionRepository struct {
	backend  *mgoBackend
	c        *mgo.Collection
	projects ProjectRepository
	group    RepositoryGroup
}

func (sr *sessionRepository) Get(id int) (Session, error) {
	result := Session{}
	q := sr.c.Find(bson.M{"id": id})
	err := q.One(&result)
	if err == mgo.ErrNotFound {
		return result, NewErrNotFound(Session{}, id)
	}
	return result, err
}

func (sr *sessionRepository) Put(session Session) (Session, error) {
	if session.Id == 0 {
		session.Id = sr.backend.nextId("session")
	}
	_, err := sr.c.Upsert(bson.M{"id": session.Id}, session)
	if err != nil {
		return Session{}, err
	}
	return sr.Get(session.Id)
}

func (sr *sessionRepository) Remove(id int) error {
	return fmt.Errorf("Not implemented")
}

func (sr *sessionRepository) Scan(projectId int, from, to int) ([]Session, error) {
	results := []Session{}

	if _, err := sr.projects.Get(projectId); err != nil {
		return results, err
	}

	query := bson.M{"id": bson.M{"$gte": from}, "projectid": projectId}
	if to > 0 {
		query["id"].(bson.M)["$lte"] = to
	}

	q := sr.c.Find(query)
	err := q.All(&results)

	return results, err
}

func (sr *sessionRepository) Group() RepositoryGroup {
	return sr.group
}
