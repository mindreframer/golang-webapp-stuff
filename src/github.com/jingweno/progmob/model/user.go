package model

import (
	"labix.org/v2/mgo/bson"
	"time"
)

type User struct {
	model ",inline"

	ID          int
	Login       string
	Email       string
	HTMLURL     string
	AvatarURL   string
	AccessToken string
	CreatedAt   time.Time
	Mobs        []bson.ObjectId
}
