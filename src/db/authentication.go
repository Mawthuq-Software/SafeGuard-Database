package db

import (
	"errors"

	"gorm.io/gorm"
)

func FindAuthFromUsername(username string) (Authentications, error) {
	db := DBSystem
	var findAuth Authentications

	authQuery := db.Where("username = ?", username).First(&findAuth)
	if errors.Is(authQuery.Error, gorm.ErrRecordNotFound) {
		return findAuth, ErrAuthNotFound
	} else if authQuery.Error != nil {
		combinedLogger.Error("Finding auth " + authQuery.Error.Error())
		return findAuth, ErrQuery
	}
	return findAuth, nil
}

func FindAuthFromUserID(userID int) (auth Authentications, err error) {
	user, err := FindUserFromUserID(userID)
	if err != nil {
		return
	}
	findAuth, err := findAuthFromAuthID(user.AuthID)
	if err != nil {
		return
	}
	return findAuth, nil
}

func FindUserFromAuthID(authID int) (Users, error) {
	db := DBSystem
	var findUser Users

	userQuery := db.Where("auth_id = ?", authID).First(&findUser)
	if errors.Is(userQuery.Error, gorm.ErrRecordNotFound) {
		return findUser, ErrUserNotFound
	} else if userQuery.Error != nil {
		combinedLogger.Error("Finding user " + userQuery.Error.Error())
		return findUser, ErrQuery
	}

	return findUser, nil
}

func findAuthFromAuthID(authID int) (Authentications, error) {
	db := DBSystem
	var findAuth Authentications

	authQuery := db.Where("id = ?", authID).First(&findAuth)
	if errors.Is(authQuery.Error, gorm.ErrRecordNotFound) {
		return findAuth, ErrUserNotFound
	} else if authQuery.Error != nil {
		combinedLogger.Error("Finding authentication " + authQuery.Error.Error())
		return findAuth, ErrQuery
	}

	return findAuth, nil
}
