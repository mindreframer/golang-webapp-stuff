package dialogue

import (
	"encoding/json"
	// SJ "github.com/bitly/go-simplejson"
	"fmt"
	// Group "github.com/jmadan/go-msgstory/circle"
	Connection "github.com/jmadan/go-msgstory/connection"
	Location "github.com/jmadan/go-msgstory/geolocation"
	// Msg "github.com/jmadan/go-msgstory/message"
	// User "github.com/jmadan/go-msgstory/user"
	RD "github.com/jmadan/go-msgstory/util"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"os"
	"strings"
	"time"
)

type Dialogue struct {
	Id    bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	Title string        `json:"title" bson:"title"`
	// Messages        []Msg.Message        `json:"messages" bson:"messages"`
	Venue Location.GeoLocation `json:"venue" bson:"venue"`
	// DialogueStarter User                 `json:"user" bson:"user"`
	IsApproved bool      `json:"is_approved" bson:"is_approved"`
	CreatedOn  time.Time `json:"created_on" bson:"created_on,omitempty"`
}

func (D *Dialogue) DialogueToJSON() string {
	dialogueJson, err := json.Marshal(D)
	if err != nil {
		log.Fatal(err.Error())
	}
	return string(dialogueJson)
}

func (D *Dialogue) JsonToDialogue(dialogue string) {
	err := json.Unmarshal([]byte(dialogue), &D)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func (D *Dialogue) CreateDialogue() (RD.ReturnData, Dialogue) {
	returnData := RD.ReturnData{}
	dbSession := Connection.GetDBSession()
	dbSession.SetMode(mgo.Monotonic, true)
	dataBase := strings.SplitAfter(os.Getenv("MONGOHQ_URL"), "/")
	c := dbSession.DB(dataBase[3]).C("Dialogue")
	D.Id = bson.NewObjectId()
	D.CreatedOn = time.Now()
	D.IsApproved = true

	err := c.Insert(&D)
	if err != nil {
		log.Print(err.Error())
		returnData.ErrorMsg = err.Error()
		returnData.Success = false
		returnData.Status = "422"
	} else {
		returnData.Success = true
		returnData.JsonData = []byte(D.DialogueToJSON())
		returnData.Status = "201"
	}

	return returnData, *D
}

func GetDialoguesForLocation(locationId string) RD.ReturnData {
	returnData := RD.ReturnData{}
	dbSession := Connection.GetDBSession()
	dbSession.SetMode(mgo.Monotonic, true)
	dataBase := strings.SplitAfter(os.Getenv("MONGOHQ_URL"), "/")
	c := dbSession.DB(dataBase[3]).C("Dialogue")

	res := []Dialogue{}
	err := c.Find(bson.M{"venue.fourid": locationId}).All(&res)
	if err != nil {
		log.Println("Found Nothing Or Something went wrong fetching the Dialogue")
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

func GetDialogue(DialogueId string) (returnData RD.ReturnData) {
	fmt.Println(DialogueId)
	dbSession := Connection.GetDBSession()
	// dbSession.SetMode(mgo.Monotonic, true)
	defer dbSession.Close()
	log.Println(os.Getenv("MONGOHQ_URL"))
	dataBase := strings.SplitAfter(os.Getenv("MONGOHQ_URL"), "/")
	c := dbSession.DB(dataBase[3]).C("Dialogue")
	res := Dialogue{}
	// err := c.FindId(bson.ObjectIdHex(DialogueId)).One(&res)
	err := c.Find(bson.M{"_id": bson.ObjectIdHex(DialogueId)}).One(&res)
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

func (D *Dialogue) DeleteDialogue() RD.ReturnData {
	returnData := RD.ReturnData{}
	dbSession := Connection.GetDBSession()
	dbSession.SetMode(mgo.Monotonic, true)
	dataBase := strings.SplitAfter(os.Getenv("MONGOHQ_URL"), "/")
	c := dbSession.DB(dataBase[3]).C("Dialogue")

	// err := c.Remove(bson.ObjectIdHex(DialogueId))
	err := c.Remove(D.Id)

	if err != nil {
		log.Println("Found Nothing. Something went wrong fetching the Dialogue document")
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
