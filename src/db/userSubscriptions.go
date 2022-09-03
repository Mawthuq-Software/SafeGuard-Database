package db

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func GetUserSubscriptionFromID(userSubID int) (userSubscription UserSubscriptions, err error) {
	db := DBSystem

	findUserSub := db.Where("id = ?", userSubID).First(&userSubscription)
	if errors.Is(findUserSub.Error, gorm.ErrRecordNotFound) {
		err = ErrUserSubscriptionNotFound
	} else if findUserSub.Error != nil {
		err = ErrQuery
	}
	return
}

func GetUserSubscriptionFromUserID(userID int) (userSubscription UserSubscriptions, err error) {
	db := DBSystem

	findUserSub := db.Where("user_id = ?", userID).First(&userSubscription)
	if errors.Is(findUserSub.Error, gorm.ErrRecordNotFound) {
		err = ErrUserSubscriptionNotFound
	} else if findUserSub.Error != nil {
		err = ErrQuery
	}
	return
}

// Gets all user subscriptions from the database
func GetAllUserSubscriptions() (userSubs []UserSubscriptions, err error) {
	db := DBSystem

	findUserSubs := db.Find(&userSubs)
	if errors.Is(findUserSubs.Error, gorm.ErrRecordNotFound) {
		err = ErrUserSubscriptionsNotFound
	} else if findUserSubs.Error != nil {
		err = ErrQuery
	}
	return
}

func AddUserSubscription(userID int, subscriptionID int, expiryTime time.Time) (err error) {
	var userSubFind UserSubscriptions
	db := DBSystem

	findUserSub := db.Where("userID = ?", userID).First(&userSubFind)
	if !errors.Is(findUserSub.Error, gorm.ErrRecordNotFound) {
		return ErrUserSubscriptionExists
	}

	userSub := UserSubscriptions{UserID: userID, SubscriptionID: subscriptionID, UsedBandwidth: 0, Expiry: expiryTime}
	creatUserSub := db.Create(&userSub)
	if creatUserSub.Error != nil {
		combinedLogger.Error("Adding user to db " + creatUserSub.Error.Error())
		return ErrCreatingUserSubscription
	}
	return nil
}

// Checks to see if a new key can be added
func checkSubscriptionKeyAddition(userID int) (err error) {
	userSubscription, err := GetUserSubscriptionFromUserID(userID)
	if err != nil {
		fmt.Println(err)

		return
	}
	if userSubscription.Expiry.After(time.Now()) {
		err = ErrSubscriptionExpired
		return
	}

	subscription, err := GetSubscription(userSubscription.SubscriptionID)
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

func ValidateUsernameUserSubscription(userID int, userSubID int) (err error) {
	findUser, findUserErr := FindUserFromUserID(userID)
	if findUserErr != nil {
		return findUserErr
	}
	userSub, userSubErr := GetUserSubscriptionFromID(userSubID)
	if userSubErr != nil {
		return userSubErr
	}

	if userSub.UserID != findUser.ID {
		return ErrUserSubscriptionValidation
	} else {
		return nil
	}
}
