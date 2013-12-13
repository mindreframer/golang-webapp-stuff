package user

import (
	// "encoding/json"
	// "labix.org/v2/mgo/bson"
	"fmt"
	"testing"
	"time"
)

var user = User{
	UserId:      11,
	Name:        "Jasdeep",
	Email:       "jasdeepm@gmail.com",
	Handle:      "JD",
	PhoneNumber: "",
	CreatedOn:   time.Now(),
}

func Test_CreatePerson(t *testing.T) {
	user := User{}
	user.SetEmail("jasdeepm@gmail.com")
	user.Handle = "JD"
	user.Name = "Jasdeep"
	user.UserId = 12
	res := user.CreateUser()
	if res.Success {
		t.Log("Test_CreateUser PASSED")
	} else {
		t.Fail()
	}
}

func Test_GetByHandle(t *testing.T) {
	res := GetByHandle("jasdeepm@gmail.com")
	if res.Name == "JD" {
		t.Log("PASSED")
	} else {
		t.Fail()
	}
}

func Test_GetUser(t *testing.T) {
	// person := User{0, 1, "Jasdeep", "jasdeepm@gmail.com", "JD", "07818912893", rels{}, time.Now()}
	temp_person := User{}
	tUser := User{}
	// temp_person.SetUserid(1)
	temp_person.SetEmail("jasdeepm@gmail.com")
	temp_person.SetUserid(12)
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

func Test_CreateUserLogin(t *testing.T) {
	useremail := "test@test.com"
	password := "password"
	user_id := CreateUserLogin(useremail, password)
	uid := getUserByEmail(useremail)

	if uid == user_id {
		t.Log("Test_CreateUserLogn PASSED")
	} else {
		t.Fail()
		t.Log("Test_CreateUserLogn FAILED")
	}
}

func Test_UserJSON(t *testing.T) {
	var str = user.UserToJSON()
	if str != "Error" {
		t.Log("Test_UserJSON PASSED")
	} else {
		t.Fail()
	}
}

func Test_GetUserById(t *testing.T) {
	res, err := GetUserById("51b8e5b62ffc2c5db5e9b213")
	if err != nil {
		t.Fail()
	} else {
		fmt.Println(res)
		t.Log("PASSED")
	}
}
