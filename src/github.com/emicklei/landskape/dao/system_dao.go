package dao

import (
	"github.com/emicklei/landskape/model"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type SystemDao struct {
	Collection *mgo.Collection
}

func (self SystemDao) Exists(scope, id string) bool {
	_, err := self.FindById(scope, id)
	return err == nil
}

func (self SystemDao) Save(app *model.System) error {
	_, err := self.Collection.Upsert(bson.M{"scope": app.Scope, "id": app.Id}, app) // ChangeInfo
	return err
}

func (self SystemDao) FindAll(scope string) ([]model.System, error) {
	query := bson.M{"scope": scope}
	result := []model.System{}
	err := self.Collection.Find(query).All(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (self SystemDao) FindById(scope, id string) (model.System, error) {
	query := bson.M{"_id": id, "scope": scope}
	result := model.System{}
	err := self.Collection.Find(query).One(&result)
	return result, err
}

func (self SystemDao) RemoveById(scope, id string) error {
	return self.Collection.Remove(bson.M{"scope": scope, "id": id})
}
