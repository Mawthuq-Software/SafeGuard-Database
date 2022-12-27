package db

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/wgmanager"
	"gorm.io/gorm"
)

// CREATE

// Adds a key to the database and returns the keyID
func createKey(serverID int, publicKey string, presharedKey string) (keyID int, err error) {
	db := DBSystem
	_, err = ReadServer(serverID)
	if err != nil {
		return
	}

	//check pub key
	err = checkKeyValidity(publicKey)
	if err != nil {
		return
	}

	//check pre key
	err = checkKeyValidity(presharedKey)
	if err != nil {
		return
	}

	//get server config
	configuration, err := ReadConfigurationFromServerID(serverID)
	if err != nil {
		return
	}

	//check keys on server
	currKeys, err := readKeysWithServerID(serverID)
	if err != nil {
		return
	}

	if len(currKeys) >= configuration.NumberOfKeys {
		err = ErrTooManyKeys
		return
	}

	octetsValid := false

	//generate main public subnet IP
	const CONFIGURATION_IP string = "10.0.0.1"
	confMask := configuration.Mask
	fullIP := CONFIGURATION_IP + "/" + strconv.Itoa(confMask)

	rand.Seed(time.Now().UnixNano())
	suitableIPs, err := Hosts(fullIP)
	var newIP string

	for i := 0; i < 10; i++ {
		index := rand.Intn(len(suitableIPs) - 1)
		newIP = suitableIPs[index]
		fmt.Println(index)
		_, err := readKeyFromServerIDAndIP(serverID, newIP)
		fmt.Println(err)
		if err == ErrKeyNotFound {
			octetsValid = true
			break
		}
	}

	if !octetsValid {
		err = ErrUnableToFindIP
		return
	}

	newKey := VPNKeys{ServerID: serverID, PublicKey: publicKey, PresharedKey: presharedKey, PrivateIPv4: newIP, PrivateIPv6: "0"}
	keyCreation := db.Create(&newKey)
	if keyCreation.Error != nil {
		err = ErrCreatingKey
		return
	}
	keyID = newKey.ID
	return
}

//Adds a user's key after checking their subscription validity
func CreateKeyAndLink(userID int, serverID int, publicKey string, presharedKey string) (err error) {
	err = checkSubscriptionKeyAddition(userID)
	if err != nil {
		return
	}

	keyID, err := createKey(serverID, publicKey, presharedKey)
	if err != nil {
		return
	}
	_, err = createUserKeyLink(userID, keyID)
	return
}

// READ

//finds a key object from a keyID
func readKey(keyID int) (key VPNKeys, err error) {
	db := DBSystem

	keyQuery := db.Where("id = ?", keyID).First(&key)
	if errors.Is(keyQuery.Error, gorm.ErrRecordNotFound) {
		err = ErrKeyNotFound
	} else if keyQuery.Error != nil {
		combinedLogger.Error("Finding key " + keyQuery.Error.Error())
		err = ErrQuery
	}
	return
}

//finds all the user's keys from their userID
func ReadUserKeys(userID int) (userKeys []UserKeys, err error) {
	db := DBSystem
	userKeysQuery := db.Where("user_id = ?", userID).Find(&userKeys)
	if errors.Is(userKeysQuery.Error, gorm.ErrRecordNotFound) {
		err = ErrKeyNotFound
		return
	} else if userKeysQuery.Error != nil {
		return
	}
	return
}

//gets all keys in database
func ReadAllKeys() (keys []VPNKeys, err error) {
	db := DBSystem

	dbResult := db.Find(&keys)
	err = dbResult.Error
	return keys, err
}

func readKeysWithServerID(serverID int) (keys []VPNKeys, err error) {
	db := DBSystem

	keyQuery := db.Where("server_id = ?", serverID).Find(&keys)
	if errors.Is(keyQuery.Error, gorm.ErrRecordNotFound) {
		err = ErrKeyNotFound
	} else if keyQuery.Error != nil {
		combinedLogger.Error("Finding key " + keyQuery.Error.Error())
		err = ErrQuery
	}
	return
}

func readKeyFromServerIDAndIP(serverID int, ipv4Address string) (key VPNKeys, err error) {
	db := DBSystem

	keyQuery := db.Where("server_id = ? AND private_ipv4 = ?", serverID, ipv4Address).First(&key)
	if errors.Is(keyQuery.Error, gorm.ErrRecordNotFound) {
		err = ErrKeyNotFound
	} else if keyQuery.Error != nil {
		combinedLogger.Error("Finding key " + keyQuery.Error.Error())
		err = ErrQuery
	}
	fmt.Println(key)
	return
}

//UPDATE

//updates a key object
func updateKey(key VPNKeys) (err error) {
	db := DBSystem

	err = db.Save(&key).Error
	return
}

// DELETE

//Deletes a key from keyID
func DeleteKey(keyID int) (err error) {
	db := DBSystem

	keyQuery, err := readKey(keyID)
	if err != nil {
		return
	}

	keyDelete := db.Delete(&keyQuery)
	if keyDelete.Error != nil {
		err = ErrDeletingKey
	}
	return
}

//Deletes a user's key and link
func DeleteKeyAndLink(keyID int) (err error) {
	err = deleteUserKeyLink(keyID)
	if err != nil {
		return
	}
	err = DeleteKey(keyID)
	return
}

//MISC

//Toggles a key usability from true to false and viceversa
func ToggleKey(keyID int) (err error) {
	key, err := readKey(keyID)
	if err != nil {
		return
	}

	key.Enabled = !key.Enabled
	err = updateKey(key)
	return
}

//checks to see if wireguard key is appropriate
func checkKeyValidity(key string) (err error) {
	_, err = wgmanager.ParseKey(key) //parse string
	if err != nil {
		err = ErrPublicKeyIncorrectForm
	}
	return
}

func checkIPInRangeValidty(mainIP string, ipToCheck string) (ipContained bool, err error) {
	_, mainIPSubnet, err := net.ParseCIDR(mainIP)
	if err != nil {
		return
	}
	ipToCheckParsed, _, err := net.ParseCIDR(ipToCheck)
	if err != nil {
		return
	}
	if mainIPSubnet.Contains(ipToCheckParsed) {
		ipContained = true
	} else {
		ipContained = false
	}
	return
}

//https://go.dev/play/p/fe-F2k6prlA
func Hosts(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}

	// remove network address and broadcast address
	lenIPs := len(ips)
	switch {
	case lenIPs < 2:
		return ips, nil

	default:
		return ips[1 : len(ips)-1], nil
	}
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
