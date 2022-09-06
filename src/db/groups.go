package db

func CreateGroup(groupName string) error {
	db := DBSystem

	newGroup := Groups{Name: groupName}
	resultGroup := db.Create(&newGroup)
	if resultGroup.Error != nil {
		combinedLogger.Warn("Adding group " + groupName + " to db " + resultGroup.Error.Error())
		return resultGroup.Error
	}
	return nil
}
