package model

import (
	"fmt"
	"labix.org/v2/mgo/bson"
)

type model struct {
	ObjectId    bson.ObjectId "_id,omitempty"
}

func (m *model) ObjectIdString() string {
	return fmt.Sprintf("%x", string(m.ObjectId))
}
