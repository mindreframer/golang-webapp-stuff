// Copyright 2012 Andreas Louca <andreas@louca.org>. All rights reserved.
// Use of this source code is governed by the 2-clause BSD license
// license that can be found in the LICENSE file.

package mongoqueue

import (
	"encoding/json"
	tl "log"
	"testing"
	"time"
)

var mq *MongoQueue
var msj *MongoScheduleJobs

func init() {
	mq = NewMongoQueue("mq", "testing", "127.0.0.1", &MongoQueueSettings{Cleanup: 30, LockLimit: 5, RetryLimit: 2, MinBackoff: 1, MaxBackoff: 3, MaxDoublings: 2, AgeLimit: 25})
	msj, _ = NewMongoScheduleJobs("mq", "127.0.0.1")
}

type Testdata struct {
	Prio int
	Test string
}

func TestInsert(t *testing.T) {
	err := mq.Truncate()

	for i := 0; i < 50; i++ {
		d, _ := json.Marshal(Testdata{i, "lalala"})
		id, err := mq.Add(string(d), "", i)
		if err != nil {
			t.Fatal(err)
		} else {
			if id == "" {
				t.Fatal("Error getting ID for inserted job")
			} else {
				t.Logf("Added new job with ID %s\n", id)
			}
		}
	}

	c, err := mq.Count()

	if err != nil {
		t.Fatal(err)
	} else if c != 50 {
		t.Fatal("Inserted 50 tasks, but count is not 50: %d", c)
	}
}

func TestRetrieve(t *testing.T) {
	stats, _ := mq.Stats()
	tl.Printf("Stats: %+v\n", stats)
	// We are expecing 50 documents to be retrieved in correct order
	for i := 0; i < 50; i++ {
		_, d, err := mq.Pop()
		if err != nil {
			t.Fatal("Error getting data: %s", err.Error())
		}

		if dstring, ok := d.(string); ok {
			var data Testdata
			err = json.Unmarshal([]byte(dstring), &data)
			if err != nil {
				t.Logf("Error unmarshalling data: %s\n", err.Error())
			} else {
				t.Logf("Got %d for %d\n", data.Prio, i)
			}
		} else {
			t.Logf("Got %v instead of string\n", d)
		}

	}
}

func TestExpire(t *testing.T) {
	d, _ := json.Marshal(Testdata{10, "testing job expiration"})
	pid := "testing-Pid"
	var id string

	aid, err := mq.Add(string(d), "", 10)
	if err != nil {
		t.Fatal(err)
	}
	if aid == "" {
		t.Fatal("Error getting ID for insert")
	}

	t.Logf("Added test job with ID: %s\n", aid)
	// Acquire the lock
	id, job, err := mq.Lock(pid)

	if err != nil {
		t.Fatalf("Error acquiring lock on test-expiration: %s\n", err.Error())
	}

	if aid != id {
		t.Fatal("Failure comparing acquired lock versus added id: %s vs %s\n", aid, id)
	}

	if sjob, ok := job.(string); ok {
		var data Testdata
		err = json.Unmarshal([]byte(sjob), &data)

		t.Logf("Lock aquired on job\n")

		// Testing to acquire a new job (should return an error)
		_, job, err = mq.Lock(pid)
		if job == nil {
			t.Logf("Got empty\n")
		} else {
			t.Fatal("MQ should be empty, but document was returned!\n")
		}

		// Upon calling expire, a job should be available to retrieve
		time.Sleep(time.Second * 10)
		mq.Cleanup()
		id, job, err = mq.Lock(pid)
		if err != nil {
			t.Fatalf("Got empty, when expecting a job: %s\n", err.Error())
		} else {
			t.Logf("Acquired lock on Job ID: %s\n", id)
			err := mq.Complete(id)

			if err != nil {
				t.Fatalf("Job failed to be marked as complete: %s\n", err.Error())
			}
		}
	}
	stats, _ := mq.Stats()
	tl.Printf("Stats: %+v\n", stats)
}

func TestFailed(t *testing.T) {
	t.Logf("Starting test for Failures\n")
	d, _ := json.Marshal(Testdata{10, "testing job failure"})
	pid := "testing-failed-Pid"
	var id string

	id, err := mq.Add(string(d), "", 10)
	if err != nil {
		t.Fatal(err)
	}

	if id == "" {
		t.Fatal("Error getting ID for insert")
	} else {
		t.Logf("Added new test job with id: %s\n", id)
	}
	retryLimit := 2

	for i := 0; i <= retryLimit; i++ {
		// Acquire the lock
		id, job, err := mq.Lock(pid)
		if sjob, ok := job.(string); ok {
			var data Testdata
			err = json.Unmarshal([]byte(sjob), &data)

			tl.Print("Lock aquired on job\n")
			err = mq.Fail(id)
			time.Sleep(time.Second * 5)
			if err != nil {
				t.Fatalf("Error during fail() call: %s\n", err.Error())
			}
		} else {
			t.Fatalf("Failure to acquire lock: %s\n", err.Error())
		}
	}

	id, job, err := mq.Lock(pid)
	if sjob, ok := job.(string); ok {
		var data Testdata
		err = json.Unmarshal([]byte(sjob), &data)

		t.Fatal("Lock aquired on job, when it should not!\n")
	} else {
		t.Logf("Correct behavior: Failed to acquire lock: %s\n", err.Error())
	}
}

func TestSchedule(t *testing.T) {
	t.Logf("Starting test for scheduling\n")

	msj.ScheduleJob("testing", map[string]string{"testing": "testing"}, 1, 30)

	ticker := time.NewTicker(time.Second * 10)

	for {
		select {
		case <-ticker.C:
			_, d, err := mq.Pop()
			if err != nil {
				t.Fatal("Error getting data: %s", err.Error())
			} else {
				t.Logf("Got data %+v\n", d)
			}
		}
	}
}

/*
func TestTruncate(t *testing.T) {
	err := mq.Truncate()
	if err != nil {
		t.Fatal(err)
	}

	c, err := mq.Count()

	if err != nil {
		t.Fatal(err)
	} else if c > 0 {
		t.Fatal("Truncated queue, but count > 0")
	}
}
*/
