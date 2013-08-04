package repository

import (
	"github.com/jingweno/progmob/model"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type user struct{}

func (u *user) ensureIndex() error {
	index := mgo.Index{
		Key:        []string{"login"},
		Unique:     true,
		DropDups:   true,
		Background: true,
	}
	err := u.C().EnsureIndex(index)
	if err != nil {
		return err
	}

	index = mgo.Index{
		Key:        []string{"id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
	}

	return u.C().EnsureIndex(index)
}

func (u *user) Insert(user *model.User) error {
	user.ObjectId = bson.NewObjectId()
	return u.C().Insert(user)
}

func (u *user) AddMob(user *model.User, mob *model.Mob) error {
	return u.C().UpdateId(user.ObjectId, Query{"$addToSet": Query{"mobs": mob.ObjectId}})
}

func (u *user) Upsert(user *model.User) (*model.User, error) {
	changed, err := u.C().Upsert(Query{"id": user.ID}, user)
	if err != nil {
		return nil, err
	}

	if changed.UpsertedId == nil {
		return u.First(Query{"id": user.ID}), nil
	} else {
		user.ObjectId = changed.UpsertedId.(bson.ObjectId)
		return user, nil
	}
}

func (u *user) First(query Query) *model.User {
	var user *model.User
	u.C().Find(query).One(&user)

	return user
}

func (u *user) FirstByObjectId(id string) *model.User {
	var user *model.User
	objectId := bson.ObjectIdHex(id)
	u.C().FindId(objectId).One(&user)

	return user
}

func (u *user) C() *mgo.Collection {
	return db.C("users")
}
