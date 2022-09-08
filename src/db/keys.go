package db

import (
	"errors"

	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/wgmanager"
	"gorm.io/gorm"
)

// CREATE

// Adds a key to the database and returns the keyID
func createKey(serverID int, publicKey string, presharedKey string) (keyID int, err error) {
	db := DBSystem
	_, err = ReadServer(serverID)
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

//Adds a user's key after checking their subscription validity
func CreateKeyAndLink(userID int, serverID int, publicKey string, presharedKey string) (err error) {
	err = checkSubscriptionKeyAddition(userID)
	if err != nil {
		return
	}

	keyID, err := createKey(serverID, publicKey, presharedKey)
	if err != nil {
		return
	}
	_, err = createUserKeyLink(userID, keyID)
	return
}

// READ

//finds a key object from a keyID
func readKey(keyID int) (key Keys, err error) {
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

//finds all the user's keys from their userID
func readUserKeys(userID int) (userKeys []UserKeys, err error) {
	db := DBSystem
	userKeysQuery := db.Where("user_id = ?", userID).Find(&userKeys)
	if !errors.Is(userKeysQuery.Error, gorm.ErrRecordNotFound) {
		err = ErrKeyNotFound
		return
	} else if userKeysQuery.Error != nil {
		return
	}
	return
}

//gets all keys in database
func ReadAllKeys() (keys []Keys, err error) {
	db := DBSystem

	dbResult := db.Find(&keys)
	err = dbResult.Error
	return keys, err
}

func readKeysWithServerID(serverID int) (keys []Keys, err error) {
	db := DBSystem

	keyQuery := db.Where("serverID = ?", serverID).Find(&keys)
	if errors.Is(keyQuery.Error, gorm.ErrRecordNotFound) {
		err = ErrKeyNotFound
	} else if keyQuery.Error != nil {
		combinedLogger.Error("Finding key " + keyQuery.Error.Error())
		err = ErrQuery
	}
	return
}

//UPDATE

//updates a key object
func updateKey(key Keys) (err error) {
	db := DBSystem

	err = db.Save(&key).Error
	return
}

// DELETE

//Deletes a key from keyID
func DeleteKey(keyID int) (err error) {
	db := DBSystem

	keyQuery, err := readKey(keyID)
	if err != nil {
		return
	}

	keyDelete := db.Delete(&keyQuery)
	if keyDelete.Error != nil {
		err = ErrDeletingKey
	}
	return
}

//Deletes a user's key and link
func DeleteKeyAndLink(keyID int) (err error) {
	err = deleteUserKeyLink(keyID)
	if err != nil {
		return
	}
	err = DeleteKey(keyID)
	return
}

//MISC

//Toggles a key usability from true to false and viceversa
func ToggleKey(keyID int) (err error) {
	key, err := readKey(keyID)
	if err != nil {
		return
	}

	key.Enabled = !key.Enabled
	err = updateKey(key)
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
