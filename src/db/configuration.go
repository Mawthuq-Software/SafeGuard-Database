package db

import (
	"net"

	"gorm.io/gorm"
)

func CreateConfiguration(name string, dns string, mask int) (err error) {
	db := DBSystem
	ipAddr := net.ParseIP(dns)
	if ipAddr == nil {
		return ErrServerIPInvalid
	}
	if mask > 32 {
		return ErrConfMaskInvalid
	}

	conf := Configurations{Name: name, DNS: dns, Mask: mask}
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
