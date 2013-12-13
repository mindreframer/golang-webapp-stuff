package conversation

import (
	"fmt"
	Circle "github.com/jmadan/go-msgstory/circle"
	Location "github.com/jmadan/go-msgstory/geolocation"
	Msg "github.com/jmadan/go-msgstory/message"
	User "github.com/jmadan/go-msgstory/user"
	"labix.org/v2/mgo/bson"
	"testing"
)

var M1 = Msg.Message{
	MsgText: "Hola! how is everyone doing?",
	UserId:  "SomeOwnerID",
}

var M2 = Msg.Message{
	MsgText: "Hola! Not bad mate",
	UserId:  "user id",
}

var Cir = Circle.Circle{
	Name: "Tapori",
}
var CirPub = Circle.Circle{
	Name: "Public",
}

var Loc = Location.GeoLocation{
	FourID:   "4bce6383ef109521bd238486",
	Name:     "City of Manchester",
	Contact:  "",
	Address:  "",
	Lat:      42,
	Lng:      -71,
	Distance: 0,
	Postcode: "03104",
	City:     "Manchester",
	State:    "NH",
	Country:  "United States",
}

var Test_User = User.User{
	Id:          bson.NewObjectId(),
	UserId:      100,
	Name:        "Test User",
	Email:       "test@test.com",
	Handle:      "testuser",
	PhoneNumber: "9008307311",
}

// var Conv = Conversation{
// 	Title:     "This is 2nd test conversation!",
// 	Messages:  []Msg.Message{},
// 	Venue:     Loc,
// 	Circles:   []Circle.Circle{Cir, CirPub},
// 	ConvOwner: "001",
// }

var Conv = Conversation{
	Title:    "This is 2nd test conversation!",
	Messages: []Msg.Message{M1, M2},
	Venue:    Loc,
	Circles:  []Circle.Circle{Cir, CirPub},
	User:     Test_User,
}

// func Test_CreateConversation(t *testing.T) {
// 	response, saved_conversation := Conv.CreateConversation()
// 	if response.Success {
// 		fmt.Println(string(response.JsonData))
// 		fmt.Println(string(saved_conversation.ConversationToJSON()))
// 		t.Log("Test_CreateConversation PASSED")
// 	} else {
// 		t.Fail()
// 		t.Log("Test_CreateConversation FAILED")
// 		fmt.Println(response.ErrorMsg)
// 	}
// }

// func Test_GetAllConversations(t *testing.T) {
// 	res := GetAllConversations()
// 	// res := GetConversationsForLocation("4bce6383ef109521bd238486")
// 	if res.Success {
// 		fmt.Print("success")
// 		t.Log("Test_GetConversationsForLocation PASSED")
// 	} else {
// 		t.Fail()
// 		t.Log("Test_GetConversationsForLocation FAILED")
// 	}
// }

func Test_GetConversationsForLocation(t *testing.T) {
	res, err := GetConversationsForLocation("4bce6383ef109521bd238486")
	// res := GetConversationsForLocation("4bce6383ef109521bd238486")
	if err == nil {
		t.Log("Test_GetConversationsForLocation PASSED")
	} else {
		fmt.Println(res)
		t.Fail()
		t.Log("Test_GetConversationsForLocation FAILED")
	}
}

// func Test_DeleteConversation(t *testing.T) {
// 	res := DeleteConversation("51bc529e2ffc2c5db5e9b215")
// 	if res.Success {
// 		t.Log("Test_DeleteConversation PASSED")
// 	} else {
// 		fmt.Println(res.ErrorMsg)
// 		t.Fail()
// 		t.Log("Test_DeleteConversation FAILED")
// 	}
// }

// func Test_GetConversation(t *testing.T) {
// 	conId := "51ec3bee5fb5e52f2860a04c"
// 	res := GetConversation(conId)
// 	if res.Success {
// 		t.Log("Test_GetConversation PASSED")
// 		fmt.Println(string(res.JsonData))
// 	} else {
// 		fmt.Println(res.ErrorMsg)
// 		t.Fail()
// 		t.Log("Test_GetConversation FAILED")
// 	}
// }

// func Test_SaveMessage(t *testing.T) {
// 	conId := "51bc529e2ffc2c5db5e9b215"
// 	json_msg, err := json.Marshal(M2)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		t.Fail()
// 		t.Log("Test_SaveMessage FAILED")
// 	}

// 	res := SaveMessage(conId, string(json_msg))
// 	if res.Success {
// 		t.Log("Test_SaveMessage PASSED")
// 		fmt.Println(string(res.JsonData))
// 	} else {
// 		fmt.Println(res.ErrorMsg)
// 		t.Fail()
// 		t.Log("Test_SaveMessage FAILED")
// 	}
// }
