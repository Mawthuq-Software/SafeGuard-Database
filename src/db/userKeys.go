package db

import (
	"errors"

	"gorm.io/gorm"
)

// CREATE

//Adds a link between the userID and keyID in userKey table
func createUserKeyLink(userID int, keyID int) (userKeyID int, err error) {
	db := DBSystem

	_, err = FindUserFromUserID(userID)
	if err != nil {
		return
	}
	_, err = readKey(keyID)
	if err != nil {
		return
	}
	userKey := UserKeys{UserID: userID, KeyID: keyID}
	userKeyCreation := db.Create(&userKey)
	if userKeyCreation.Error != nil {
		err = ErrUserKeyLink
		return
	}
	userKeyID = userKey.ID
	return
}

// READ

//finds a user key link from the keyID
func readUserKey(keyID int) (userKeys UserKeys, err error) {
	db := DBSystem

	keyQuery := db.Where("id = ?", keyID).First(&userKeys)
	if errors.Is(keyQuery.Error, gorm.ErrRecordNotFound) {
		err = ErrUserKeyNotFound
	} else if keyQuery.Error != nil {
		combinedLogger.Error("Finding user key " + keyQuery.Error.Error())
		err = ErrQuery
	}
	return
}

// DELETE

//Deletes a link between a user and key
func deleteUserKeyLink(keyID int) (err error) {
	db := DBSystem

	userKey, err := readUserKey(keyID)
	if err != nil {
		return
	}
	userKeyDelete := db.Delete(&userKey)
	if userKeyDelete.Error != nil {
		err = ErrDeletingUserKey
	}
	return
}
