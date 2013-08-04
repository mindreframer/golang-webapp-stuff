package user

import (
	"encoding/json"
	"labix.org/v2/mgo/bson"
	"testing"
	"time"
)

func Test_CreatePerson(t *testing.T) {
	user := User{bson.NewObjectId(), 12, "Jasdeep", "jasdeepm@gmail.com", "JD", "07818912893", rels{}, time.Now()}
	res := user.CreateUser()
	if res {
		t.Log("Test_CreateUser PASSED")
	} else {
		t.Fail()
	}
}

// func Test_GetByHandle(t *testing.T) {
// 	res := GetByHandle("jasdeepm@gmail.com")
// 	if res.Name == "JD" {
// 		t.Log("PASSED")
// 	} else {
// 		t.Fail()
// 	}
// }

func Test_GetUser(t *testing.T) {
	// person := User{0, 1, "Jasdeep", "jasdeepm@gmail.com", "JD", "07818912893", rels{}, time.Now()}
	temp_person := User{}
	tUser := User{}
	// temp_person.SetUserid(1)
	temp_person.SetEmail("jasdeepm@gmail.com")
	str := temp_person.GetUser()
	err := json.Unmarshal([]byte(str), &tUser)
	if err != nil {
		t.Fail()
		t.Log(err)
	} else {
		if tUser.GetName() == "Jasdeep" {
			t.Log("Test_GetUser PASSED")
		} else {
			t.Failed()
		}
	}
}

// func Test_CreateUserLogin(t *testing.T) {
// 	useremail := "test@test.com"
// 	password := "password"
// 	user_id := CreateUserLogin(useremail, password)
// 	uid := getUserByEmail(useremail)

// 	if uid == user_id {
// 		t.Log("Test_CreateUserLogn PASSED")
// 	} else {
// 		t.Fail()
// 		t.Log("Test_CreateUserLogn FAILED")
// 	}
// }
