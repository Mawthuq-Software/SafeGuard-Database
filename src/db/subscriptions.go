package db

import (
	"errors"

	"gorm.io/gorm"
)

// CREATE

//Creates a new subscription
func CreateSubscription(name string, numKeys int, totalBandwidth int) (err error) {
	db := DBSystem
	_, err = ReadSubscriptionByName(name)
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

// READ

// Gets a subscription from subscriptionID
func ReadSubscription(subscriptionID int) (subscription Subscriptions, err error) {
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

func ReadSubscriptionByName(subscriptionName string) (subscription Subscriptions, err error) {
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

func ReadAllSubscriptions() (subscriptions []Subscriptions, err error) {
	db := DBSystem

	findSub := db.Find(&subscriptions)
	if errors.Is(findSub.Error, gorm.ErrRecordNotFound) {
		err = ErrSubscriptionNotFound
		return
	} else if findSub.Error != nil {
		err = ErrQuery
	}
	return
}

// UPDATE

//Updates a subscription. Use -1 for numKeys or totalBandwidth to keep current value.
func UpdateSubscription(name string, numKeys int, totalBandwidth int) (err error) {
	db := DBSystem
	subs, err := ReadSubscriptionByName(name)
	if err != nil {
		return
	}

	if numKeys > -1 {
		subs.NumberOfKeys = numKeys
	}
	if totalBandwidth > -1 {
		subs.TotalBandwidth = totalBandwidth
	}
	saveErrs := db.Save(&subs)
	if saveErrs.Error != nil {
		err = saveErrs.Error
		return
	}
	return nil
}

// DELETE

func DeleteSubscription(subscriptionID int) (err error) {
	db := DBSystem
	var userSubs []UserSubscriptions

	userSubs, err = ReadUserSubscriptionWithSubscriptionID(subscriptionID)
	if err != nil && err != ErrUserSubscriptionsNotFound { // check error is valid
		return err
	} else if len(userSubs) > 0 {
		return ErrSubscriptionUserSubExists
	}

	delErr := db.Delete(&userSubs)
	if delErr.Error != nil {
		err = delErr.Error
		return
	}
	return nil
}

func DeleteSubscriptionByName(name string) (err error) {
	subs, err := ReadSubscriptionByName(name)
	if err != nil {
		return
	}
	err = DeleteSubscription(subs.ID)
	return
}
