package entities

import (
	"time"
)

type Session struct {
	Id               int
	ProjectId        int
	RelatedSessionId int
	Created          time.Time
	Completed        bool
	UserAgent        string
	FormValues       []FormFieldValue
}

type FormFieldValue struct {
	FieldSlug string
	Value     interface{}
}
