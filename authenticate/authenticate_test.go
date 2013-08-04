package authenticate

import (
	"testing"
)

func Test_Authorize_Incorrect_Credentials(t *testing.T) {

	authenticate := Authenticate{"a@a.com", "something", 0, false}
	authenticate.Authorize()
	if authenticate.IsAuthenticated {
		t.Log("Test_Authorize_Incorrect_Credentials Failed")
		t.Fail()
	} else {
		t.Log("Test_Authorize_Incorrect_Credentials PASSED")
	}
}

func Test_Authorize_Correct_Credentials(t *testing.T) {
	authenticate := Authenticate{"jasdeepm@gmail.com", "98036054", 0, false}
	authenticate.Authorize()
	if authenticate.IsAuthenticated {
		t.Log("PASSED")
	} else {
		t.Log("Test_Authorize_Correct_Credentials Failed")
		t.Fail()
	}
}

func Test_Login_Method(t *testing.T) {
	auth := Login("jasdeepm@gmail.com", "98036054")
	if auth.IsAuthenticated {
		t.Log("Login Method PASSED")
	} else {
		t.Log("Login Method FAILED")
		t.Fail()
	}
}
