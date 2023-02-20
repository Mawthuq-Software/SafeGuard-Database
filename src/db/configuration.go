package db

import (
	"net"

	"gorm.io/gorm"
)

func CreateConfiguration(name string, dns string, mask int, numOfKeys int) (err error) {
	db := DBSystem
	confErr := CheckConfig(dns, mask)
	if confErr != nil {
		return confErr
	}

	conf := Configurations{Name: name, DNS: dns, Mask: mask, NumberOfKeys: numOfKeys}
	createErr := db.Create(&conf)
	if createErr.Error != nil {
		return ErrCreatingConf
	}
	return nil
}

func ReadConfiguration(confID int) (conf Configurations, err error) {
	db := DBSystem
	findConf := db.Where("id = ?", confID).First(&conf)
	if findConf.Error == gorm.ErrRecordNotFound {
		err = ErrConfNotFound
		return
	}
	return
}

func ReadConfigurationFromName(name string) (conf Configurations, err error) {
	db := DBSystem
	findConf := db.Where("name = ?", name).First(&conf)
	if findConf.Error == gorm.ErrRecordNotFound {
		err = ErrConfNotFound
		return
	}
	return
}

func ReadConfigurationFromServerID(serverID int) (conf Configurations, err error) {
	serverConfig, err := ReadServerConfigFromServerID(serverID)
	if err != nil {
		return
	}

	conf, err = ReadConfiguration(serverConfig.ConfigID)
	return
}

func ReadAllConfigurations() (configs []Configurations, err error) {
	db := DBSystem
	findConf := db.Find(&configs)
	if findConf.Error == gorm.ErrRecordNotFound {
		err = ErrConfNotFound
		return
	}
	return
}

func UpdateConfiguration(conf Configurations) (err error) {
	db := DBSystem

	saveErr := db.Save(&conf)
	if saveErr.Error != nil {
		return saveErr.Error
	}
	return
}

func DeleteConfiguration(confID int) (err error) {
	db := DBSystem

	conf, confErr := ReadConfiguration(confID)
	if confErr != nil {
		return confErr
	}

	deleteErr := db.Delete(&conf)
	if deleteErr.Error != nil {
		return deleteErr.Error
	}

	return
}

func CheckConfig(dns string, mask int) (err error) {
	if dns != "" {
		ipAddr := net.ParseIP(dns)
		if ipAddr == nil {
			return ErrConfDNSInvalid
		}
	}
	if mask != -1 && mask > 32 {
		return ErrConfMaskInvalid
	}
	return
}
