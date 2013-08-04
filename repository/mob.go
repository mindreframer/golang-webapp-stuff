package repository

import (
	"github.com/jingweno/progmob/model"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type mob struct{}

func (m *mob) ensureIndex() error {
	index := mgo.Index{
		Key:        []string{"repofullname"},
		Unique:     true,
		DropDups:   true,
		Background: true,
	}

	return m.C().EnsureIndex(index)
}

func (m *mob) ForUser(user *model.User) ([]model.Mob, error) {
	query := Query{"_id": Query{"$in": user.Mobs}}
	var mobs []model.Mob
	err := m.C().Find(query).All(&mobs)
	if err != nil {
		return nil, err
	}

	return mobs, nil
}

func (m *mob) First(query Query) *model.Mob {
	var mob *model.Mob
	m.C().Find(query).One(&mob)

	return mob
}

func (m *mob) Upsert(mob *model.Mob) (*model.Mob, error) {
	query := Query{"repofullname": mob.RepoFullName}
	changed, err := m.C().Upsert(query, mob)
	if err != nil {
		return nil, err
	}

	if changed.UpsertedId == nil {
		return m.First(query), nil
	}

	mob.ObjectId = changed.UpsertedId.(bson.ObjectId)
	return mob, nil
}

func (m *mob) C() *mgo.Collection {
	return db.C("mobs")
}
