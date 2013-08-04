package circle

import (
	"fmt"
	"testing"
	"time"
)

var grp = Circle{
	Name:        "Tapori",
	Description: "Test Group",
	CreatorID:   "011",
	CreatedOn:   time.Now().String(),
	Members:     []string{"011"},
}

func Test_MakeCircle(t *testing.T) {
	test_circle, err := grp.makeCircle()
	if test_circle {
		t.Log("Test_MakeCircle PASSED")
	} else {
		t.Log(err.Error())
		t.Log("Test_MakeCircle Failed")
		t.Fail()
	}
}

func Test_CircleExists(t *testing.T) {
	exists := grp.CircleExists()
	if exists {
		fmt.Println(exists)
	}
}
