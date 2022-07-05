package db

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

func getUserSubscriptionFromUserID(userID int) (userSubscription UserSubscriptions, err error) {
	db := DBSystem

	findUserSub := db.Where("user_id = ?", userID).First(&userSubscription)
	if !errors.Is(findUserSub.Error, gorm.ErrRecordNotFound) {
		err = ErrUserSubscriptionNotFound
	} else if findUserSub.Error != nil {
		err = ErrQuery
	}
	return
}

// Gets a subscription from subscriptionID
func getSubscription(subscriptionID int) (subscription Subscriptions, err error) {
	db := DBSystem

	findSub := db.Where("id = ?", subscriptionID).First(&subscription)
	if !errors.Is(findSub.Error, gorm.ErrRecordNotFound) {
		err = ErrSubscriptionNotFound
		return
	} else if findSub.Error != nil {
		err = ErrQuery
	}
	return
}

// Checks to see if a new key can be added
func checkSubscriptionKeyAddition(userID int) (err error) {
	userSubscription, err := getUserSubscriptionFromUserID(userID)
	if err != nil {
		return
	}

	if userSubscription.Expiry.After(time.Now()) {
		err = ErrSubscriptionExpired
		return
	}

	subscription, err := getSubscription(userSubscription.SubscriptionID)
	if err != nil {
		return
	}

	numKeys := subscription.NumberOfKeys
	userKeys, err := findUserKeys(userID)
	if err != nil {
		return
	}
	if len(userKeys) >= numKeys {
		err = ErrSubscriptionKeyMaxed
		return
	}
	return
}
