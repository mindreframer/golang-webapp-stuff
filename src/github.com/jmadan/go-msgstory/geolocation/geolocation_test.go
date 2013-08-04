package geolocation

import (
	"fmt"
	SJ "github.com/bitly/go-simplejson"
	"testing"
)

func Test_GetVenuesWithLatitudeAndLongitude(t *testing.T) {
	// location := GeoLocation{}
	resbody := GetVenuesWithLatitudeAndLongitude("37.422005", "-122.084095")
	js, err := SJ.NewJson([]byte(resbody))
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	arr, _ := js.Get("meta").Get("code").Int()
	if arr == 200 {
		t.Log("Test_GetVenuesWithLatitudeAndLongitude PASSED")
	} else {
		t.Fail()
	}
}

// func Test_GetVenues(t *testing.T) {
// 	location := GeoLocation{}
// 	res := location.GetVenues("Manchester")
// 	if res == "200" {
// 		t.Log("Passed")
// 	}

// }
