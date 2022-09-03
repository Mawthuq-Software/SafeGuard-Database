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

//Deletes a key from keyID
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

//Adds a link between the userID and keyID in userKey table
func addUserKeyLink(userID int, keyID int) (userKeyID int, err error) {
	db := DBSystem

	_, err = FindUserFromUserID(userID)
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

//Deletes a link between a user and key
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

//finds a key object from a keyID
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

//finds a user key link from the keyID
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

//finds all the user's keys from their userID
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

//checks to see if wireguard key is appropriate
func checkKeyValidity(key string) (err error) {
	_, err = wgmanager.ParseKey(key) //parse string
	if err != nil {
		err = ErrPublicKeyIncorrectForm
	}
	return
}

//updates a key object
func updateKey(key Keys) (err error) {
	db := DBSystem

	err = db.Save(&key).Error
	return
}

//gets all keys in database
func getAll() (keys []Keys, err error) {
	db := DBSystem

	dbResult := db.Find(&keys)
	err = dbResult.Error
	return keys, err
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
	return
}

//Deletes a user's key
func DeleteUserKey(keyID int) (err error) {
	err = deleteUserKeyLink(keyID)
	if err != nil {
		return
	}
	err = deleteKey(keyID)
	return
}

//Toggles a key usability from true to false and viceversa
func ToggleKey(keyID int) (err error) {
	key, err := findKeyFromKeyID(keyID)
	if err != nil {
		return
	}

	key.Enabled = !key.Enabled
	err = updateKey(key)
	return
}

func GetAllKeys() (keys []Keys, err error) {
	keys, err = getAll()
	return
}
