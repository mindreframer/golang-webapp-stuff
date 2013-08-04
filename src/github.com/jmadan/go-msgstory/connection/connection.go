package connection

import (
	"labix.org/v2/mgo"
	"log"
	"os"
)

var (
	mgoSession   *mgo.Session
	databaseName = "app15287973"
)

func GetDBSession() *mgo.Session {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial(os.Getenv("MONGOHQ_URL"))
		if err != nil {
			log.Println(err) // no, not really
		}
	}
	return mgoSession.Clone()
}

func WithCollection(collection string, s func(*mgo.Collection) error) error {
	dbSession := GetDBSession()
	defer dbSession.Close()
	c := dbSession.DB(databaseName).C(collection)
	return s(c)
}

// func ExecQuery (collection string, q interface{}, skip int, limit int) (searchResults []string, searchErr string) {
//     searchErr     = ""
//     searchResults = []Person{
//     query := func(c *mgo.Collection) error {
//         fn := c.Find(q).Skip(skip).Limit(limit).All(&searchResults)
//         if limit < 0 {
//             fn = c.Find(q).Skip(skip).All(&searchResults)
//         }
//         return fn
//     }
//     search := func() error {
//         return withCollection("person", query)
//     }
//     err := search()
//     if err != nil {
//         searchErr = "Database Error"
//     }
//     return
// }
