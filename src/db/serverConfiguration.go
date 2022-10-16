package db

import (
	"errors"

	"gorm.io/gorm"
)

func CreateServerConfig(serverID int, configID int) (err error) {
	db := DBSystem

	_, confErr := ReadConfiguration(configID)
	if confErr != nil {
		return confErr
	}

	_, serverErr := ReadServer(serverID)
	if serverErr != nil {
		return serverErr
	}

	_, serverConfErr := ReadServerConfigFromServerID(serverID)
	if serverConfErr != ErrServerConfNotFound {
		return serverConfErr
	} else if serverConfErr != nil {
		return ErrServerConfAlreadyExists
	}
	newServerConf := ServerConfigurations{ServerID: serverID, ConfigID: configID}
	createErr := db.Create(&newServerConf)
	if createErr.Error != nil {
		return createErr.Error
	}

	return
}

func ReadServerConfig(serverConfigID int) (serverConf ServerConfigurations, err error) {
	db := DBSystem

	findErr := db.Where("id = ?", serverConfigID).First(&serverConf)
	if errors.Is(findErr.Error, gorm.ErrRecordNotFound) {
		err = ErrServerConfNotFound
		return
	} else if findErr.Error != nil {
		err = ErrQuery
		return
	}
	return
}

func ReadServerConfigFromServerID(serverID int) (serverConf ServerConfigurations, err error) {
	db := DBSystem

	findErr := db.Where("server_id = ?", serverID).First(&serverConf)
	if errors.Is(findErr.Error, gorm.ErrRecordNotFound) {
		err = ErrServerConfNotFound
		return
	} else if findErr.Error != nil {
		err = ErrQuery
		return
	}
	return
}

func UpdateServerConfig(serverConf ServerConfigurations) (err error) {
	db := DBSystem

	saveErr := db.Save(&serverConf)
	if saveErr.Error != nil {
		err = saveErr.Error
	}
	return
}

func DeleteServerConfig(serverConfigID int) (err error) {
	db := DBSystem

	serverConf, errFind := ReadServerConfig(serverConfigID)
	if errFind != nil {
		return errFind
	}
	errDelete := db.Delete(&serverConf)
	if errDelete.Error != nil {
		return errDelete.Error
	}
	return
}
