package testsetup

import (
	"log"
	"math/rand"
	"net/http"
	"strings"
	"testing"
	"time"
)

func RandomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	charSet := "abcdedfghijklmnopqrstABCDEFGHIJKLMNOPQRST"
	var output strings.Builder
	for i := 0; i < n; i++ {
		random := rand.Intn(len(charSet))
		randomChar := charSet[random]
		output.WriteString(string(randomChar))
	}
	return output.String()
}

func CheckResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func CheckErr(t *testing.T, expected, actual error) {
	if expected != actual {
		t.Errorf("Expected error code %s. Got %s\n", expected.Error(), actual.Error())
	}
}

func CheckString(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Errorf("Expected string %s. Got %s\n", expected, actual)
	}
}

func MakeReq(req *http.Request) *http.Response {
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	return resp
}
