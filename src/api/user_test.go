package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/api/router/responses"
	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/api/router/routes"
	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/testsetup"
)

var url = "https://eu1.raspberrytech.one:8446"

func newRandomUser() (email string, username string, password string) {
	email = testsetup.RandomString(9) + "@fakemail.com"
	username = testsetup.RandomString(11)
	password = testsetup.RandomString(13)
	return
}

func addNewRandomUser(randomUserEmail string, randomUserUsername string, randomUserPassword string) (res *http.Response) {
	var newRandomUser routes.User

	newRandomUser.Email = randomUserEmail
	newRandomUser.Username = randomUserUsername
	newRandomUser.Password = randomUserPassword

	jsonResp, _ := json.Marshal(newRandomUser)
	reqReader := bytes.NewReader(jsonResp)

	req, _ := http.NewRequest("POST", url+"/user/register", reqReader)
	res = testsetup.MakeReq(req)
	return
}

func loginUser(existingUser routes.User) (res *http.Response) {
	jsonResp, _ := json.Marshal(existingUser)
	reqReader := bytes.NewReader(jsonResp)

	req, _ := http.NewRequest("POST", url+"/user/login", reqReader)
	res = testsetup.MakeReq(req)
	return
}

func checkToken(userToken string) (res *http.Response) {
	req, _ := http.NewRequest("POST", url+"/token/validate", nil)
	req.Header.Set("Bearer", userToken)
	res = testsetup.MakeReq(req)
	return
}

func Test_AddUser(t *testing.T) {
	email, username, password := newRandomUser()
	res := addNewRandomUser(email, username, password)
	testsetup.CheckResponseCode(t, http.StatusAccepted, res.StatusCode)
}

func Test_AddExistingUser(t *testing.T) {
	email, username, password := newRandomUser()
	res := addNewRandomUser(email, username, password)
	if res.StatusCode == http.StatusAccepted {
		res = addNewRandomUser(email, username, password)
		testsetup.CheckResponseCode(t, http.StatusBadRequest, res.StatusCode)
	} else {
		t.Error("Random user did not come back as correct status code")
	}
}

//This test must run for token to be generated
func Test_LoginExistingUser(t *testing.T) {
	// create new user
	var existingUser routes.User
	email, username, password := newRandomUser()
	res := addNewRandomUser(email, username, password)

	existingUser.Email = email
	existingUser.Username = username
	existingUser.Password = password

	if res.StatusCode == http.StatusAccepted {
		res := loginUser(existingUser)
		testsetup.CheckResponseCode(t, http.StatusAccepted, res.StatusCode)
	} else {
		t.FailNow()
	}
}

func Test_LoginNonExistentUser(t *testing.T) {
	// create new user
	var existingUser routes.User
	email, username, password := newRandomUser()

	existingUser.Email = email
	existingUser.Username = username
	existingUser.Password = password

	res := loginUser(existingUser)
	testsetup.CheckResponseCode(t, http.StatusBadRequest, res.StatusCode)
}

func Test_ValidateToken(t *testing.T) {
	// create new user
	var existingUser routes.User
	email, username, password := newRandomUser()
	res := addNewRandomUser(email, username, password)

	existingUser.Email = email
	existingUser.Username = username
	existingUser.Password = password

	if res.StatusCode == http.StatusAccepted {
		res := loginUser(existingUser)
		if res.StatusCode == http.StatusAccepted {
			var token responses.TokenResponse
			err := routes.ParseResponse(res, &token)
			if err != nil {
				t.Errorf(err.Error())
			}
			res = checkToken(token.Token)
			testsetup.CheckResponseCode(t, http.StatusOK, res.StatusCode)
		} else {
			t.FailNow()
		}
	} else {
		t.FailNow()
	}
}
