package message

import (
	"encoding/json"
	Connection "github.com/jmadan/go-msgstory/connection"
	// User "github.com/jmadan/go-msgstory/user"
	RD "github.com/jmadan/go-msgstory/util"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"os"
	"strings"
	"time"
)

type Message struct {
	Id         bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	MsgText    string        `json:"msg_text" bson:"msg_text"`
	UserId     string        `json:"user_id" bson:"user_id"`
	UserHandle string        `json:"user_handle" bson:"user_handle"`
	CreatedOn  time.Time     `json:"created_on" bson:"created_on"`
	// User      User.User     `json:"user" bson:"user"`
}

type Messages struct {
	Messages []Message `json:"messages" bson:"messages"`
}

func (M *Message) MsgToJSON() string {
	mjson, err := json.Marshal(M)
	if err != nil {
		log.Fatal(err.Error())
	}

	return string(mjson)
}

func (M *Message) JsonToMsg(msgtext string) {
	err := json.Unmarshal([]byte(msgtext), &M)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func (msg *Message) SaveMessage(conversationId string) RD.ReturnData {
	returnData := RD.ReturnData{}
	dbSession := Connection.GetDBSession()
	dbSession.SetMode(mgo.Monotonic, true)
	dataBase := strings.SplitAfter(os.Getenv("MONGOHQ_URL"), "/")
	c := dbSession.DB(dataBase[3]).C("conversation")
	msg.CreatedOn = time.Now()

	err := c.Update(bson.M{"_id": bson.ObjectIdHex(conversationId)}, bson.M{
		"$push": bson.M{"messages": bson.M{
			"_id":         bson.NewObjectId(),
			"msg_text":    msg.MsgText,
			"user_id":     msg.UserId,
			"user_handle": msg.UserHandle,
			"created_on":  msg.CreatedOn,
		}}})

	if err != nil {
		log.Println(err.Error())
		returnData.ErrorMsg = err.Error()
		returnData.Success = false
		returnData.Status = "422"
	} else {
		jsonData := []byte("{}")
		returnData.Success = true
		returnData.JsonData = jsonData
		returnData.Status = "201"
	}
	return returnData
}

func GetMessages(conversationId string) RD.ReturnData {
	returnData := RD.ReturnData{}
	dbSession := Connection.GetDBSession()
	dbSession.SetMode(mgo.Monotonic, true)
	dataBase := strings.SplitAfter(os.Getenv("MONGOHQ_URL"), "/")
	c := dbSession.DB(dataBase[3]).C("conversation")

	Msgs := []Message{}
	m := Messages{}
	err := c.Find(bson.M{"_id": bson.ObjectIdHex(conversationId)}).Select(bson.M{"messages": 1}).One(&m)
	if err != nil {
		log.Println(err.Error())
		returnData.ErrorMsg = err.Error()
		returnData.Success = false
		returnData.Status = "422"
	} else {
		log.Println(Msgs)
		jsonData, _ := json.Marshal(&m)
		returnData.Success = true
		returnData.JsonData = jsonData
		returnData.Status = "201"
	}
	return returnData
}

func GetUserMessages(userId string) RD.ReturnData {
	returnData := RD.ReturnData{}
	dbSession := Connection.GetDBSession()
	dbSession.SetMode(mgo.Monotonic, true)
	dataBase := strings.SplitAfter(os.Getenv("MONGOHQ_URL"), "/")
	c := dbSession.DB(dataBase[3]).C("conversation")

	Msgs := []Message{}
	m := Messages{}
	err := c.Find(bson.M{"messages.user_id": bson.ObjectIdHex(userId)}).Select(bson.M{"messages": 1}).One(&m)
	if err != nil {
		log.Println(err.Error())
		returnData.ErrorMsg = err.Error()
		returnData.Success = false
		returnData.Status = "422"
	} else {
		log.Println(Msgs)
		jsonData, _ := json.Marshal(&m)
		returnData.Success = true
		returnData.JsonData = jsonData
		returnData.Status = "201"
	}
	return returnData
}

func GetUserMessagesList(userId string) (string, error) {
	var response []byte
	dbSession := Connection.GetDBSession()
	dbSession.SetMode(mgo.Monotonic, true)
	dataBase := strings.SplitAfter(os.Getenv("MONGOHQ_URL"), "/")
	c := dbSession.DB(dataBase[3]).C("conversation")

	Msgs := []Messages{}
	err := c.Find(bson.M{"messages.user_id": userId}).Select(bson.M{"messages": 1}).All(&Msgs)
	if err != nil {
	} else {
		response, err = json.Marshal(Msgs)
	}
	log.Println(string(response))

	return string(response), err
}
