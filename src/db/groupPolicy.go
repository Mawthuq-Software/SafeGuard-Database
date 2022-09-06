package db

import (
	"errors"

	"gorm.io/gorm"
)

// CREATE

func CreateGroupPolicies(groupName string, policyNames []string) error {
	db := DBSystem
	var findGroup Groups

	resFindGroup := db.Where("name = ?", groupName).First(&findGroup)
	if resFindGroup.Error != nil {
		combinedLogger.Warn("Finding group " + resFindGroup.Error.Error())
		return resFindGroup.Error
	}

	for i := 0; i < len(policyNames); i++ {
		var findPolicy Policies
		resFindPol := db.Where("name = ?", policyNames[i]).First(&findPolicy)
		if resFindPol.Error != nil {
			combinedLogger.Warn("Finding policy " + groupName + "to db" + resFindPol.Error.Error())
			return resFindPol.Error
		}

		groupPolicy := GroupPolicies{GroupID: findGroup.ID, PolicyID: findPolicy.ID}
		createGroupPolicy := db.Create(&groupPolicy)
		if createGroupPolicy.Error != nil {
			combinedLogger.Error("Creating group " + groupName + " policy " + policyNames[i] + "to db" + resFindPol.Error.Error())
			return createGroupPolicy.Error
		}
	}
	return nil
}

// READ

func ReadGroupPolicies(groupID int) ([]GroupPolicies, error) {
	db := DBSystem
	var findGroupPol []GroupPolicies

	errUserGroupQuery := db.Where("group_id = ?", groupID).Find(&findGroupPol)
	if errors.Is(errUserGroupQuery.Error, gorm.ErrRecordNotFound) {
		return nil, ErrGroupPolicyNotFound
	} else if errUserGroupQuery.Error != nil {
		return nil, ErrQuery
	}
	return findGroupPol, nil
}
