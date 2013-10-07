// Package dummy implements dummy storage for API entities
package dummy

import (
	. "github.com/ku-ovdp/api/entities"
	. "github.com/ku-ovdp/api/repository"
	"time"
)

func (d dummyBackend) NewSampleRepository(repositories RepositoryGroup) VoiceSampleRepository {
	return sampleRepository{
		dummySampleData,
		repositories["sessions"].(SessionRepository),
		repositories,
	}
}

type sampleRepo map[int]VoiceSample

type sampleRepository struct {
	sampleRepo
	sessions SessionRepository
	group    RepositoryGroup
}

var dummySampleData = map[int]VoiceSample{
	1: {Id: 1, SessionId: 1,
		Created:  time.Now().Add(time.Hour * -24 * 14),
		Bitrate:  24000,
		MimeType: "audio/mpeg",
		Size:     12537,
		AudioURL: "http://watbutton.com/wat3.mp3",
	},
}

func (sr sampleRepository) Get(sessionId, id int) (VoiceSample, error) {
	if obj, ok := sr.sampleRepo[id]; ok {
		return obj, nil
	} else {
		return VoiceSample{}, NewErrNotFound(VoiceSample{}, id)
	}
}

func (sr sampleRepository) Put(sample VoiceSample) (VoiceSample, error) {
	if sample.Id == 0 {
		sample.Id = len(sr.sampleRepo) + 1
	}
	sr.sampleRepo[sample.Id] = sample
	return sample, nil
}

func (sr sampleRepository) Remove(sessionId, id int) error {
	delete(sr.sampleRepo, id)
	return nil
}

func (sr sampleRepository) Scan(sessionId int, from, to int) ([]VoiceSample, error) {
	results := []VoiceSample{}
	if _, err := sr.sessions.Get(sessionId); err != nil {
		return results, err
	}
	for id, value := range sr.sampleRepo {
		if id < from {
			continue
		}
		if id > to && to != 0 {
			continue
		}
		if value.SessionId != sessionId {
			continue
		}
		results = append(results, value)
	}
	return results, nil
}

func (sr sampleRepository) Group() RepositoryGroup {
	return sr.group
}
