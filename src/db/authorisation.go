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
