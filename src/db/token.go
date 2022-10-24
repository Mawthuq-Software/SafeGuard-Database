package db

import (
	"errors"
	"os/exec"
	"strings"

	"gorm.io/gorm"
)

func CreateToken(name string) (uuid string, err error) {
	db := DBSystem

	newUUIDStr, hashedUUID, err := generateToken()
	newToken := Tokens{Name: name, AccessToken: hashedUUID}
	tokenCreation := db.Create(&newToken)
	if tokenCreation.Error != nil {
		err = ErrAddingToken
		return
	}

	return newUUIDStr, err
}

func ReadToken(tokenID int) (token Tokens, err error) {
	db := DBSystem

	findToken := db.Where("id = ?", tokenID).First(&token)
	if errors.Is(findToken.Error, gorm.ErrRecordNotFound) {
		err = ErrTokenNotFound
		return
	}
	return
}

func UpdateToken(token Tokens) (err error) {
	db := DBSystem

	saveErr := db.Save(&token)
	if saveErr.Error != nil {
		return saveErr.Error
	}
	return
}

func DeleteToken(tokenID int) (err error) {
	db := DBSystem

	token, err := ReadToken(tokenID)
	if err != nil {
		return err
	}
	deleteErr := db.Delete(&token)
	if deleteErr.Error != nil {
		return deleteErr.Error
	}
	return
}

func generateToken() (uuid string, hashedUUID string, err error) {
	newUUIDByte, genErr := exec.Command("uuidgen").Output()
	if genErr != nil {
		err = ErrGeneratingToken
		return
	}

	newUUIDStr := string(newUUIDByte[:])
	newUUIDStr = strings.Replace(newUUIDStr, "\n", "", -1) //Replace new lines
	var passHash Hash
	hashedUUID, hashErr := passHash.Generate(newUUIDStr)
	if hashErr != nil {
		combinedLogger.Error("Hashing password " + hashErr.Error())
		err = ErrHashing
		return
	}
	return newUUIDStr, hashedUUID, err
}
