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
