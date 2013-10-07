package user

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	Connection "github.com/jmadan/go-msgstory/connection"
	Message "github.com/jmadan/go-msgstory/message"
	RD "github.com/jmadan/go-msgstory/util"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"os"
	"strings"
	"time"
)

type User struct {
	Id          bson.ObjectId `json:"_id" bson:"_id"`
	UserId      int           `json:"userid" bson:"userid"`
	Name        string        `json:"name" bson:"name"`
	Email       string        `json:"email" bson:"email"`
	Handle      string        `json:"handle" bson:"handle"`
	PhoneNumber string        `json:"phone" bson:"phone"`
	Relations   rels          `json:"relations" bson:"relations"`
	Created_on  time.Time     `json:"created_on" bson:"created_on"`
}

type rels struct {
	Messages []Message.Message `json:"messages" bson:"messages"`
}

func (u *User) SetEmail(email string) {
	u.Email = email
}

func (u *User) SetUserid(id int) {
	u.UserId = id
}

func (u *User) GetName() string {
	return u.Name
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) GetHandle() string {
	return u.Handle
}

func (u *User) GetMessages() string {
	str, err := json.Marshal(u.Relations.Messages)
	if err != nil {
		fmt.Println("what the fuck!")
	}
	return string(str)
}

func (u *User) GetUser() string {
	dbSession := Connection.GetDBSession()
	dbSession.SetMode(mgo.Monotonic, true)
	dataBase := strings.SplitAfter(os.Getenv("MONGOHQ_URL"), "/")
	c := dbSession.DB(dataBase[3]).C("jove")

	result := User{}
	err := c.Find(bson.M{"email": u.Email, "userid": u.UserId}).One(&result)
	if err != nil {
		log.Println(err.Error())
	}
	js, _ := json.Marshal(result)
	return string(js)
}

func GetAll() string {
	dbSession := Connection.GetDBSession()
	dbSession.SetMode(mgo.Monotonic, true)
	dataBase := strings.SplitAfter(os.Getenv("MONGOHQ_URL"), "/")
	c := dbSession.DB(dataBase[3]).C("jove")

	result := []User{}
	err := c.Find(nil).Limit(10).All(&result)
	if err != nil {
		panic(err.Error())
	}

	return "hello"
}

func GetByEmailAndUserId(email string, user_id int) (User, error) {
	dbSession := Connection.GetDBSession()
	dbSession.SetMode(mgo.Monotonic, true)
	dataBase := strings.SplitAfter(os.Getenv("MONGOHQ_URL"), "/")
	c := dbSession.DB(dataBase[3]).C("jove")

	result := User{}
	err := c.Find(bson.M{"email": email, "userid": user_id}).One(&result)
	if err != nil {
		log.Fatal(err)
	}
	return result, err
}

func (u *User) GetByHandle() User {
	dbSession := Connection.GetDBSession()
	dbSession.SetMode(mgo.Monotonic, true)
	dataBase := strings.SplitAfter(os.Getenv("MONGOHQ_URL"), "/")
	c := dbSession.DB(dataBase[3]).C("jove")

	result := User{}
	err := c.Find(bson.M{"handle": u.Handle}).One(&result)
	if err != nil {
		log.Fatal(err)
	}

	return result
}

func (u *User) CreateUser() RD.ReturnData {
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

func CreateUserLogin(useremail, password string) string {
	dburl := os.Getenv("DATABASE_URL")
	db, err := sql.Open("mysql", dburl[8:])
	if err != nil {
		log.Fatal("Phat Gayee : " + err.Error())
	}
	defer db.Close()

	stmtIns, err := db.Prepare("INSERT INTO USERS (USEREMAIL,PASSWORD) VALUES (?,?)")
	if err != nil {
		log.Fatal("stmtError :" + err.Error())
	}
	defer stmtIns.Close()

	// err = stmtOut.QueryRow(useremail, userpassword).Scan(&authorize.user_id, &authorize.email)
	_, err = stmtIns.Exec(useremail, password)
	if err != nil {
		log.Print("stmtExecution: " + err.Error())
	}

	var userid string
	stmtOut, err := db.Prepare("SELECT USER_ID FROM USERS WHERE USEREMAIL=?")
	if err != nil {
		log.Println("stmtError: " + err.Error())
	}

	err = stmtOut.QueryRow(useremail).Scan(&userid)
	if err != nil {
		log.Println(err.Error())
	}

	return userid
}

func GetUserByEmail(user_email string) string {
	dburl := os.Getenv("DATABASE_URL")
	db, err := sql.Open("mysql", dburl[8:])
	if err != nil {
		log.Fatal("Phat Gayee : " + err.Error())
	}
	defer db.Close()

	stmtOut, err := db.Prepare("SELECT USER_ID FROM USERS WHERE USEREMAIL = ?")
	if err != nil {
		log.Fatal("stmtError :" + err.Error())
	}
	defer stmtOut.Close()

	var uid string

	// err = stmtOut.QueryRow(useremail, userpassword).Scan(&authorize.user_id, &authorize.email)
	err = stmtOut.QueryRow(user_email).Scan(&uid)

	if err != nil {
		log.Print("stmtExecution: " + err.Error())
		return err.Error()
	} else {
		return uid
	}

	return uid
}
