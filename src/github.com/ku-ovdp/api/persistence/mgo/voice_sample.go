package mgo

import (
	"fmt"
	. "github.com/ku-ovdp/api/entities"
	. "github.com/ku-ovdp/api/repository"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

func (m *mgoBackend) NewSampleRepository(repositories RepositoryGroup) VoiceSampleRepository {
	if m.db == nil {
		panic("Not initialized! Missing call to Init()?")
	}
	sr := &sampleRepository{
		m,
		m.db.C("sessions"),
		repositories["sessions"].(SessionRepository),
		repositories,
	}

	return sr
}

type sampleRepository struct {
	backend  *mgoBackend
	c        *mgo.Collection
	sessions SessionRepository
	group    RepositoryGroup
}

func (sr *sampleRepository) Get(sessionId, id int) (VoiceSample, error) {
	result := VoiceSample{}
	q := sr.c.Find(bson.M{"id": id})
	err := q.One(&result)
	if err == mgo.ErrNotFound {
		return result, NewErrNotFound(Session{}, id)
	}
	return result, err
}

func (sr *sampleRepository) Put(sample VoiceSample) (VoiceSample, error) {
	if sample.SessionId == 0 {
		return VoiceSample{}, fmt.Errorf("Invalid SessionId")
	}
	if sample.Id == 0 {
		sample.Id = sr.backend.nextId(fmt.Sprintf("voice_samples:%d", sample.SessionId))
	}
	_, err := sr.c.Upsert(bson.M{"id": sample.Id}, sample)
	if err != nil {
		return VoiceSample{}, err
	}
	return sr.Get(sample.SessionId, sample.Id)
}

func (sr *sampleRepository) Remove(sessionId, id int) error {
	return fmt.Errorf("Not implemented")
}

func (sr *sampleRepository) Scan(sessionId int, from, to int) ([]VoiceSample, error) {
	results := []VoiceSample{}

	if _, err := sr.sessions.Get(sessionId); err != nil {
		return results, err
	}

	query := bson.M{"id": bson.M{"$gte": from}, "sessionid": sessionId}
	if to > 0 {
		query["id"].(bson.M)["$lte"] = to
	}

	q := sr.c.Find(query)
	err := q.All(&results)

	return results, err
}

func (sr *sampleRepository) Group() RepositoryGroup {
	return sr.group
}
