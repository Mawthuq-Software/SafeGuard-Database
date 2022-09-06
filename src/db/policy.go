package db

import "strconv"

func CreatePolicy(policyName string, perms []int) error {
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
