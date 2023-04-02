package db

import (
	"testing"

	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/config"
	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/testsetup"
	"gorm.io/gorm"
)

var db *gorm.DB

func setup() {
	config.LoadConfig()

	DBStart(true)
	db = DBSystem
}

func Test_AddUser(t *testing.T) {
	// variable setup
	setup()
	username := testsetup.RandomString(8)
	password := testsetup.RandomString(12)
	email := testsetup.RandomString(10) + "@abc.com"

	// check no issues when creating user
	err := CreateUser(username, password, email)
	testsetup.CheckErr(t, nil, err)

	auth, err := FindAuthFromUsername(username)

	testsetup.CheckErr(t, nil, err)
	testsetup.CheckString(t, auth.Username, username)
	// We can't compared password directly, its hashed in db
	// testsetup.CheckString(t, auth.Password, password)
	testsetup.CheckString(t, auth.Email, email)
}
