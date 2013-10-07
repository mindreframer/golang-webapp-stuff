// Package mgo implements a MongoDB storage backend for API entities
package mgo

import (
	"fmt"
	"github.com/ku-ovdp/api/persistence"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"os"
)

type mgoBackend struct {
	session *mgo.Session
	db      *mgo.Database
}

var MongoURIEnvVar = "MONGOURI"

func (m *mgoBackend) Init(args ...interface{}) {
	var err error
	uri := os.Getenv(MongoURIEnvVar)
	if uri == "" {
		panic(fmt.Sprintln("Environment variable", MongoURIEnvVar, "is empty!"))
	}
	log.Println("Connecting to", uri, "...")
	m.session, err = mgo.Dial(uri)
	if err != nil {
		panic(err)
	}
	log.Println("Connected.")
	m.db = m.session.DB("")
}

func (m *mgoBackend) nextId(namespace string) int {
	ctr := struct {
		Id   string `bson:"_id"`
		Next int
	}{namespace, 1}

	m.db.C("counters").EnsureIndex(mgo.Index{
		Key:    []string{"_id", "next"},
		Unique: true,
	})

	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"next": 1}},
		ReturnNew: true,
	}
	_, err := m.db.C("counters").Find(bson.M{"_id": namespace}).Apply(change, &ctr)
	if err == mgo.ErrNotFound {
		// @todo this is a race
		err = m.db.C("counters").Insert(ctr)
		log.Println(err)
	}
	return ctr.Next
}

func init() {
	persistence.Register("mgo", new(mgoBackend))
}
