MGORX
=====

This is an experimental Go lib to wrap mgo to minimize code redundancy

Usage
-----

define your models like this

```go
package models

import (
  "github.com/robfig/revel"
  "labix.org/v2/mgo/bson"
  "github.com/tanema/mgorx"
  "time"
)

type User struct {
  mgorx.Document    "-"
  Id          bson.ObjectId "_id,omitempty"
  Username    string
  Email       string
  Password    string
  Created_at  time.Time
}

func Users() *mgorx.Collection{
  return mgorx.GetCollection(User{})
}

func (user *User) Validate(v *revel.Validation) {
  v.Required(user.Username).Message("Your Username is required!")
  v.Required(user.Email).Message("Your Email is required!")
  v.Required(user.Password).Message("Your Password is required!")
}
```

And then you are able to do something like this in your controllers for a more natural CRUD

```go
package controllers

import (
  "github.com/robfig/revel"
  "yourapp/app/models"
)

type UsersController struct {
  *revel.Controller
}

func (c UsersController) New() revel.Result{
  return c.Render()
}

func (c UsersController) Create(user models.User) revel.Result {
  saved := models.Signatures().Create(&signature, c.Validation)
  if !saved || c.Validation.HasErrors() {
    return c.Render(user)
  }
  return c.Redirect(App.Index)
}

func (c UsersController) Edit(id string) revel.Result {
  var user models.User
  models.Users().Find(&user, id)
  return c.Render(user)
}

func (c UsersController) Update(id string, updates models.User) revel.Result {
  var user models.User
  models.Users().Find(&user, id)
  saved := user.Update(updates, c.Validation)
  if !saved || c.Validation.HasErrors() {
		return c.Render(user)
	}
  return c.Redirect(App.Index)
}

func (c UsersController) Delete(id string) revel.Result {
  models.Users().Delete(id)
  return c.Redirect(App.Index)
}
```

woo I think that is a bit more sane

Last you can do other stuff cool like this

```go
func (c App) Index() revel.Result {
  users := []models.User{}
  models.Users().All(&users, bson.M{"order": "-_id"})
  count, _ := models.Users().Count(nil)
  return c.Render(users, count)
}
```

It provides methods such as Collection.Where(result, query, options) and Collection.All(result, options), Collection.Count(query)

All the queries can be done using mgo.bson.M{} and same with the options

