package models

import (
	"errors"
	"time"
	"appengine"
	"appengine/datastore"
	"net/http"
	"fmt"
	"encoding/json"
)

type Header struct {
	Key,  Value string
}

type Code struct {
	Title              string
	UserId             int64
	Read,  Started     bool
	Id                 int64
	CreatedAt,  ReadAt time.Time
	Headers            []Header
}

func FindCodesByUserId(ctx appengine.Context, userId int64) ([]Code, error) {
	q := datastore.NewQuery("code").Filter("UserId =", userId).Order("-CreatedAt")

	var codes []Code
	keys, e := q.GetAll(ctx, &codes)

	if e == nil {
		for i, key := range keys {
			codes[i].Id = key.IntID()
		}
	}

	return codes, e
}

func DestroyCodeBy(ctx appengine.Context, id int64) (bool, error) {
	k := datastore.NewKey(ctx, "code", "", id, nil)
	e := datastore.Delete(ctx, k)

	if e != nil {
		return false, e
	} else {
		return true, nil
	}
}

func UpdateCode(ctx appengine.Context, id int64, headers []Header) (bool, error) {
	k := datastore.NewKey(ctx, "code", "", id, nil)
	var code Code
	e := datastore.Get(ctx, k, &code)

	if e != nil {
		return false, errors.New(fmt.Sprintf("Not found: %s", e.Error()))
	} else {
		if code.Started {
			code.Read = true
			code.ReadAt = time.Now()
			code.Headers = headers
			_, err := datastore.Put(ctx, k, &code)

			if err != nil {
				return false, errors.New(fmt.Sprintf("Failed to update: %s", err.Error()))
			} else {
				return true, nil
			}
		} else {
			return false, errors.New("Code tracking was not started!")
		}
	}
}

func StartTrackingCode(ctx appengine.Context, id int64) (bool, error) {
	k := datastore.NewKey(ctx, "code", "", id, nil)
	var code Code
	e := datastore.Get(ctx, k, &code)

	if e != nil {
		return false, errors.New(fmt.Sprintf("Not found: %s", e.Error()))
	} else {
		code.Started = true
		_, err := datastore.Put(ctx, k, &code)

		if err != nil {
			return false, errors.New(fmt.Sprintf("Failed to update: %s", err.Error()))
		} else {
			return true, nil
		}
	}

}

func (u *Code) IsValid() (bool, error) {
	valid := u.Title != "" && u.UserId != 0

	if valid == false {
		if u.Title == "" {
			return valid, errors.New("Name is empty")
		} else if u.UserId == 0 {
			return valid, errors.New("UserId is empty")
		}
	}

	return valid, nil
}

func (c *Code) Save(ctx appengine.Context) (bool, error, error) {
	c.CreatedAt = time.Now()

	valid, validationErr := c.IsValid()
	if valid {
		gKey := datastore.NewIncompleteKey(ctx, "code", nil)
		key, err := datastore.Put(ctx, gKey, c)

		if err == nil {
			c.Id = key.IntID()
			return true, nil, nil
		} else {
			return false, nil, err
		}
	} else {
		return false, validationErr, nil
	}
}

func (c Code) AsJson() string {
	d, _ := json.Marshal(c)
	return string(d)
}

func PopulateCode(r *http.Request, postReq bool) *Code {
	if postReq {
		return &Code{
			Title: r.PostFormValue("title"),
		}
	} else {
		return &Code{
			Title: r.FormValue("title"),
		}
	}
}
