package usermetadata

import (
	// "database/sql"
	"encoding/json"
	// "fmt"
	// _ "github.com/go-sql-driver/mysql"
	Connection "github.com/jmadan/go-msgstory/connection"
	RD "github.com/jmadan/go-msgstory/util"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"os"
	"strings"
	"time"
)

type UserMetaData struct {
	Id          bson.ObjectId `json:"_id" bson:"_id"`
	UserId      string        `json:"userid" bson:"userid"`
	PhoneNumber string        `json:"phone" bson:"phone"`
	Age         int           `json:"age" bson:"age"`
	Created_on  time.Time     `json:"created_on" bson:"created_on"`
}

func (u *UserMetaData) SetUserid(id string) {
	u.UserId = id
}

func (u *UserMetaData) SaveUserMetaData() RD.ReturnData {
	returnData := RD.ReturnData{}
	dbSession := Connection.GetDBSession()
	dbSession.SetMode(mgo.Monotonic, true)
	dataBase := strings.SplitAfter(os.Getenv("MONGOHQ_URL"), "/")
	c := dbSession.DB(dataBase[3]).C("jove")

	u.Id = bson.NewObjectId()
	u.Created_on = time.Now()

	err := c.Insert(u)
	if err != nil {
		log.Print(err.Error())
		returnData.ErrorMsg = err.Error()
		returnData.Success = false
		returnData.Status = "422"
	} else {
		returnData.Success = true
		jsonData, _ := json.Marshal(&u)
		returnData.JsonData = jsonData
		returnData.Status = "201"
	}

	return returnData
}

func GetUserById(userId string) string {
	var response string
	dbSession := Connection.GetDBSession()
	dbSession.SetMode(mgo.Monotonic, true)
	dataBase := strings.SplitAfter(os.Getenv("MONGOHQ_URL"), "/")
	c := dbSession.DB(dataBase[3]).C("jove")

	result := UserMetaData{}
	err := c.Find(bson.M{"_id": bson.ObjectIdHex(userId)}).One(&result)
	if err != nil {
		log.Println(err)
		response = err.Error()
		// returnData.ErrorMsg = err.Error()
		// returnData.Status = "400"
		// returnData.Success = false
	} else {
		log.Println(result)
		// returnData.ErrorMsg = "All is well"
		// returnData.Status = "200"
		// returnData.Success = true
		response, _ := json.Marshal(result)
		log.Println(response)
		// returnData.JsonData = jsonRes
	}
	return response
}
