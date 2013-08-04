package mgorx

import (
  "github.com/robfig/revel"
  "labix.org/v2/mgo"
)

var (
    mgoSession    *mgo.Session
    database_name string
)

func getSession () *mgo.Session {
    if mgoSession == nil {
        var err error
        url, found := revel.Config.String("db.url")
        if !found {
             panic("db.url not set in config")
        }
        database_name, found = revel.Config.String("db.name")
        if !found {
             panic("db.name not set in config")
        }
        mgoSession, err = mgo.Dial(url)
        if err != nil {
             panic(err) // no, not really
        }
    }
    return mgoSession.Clone()
}
