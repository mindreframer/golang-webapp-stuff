package geolocation

import (
	"encoding/json"
	SJ "github.com/bitly/go-simplejson"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type GeoLocation struct {
	FourID   string  `json:"four_id"	bson:"four_id"`
	Name     string  `json:"name"	bson:"name"`
	Contact  string  `json:"contact"	bson:"contact"`
	Address  string  `json:"address"	bson:"address"`
	Lat      float64 `json:"lat"	bson:"lat"`
	Lng      float64 `json:"lng"	bson:"lng"`
	Distance int     `json:"distance"	bson:"distance"`
	Postcode string  `json:"postcode"	bson:"postcode"`
	City     string  `json:"city"	bson:"city"`
	State    string  `json:"state"	bson:"state"`
	Country  string  `json:"country"	bson:"country"`
}

type Feed struct {
	Gmeta     Metaf    `json:"meta"`
	Gresponse Response `json:"response"`
}

type Response struct {
	Rvenue []Venue `json:"venues"`
}

type Venue struct {
	Id        string   `xml:"id"`
	Name      string   `xml:"name"`
	Gcontact  Contact  `json:"contact"`
	Glocation Location `json:"location"`
	// canonicalUrl string
	// Categories   []Category
	// Verified     bool
	// Restricted   bool
	// Stats        Stat
	// Url          string
	// ReferralId   string
}

type Contact struct {
	Phone          float64 `xml:"phone"`
	FormattedPhone string  `xml:"formatedPhone"`
}

type Location struct {
	Address    string  `xml:"address"`
	Lat        float64 `xml:"lat"`
	Lng        float64 `xml:"lng"`
	PostalCode string  `xml:"postalCode"`
	City       string  `xml:"city"`
	State      string  `xml:"state"`
	Country    string  `xml:"country"`
	CC         string  `xml:"cc"`
}

type Stat struct {
	CheckinsCount float64
	UserCount     float64
	TipCount      float64
}

type Special struct {
	Count float64
	Items []string
}

type Metaf struct {
	Code float64 `xml:"code"`
}
type Category struct {
	Id         string
	Name       string
	PluralName string
	ShortName  string
	Logo       Icon
	Primary    bool
}

type Icon struct {
	Prefix string
	Suffix string
}

func GetVenues(near string) string {
	FSqrUrl := "https://api.foursquare.com/v2/venues/search?v=20130417&near=<nearLocation>&client_id=" + os.Getenv("FSQR_CLIENT_ID") + "&client_secret=" + os.Getenv("FSQR_CLIENT_SECRET")
	FSqrUrl = strings.Replace(FSqrUrl, "<nearLocation>", near, -1)
	log.Println(FSqrUrl)

	return getLocations(FSqrUrl)
}

func GetVenuesWithLatitudeAndLongitude(lt, lg string) string {
	var FSqrUrl string
	FSqrUrl = "https://api.foursquare.com/v2/venues/search?v=20130417&ll=" + lt + "," + lg + "&client_id=" + os.Getenv("FSQR_CLIENT_ID") + "&client_secret=" + os.Getenv("FSQR_CLIENT_SECRET")
	// FSqrUrl = strings.Replace(FSqrUrl, "<lat>", lt, -1)
	// FSqrUrl = strings.Replace(FSqrUrl, "<lng>", lg, -1)

	return getLocations(FSqrUrl)
}

func GetVenueWithId(venueId string) string {
	var FSqrUrl string
	FSqrUrl = "https://api.foursquare.com/v2/venues/" + venueId + "&client_id=" + os.Getenv("FSQR_CLIENT_ID") + "&client_secret=" + os.Getenv("FSQR_CLIENT_SECRET")
	return getLocations(FSqrUrl)
}

func getLocations(FSqrUrl string) string {
	res, err := http.Get(FSqrUrl)
	if err != nil {
		log.Println("getVenues error: " + err.Error())
	}
	defer res.Body.Close()
	str, _ := ioutil.ReadAll(res.Body)
	// Venues := getJsonVenues(str)
	return string(str)
}

func getJsonVenues(resp []byte) string {
	var str string
	var flt float64
	locs := "["

	js, err := SJ.NewJson(resp)
	if err != nil {
		log.Println(err.Error())
		return err.Error()
	}
	arr := js.Get("response").Get("venues")
	resLen, _ := arr.Array()

	// fmt.Println(len(resLen))
	places := make([]Venue, len(resLen))
	for i := 0; i < len(resLen); i++ {
		if id, ok := arr.GetIndex(i).CheckGet("id"); ok {
			str, _ = id.String()
			places[i].Id = str
		}
		if name, ok := arr.GetIndex(i).CheckGet("name"); ok {
			str, _ = name.String()
			places[i].Name = str
		}
		if phone, ok := arr.GetIndex(i).Get("contact").CheckGet("phone"); ok {
			flt, _ = phone.Float64()
			places[i].Gcontact.Phone = flt
		}
		if formattedPhone, ok := arr.GetIndex(i).Get("contact").CheckGet("formattedPhone"); ok {
			str, _ = formattedPhone.String()
			places[i].Gcontact.FormattedPhone = str
		}
		if address, ok := arr.GetIndex(i).Get("location").CheckGet("address"); ok {
			str, _ = address.String()
			places[i].Glocation.Address = str
		}
		if lat, ok := arr.GetIndex(i).Get("location").CheckGet("lat"); ok {
			flt, _ = lat.Float64()
			places[i].Glocation.Lat = flt
		}
		if lng, ok := arr.GetIndex(i).Get("location").CheckGet("lng"); ok {
			flt, _ = lng.Float64()
			places[i].Glocation.Lng = flt
		}
		if postalCode, ok := arr.GetIndex(i).Get("location").CheckGet("postalCode"); ok {
			str, _ = postalCode.String()
			places[i].Glocation.PostalCode = str
		}
		if city, ok := arr.GetIndex(i).Get("location").CheckGet("city"); ok {
			str, _ = city.String()
			places[i].Glocation.City = str
		}
		if state, ok := arr.GetIndex(i).Get("location").CheckGet("state"); ok {
			str, _ = state.String()
			places[i].Glocation.State = str
		}
		if country, ok := arr.GetIndex(i).Get("location").CheckGet("country"); ok {
			str, _ = country.String()
			places[i].Glocation.Country = str
		}
		if cc, ok := arr.GetIndex(i).Get("location").CheckGet("cc"); ok {
			str, _ = cc.String()
			places[i].Glocation.CC = str
		}
		temp, _ := json.Marshal((places[i]))
		locs += string(temp)

		if i < len(resLen)-1 {
			locs += ","
		}
	}
	locs += "]"
	return string(locs)
}
