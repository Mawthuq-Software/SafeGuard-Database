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
	SERVER_VIEW                       int = 1012

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
	SERVER_ADD_ALL               int = 2013
	SERVER_MODIFY_ALL            int = 2014
	CONFIGURATION_ADD            int = 2020
	CONFIGURATION_MODIFY         int = 2021
	CONFIGURATION_READ           int = 2022
	CONFIGURATION_DELETE         int = 2023
	SERVER_CONFIGURATION_ADD     int = 2030
	SERVER_CONFIGURATION_READ    int = 2031
	SERVER_CONFIGURATION_MODIFY  int = 2032
	SERVER_CONFIGURATION_DELETE  int = 2033
	TOKEN_ADD                    int = 2040
	TOKEN_READ                   int = 2041
	TOKEN_MODIFY                 int = 2042
	TOKEN_DELETE                 int = 2043
	SERVER_TOKEN_ADD             int = 2050
	SERVER_TOKEN_READ            int = 2051
	SERVER_TOKEN_MODIFY          int = 2052
	SERVER_TOKEN_DELETE          int = 2053
	WIREGUARD_INSTANCE_CREATE    int = 2060
	WIREGUARD_INSTANCE_READ      int = 2061
	WIREGUARD_INSTANCE_MODIFY    int = 2062
	WIREGUARD_INSTANCE_DELETE    int = 2063
)

func startupCreation() {

	standardVPNPerms := []int{PERSONAL_KEYS_VIEW, PERSONAL_KEYS_ADD, PERSONAL_PASSWORD_RESET, PERSONAL_LOGIN, PERSONAL_USER_SUBSCRIPTION_VIEW,
		PERSONAL_USER_SUBSCRIPTION_ADD, PERSONAL_USER_SUBSCRIPTION_DELETE, PERSONAL_USER_SUBSCRIPTION_MODIFY, PERSONAL_SUBSCRIPTION_VIEW, SERVER_VIEW}
	CreatePolicy("STANDARD_USER_VPN", standardVPNPerms)

	userSettingPerms := []int{PERSONAL_PASSWORD_RESET, PERSONAL_LOGIN}
	CreatePolicy("STANDARD_USER_SETTINGS", userSettingPerms)

	advancedUserVPNPerms := []int{PERSONAL_KEYS_MODIFY, PERSONAL_KEYS_DELETE}
	CreatePolicy("ADVANCED_USER_VPN", advancedUserVPNPerms)

	adminPerms := []int{KEYS_VIEW_ALL, KEYS_ADD_ALL, KEYS_ADD_OVERRIDE, KEYS_DELETE_ALL, KEYS_MODIFY_ALL, PASSWORD_RESET_ALL, SUBSCRIPTION_MODIFY_ALL,
		SUBSCRIPTION_VIEW_ALL, USER_SUBSCRIPTION_MODIFY_ALL, USER_SUBSCRIPTION_VIEW_ALL, SERVER_ADD_ALL, SERVER_MODIFY_ALL, CONFIGURATION_ADD,
		CONFIGURATION_MODIFY, CONFIGURATION_READ, CONFIGURATION_DELETE, SERVER_CONFIGURATION_ADD, SERVER_CONFIGURATION_MODIFY, SERVER_CONFIGURATION_READ,
		SERVER_CONFIGURATION_DELETE, TOKEN_ADD, TOKEN_READ, TOKEN_MODIFY, TOKEN_DELETE, SERVER_TOKEN_ADD, SERVER_TOKEN_READ, SERVER_TOKEN_MODIFY,
		SERVER_TOKEN_DELETE, WIREGUARD_INSTANCE_CREATE, WIREGUARD_INSTANCE_READ, WIREGUARD_INSTANCE_MODIFY, WIREGUARD_INSTANCE_DELETE}
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
