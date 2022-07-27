package db

import (
	"errors"

	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/wgmanager"
	"gorm.io/gorm"
)

// Adds a key to the database and returns the keyID
func addKey(serverID int, publicKey string, presharedKey string) (keyID int, err error) {
	db := DBSystem
	_, err = findServerFromServerID(serverID)
	if err != nil {
		return
	}

	//check pub key
	err = checkKeyValidity(publicKey)
	if err != nil {
		return
	}

	//check pre key
	err = checkKeyValidity(presharedKey)
	if err != nil {
		return
	}

	newKey := Keys{ServerID: serverID, PublicKey: publicKey, PresharedKey: presharedKey}
	keyCreation := db.Create(&newKey)
	if keyCreation.Error != nil {
		err = ErrCreatingKey
		return
	}
	keyID = newKey.ID
	return
}

func deleteKey(keyID int) (err error) {
	db := DBSystem

	keyQuery, err := findKeyFromKeyID(keyID)
	if err != nil {
		return
	}

	keyDelete := db.Delete(&keyQuery)
	if keyDelete.Error != nil {
		err = ErrDeletingKey
	}
	return
}

func addUserKeyLink(userID int, keyID int) (userKeyID int, err error) {
	db := DBSystem

	_, err = findUserFromUserID(userID)
	if err != nil {
		return
	}
	_, err = findKeyFromKeyID(keyID)
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

func deleteUserKeyLink(keyID int) (err error) {
	db := DBSystem

	userKey, err := findUserKeysFromKeyID(keyID)
	if err != nil {
		return
	}
	userKeyDelete := db.Delete(&userKey)
	if userKeyDelete.Error != nil {
		err = ErrDeletingUserKey
	}
	return
}

func findKeyFromKeyID(keyID int) (key Keys, err error) {
	db := DBSystem

	keyQuery := db.Where("id = ?", keyID).First(&key)
	if errors.Is(keyQuery.Error, gorm.ErrRecordNotFound) {
		err = ErrKeyNotFound
	} else if keyQuery.Error != nil {
		combinedLogger.Error("Finding key " + keyQuery.Error.Error())
		err = ErrQuery
	}
	return
}

func findUserKeysFromKeyID(keyID int) (userKeys UserKeys, err error) {
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

func findUserKeys(userID int) (userKeys []UserKeys, err error) {
	db := DBSystem
	userKeysQuery := db.Where("user_id = ?", userID).Find(&userKeys)
	if !errors.Is(userKeysQuery.Error, gorm.ErrRecordNotFound) {
		err = ErrSubscriptionNotFound
		return
	} else if userKeysQuery.Error != nil {
		return
	}
	return
}

func checkKeyValidity(key string) (err error) {
	_, err = wgmanager.ParseKey(key) //parse string
	if err != nil {
		err = ErrPublicKeyIncorrectForm
	}
	return
}

//Adds a user's key after checking their subscription validity
func AddUserKey(userID int, serverID int, publicKey string, presharedKey string) (err error) {
	err = checkSubscriptionKeyAddition(userID)
	if err != nil {
		return
	}

	keyID, err := addKey(serverID, publicKey, presharedKey)
	if err != nil {
		return
	}
	_, err = addUserKeyLink(userID, keyID)
	// if err != nil {
	// 	return
	// }
	return
}

func DeleteUserKey(keyID int) (err error) {
	err = deleteUserKeyLink(keyID)
	if err != nil {
		return
	}
	err = deleteKey(keyID)
	return
}
