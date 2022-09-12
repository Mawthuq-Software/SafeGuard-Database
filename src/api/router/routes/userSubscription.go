package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/api/router/responses"
	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/db"
)

type UserSubscription struct {
	UserSubscriptionID int `json:"userSubscriptionID"`
}

type EditUserSubscription struct {
	UserSubscription
	Subscription
	UsedBandwidth int    `json:"usedBandwidth"`
	Expiry        string `json:"expiry"`
}

type AddUserSubscription struct {
	User
	Subscription
}

// CREATE

//Ties a user to a subscription.
func CreateUserSubscription(res http.ResponseWriter, req *http.Request) {
	var bodyReq AddUserSubscription
	bodyRes := responses.StandardResponse{}

	err := ParseRequest(req, &bodyReq)
	if err != nil {
		combinedLogger.Error("Parsing request " + err.Error())
		bodyRes.Response = "Error parsing request"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	if bodyReq.UserID <= 0 {
		bodyRes.Response = "userID must be > 0"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	} else if bodyReq.SubscriptionID <= 0 {
		bodyRes.Response = "subscriptionID must be > 0"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	//check perms
	bearerToken := req.Header.Get("Bearer")

	userPerms := []int{db.PERSONAL_USER_SUBSCRIPTION_ADD}
	userID, validUserErr := db.ValidatePerms(bearerToken, userPerms)

	adminPerms := []int{db.USER_SUBSCRIPTION_MODIFY_ALL}
	_, validAdminErr := db.ValidatePerms(bearerToken, adminPerms)

	if validUserErr != nil && validAdminErr != nil {
		bodyRes.Response = "user does not have permission or an error occurred"
		responses.Standard(res, bodyRes, http.StatusForbidden)
		return
	} else if validAdminErr == nil { //If request is from admin with perms
		expTime := time.Now().Add(time.Hour * 24000)
		userSubErr := db.CreateUserSubscription(bodyReq.UserID, bodyReq.SubscriptionID, expTime)
		if userSubErr != nil {
			bodyRes.Response = userSubErr.Error()
			responses.Standard(res, bodyRes, http.StatusBadRequest)
			return
		}
		bodyRes.Response = "added user subscription successfully"
		responses.Standard(res, bodyRes, http.StatusAccepted)
		return

	} else if validUserErr == nil { //If incoming request is not from admin, run this
		if userID != bodyReq.UserID {
			bodyRes.Response = "you do not have permission to do this action"
			responses.Standard(res, bodyRes, http.StatusForbidden)
			return
		}
		expTime := time.Now().Add(time.Hour * 24000)
		userSubErr := db.CreateUserSubscription(bodyReq.UserID, bodyReq.SubscriptionID, expTime)
		if userSubErr != nil {
			bodyRes.Response = userSubErr.Error()
			responses.Standard(res, bodyRes, http.StatusBadRequest)
			return
		}
		bodyRes.Response = "added user subscription successfully"
		responses.Standard(res, bodyRes, http.StatusAccepted)
		return

	} else {
		bodyRes.Response = "an error occurred"
		responses.Standard(res, bodyRes, http.StatusInternalServerError)
		return
	}
}

// READ

//Get a user's subscription
func ReadUserSubscription(res http.ResponseWriter, req *http.Request) {
	bodyRes := responses.UserSubscriptionResponse{}
	queryVars := req.URL.Query()

	userSubIDStr := queryVars.Get("id")

	//check perms
	bearerToken := req.Header.Get("Bearer")
	userPerms := []int{db.PERSONAL_USER_SUBSCRIPTION_VIEW}
	userID, validUserErr := db.ValidatePerms(bearerToken, userPerms)

	adminPerms := []int{db.USER_SUBSCRIPTION_VIEW_ALL}
	_, validAdminErr := db.ValidatePerms(bearerToken, adminPerms)

	if userSubIDStr != "" { //check param
		userSubscriptionID, convErr := strconv.Atoi(userSubIDStr) //convert param to int
		if convErr != nil {
			bodyRes.Response = "id was unable to be converted to int"
			responses.UserSubscription(res, bodyRes, http.StatusBadRequest)
			return
		} else if validUserErr != nil && validAdminErr != nil {
			fmt.Println(validUserErr.Error())
			bodyRes.Response = "user does not have permission or an error occurred"
			responses.UserSubscription(res, bodyRes, http.StatusForbidden)
			return
		} else if validUserErr == nil { //If incoming request is not from admin, run this
			validationErr := db.ValidateUsernameUserSubscription(userID, userSubscriptionID)
			if validationErr != nil {
				bodyRes.Response = validationErr.Error()
				responses.UserSubscription(res, bodyRes, http.StatusBadRequest)
				return
			}
		}
		//assume permission is admin or user is allowed to
		userSub, userSubErr := db.ReadUserSubscriptionFromID(userSubscriptionID)
		if userSubErr != nil {
			bodyRes.Response = userSubErr.Error()
			responses.UserSubscription(res, bodyRes, http.StatusBadRequest)
			return
		}
		bodyRes.UserSubscription = userSub
		bodyRes.Response = "pulled user subscription successfully"
		responses.UserSubscription(res, bodyRes, http.StatusAccepted)
		return
	} else { //param is empty
		bodyRes.Response = "id needs to be filled"
		responses.UserSubscription(res, bodyRes, http.StatusBadRequest)
	}
}

//Get a user's subscription from their user ID
func ReadUserSubscriptionFromUserID(res http.ResponseWriter, req *http.Request) {
	bodyRes := responses.UserSubscriptionResponse{}
	queryVars := req.URL.Query()

	userSubIDStr := queryVars.Get("id")

	//check perms
	bearerToken := req.Header.Get("Bearer")
	userPerms := []int{db.PERSONAL_USER_SUBSCRIPTION_VIEW}
	userID, validUserErr := db.ValidatePerms(bearerToken, userPerms)

	adminPerms := []int{db.USER_SUBSCRIPTION_VIEW_ALL}
	_, validAdminErr := db.ValidatePerms(bearerToken, adminPerms)

	if userSubIDStr != "" { //check param
		userIDInt, convErr := strconv.Atoi(userSubIDStr) //convert param to int
		if convErr != nil {
			bodyRes.Response = "id was unable to be converted to int"
			responses.UserSubscription(res, bodyRes, http.StatusBadRequest)
			return
		} else if validUserErr != nil && validAdminErr != nil {
			bodyRes.Response = "user does not have permission or an error occurred"
			responses.UserSubscription(res, bodyRes, http.StatusForbidden)
			return
		} else if validUserErr == nil { //If incoming request is not from admin, run this
			if userIDInt != userID {
				bodyRes.Response = "user does not have permission or an error occurred"
				responses.UserSubscription(res, bodyRes, http.StatusForbidden)
				return
			}
		}
		//assume admin is present or user has permission
		userSub, userSubErr := db.ReadUserSubscriptionFromUserID(userID)
		if userSubErr != nil {
			bodyRes.Response = userSubErr.Error()
			responses.UserSubscription(res, bodyRes, http.StatusBadRequest)
			return
		}
		bodyRes.UserSubscription = userSub
		bodyRes.Response = "pulled user subscription successfully"
		responses.UserSubscription(res, bodyRes, http.StatusAccepted)
	} else { //param is empty
		bodyRes.Response = "id needs to be filled"
		responses.UserSubscription(res, bodyRes, http.StatusBadRequest)
	}
}

//Get all users subscriptions
func ReadAllUserSubscriptions(res http.ResponseWriter, req *http.Request) {
	bodyRes := responses.DumpUserSubscriptionResponse{}
	bearerToken := req.Header.Get("Bearer")
	adminPerms := []int{db.USER_SUBSCRIPTION_VIEW_ALL}

	_, validAdminErr := db.ValidatePerms(bearerToken, adminPerms)
	if validAdminErr != nil {
		bodyRes.Response = "user does not have permission or an error occurred"
		responses.DumpUserSubscription(res, bodyRes, http.StatusBadRequest)
		return
	}

	userSubs, dumpErr := db.ReadAllUserSubscriptions()
	if dumpErr != nil {
		bodyRes.Response = dumpErr.Error()
		responses.DumpUserSubscription(res, bodyRes, http.StatusBadRequest)
		return
	}
	bodyRes.Response = "dumped database table successfully"
	bodyRes.UserSubscription = userSubs
	responses.DumpUserSubscription(res, bodyRes, http.StatusAccepted)
}

// UPDATE

//Edits a user's subscription
func UpdateUserSubscription(res http.ResponseWriter, req *http.Request) {
	bodyRes := responses.StandardResponse{}
	queryVars := req.URL.Query()

	userSubIDStr := queryVars.Get("id")
	usedBWStr := queryVars.Get("usedBandwidth")
	expiryStr := queryVars.Get("expiry")

	bearerToken := req.Header.Get("Bearer") // Bearer token

	//NEEDS TO BE REMOVED FOR PAYMENT!!!
	userPerms := []int{db.PERSONAL_USER_SUBSCRIPTION_MODIFY} //db.SUBSCRIPTION_MODIFY_ALL
	username, validUserErr := db.ValidatePerms(bearerToken, userPerms)
	// REMOVE FOR PAYMENT!! WILL CAUSE SECURITY ISSUE

	adminPerms := []int{db.USER_SUBSCRIPTION_MODIFY_ALL} //db.SUBSCRIPTION_MODIFY_ALL
	_, validAdminErr := db.ValidatePerms(bearerToken, adminPerms)

	if userSubIDStr == "" {
		bodyRes.Response = "id needs to be filled"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	} else if usedBWStr == "" && expiryStr == "" {
		bodyRes.Response = "usedBandwidth or expiry needs to be filled"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	userSubID, errConv := strconv.Atoi(userSubIDStr)
	if errConv != nil {
		bodyRes.Response = "could not convert id to integer"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	usedBW, errConv := strconv.Atoi(usedBWStr)
	if errConv != nil {
		bodyRes.Response = "could not convert usedBandwidth to integer"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	expiry, timeErr := time.Parse(time.RFC822, expiryStr)
	if timeErr != nil {
		bodyRes.Response = "expiry is not in correct time format"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	if errConv != nil {
		bodyRes.Response = "Could not convert usedBandwidth to integer"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	if validUserErr != nil && validAdminErr != nil {
		bodyRes.Response = "user does not have permission or an error occurred"
		responses.Standard(res, bodyRes, http.StatusForbidden)
		return
	} else if validUserErr == nil { //If incoming request is not from admin, run this
		validationErr := db.ValidateUsernameUserSubscription(username, userSubID)
		if validationErr != nil {
			bodyRes.Response = validationErr.Error()
			responses.Standard(res, bodyRes, http.StatusBadRequest)
			return
		}
	}
	// No need to check admin, it is assumed from logic they are allowed.
	errUpdate := db.UpdateUserSubscription(userSubID, usedBW, expiry)
	if errUpdate != nil {
		bodyRes.Response = errUpdate.Error()
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}
	bodyRes.Response = "updated userSubscription successfully"
	responses.Standard(res, bodyRes, http.StatusAccepted)
}

// DELETE

//Deletes a user's usersubscription
func DeleteUserSubscription(res http.ResponseWriter, req *http.Request) {
	bodyRes := responses.StandardResponse{}
	queryVars := req.URL.Query()

	userSubIDStr := queryVars.Get("id")

	bearerToken := req.Header.Get("Bearer") // Bearer token

	//NEED TO CANCEL SUBSCRIPTION PAYMENT
	userPerms := []int{db.PERSONAL_USER_SUBSCRIPTION_DELETE} //db.SUBSCRIPTION_MODIFY_ALL
	username, validUserErr := db.ValidatePerms(bearerToken, userPerms)

	adminPerms := []int{db.USER_SUBSCRIPTION_MODIFY_ALL} //db.SUBSCRIPTION_MODIFY_ALL
	_, validAdminErr := db.ValidatePerms(bearerToken, adminPerms)

	if userSubIDStr == "" {
		bodyRes.Response = "id needs to be filled"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	userSubID, errConv := strconv.Atoi(userSubIDStr)
	if errConv != nil {
		bodyRes.Response = "could not convert id to integer"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	if validUserErr != nil && validAdminErr != nil {
		bodyRes.Response = "user does not have permission or an error occurred"
		responses.Standard(res, bodyRes, http.StatusForbidden)
		return
	} else if validUserErr == nil { //If incoming request is not from admin, run this
		validationErr := db.ValidateUsernameUserSubscription(username, userSubID)
		if validationErr != nil {
			bodyRes.Response = validationErr.Error()
			responses.Standard(res, bodyRes, http.StatusBadRequest)
			return
		}
	}
	// No need to check admin, it is assumed from logic they are allowed.
	errDelete := db.DeleteUserSubscription(userSubID)
	if errDelete != nil {
		bodyRes.Response = errDelete.Error()
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}
	bodyRes.Response = "Deleted usersubscription successfully"
	responses.Standard(res, bodyRes, http.StatusAccepted)
}
