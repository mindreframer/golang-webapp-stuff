package main

import (
	"code.google.com/p/gorest"
	"encoding/json"
	"fmt"
	Authenticate "github.com/jmadan/go-msgstory/authenticate"
	Circle "github.com/jmadan/go-msgstory/circle"
	Conversation "github.com/jmadan/go-msgstory/conversation"
	Glocation "github.com/jmadan/go-msgstory/geolocation"
	Msg "github.com/jmadan/go-msgstory/message"
	Register "github.com/jmadan/go-msgstory/register"
	User "github.com/jmadan/go-msgstory/user"
	ReturnData "github.com/jmadan/go-msgstory/util"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type AppService struct {
	gorest.RestService `root:"/api/" consumes:"application/json" produces:"application/json"`
	getApp             gorest.EndPoint `method:"GET" path:"/" output:"string"`
}

type UserService struct {
	gorest.RestService `root:"/api/user/" consumes:"application/json" produces:"application/json"`

	getUser      gorest.EndPoint `method:"GET" path:"/{userid:string}" output:"string"`
	getAll       gorest.EndPoint `method:"GET" path:"/" output:"string"`
	registerUser gorest.EndPoint `method:"POST" path:"/" postdata:"string"`
	// createUser         gorest.EndPoint `method:"GET" path:"/new/{uemail:string}/{pass:string}" output:"string"`
}

type ConversationService struct {
	gorest.RestService `root:"/api/conversation/" consumes:"application/json" produces:"application/json"`

	createConversation          gorest.EndPoint `method:"POST" path:"/" postdata:"string"`
	getConversationsForLocation gorest.EndPoint `method:"GET" path:"/all/{locationId:string}" output:"string"`
	getConversation             gorest.EndPoint `method:"GET" path:"/{convoId:string}" output:"string"`
	deleteConversation          gorest.EndPoint `method:"DELETE" path:"/{convoId:string}/"`
	deleteMessage               gorest.EndPoint `method:"DELETE" path:"/{convoId:string}/messages/{msgId:string}"`
}

type MsgService struct {
	gorest.RestService `root:"/api/message/" consumes:"application/json" produces:"application/json"`

	saveMessage gorest.EndPoint `method:"POST" path:"/conversation/{convoId:string}/" postdata:"string"`
	getMessage  gorest.EndPoint `method:"GET" path:"/{msgId:string}" output:"string"`
	getMessages gorest.EndPoint `method:"GET" path:"/conversation/{convoId:string}" output:"string"`
}

type AuthenticateService struct {
	gorest.RestService `root:"/api/auth/" consumes:"application/json" produces:"application/json"`

	loginUser gorest.EndPoint `method:"POST" path:"/login/" postdata:"string"`
}

type CircleService struct {
	gorest.RestService `root:"/api/circle/" consumes:"application/json" produces:"application/json"`

	createCircle gorest.EndPoint `method:"POST" path:"/new/" postdata:"string"`
	getCircles   gorest.EndPoint `method:"GET" path:"/circles/" output:"string"`
}

type LocationService struct {
	gorest.RestService `root:"/api/location/" consumes:"application/json" produces:"application/json"`

	getLocations            gorest.EndPoint `method:"GET" path:"/near/{place:string}" output:"string"`
	getLocationsWithLatLng  gorest.EndPoint `method:"GET" path:"/coordinates/{lat:string}/{lng:string}" output:"string"`
	getLocationDetailWithId gorest.EndPoint `method:"GET" path:"/{venueId:string}" output:"string"`
}

//*************Conversation Service Methods ***********
func (serv ConversationService) CreateConversation(posted string) {
	var returnData ReturnData.ReturnData
	var formData []string
	formData = strings.Split(posted, "=")
	conv := Conversation.Conversation{}
	err := json.Unmarshal([]byte(formData[1]), &conv)
	if err != nil {
		log.Println("conversation marshelling error>>>>>" + err.Error())
		serv.ResponseBuilder().SetResponseCode(400).WriteAndOveride([]byte(err.Error()))
		return
	} else {
		returnData, _ = conv.CreateConversation()
	}

	if returnData.Success {
		serv.ResponseBuilder().SetResponseCode(201).Write([]byte(returnData.ToString()))
	} else {
		serv.ResponseBuilder().SetResponseCode(400).WriteAndOveride([]byte(returnData.ToString()))
	}

}

func (serv ConversationService) GetConversationsForLocation(locationId string) string {
	var data ReturnData.ReturnData
	data = Conversation.GetConversationsForLocation(locationId)
	if data.Success {
		serv.ResponseBuilder().SetResponseCode(200)
	} else {
		serv.ResponseBuilder().SetResponseCode(400).WriteAndOveride([]byte(data.ToString()))
	}
	return string(data.ToString())
}

func (serv ConversationService) GetConversation(convoId string) string {
	var data ReturnData.ReturnData
	data = Conversation.GetConversation(convoId)
	if data.Success {
		serv.ResponseBuilder().SetResponseCode(200)
	} else {
		serv.ResponseBuilder().SetResponseCode(400).WriteAndOveride([]byte(data.ToString()))
	}
	return string(data.ToString())
}

func (serv ConversationService) DeleteConversation(convoId string) {}

func (serv ConversationService) DeleteMessage(convoId, msgId string) {}

// ************Location Service Methods ***********
func (serv LocationService) GetLocations(place string) string {
	fmt.Println(place)
	resp := Glocation.GetVenues("Chelsea,London")
	serv.ResponseBuilder().SetResponseCode(200)
	return resp
}

func (serv LocationService) GetLocationsWithLatLng(lat, lng string) string {
	str := Glocation.GetVenuesWithLatitudeAndLongitude(lat, lng)
	serv.ResponseBuilder().SetResponseCode(200)
	return str
}

func (serv LocationService) GetLocationDetailWithId(venueId string) string {
	str := Glocation.GetVenueWithId(venueId)
	serv.ResponseBuilder().SetResponseCode(200)
	return str
}

//*************Circle Service Methods ***************
func (serv CircleService) CreateCircle(posted string) {
	// var str []string
	// str = strings.Split(posted, "=")
	// msg, err := Circle.CreateCircle(str[1], "", "", nil)
	// if err != nil {
	// 	log.Println(err)
	// } else {
	// 	fmt.Println(msg)
	// 	serv.ResponseBuilder().SetResponseCode(200)
	// }
}

func (serv CircleService) GetCircles() string {
	return Circle.GetUserCircles("")[0]
}

//*************Authentication Service Methods ***************

func (serv AuthenticateService) LoginUser(posted string) {
	fmt.Println(posted)
	var str []string
	str = strings.Split(posted, "=")
	auth := Authenticate.Authenticate{}
	user := User.User{}
	err := json.Unmarshal([]byte(str[1]), &auth)
	if err != nil {
		log.Println(err.Error())
		serv.ResponseBuilder().SetResponseCode(404).WriteAndOveride(nil)
		return
	} else {
		auth.Authorize()
		user.SetEmail(auth.Email)
		user.SetUserid(auth.User_id)
		serv.ResponseBuilder().SetResponseCode(200).Write([]byte(user.GetUser()))
		return
	}
}

//*************Message Service Methods ***************
func (serv MsgService) GetMessages(convoId string) string {
	var data ReturnData.ReturnData
	data = Msg.GetMessages(convoId)
	if data.Success {
		serv.ResponseBuilder().SetResponseCode(200)
	} else {
		serv.ResponseBuilder().SetResponseCode(400).WriteAndOveride([]byte(data.ToString()))
	}
	return string(data.ToString())
}

func (serv MsgService) GetMessage(msgId string) string {
	var data ReturnData.ReturnData
	data.Success = true
	data.JsonData = []byte("Get Message call")
	data.Status = "200"
	if data.Success {
		serv.ResponseBuilder().SetResponseCode(200)
	} else {
		serv.ResponseBuilder().SetResponseCode(400).WriteAndOveride([]byte(data.ToString()))
	}
	return string(data.ToString())
}

func (serv MsgService) SaveMessage(posted, convoId string) {
	var data ReturnData.ReturnData
	var str []string
	str = strings.Split(posted, "=")
	msg := Msg.Message{}
	err := json.Unmarshal([]byte(str[1]), &msg)
	if err != nil {
		log.Println(err.Error())
		serv.ResponseBuilder().SetResponseCode(400).WriteAndOveride(nil)
		return
	} else {
		data = msg.SaveMessage(convoId)
	}
	if data.Success {
		serv.ResponseBuilder().SetResponseCode(201).Write([]byte(data.ToString()))
	} else {
		serv.ResponseBuilder().SetResponseCode(400).WriteAndOveride([]byte(data.ToString()))
	}

}

//*************User Service Methods ***************
func (serv UserService) RegisterUser(posted string) {

	type newUser struct {
		Name     string `json:"name" bson:"name"`
		Email    string `json:"email" bson:"email"`
		Handle   string `json:"handle" bson:"handle"`
		Password string `json:"password" bson:"password"`
	}

	var data ReturnData.ReturnData
	var formData []string
	formData = strings.Split(posted, "=")
	user := User.User{}
	tempUser := newUser{}
	err := json.Unmarshal([]byte(formData[1]), &tempUser)

	if err != nil {
		log.Println(err.Error())
		serv.ResponseBuilder().SetResponseCode(400).WriteAndOveride(nil)
		return
	} else {
		user_id := User.CreateUserLogin(tempUser.Email, tempUser.Password)
		user.UserId, _ = strconv.Atoi(user_id)
		user.Name = tempUser.Name
		user.Email = tempUser.Email
		user.Handle = tempUser.Handle
		data = user.CreateUser()
	}
	if data.Success {
		serv.ResponseBuilder().SetResponseCode(201).Write([]byte(data.ToString()))
	} else {
		serv.ResponseBuilder().SetResponseCode(400).WriteAndOveride([]byte(data.ToString()))
	}
}

func (serv UserService) CreateUser(uemail, pass string) string {
	Register.Register(uemail, pass)
	return "Executed!!!"
}

func (serv UserService) GetUser(userid string) string {
	// user := User.User{}
	// per := "{User:[" + User.User.GetUser() + "]}"
	// serv.ResponseBuilder().SetResponseCode(404).Overide(true)
	return "Some User"
}

func (serv UserService) GetAll() string {
	fmt.Print("I am here")
	per := "User:[" + User.GetAll() + "]"
	return per
}

//*************App Service Methods ***************
func (serv AppService) GetApp() string {
	m := "{\"Message\": \"Welcome to Mesiji\"}"
	return m
}

func getData(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.FormValue("inputEmail"))
}

func main() {
	log.Println(os.Getenv("PORT"))
	gorest.RegisterService(new(AppService))
	gorest.RegisterService(new(UserService))
	gorest.RegisterService(new(ConversationService))
	gorest.RegisterService(new(MsgService))
	gorest.RegisterService(new(AuthenticateService))
	gorest.RegisterService(new(CircleService))
	gorest.RegisterService(new(LocationService))
	http.Handle("/", gorest.Handle())
	// http.HandleFunc("/tempurl", getData)
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	// User.GetAll()
}

func getResponse() string {
	log.Println("something works")
	return "All is Well"
}
