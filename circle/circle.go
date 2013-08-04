package circle

import (
	"encoding/json"
	"errors"
	"fmt"
	Connection "github.com/jmadan/go-msgstory/connection"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
)

type Circle struct {
	Name        string   `json:"name" bson:"name"`
	Description string   `json:"description" bson:"description"`
	CreatorID   string   `json:"creator" bson:"creator"`
	CreatedOn   string   `json:"createdon" bson:"createdon"`
	Members     []string `json:"members" bson:"members"`
}

// type JsonCircle struct {
// 	Name string `json:"name" bson:"name"`
// }

func (c *Circle) GetName() string {
	return c.Name
}

func GetUserCircles(userID string) []string {
	var userCircles []string
	searchResults := []Circle{}
	query := func(c *mgo.Collection) error {
		fn := c.Find(bson.M{"members": userID}).All(&searchResults)
		return fn
	}
	search := func() error {
		return Connection.WithCollection("circle", query)
	}
	err := search()
	if err != nil {
		searchErr := "Database Error"
		log.Println(searchErr)
	}

	for i, v := range searchResults {
		userCircles[i] = v.Name
	}
	return userCircles
}

func GetCircleMembers(circleName string) []string {
	var circleMembers []string
	searchResults := []Circle{}
	query := func(c *mgo.Collection) error {
		fn := c.Find(bson.M{"name": circleName}).All(&searchResults)
		return fn
	}
	search := func() error {
		return Connection.WithCollection("circle", query)
	}
	err := search()
	if err != nil {
		searchErr := "Database Error"
		log.Println(searchErr)
	}

	for i, v := range searchResults {
		circleMembers[i] = v.Name
	}
	return circleMembers
}

func (c *Circle) GetJson() string {
	str, err := json.Marshal(&c)
	if err != nil {
		fmt.Println(err.Error())
	}
	return string(str)
}

func (cir *Circle) makeCircle() (bool, error) {
	status := true
	err := errors.New("")
	dbSession := Connection.GetDBSession()
	dbSession.SetMode(mgo.Monotonic, true)
	c := dbSession.DB("msgme").C("circle")

	if cir.CircleExists() {
		status = false
		err = errors.New("Circle already exists with this name")
	} else {
		err = c.Insert(&cir)
		if err != nil {
			status = false
			log.Println(err.Error())
		}
	}

	return status, err
}

// func CircleExists(name string, owner User) (exists bool, msg string) {
func (circle *Circle) CircleExists() (exists bool) {
	dbSession := Connection.GetDBSession()
	dbSession.SetMode(mgo.Monotonic, true)
	c := dbSession.DB("msgme").C("circle")

	result := Circle{}
	err := c.Find(bson.M{"name": circle.Name}).One(&result)
	if err == nil {
		exists = true
	} else {
		exists = false
	}

	return exists
}
