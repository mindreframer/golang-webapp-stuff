package dao

import (
	"github.com/emicklei/landskape/model"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type ConnectionDao struct {
	Collection *mgo.Collection
}

func (self ConnectionDao) FindAllMatching(scope string, filter model.ConnectionsFilter) ([]model.Connection, error) {
	query := bson.M{"scope": scope}
	if len(filter.Types) > 0 {
		query["type"] = bson.M{"$in": filter.Types}
	}
	if len(filter.Centers) > 0 {
		froms := bson.M{"from": bson.M{"$in": filter.Centers}}
		tos := bson.M{"to": bson.M{"$in": filter.Centers}}
		query["$or"] = []bson.M{froms, tos}
	} else {
		if len(filter.Froms) > 0 {
			query["from"] = bson.M{"$in": filter.Froms}
		}
		if len(filter.Tos) > 0 {
			query["to"] = bson.M{"$in": filter.Tos}
		}
	}
	//model.Debug("query", query)
	result := []model.Connection{}
	err := self.Collection.Find(query).All(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (self ConnectionDao) Save(con model.Connection) error {
	query := bson.M{"scope": con.Scope, "from": con.From, "to": con.To, "type": con.Type}
	_, err := self.Collection.Upsert(query, con) // ChangeInfo
	return err
}

func (self ConnectionDao) Remove(con model.Connection) error {
	query := bson.M{"scope": con.Scope, "from": con.From, "to": con.To, "type": con.Type}
	return self.Collection.Remove(query)
}

func (self ConnectionDao) RemoveAllToOrFrom(scope, toOrFrom string) error {
	return nil
}
