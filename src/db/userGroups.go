package db

import (
	"errors"

	"gorm.io/gorm"
)

func ReadAllUserGroups(userID int) ([]UserGroups, error) {
	db := DBSystem
	var findUserGroup []UserGroups

	errUserGroupQuery := db.Where("user_id = ?", userID).Find(&findUserGroup)
	if errors.Is(errUserGroupQuery.Error, gorm.ErrRecordNotFound) {
		return nil, ErrUserGroupNotFound
	} else if errUserGroupQuery.Error != nil {
		return nil, ErrQuery
	}
	return findUserGroup, nil
}
