package db

import (
	"errors"

	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/network"
	"gorm.io/gorm"
)

func CreateServerWireguardInterface(serverID int, listenPort int, publicKey string, ipv4Address string, ipv6Address string) (err error) {
	_, err = ReadServer(serverID)
	if err != nil {
		return
	}

	//need to check ips and public key
	err = checkKeyValidity(publicKey)
	if err != nil {
		return
	}
	//check IPv4
	err = checkIPValidty(ipv4Address)
	if err != nil && ipAddrType(ipv4Address) != IPv4 {
		return ErrIPv4AddressNotValid
	}
	err = checkIPValidty(ipv6Address)
	if err != nil && ipAddrType(ipv6Address) != IPv6 {
		return ErrIPv6AddressNotValid
	}

	instanceID, err := createWireguardInterface(listenPort, publicKey, ipv4Address, ipv6Address)
	if err != nil {
		return
	}
	err = createServerInterfaceLink(serverID, instanceID)
	// final line so no need for error check
	return
}

func createServerInterfaceLink(serverID int, wireguardInterfaceID int) (err error) {
	db := DBSystem

	newServerInterface := ServerInterfaces{ServerID: serverID, InterfaceID: wireguardInterfaceID}
	createInfo := db.Create(&newServerInterface)
	if createInfo.Error != nil {
		return ErrCreatingServerInterface
	}
	return
}

func ReadServerInterface(serverInterfaceID int) (serverInterface ServerInterfaces, err error) {
	db := DBSystem

	findInterface := db.Where("id = ?", serverInterfaceID).First(&serverInterface)
	if errors.Is(findInterface.Error, gorm.ErrRecordNotFound) {
		err = ErrFindingServerInterface
	} else if findInterface.Error != nil {
		err = ErrQuery
	}
	return
}

func ReadServerInterfaceFromServerID(serverID int) (serverInterface ServerInterfaces, err error) {
	db := DBSystem

	findInterface := db.Where("server_id = ?", serverID).First(&serverInterface)
	if errors.Is(findInterface.Error, gorm.ErrRecordNotFound) {
		err = ErrFindingServerInterface
	} else if findInterface.Error != nil {
		err = ErrQuery
	}
	return
}

func DeleteServerInterface(serverID int) (err error) {
	// make checks before deleting the interface
	serverInterface, err := ReadServerInterfaceFromServerID(serverID)
	if err != nil {
		return err
	}

	keys, err := readKeysWithServerID(serverInterface.ID)
	if err != nil {
		return err
	}
	if len(keys) > 0 {
		return ErrKeysExistOnServerInterface
	}

	err = deleteServerInterfaceLink(serverInterface.ID)
	if err != nil {
		return
	}

	err = deleteWireguardInterface(serverInterface.InterfaceID)
	return
}

func deleteServerInterfaceLink(serverInterfaceID int) (err error) {
	db := DBSystem

	findInterface, err := ReadServerInterface(serverInterfaceID)
	if err != nil {
		return
	}

	deleteInfo := db.Delete(&findInterface)
	if deleteInfo.Error != nil {
		return ErrDeletingServerInterface
	}

	return
}

//MISC

func checkIPValidty(ipAddress string) (err error) {
	_, err = network.ParseIP(ipAddress)
	if err != nil {
		err = ErrPublicKeyIncorrectForm
	}
	return
}

type IPAddressType string

const (
	IPv4 IPAddressType = "v4"
	IPv6 IPAddressType = "v6"
)

func ipAddrType(ipAddress string) IPAddressType {
	for i := 0; i < len(ipAddress); i++ {
		switch ipAddress[i] {
		case '.':
			return IPv4
		case ':':
			return IPv6
		}
	}
	return ""
}
