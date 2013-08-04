package util

import (
	"fmt"
	"testing"
)

var msg = []byte("{\"name\":\"something\"}")

// var msg = []byte("eyJuYW1lIjoic29tZXRoaW5nIn0=")
var some_data = ReturnData{true, "No Error", msg, "200"}

func Test_ToString(t *testing.T) {
	testMsg := some_data.ToString()
	fmt.Println(testMsg)
}

func Test_ToString_With_Empty_Data(t *testing.T) {
	var some_data1 = ReturnData{true, "No Error", []byte(""), "200"}
	testMsg := some_data1.ToString()
	fmt.Println(testMsg)
}
