package message

import (
	"fmt"
	// "labix.org/v2/mgo/bson"
	"testing"
)

var Msg = "{\"msg_text\":\"this is a text message\",\"user_id\":\"1234567\"}"

func Test_SaveMessage(t *testing.T) {
	var m Message
	m.JsonToMsg(Msg)
	res := m.SaveMessage("51f0f90fd323a50002000001")
	if res.Success {
		fmt.Println(string(res.JsonData))
		t.Log("Test_SaveMessage PASSED")
	} else {
		t.Fail()
		t.Log("Test_SaveMessage FAILED")
	}
}

func Test_GetMessages(t *testing.T) {
	res := GetMessages("51f0f90fd323a50002000001")
	if res.Success {
		fmt.Println(string(res.JsonData))
		t.Log("Test_GetConversationsForLocation PASSED")
	} else {
		t.Fail()
		t.Log("Test_GetConversationsForLocation FAILED")
	}
}
