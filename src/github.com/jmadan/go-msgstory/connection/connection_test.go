package connection

import (
	User "github.com/jmadan/go-msgstory/user"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"testing"
)

func Test_DatabaseConnection(t *testing.T) {
	searchResults := *User.User{}
	query := func(c *mgo.Collection) error {
		fn := c.Find(bson.M{"_id": "516f921a8a00b54a977f5fa7"}).One(&searchResults)
		return fn
	}
	search := func() error {
		return WithCollection("circle", query)
	}
	err := search()
	if err != nil {
		t.Error("Database Error")
		t.Fail()
	}

	if searchResults != nil {
		t.Log("DATABASE CONNECTION PASSED")
	} else {
		t.Fail()
	}
}
