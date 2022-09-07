package db

import (
	"errors"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

const (
	//1000 for user permissions
	//2000 for admin permissions
	//3000 for server permissions

	//User Permissions
	PERSONAL_KEYS_VIEW                int = 1000
	PERSONAL_KEYS_ADD                 int = 1001
	PERSONAL_PASSWORD_RESET           int = 1002
	PERSONAL_LOGIN                    int = 1003
	PERSONAL_USER_SUBSCRIPTION_VIEW   int = 1004
	PERSONAL_USER_SUBSCRIPTION_ADD    int = 1005
	PERSONAL_USER_SUBSCRIPTION_DELETE int = 1006
	PERSONAL_USER_SUBSCRIPTION_MODIFY int = 1007
	PERSONAL_SUBSCRIPTION_VIEW        int = 1008

	//Advanced User Perms
	PERSONAL_KEYS_DELETE         int = 1100
	PERSONAL_KEYS_MODIFY         int = 1101
	PERSONAL_SUBSCRIPTION_MODIFY int = 1102

	//Admin User Permissions
	KEYS_VIEW_ALL                int = 2000
	KEYS_ADD_ALL                 int = 2001
	KEYS_ADD_OVERRIDE            int = 2002
	KEYS_DELETE_ALL              int = 2003
	KEYS_MODIFY_ALL              int = 2004
	PASSWORD_RESET_ALL           int = 2005
	SUBSCRIPTION_MODIFY_ALL      int = 2007
	SUBSCRIPTION_VIEW_ALL        int = 2008
	USER_SUBSCRIPTION_MODIFY_ALL int = 2009
	USER_SUBSCRIPTION_VIEW_ALL   int = 2010
)

func startupCreation() {

	standardVPNPerms := []int{PERSONAL_KEYS_VIEW, PERSONAL_KEYS_ADD, PERSONAL_PASSWORD_RESET, PERSONAL_LOGIN, PERSONAL_USER_SUBSCRIPTION_VIEW,
		PERSONAL_USER_SUBSCRIPTION_ADD, PERSONAL_USER_SUBSCRIPTION_DELETE, PERSONAL_USER_SUBSCRIPTION_MODIFY, PERSONAL_SUBSCRIPTION_VIEW}
	CreatePolicy("STANDARD_USER_VPN", standardVPNPerms)

	userSettingPerms := []int{PERSONAL_PASSWORD_RESET, PERSONAL_LOGIN}
	CreatePolicy("STANDARD_USER_SETTINGS", userSettingPerms)

	advancedUserVPNPerms := []int{PERSONAL_KEYS_MODIFY, PERSONAL_KEYS_DELETE}
	CreatePolicy("ADVANCED_USER_VPN", advancedUserVPNPerms)

	adminPerms := []int{KEYS_VIEW_ALL, KEYS_ADD_ALL, KEYS_ADD_OVERRIDE, KEYS_DELETE_ALL, KEYS_MODIFY_ALL, PASSWORD_RESET_ALL}
	CreatePolicy("ADMIN_USER", adminPerms)

	CreateGroup("User")
	userPolicies := []string{"STANDARD_USER_VPN", "STANDARD_USER_SETTINGS"}
	CreateGroupPolicies("User", userPolicies)

	CreateGroup("Admin")
	adminPolicies := []string{"STANDARD_USER_VPN", "STANDARD_USER_SETTINGS", "ADVANCED_USER_VPN", "ADMIN_USER"}
	CreateGroupPolicies("Admin", adminPolicies)
}

func CheckPermission(userID int, permID int) (bool, error) {
	//REWORK THIS TRASH FUNCTION
	db := DBSystem
	var findUser Users

	errUserQuery := db.Where("id = ?", userID).First(&findUser)
	if errors.Is(errUserQuery.Error, gorm.ErrRecordNotFound) {
		return false, ErrUserNotFound
	} else if errUserQuery.Error != nil {
		return false, ErrQuery
	}

	userGroups, userGroupErr := ReadAllUserGroups(userID)
	if userGroupErr != nil {
		return false, userGroupErr
	}

	for x := 0; x < len(userGroups); x++ {
		groupID := userGroups[x].GroupID
		groupPol, errGroupPol := ReadGroupPolicies(groupID)
		if errGroupPol != nil {
			return false, errGroupPol
		}

		for y := 0; y < len(groupPol); y++ {
			policyID := groupPol[y].PolicyID
			permsArr, errPerms := ReadPermissions(policyID)
			if errPerms != nil {
				return false, errPerms
			}
			for z := 0; z < len(permsArr); z++ {
				perm, convErr := strconv.Atoi(permsArr[z])
				if convErr != nil {
					return false, convErr
				}
				if perm == permID {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

func ReadPermissions(policyID int) ([]string, error) {
	db := DBSystem
	var findPolicy Policies

	errPolicyQuery := db.Where("id = ?", policyID).First(&findPolicy)
	if errors.Is(errPolicyQuery.Error, gorm.ErrRecordNotFound) {
		return nil, ErrPolicyNotFound
	} else if errPolicyQuery.Error != nil {
		return nil, ErrQuery
	}
	permsStr := findPolicy.Permissions
	perms := strings.Split(permsStr, ";")
	perms = perms[0:(len(perms) - 1)]
	return perms, nil
}
