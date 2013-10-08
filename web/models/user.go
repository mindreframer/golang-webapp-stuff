package models

import (
	"errors"
	"net/http"
	"appengine"
	"appengine/datastore"
)

type User struct {
	Name, Email, Password string
	Id int64
}

func (u *User) IsValid(ctx appengine.Context) (bool, error) {
	valid := u.Name != "" && u.Email != "" && u.Password != ""
	if valid == false {
		if u.Name == "" {
			return  valid, errors.New("Name is empty")
		} else if u.Email == "" {
			return  valid, errors.New("Email is empty")
		} else if u.Password == "" {
			return  valid, errors.New("Password is empty")
		}
	} else {
		q := datastore.NewQuery("user").Filter("Email =", u.Email)
		var users []User
		keys, _ := q.GetAll(ctx, &users)

		if len(keys) > 0 {
			valid = false
			return valid, errors.New("Email address is already in use")
		}
	}

	return valid, nil
}

func (u *User) Save(ctx appengine.Context) (bool, error, error) {
	valid, validationErr := u.IsValid(ctx)

	if valid {
		gKey := datastore.NewIncompleteKey(ctx, "user", nil)
		key, err := datastore.Put(ctx, gKey, u)
		u.Id = key.IntID()

		if err == nil {
			return true, nil, nil
		} else {
			return false, validationErr, err
		}
	} else {
		return false, validationErr, nil
	}
}

func Authenticate(ctx appengine.Context, r *http.Request) (*User, error) {
	email, password := r.PostFormValue("email"), r.PostFormValue("password")

	if email != "" && password != "" {
		var users []User
		q := datastore.NewQuery("user").Filter("Email =", email)
		keys, err := q.GetAll(ctx, &users)

		if err != nil {
			return nil, err
		}

		if len(users) == 0 {
			return nil, errors.New("Invalid Email or Password.")
		}

		user := users[0];
		user.Id = keys[0].IntID()

		valid := user.Password == password

		if valid {
			return &user, nil
		} else {
			return nil, errors.New("Invalid Email or Password.")
		}
	} else {
		return nil, errors.New("Email or Password field is empty")
	}
}

func PopulateUser(r *http.Request, postReq bool) *User {
	if postReq {
		return &User{
			Name: r.PostFormValue("name"),
			Email: r.PostFormValue("email"),
			Password: r.PostFormValue("password"),
		}
	} else {
		return &User{
			Name: r.FormValue("name"),
			Email: r.FormValue("email"),
			Password: r.FormValue("password"),
		}
	}
}
