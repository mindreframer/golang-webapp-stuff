package database

import (
	"labix.org/v2/mgo"
)

type Database struct {
	name    string
	session *mgo.Session
}

func (db *Database) Close() {
	db.session.Close()
}

func (db *Database) C(name string) *mgo.Collection {
	return db.DB().C(name)
}

func (db *Database) DropDatabase() error {
	return db.DB().DropDatabase()
}

func (db *Database) DB() *mgo.Database {
	return db.session.DB(db.name)
}

func New(url string) (*Database, error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}
	session.SetMode(mgo.Monotonic, true)

	return &Database{session: session}, nil
}
