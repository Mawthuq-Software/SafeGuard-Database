package db

import (
	"errors"

	"gorm.io/gorm"
)

func findServerFromServerID(serverID int) (server Servers, err error) {
	db := DBSystem

	serverQuery := db.Where("id = ?", serverID).First(&server)
	if errors.Is(serverQuery.Error, gorm.ErrRecordNotFound) {
		err = ErrServerNotFound
		return
	} else if serverQuery.Error != nil {
		err = ErrQuery
		return
	}
	return
}
