# MongoQueue

MongoQueue is Job queue written in Go, which uses Mongo as a backend storage engine. It supports a sophisticated feature set,  facilitating fine-grained job queueing.

It supports job prioritisation, locking, retries for failed jobs, retry timers, age limits and failure limits.

MongoQueue is under the BSD license, found in the LICENSE file.

Copyright © 2012 Andreas Louca <andreas@louca.org>

## Example usage

To create a new MongoQueue:

`mq := NewMongoQueue("mq", "testing", "127.0.0.1", &MongoQueueSettings{Cleanup: 30, LockLimit: 5, RetryLimit: 2, MinBackoff: 1, MaxBackoff: 3, MaxDoublings: 2, AgeLimit: 25})`

mq is the database, testing is the collection and 127.0.0.1 is the mongo server. The MongoQueueSettings is a struct, which contains all the necessary queue behaviour parameters. All of the time parameters are specified in seconds.

Parameter description:

* **Cleanup**: The number of seconds between the cleanup process runs.
* **LockLimit**: The maximum number of seconds a job can remain locked to a pid
* **RetryLimit**: The maximum number of retry attempts for a failed task
* **MinBackoff**: The minimum number of seconds to wait before retrying a task after it fails.
* **MaxBackoff**: The minimum number of seconds to wait before retrying a task after it fails.
* **MaxDoublings**: The maximum number of times that the interval between failed task retries will be doubled before the increase becomes constant. The constant is: 2**(max_doublings - 1) * min_backoff_seconds.
* **AgeLimit**: The time limit for retrying a failed task, in seconds, measured from the time the task was created.

### Adding a Job

MQ is data-agnostic, so when adding a job an `interface{}` is passed to it. With the job, a priority is also passed, which indicates the priority level of the job being added to the queue. The higher the number, the highest the priority. MQ also allows the programmer to specify its own ID to a job, by passing the ID parameter. If left empty, a UUID will be automatically generated for the queued job.

If you do not wish to use the priority in jobs, specify 0, and MQ will act as FIFO queue.

Adding a job is really simple:

`id, err := mq.Add(map[string][int]{"testing": 1}, "testing-id-1" 10)`

The `Add` calls returns the Job ID, which is used to identify the job in later calls.

### Retrieving Jobs from the queue

There are two ways of retrieving queued jobs:

* With job locking: Provides support for job retries, if the job fails to execute for whatever reason. Locking also denies other instances from de-queueing that job by locking it to a PID, until the locks are expired.
* Without job locking: Jobs are popped from the queue directly according to priority (or FIFO), and are deleted from the queue. No job retries or locking is supported in this mode.

#### Retrieving with Job Locking

When operating in job locking mode, each thread or program that will interact with MQ, must have its own unique program identifier. A PID can be any arbitrary PID.

To lock a job to a PID:

`id, job, err := mq.Lock(pid)`

This will return the job ID as well as the data of the job inserted before. The call might fail if no jobs are available in the queue, and an error will return.

When the program finishes processing the job, it must notify MQ that processing was finished successfully, so the following call must be made:

`err := mq.Complete(id)`

This marks the job as successfully completed, and deletes it from the job queue. 

If for any reason the execution has failed (eg. the remote service was unavailable at that time), and you want to retry the execution of the task at a later time, you must mark the job as failed:

`err := mq.Fail(id)`

This will mark the job as failed, and will be re-queued for execution according to the queue behaviour parameters.

The ids passed in `Complete()` and `Fail()` are the Ids received when `Lock()` is called.

#### Retrieving without job locking

To retrieve a job from MQ without locking the following call must be made:

`d, err := mq.Pop()`

This will remove a job from MQ according to priority (or FIFO, depending what is used during job addition), and return the data to the caller. An error can be generated if the queue is empty. This deletes the job directly from MQ, and cannot use the extra mechanisms for retries.

### Statistics

MQ can provide statistics for pending jobs using:

`stats, err := mq.Stats()`

A struct is returned, with the total jobs currently in queue, in progress and failed jobs:

`type MongoQueueStats struct {
	Total      int
	InProgress int
	Failed     int
}`

### Truncate

To cleanup and remove all the jobs from MQ use the call:

`mq.Truncate()`

### Cleanup process

The cleanup runs at intervals, specified in settings, which releases the locks for jobs that are locked for a period over the specified and marks jobs as permanently failed if necessary. 

### Changelog

* 2013-3-2: Added support for specifying custom Job ID, instead of relaying on Mongo ObjectID field _id. Allows for more flexibility when integrating with existing systems.