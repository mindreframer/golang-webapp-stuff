package model

import (
	"time"
)

type Repo struct {
	model ",inline"

	ID        int
	Owner     string
	Name      string
	FullName  string
	Private   bool
	HTMLURL   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
