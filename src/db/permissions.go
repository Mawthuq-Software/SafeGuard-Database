package db

import (
	"errors"
	"strconv"

	"gorm.io/gorm"
)

const (
	//Standard User Permissions
	PERSONAL_KEYS_VIEW      int = 1
	PERSONAL_KEYS_ADD       int = 2
	PERSONAL_KEYS_DELETE    int = 3
	PERSONAL_KEYS_MODIFY    int = 4
	PERSONAL_PASSWORD_RESET int = 5
	PERSONAL_LOGIN          int = 6

	//Admin User Permissions
	KEYS_VIEW_ALL      int = 7
	KEYS_ADD_ALL       int = 8
	KEYS_ADD_OVERRIDE  int = 12
	KEYS_DELETE_ALL    int = 9
	KEYS_MODIFY_ALL    int = 10
	PASSWORD_RESET_ALL int = 11
)

func CheckPermission(userID int, policyID int) (bool, error) {
	db := DBSystem

	var findUser Users

	errUserQuery := db.Where("id = ?", userID).First(&findUser)
	if errors.Is(errUserQuery.Error, gorm.ErrRecordNotFound) {
		return false, ErrUserNotFound
	} else if errUserQuery != nil {
		return false, ErrQuery
	}

	userGroups, userGroupErr := GetAllUserGroups(userID)
	if userGroupErr != nil {
		return false, userGroupErr
	}

	for i := 0; i < len(userGroups); i++ {
		groupID := userGroups[i].GroupID
		groupPol, errGroupPol := GetGroupPolicies(groupID)
		if errGroupPol != nil {
			return false, errGroupPol
		}

		for x := 0; x < len(groupPol); x++ {
			if groupPol[x].PolicyID == policyID {
				return true, nil
			}
		}
	}
	return false, nil
}

func GetAllUserGroups(userID int) ([]UserGroups, error) {
	db := DBSystem
	var findUserGroup []UserGroups

	errUserGroupQuery := db.Where("user_id = ?", userID).Find(&findUserGroup)
	if errors.Is(errUserGroupQuery.Error, gorm.ErrRecordNotFound) {
		return nil, ErrUserGroupNotFound
	} else if errUserGroupQuery != nil {
		return nil, ErrQuery
	}
	return findUserGroup, nil
}

func GetGroupPolicies(groupID int) ([]GroupPolicies, error) {
	db := DBSystem
	var findGroupPol []GroupPolicies

	errUserGroupQuery := db.Where("group_id = ?", groupID).Find(&findGroupPol)
	if errors.Is(errUserGroupQuery.Error, gorm.ErrRecordNotFound) {
		return nil, ErrGroupPolicyNotFound
	} else if errUserGroupQuery != nil {
		return nil, ErrQuery
	}
	return findGroupPol, nil
}

func startupCreation() {

	standardVPNPerms := []int{PERSONAL_KEYS_VIEW, PERSONAL_KEYS_ADD}
	AddPolicy("STANDARD_USER_VPN", standardVPNPerms)

	advancedUserVPNPerms := []int{PERSONAL_KEYS_MODIFY, PERSONAL_KEYS_DELETE}
	AddPolicy("ADVANCED_USER_VPN", advancedUserVPNPerms)

	userSettingPerms := []int{PERSONAL_PASSWORD_RESET, PERSONAL_LOGIN}
	AddPolicy("STANDARD_USER_SETTINGS", userSettingPerms)

	adminPerms := []int{KEYS_VIEW_ALL, KEYS_ADD_ALL, KEYS_ADD_OVERRIDE, KEYS_DELETE_ALL, KEYS_MODIFY_ALL, PASSWORD_RESET_ALL}
	AddPolicy("ADMIN_USER", adminPerms)
	AddGroup("User")
	userPolicies := []string{"STANDARD_USER_VPN", "STANDARD_USER_SETTINGS"}
	AddGroupPolicies("User", userPolicies)

	adminPolicies := []string{"STANDARD_USER_VPN", "STANDARD_USER_SETTINGS", "ADVANCED_USER_VPN", "ADMIN_USER"}
	AddGroupPolicies("Admin", adminPolicies)

}

func AddPolicy(policyName string, perms []int) error {
	db := DBSystem

	totalPerms := ""
	for i := 0; i < len(perms); i++ {
		perm := strconv.Itoa(perms[i])
		totalPerms += perm + ";"
	}

	newPerms := Policies{Name: policyName, Permissions: totalPerms}
	result := db.Create(&newPerms)
	if result.Error != nil {
		combinedLogger.Warn("Adding policy " + policyName + " to db " + result.Error.Error())
		return result.Error
	}
	return nil
}

func AddGroup(groupName string) error {
	db := DBSystem

	newGroup := Groups{Name: groupName}
	resultGroup := db.Create(&newGroup)
	if resultGroup.Error != nil {
		combinedLogger.Warn("Adding group " + groupName + " to db " + resultGroup.Error.Error())
		return resultGroup.Error
	}
	return nil
}

func AddGroupPolicies(groupName string, policyNames []string) error {
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
