package db

import (
	"errors"

	"gorm.io/gorm"
)

// creates a new wireguard interface
func createWireguardInterface(listenPort int, publicKey string, ipv4Address string, ipv6Address string) (id int, err error) {
	db := DBSystem

	newInterface := WireguardInterfaces{ListenPort: listenPort, PublicKey: publicKey, IPv4Address: ipv4Address, IPv6Address: ipv6Address}
	createInterface := db.Create(&newInterface)

	if createInterface.Error != nil {
		err = ErrCreatingWireguardInterface
	}
	id = newInterface.ID
	return
}

func readWireguardInterface(wireguardInterfaceID int) (wgInterface WireguardInterfaces, err error) {
	db := DBSystem

	findInterface := db.Where("id = ?", wireguardInterfaceID).First(&wgInterface)

	if errors.Is(findInterface.Error, gorm.ErrRecordNotFound) {
		err = ErrFindingWireguardInterface
	} else if findInterface.Error != nil {
		err = ErrQuery
	}
	return
}

func ReadWireguardInstanceFromServerID(serverID int) (wgInterface WireguardInterfaces, err error) {
	db := DBSystem

	serverInterface, err := ReadServerInterfaceFromServerID(serverID)
	if err != nil {
		return
	}

	findInterface := db.Where("id = ?", serverInterface.InterfaceID).First(&wgInterface)

	if errors.Is(findInterface.Error, gorm.ErrRecordNotFound) {
		err = ErrFindingWireguardInterface
	} else if findInterface.Error != nil {
		err = ErrQuery
	}
	return
}

func deleteWireguardInterface(wireguardInterfaceID int) (err error) {
	db := DBSystem

	wgInterface, err := readWireguardInterface(wireguardInterfaceID)
	if err != nil {
		return err
	}

	deleteInfo := db.Delete(&wgInterface)
	if deleteInfo.Error != nil {
		return ErrDeletingWireguardInterface
	}
	return
}
