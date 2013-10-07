package conversation

import (
	"encoding/json"
	// SJ "github.com/bitly/go-simplejson"
	"fmt"
	Group "github.com/jmadan/go-msgstory/circle"
	Connection "github.com/jmadan/go-msgstory/connection"
	Location "github.com/jmadan/go-msgstory/geolocation"
	Msg "github.com/jmadan/go-msgstory/message"
	User "github.com/jmadan/go-msgstory/user"
	RD "github.com/jmadan/go-msgstory/util"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"os"
	"strings"
	"time"
)

// type Conversation struct {
// 	Id          bson.ObjectId        `json:"_id" bson:"_id,omitempty"`
// 	Title       string               `json:"title" bson:"title"`
// 	Messages    []Msg.Message        `json:"messages" bson:"messages"`
// 	Venue       Location.GeoLocation `json:"venue" bson:"venue"`
// 	Circles     []Group.Circle       `json:"circles" bson:"circles"`
// 	ConvOwner   string               `json:"creator" bson:"creator"`
// 	Is_Approved bool                 `json:"is_approved" bson:"is_approved"`
// 	Created_On  time.Time            `json:"created_on" bson:"created_on,omitempty"`
// }

type Conversation struct {
	Id          bson.ObjectId        `json:"_id" bson:"_id,omitempty"`
	Title       string               `json:"title" bson:"title"`
	Messages    []Msg.Message        `json:"messages" bson:"messages"`
	Venue       Location.GeoLocation `json:"venue" bson:"venue"`
	Circles     []Group.Circle       `json:"circles" bson:"circles"`
	ConvOwner   User.User            `json:"user" bson:"user"`
	Is_Approved bool                 `json:"is_approved" bson:"is_approved"`
	Created_On  time.Time            `json:"created_on" bson:"created_on,omitempty"`
}

func (C *Conversation) ConversationToJSON() string {
	cjson, err := json.Marshal(C)
	if err != nil {
		log.Fatal(err.Error())
	}

	return string(cjson)
}

func (C *Conversation) JsonToConversation(conv string) {
	err := json.Unmarshal([]byte(conv), &C)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func (conv *Conversation) CreateConversation() (RD.ReturnData, Conversation) {
	returnData := RD.ReturnData{}
	dbSession := Connection.GetDBSession()
	dbSession.SetMode(mgo.Monotonic, true)
	dataBase := strings.SplitAfter(os.Getenv("MONGOHQ_URL"), "/")
	c := dbSession.DB(dataBase[3]).C("conversation")
	conv.Id = bson.NewObjectId()
	conv.Created_On = time.Now()
	conv.Is_Approved = true

	err := c.Insert(&conv)
	if err != nil {
		log.Print(err.Error())
		returnData.ErrorMsg = err.Error()
		returnData.Success = false
		returnData.Status = "422"
	} else {

		returnData.Success = true
		jsonData, _ := json.Marshal(&conv)
		returnData.JsonData = jsonData
		returnData.Status = "201"
	}

	return returnData, *conv
}

func GetConversationsForLocation(locationId string) RD.ReturnData {
	returnData := RD.ReturnData{}
	dbSession := Connection.GetDBSession()
	dbSession.SetMode(mgo.Monotonic, true)
	dataBase := strings.SplitAfter(os.Getenv("MONGOHQ_URL"), "/")
	c := dbSession.DB(dataBase[3]).C("conversation")

	res := []Conversation{}
	err := c.Find(bson.M{"venue.fourid": locationId}).All(&res)
	if err != nil {
		log.Println("Found Nothing Or Something went wrong fetching the Conversation document")
		returnData.ErrorMsg = err.Error()
		returnData.Status = "400"
		returnData.Success = false
	} else {
		log.Println(res)
		returnData.ErrorMsg = "All is well"
		returnData.Status = "200"
		returnData.Success = true
		jsonRes, _ := json.Marshal(res)
		returnData.JsonData = jsonRes
		log.Println(string(jsonRes))
	}
	return returnData
}

func GetConversation(conversationId string) (returnData RD.ReturnData) {
	fmt.Println(conversationId)
	dbSession := Connection.GetDBSession()
	// dbSession.SetMode(mgo.Monotonic, true)
	defer dbSession.Close()
	log.Println(os.Getenv("MONGOHQ_URL"))
	dataBase := strings.SplitAfter(os.Getenv("MONGOHQ_URL"), "/")
	c := dbSession.DB(dataBase[3]).C("conversation")
	res := Conversation{}
	// err := c.FindId(bson.ObjectIdHex(conversationId)).One(&res)
	err := c.Find(bson.M{"_id": bson.ObjectIdHex(conversationId)}).One(&res)
	if err != nil {
		log.Println(err)
		returnData.ErrorMsg = err.Error()
		returnData.Status = "400"
		returnData.Success = false
	} else {
		log.Println(res)
		returnData.ErrorMsg = "All is well"
		returnData.Status = "200"
		returnData.Success = true
		jsonRes, _ := json.Marshal(res)
		returnData.JsonData = jsonRes
	}
	return
}

func DeleteConversation(conversationId string) RD.ReturnData {
	returnData := RD.ReturnData{}
	dbSession := Connection.GetDBSession()
	dbSession.SetMode(mgo.Monotonic, true)
	dataBase := strings.SplitAfter(os.Getenv("MONGOHQ_URL"), "/")
	c := dbSession.DB(dataBase[3]).C("conversation")

	err := c.Remove(bson.ObjectIdHex(conversationId))
	// err := c.Find(bson.M{"venue.fourid": locationId}).One(&res)
	if err != nil {
		log.Println("Found Nothing. Something went wrong fetching the Conversation document")
		log.Println(err)
		returnData.ErrorMsg = err.Error()
		returnData.Status = "400"
		returnData.Success = false
	} else {
		returnData.ErrorMsg = "All is well"
		returnData.Status = "200"
		returnData.Success = true
		returnData.JsonData = nil
	}
	return returnData
}

func GetAllConversations() RD.ReturnData {
	returnData := RD.ReturnData{}
	dbSession := Connection.GetDBSession()
	dbSession.SetMode(mgo.Monotonic, true)
	dataBase := strings.SplitAfter(os.Getenv("MONGOHQ_URL"), "/")
	c := dbSession.DB(dataBase[3]).C("conversation")

	res := []Conversation{}
	err := c.Find(bson.M{}).All(&res)
	if err != nil {
		log.Println("Found Nothing Or Something went wrong fetching the Conversation document")
		returnData.ErrorMsg = err.Error()
		returnData.Status = "400"
		returnData.Success = false
	} else {
		log.Println(res)
		returnData.ErrorMsg = "All is well"
		returnData.Status = "200"
		returnData.Success = true
		jsonRes, _ := json.Marshal(res)
		returnData.JsonData = jsonRes
		log.Println(string(jsonRes))
	}
	return returnData
}
