package db

import (
	"errors"

	"gorm.io/gorm"
)

func CreateServerInterfaces(serverID int, wireguardInterfaceID int) (err error) {
	db := DBSystem

	newServerInterface := ServerInterfaces{ServerID: serverID, InterfaceID: wireguardInterfaceID}
	createInfo := db.Create(&newServerInterface)
	if createInfo.Error != nil {
		return ErrCreatingServerInterface
	}

	return
}

func FindServerInterface(serverInterfaceID int) (serverInterface ServerInterfaces, err error) {
	db := DBSystem

	findInterface := db.Where("id = ?", serverInterfaceID).First(&serverInterface)
	if errors.Is(findInterface.Error, gorm.ErrRecordNotFound) {
		err = ErrFindingServerInterface
	} else if findInterface.Error != nil {
		err = ErrQuery
	}
	return
}

func DeleteServerInterface(serverInterfaceID int) (err error) {
	// make checks before deleting the interface
	serverInterface, err := FindServerInterface(serverInterfaceID)
	if err != nil {
		return err
	}

	keys, err := readKeysWithServerID(serverInterface.ID)
	if err != nil {
		return err
	}
	if len(keys) < 0 {
		return ErrKeysExistOnServerInterface
	}

	err = deleteServerInterfaceLink(serverInterfaceID)
	return
}

func deleteServerInterfaceLink(serverInterfaceID int) (err error) {
	db := DBSystem

	findInterface, err := FindServerInterface(serverInterfaceID)
	if err != nil {
		return
	}

	deleteInfo := db.Delete(&findInterface)
	if deleteInfo.Error != nil {
		return ErrDeletingServerInterface
	}

	return
}
