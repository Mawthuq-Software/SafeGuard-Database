package db

import (
	"errors"

	"gorm.io/gorm"
)

func ReadUserGroup(userGroupID int) (userGroup UserGroups, err error) {
	db := DBSystem

	errUserGroupQuery := db.Where("id = ?", userGroupID).Find(&userGroup)
	if errors.Is(errUserGroupQuery.Error, gorm.ErrRecordNotFound) {
		err = ErrUserGroupNotFound
	} else if errUserGroupQuery.Error != nil {
		err = ErrQuery
	}
	return
}

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

func DeleteUserGroup(userGroupID int) (err error) {
	db := DBSystem

	userGroup, err := ReadUserGroup(userGroupID)
	if err != nil {
		return
	}

	delete := db.Delete(&userGroup)
	return delete.Error
}
