package db

import (
	"errors"
	"time"

	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/token"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Adds a new user
func AddUser(username string, password string, email string) error {
	var userAuthStruct Authentications
	var groupStruct Groups
	db := DBSystem

	findAuth := db.Where("username = ?", username).Or("email = ?", email).First(&userAuthStruct) //find authentication in database
	if !errors.Is(findAuth.Error, gorm.ErrRecordNotFound) {
		return ErrUserExists
	}

	findGroup := db.Where("name = ?", "User").First(&groupStruct)
	if errors.Is(findGroup.Error, gorm.ErrRecordNotFound) {
		return ErrGroupNotFound
	} else if findGroup.Error != nil {
		combinedLogger.Warn("Finding group user in db " + findGroup.Error.Error())
		return ErrQuery
	}

	var passHash Hash
	hashedPassword, hashErr := passHash.Generate(password)
	if hashErr != nil {
		combinedLogger.Error("Hashing password " + hashErr.Error())
		return ErrHashing
	}

	newAuth := Authentications{Username: username, Password: hashedPassword, Email: email}
	resultAuthCreation := db.Create(&newAuth)
	if resultAuthCreation.Error != nil {
		combinedLogger.Error("Adding authentication to db " + resultAuthCreation.Error.Error())
		return ErrCreatingAuth
	}

	newUser := Users{AuthID: newAuth.ID}
	userCreation := db.Create(&newUser)
	if userCreation.Error != nil {
		combinedLogger.Error("Adding user to db " + userCreation.Error.Error())
		return ErrCreatingUser
	}

	newUserGroup := UserGroups{GroupID: groupStruct.ID, UserID: newUser.ID}
	userGroupCreation := db.Create(&newUserGroup)
	if userGroupCreation.Error != nil {
		combinedLogger.Error("Adding group to db " + userGroupCreation.Error.Error())
		return ErrCreatingUserGroup
	}
	return nil
}

//Logs in the user by their username and password
func LoginWithUsername(username string, password string) (string, error) {
	var passHash Hash

	findAuth, authErr := FindAuthFromUsername(username)
	if authErr != nil {
		return "", authErr
	}
	if passHash.Compare(findAuth.Password, password) != nil {
		return "", ErrIncorrectPass
	}

	userQuery, userErr := FindUserFromAuthID(findAuth.ID)
	if userErr != nil {
		return "", userErr
	}
	hasPerms, permsErr := CheckPermission(userQuery.ID, PERSONAL_LOGIN)
	if permsErr != nil {
		return "", permsErr
	} else if !hasPerms {
		return "", ErrMissingPermission
	}
	tokenLifetime := time.Now().AddDate(0, 0, 7)
	generatedToken, tokenErr := token.GenerateUser(username, tokenLifetime)
	if tokenErr != nil {
		combinedLogger.Error("Generating token " + tokenErr.Error())
		return "", ErrCreatingToken
	}
	return generatedToken, nil
}

// Changes the users password
func ChangeUserPassword(username string, oldPassword string, newPassword string) (err error) {
	db := DBSystem

	var passHash Hash
	findAuth, err := FindAuthFromUsername(username)
	if err != nil {
		return err
	}
	if passHash.Compare(findAuth.Password, oldPassword) != nil {
		return ErrIncorrectPass
	}

	hashedPassword, err := passHash.Generate(newPassword)
	if err != nil {
		combinedLogger.Error("Hashing password " + err.Error())
		return ErrHashing
	}

	findAuth.Password = hashedPassword
	hashSave := db.Save(&findAuth)
	if hashSave.Error != nil {
		return ErrSavingPassword
	}
	return nil
}

// Finds user from the userID
func findUserFromUserID(userID int) (user Users, err error) {
	db := DBSystem

	userQuery := db.Where("id = ?", userID).First(&user)
	if errors.Is(userQuery.Error, gorm.ErrRecordNotFound) {
		err = ErrUserNotFound
	} else if userQuery.Error != nil {
		combinedLogger.Error("Finding user " + userQuery.Error.Error())
		err = ErrQuery
	}
	return
}
func FindUserFromUsername(username string) (user Users, err error) {
	authentication, err := FindAuthFromUsername(username)
	if err != nil {
		return
	}

	user, err = FindUserFromAuthID(authentication.ID)
	return
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
