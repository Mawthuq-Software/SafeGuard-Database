package main_test

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/api/router/responses"
	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/api/router/routes"
)

var url = "https://eu1.raspberrytech.one:8446"
var userEmail = randomString(9) + "@fakemail.com"
var userUsername = randomString(11)
var userPassword = randomString(13)
var userToken string

func Test_AddUser(t *testing.T) {
	var newUser routes.User
	newUser.Email = userEmail
	newUser.Username = userUsername
	newUser.Password = userPassword

	jsonResp, _ := json.Marshal(newUser)
	reqReader := bytes.NewReader(jsonResp)

	req, _ := http.NewRequest("POST", url+"/user/add", reqReader)
	res := makeReq(req)

	checkResponseCode(t, http.StatusAccepted, res.StatusCode)
}

func Test_AddExistingUser(t *testing.T) {
	var existingUser routes.User
	existingUser.Email = userEmail
	existingUser.Username = userUsername
	existingUser.Password = userPassword

	jsonResp, _ := json.Marshal(existingUser)
	reqReader := bytes.NewReader(jsonResp)

	req, _ := http.NewRequest("POST", url+"/user/add", reqReader)
	res := makeReq(req)

	checkResponseCode(t, http.StatusBadRequest, res.StatusCode)
}

//This test must run for token to be generated
func Test_LoginExistingUser(t *testing.T) {
	var existingUser routes.User
	existingUser.Username = userUsername
	existingUser.Password = userPassword

	jsonResp, _ := json.Marshal(existingUser)
	reqReader := bytes.NewReader(jsonResp)

	req, _ := http.NewRequest("POST", url+"/user/login", reqReader)
	res := makeReq(req)

	var user responses.TokenResponse
	err := routes.ParseResponse(res, &user)
	if err != nil {
		t.Errorf(err.Error())
	}
	userToken = user.Token
	checkResponseCode(t, http.StatusAccepted, res.StatusCode)
}

func Test_LoginNonExistentUser(t *testing.T) {
	var existingUser routes.User
	existingUser.Username = randomString(20)
	existingUser.Password = randomString(20)

	jsonResp, _ := json.Marshal(existingUser)
	reqReader := bytes.NewReader(jsonResp)

	req, _ := http.NewRequest("POST", url+"/user/login", reqReader)
	res := makeReq(req)

	checkResponseCode(t, http.StatusBadRequest, res.StatusCode)
}

func Test_ValidateToken(t *testing.T) {
	req, _ := http.NewRequest("POST", url+"/token/validate", nil)
	req.Header.Set("Bearer", userToken)
	res := makeReq(req)

	checkResponseCode(t, http.StatusOK, res.StatusCode)
}

func randomString(n int) string {
	rand.Seed(time.Now().Unix())
	charSet := "abcdedfghijklmnopqrstABCDEFGHIJKLMNOPQRST"
	var output strings.Builder
	for i := 0; i < n; i++ {
		random := rand.Intn(len(charSet))
		randomChar := charSet[random]
		output.WriteString(string(randomChar))
	}
	return output.String()
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
func makeReq(req *http.Request) *http.Response {
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	return resp
}
