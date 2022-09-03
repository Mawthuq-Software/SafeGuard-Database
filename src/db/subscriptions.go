package db

import (
	"errors"

	"gorm.io/gorm"
)

// Gets a subscription from subscriptionID
func GetSubscription(subscriptionID int) (subscription Subscriptions, err error) {
	db := DBSystem

	findSub := db.Where("id = ?", subscriptionID).First(&subscription)
	if errors.Is(findSub.Error, gorm.ErrRecordNotFound) {
		err = ErrSubscriptionNotFound
		return
	} else if findSub.Error != nil {
		err = ErrQuery
	}
	return
}

func GetSubscriptionByName(subscriptionName string) (subscription Subscriptions, err error) {
	db := DBSystem

	findSub := db.Where("name = ?", subscriptionName).First(&subscription)
	if errors.Is(findSub.Error, gorm.ErrRecordNotFound) {
		err = ErrSubscriptionNotFound
		return
	} else if findSub.Error != nil {
		err = ErrQuery
	}
	return
}

func AddSubscription(name string, numKeys int, totalBandwidth int) (err error) {
	db := DBSystem
	_, err = GetSubscriptionByName(name)
	if err != ErrSubscriptionNotFound {
		newSub := Subscriptions{Name: name, NumberOfKeys: numKeys, TotalBandwidth: totalBandwidth}
		createErr := db.Create(&newSub)
		if createErr != nil {
			err = createErr.Error
			return
		}
		return nil
	}
	return
}
