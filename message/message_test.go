package message

import (
	"fmt"
	User "github.com/jmadan/go-msgstory/user"
	"labix.org/v2/mgo/bson"
	"testing"
	"time"
)

var Test_User = User.User{
	Id:          bson.NewObjectId(),
	UserId:      100,
	Name:        "Test User",
	Email:       "test@test.com",
	Handle:      "testuser",
	PhoneNumber: "9008307311",
}

var Msg = "{\"msg_text\":\"this is a text message\",\"user_id\":\"1234567\"}"

var Test_Msg = Message{
	MsgText:    "Hello Delhi",
	UserId:     "123456",
	UserHandle: "handle",
	CreatedOn:  time.Now(),
}

// func Test_SaveMessage(t *testing.T) {
// 	res := Test_Msg.SaveMessage("525e63ade7c5d90002000001")
// 	if res.Success {
// 		fmt.Println(string(res.JsonData))
// 		t.Log("Test_SaveMessage PASSED")
// 	} else {
// 		t.Fail()
// 		t.Log("Test_SaveMessage FAILED")
// 	}
// }

// func Test_SaveMessage(t *testing.T) {
// 	// var m Message
// 	// m.JsonToMsg(Msg)

// 	res := Test_Msg.SaveMessage("51fe91cc552d640002000007")
// 	if res.Success {
// 		fmt.Println(string(res.JsonData))
// 		t.Log("Test_SaveMessage PASSED")
// 	} else {
// 		t.Fail()
// 		t.Log("Test_SaveMessage FAILED")
// 	}
// }

// func Test_GetMessages(t *testing.T) {
// 	res := GetMessages("51fe91cc552d640002000007")
// 	if res.Success {
// 		fmt.Println(string(res.JsonData))
// 		t.Log("Test_GetConversationsForLocation PASSED")
// 	} else {
// 		t.Fail()
// 		t.Log("Test_GetConversationsForLocation FAILED")
// 	}
// }

// func Test_GetUserMessages(t *testing.T) {
// 	response := GetUserMessages("51b8e5b62ffc2c5db5e9b213")
// 	fmt.Println(response)
// }

func Test_GetUserMessagesList(t *testing.T) {
	_, err := GetUserMessagesList("51b8e5b62ffc2c5db5e9b213")
	// GetUserMessageList("51b8e5b62ffc2c5db5e9b213")
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	} else {
		fmt.Println("Test_GetUserMessagesList PASSED")
		t.Log("Test_GetUserMessagesList PASSED")

	}
}
