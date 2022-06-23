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
