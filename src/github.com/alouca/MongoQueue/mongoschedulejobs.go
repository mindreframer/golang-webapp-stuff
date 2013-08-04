// Copyright 2012-2013 Andreas Louca <andreas@louca.org>. All rights reserved.
// Use of this source code is governed by the 2-clause BSD license
// license that can be found in the LICENSE file.

/*
Package mongoqueue provides a job queue, which uses Mongo as a backend storage engine. 
It supports a sophisticated feature set,  facilitating fine-grained job queueing.

See: https://github.com/alouca/MongoQueue
*/

package mongoqueue

import (
	"code.google.com/p/go-uuid/uuid"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

type MongoScheduleJobs struct {
	MongoSession *mgo.Session
	Database     string
}

func NewMongoScheduleJobs(database, server string) (*MongoScheduleJobs, error) {
	var err error
	msj := new(MongoScheduleJobs)
	msj.MongoSession, err = mgo.Dial(server)
	if err != nil {
		return nil, err
	}
	msj.Database = database

	go msj.Start()

	return msj, nil
}

func (m *MongoScheduleJobs) ScheduleJob(name string, queue string, x interface{}, p int, interval int64) {
	coll := m.MongoSession.DB(m.Database).C("mongoschedulejobs")

	coll.Insert(bson.M{
		"name":     name,
		"queue":    queue,
		"next-run": time.Now().Unix(),
		"priority": p,
		"interval": interval,
		"data":     x,
	})
}

func (m *MongoScheduleJobs) DeleteJob(name string) {
	coll := m.MongoSession.DB(m.Database).C("mongoschedulejobs")

	coll.Remove(bson.M{"name": name})
}

func (m *MongoScheduleJobs) Start() {
	t := time.NewTicker(time.Second * 10)
	coll := m.MongoSession.DB(m.Database).C("mongoschedulejobs")

	for {
		select {
		case <-t.C:
			now := time.Now().Unix()

			q := coll.Find(bson.M{"next-run": bson.M{"$lte": now}}).Iter()

			var res bson.M

			for q.Next(&res) {
				l.Info("Adding job for execution in queue %s\n", res["queue"])
				// Queue the job
				m.addJob(res["queue"].(string), res["data"], res["priority"].(int))

				// Set next-run
				coll.Update(bson.M{"_id": res["_id"]}, bson.M{"$set": bson.M{"next-run": now + res["interval"].(int64)}})
			}

		}

	}
}

func (m *MongoScheduleJobs) addJob(queue string, x interface{}, p int) (string, error) {
	coll := m.MongoSession.DB(m.Database).C(queue)

	id := uuid.NewRandom().String()

	now := time.Now().Unix()
	err := coll.Insert(bson.M{
		"id":         id,
		"inprogress": false,
		"failed":     false,
		"priority":   p,
		"retries":    0,
		"added":      now,
		"runat":      now,
		"data":       x})

	if err != nil {
		l.Fatal("Error inserting new task: %s\n", err.Error())
	}

	return id, err

}
