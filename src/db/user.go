package db

import (
	"errors"
	"log"

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
	}

	newAuth := Authentications{Username: username, Password: password, Email: email}
	resultAuthCreation := db.Create(newAuth)
	if resultAuthCreation.Error != nil {
		log.Println("Error - Adding authentication to db", resultAuthCreation.Error)
		processed.Proccessed = false
		processed.Response = "Error when adding authentication to database"
		return processed
	}

	findGroup := db.Where("name = ?", "User").First(&groupStruct)
	if errors.Is(findGroup.Error, gorm.ErrRecordNotFound) {
		processed.Proccessed = false
		processed.Response = "Group was not found on the server"
		return processed

	} else if findGroup.Error != nil {
		log.Println("Warning - Finding group user in db ", findGroup.Error)
		processed.Proccessed = false
		processed.Response = "Error finding group"
		return processed
	}

	userCreation := db.Create(Users{AuthID: newAuth.ID, GroupID: groupStruct.ID})
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
