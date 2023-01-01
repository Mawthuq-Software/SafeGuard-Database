package db

import (
	"errors"
	"net"

	"gorm.io/gorm"
)

// CREATE

func CreateServer(name string, region string, country string, ipAddressStr string) (err error) {
	db := DBSystem
	if net.ParseIP(ipAddressStr) == nil {
		return ErrServerIPInvalid
	}
	_, findErr := ReadServerFromServerName(name)
	if findErr != ErrServerNotFound || findErr == nil {
		return ErrServerNameExist
	}

	server := Servers{Name: name, Region: region, Country: country, IPAddress: ipAddressStr}
	createErr := db.Create(&server)
	if createErr.Error != nil {
		return createErr.Error
	}
	return nil
}

// READ

func ReadServer(serverID int) (server Servers, err error) {
	db := DBSystem

	serverQuery := db.Where("id = ?", serverID).First(&server)
	if errors.Is(serverQuery.Error, gorm.ErrRecordNotFound) {
		err = ErrServerNotFound
		return
	} else if serverQuery.Error != nil {
		err = ErrQuery
		return
	}
	return
}

func ReadServerFromServerName(serverName string) (server Servers, err error) {
	db := DBSystem

	serverQuery := db.Where("name = ?", serverName).First(&server)
	if errors.Is(serverQuery.Error, gorm.ErrRecordNotFound) {
		err = ErrServerNotFound
		return
	} else if serverQuery.Error != nil {
		err = ErrQuery
		return
	}
	return
}

func ReadAllServers() (servers []Servers, err error) {
	db := DBSystem

	serverQuery := db.Find(&servers)
	if errors.Is(serverQuery.Error, gorm.ErrRecordNotFound) {
		err = ErrServerNotFound
		return
	} else if serverQuery.Error != nil {
		err = ErrQuery
		return
	}
	return
}

// UPDATE

func UpdateServer(server Servers) (err error) {
	db := DBSystem

	saveErr := db.Save(&server)
	if saveErr.Error != nil {
		return ErrServerSaving
	}
	return nil
}

func DeleteServer(serverID int) (err error) {
	db := DBSystem
	server, readErr := ReadServer(serverID)
	if readErr != nil {
		return readErr
	}

	keys, err := ReadKeysWithServerID(serverID)
	if err != nil && err != ErrKeyNotFound {
		return err
	} else if len(keys) > 0 {
		return ErrServerKeysExist
	}

	deleteErr := db.Delete(&server)
	if deleteErr.Error != nil {
		return ErrDeletingServer
	}
	return nil
}
