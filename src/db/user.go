package db

import (
	"errors"
	"log"
	"strconv"
	"time"

	"gitlab.com/mawthuq-software/wireguard-manager-authenticator/src/token"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func AddUser(username string, password string, email string) DatabaseResponse {
	var userAuthStruct Authentications
	var groupStruct Groups
	processed := DatabaseResponse{}
	db := DBSystem

	findUser := db.Where("username = ?", username).Or("email = ?", email).First(&userAuthStruct) //find authentication in database
	if !errors.Is(findUser.Error, gorm.ErrRecordNotFound) {
		processed.Proccessed = false
		processed.Response = "Username or email already exists"
		return processed
	} else if findUser.Error != nil {
		processed.Proccessed = false
		processed.Response = "Error finding user"
	}

	findGroup := db.Where("name = ?", "User").First(&groupStruct)
	if errors.Is(findGroup.Error, gorm.ErrRecordNotFound) {
		processed.Proccessed = false
		processed.Response = "Group was not found"
		return processed

	} else if findGroup.Error != nil {
		log.Println("Warning - Finding group user in db ", findGroup.Error)
		processed.Proccessed = false
		processed.Response = "Error finding group"
		return processed
	}

	var passHash Hash
	hashedPassword, hashErr := passHash.Generate(password)
	if hashErr != nil {
		log.Println("Error - Hashing password", hashErr)
		processed.Proccessed = false
		processed.Response = "Error when adding password"
		return processed
	}

	newAuth := Authentications{Username: username, Password: hashedPassword, Email: email}
	resultAuthCreation := db.Create(&newAuth)
	if resultAuthCreation.Error != nil {
		log.Println("Error - Adding authentication to db", resultAuthCreation.Error)
		processed.Proccessed = false
		processed.Response = "Error when adding authentication to database"
		return processed
	}

	newUser := Users{AuthID: newAuth.ID, GroupID: groupStruct.ID}
	userCreation := db.Create(&newUser)
	if userCreation.Error != nil {
		log.Println("Error - Adding user to db", userCreation.Error)
		processed.Proccessed = false
		processed.Response = "Error when adding user to database"
		return processed
	}

	processed.Proccessed = true
	processed.Response = "Added user successfully"
	return processed
}

func LoginWithUsername(username string, password string) GenToken {
	db := DBSystem
	processed := GenToken{}
	var findAuth Authentications
	var findUser Users

	authQuery := db.Where("username = ?", username).First(&findAuth)
	if errors.Is(authQuery.Error, gorm.ErrRecordNotFound) {
		processed.Proccessed = false
		processed.Response = "User was not found on the system"
		return processed
	} else if authQuery.Error != nil {
		processed.Proccessed = false
		processed.Response = "Error finding user"
		return processed
	}

	var passHash Hash
	hashedPassword, hashErr := passHash.Generate(password)
	if hashErr != nil {
		log.Println("Error - Hashing password", hashErr)
		processed.Proccessed = false
		processed.Response = "Error when finding user"
		return processed
	}

	if hashedPassword != findAuth.Password {
		processed.Proccessed = false
		processed.Response = "Username or password is incorrect"
		return processed
	}

	userQuery := db.Where("auth_id = ?", findAuth.ID).First(&findUser)
	if errors.Is(userQuery.Error, gorm.ErrRecordNotFound) {
		processed.Proccessed = false
		processed.Response = "User was not found on the system"
		return processed
	} else if userQuery.Error != nil {
		processed.Proccessed = false
		processed.Response = "Error finding user"
		return processed
	}

	tokenLifetime := time.Now().AddDate(0, 0, 7)
	generatedToken, tokenErr := token.GenerateNew(strconv.Itoa(findUser.ID), tokenLifetime)
	if tokenErr != nil {
		processed.Proccessed = false
		processed.Response = "Error creating token"
		return processed
	}

	processed.Token = generatedToken
	processed.Proccessed = true
	processed.Response = "Created token successfully"
	return processed
}

//https://hackernoon.com/how-to-store-passwords-example-in-go-62712b1d2212
type Hash struct{}

//Generate a salted hash for the input string
func (c *Hash) Generate(s string) (string, error) {
	saltedBytes := []byte(s)
	hashedBytes, err := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	hash := string(hashedBytes[:])
	return hash, nil
}

//Compare string to generated hash
func (c *Hash) Compare(hash string, s string) error {
	incoming := []byte(s)
	existing := []byte(hash)
	return bcrypt.CompareHashAndPassword(existing, incoming)
}
