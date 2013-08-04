package repository

import (
	"fmt"
	"github.com/jingweno/progmob/model"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"strings"
)

type repo struct{}

func (r *repo) ensureIndex() error {
	index := mgo.Index{
		Key:        []string{"id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
	}

	return r.C().EnsureIndex(index)
}

func (r *repo) Insert(repos ...*model.Repo) error {
	for _, repo := range repos {
		repo.ObjectId = bson.NewObjectId()
	}

	return r.C().Insert(repos)
}

func (r *repo) UpsertAll(repos ...model.Repo) ([]model.Repo, error) {
	errMsg := []string{}
	ids := []int{}
	for _, repo := range repos {
		err := r.Upsert(&repo)
		if err != nil {
			errMsg = append(errMsg, err.Error())
		}

		ids = append(ids, repo.ID)
	}

	result := []model.Repo{}
	if len(errMsg) == 0 {
		query := Query{"id": Query{"$in": ids}}
		err := r.All(query, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	}

	return nil, fmt.Errorf(strings.Join(errMsg, ", "))
}

func (r *repo) Upsert(repo *model.Repo) error {
	_, err := r.C().Upsert(Query{"id": repo.ID}, repo)

	return err
}

func (r *repo) All(query Query, all *[]model.Repo) error {
	return r.C().Find(query).All(all)
}

func (r *repo) C() *mgo.Collection {
	return db.C("repos")
}
