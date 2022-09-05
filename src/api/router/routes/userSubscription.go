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

//Edits a user's subscription
// func EditingUserSubscription(res http.ResponseWriter, req *http.Request) {
// 	//add json

// 	bodyRes := responses.StandardResponse{}
// 	bearerToken := req.Header.Get("Bearer")
// 	userPerms := []int{db.PERSONAL_SUBSCRIPTION_MODIFY} //db.SUBSCRIPTION_MODIFY_ALL
// 	username, validUserErr := db.ValidatePerms(bearerToken, userPerms)

// 	adminPerms := []int{db.SUBSCRIPTION_MODIFY_ALL} //db.SUBSCRIPTION_MODIFY_ALL
// 	_, validAdminErr := db.ValidatePerms(bearerToken, adminPerms)

// 	if validUserErr != nil && validAdminErr != nil {
// 		bodyRes.Response = "user does not have permission or an error occurred"
// 		responses.Standard(res, bodyRes, http.StatusBadRequest)
// 	} else if validAdminErr == nil {
// 		// Do logic ok

// 	} else if validUserErr == nil { //If incoming request is not from admin, run this
// 		validationErr := db.ValidateUsernameUserSubscription(username, UserSubscriptionID)
// 		if validationErr != nil {
// 			bodyRes.Response = validationErr.Error()
// 			responses.Standard(res, bodyRes, http.StatusBadRequest)
// 			return
// 		}
// 	}
// }

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

	userPerms := []int{db.PERSONAL_SUBSCRIPTION_MODIFY}
	userID, validUserErr := db.ValidatePerms(bearerToken, userPerms)

	adminPerms := []int{db.USER_SUBSCRIPTION_MODIFY_ALL}
	_, validAdminErr := db.ValidatePerms(bearerToken, adminPerms)

	if validUserErr != nil && validAdminErr != nil {
		bodyRes.Response = "user does not have permission or an error occurred"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
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

//Get a user's subscription
func GetUserSubscription(res http.ResponseWriter, req *http.Request) {
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
		}
		if validUserErr != nil && validAdminErr != nil {
			fmt.Println(validUserErr.Error())
			bodyRes.Response = "user does not have permission or an error occurred"
			responses.UserSubscription(res, bodyRes, http.StatusBadRequest)
		} else if validAdminErr == nil { //If request is from admin with perms
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
		} else if validUserErr == nil { //If incoming request is not from admin, run this
			validationErr := db.ValidateUsernameUserSubscription(userID, userSubscriptionID)
			if validationErr != nil {
				bodyRes.Response = validationErr.Error()
				responses.UserSubscription(res, bodyRes, http.StatusBadRequest)
				return
			}
		} else {
			bodyRes.Response = "an error occurred"
			responses.UserSubscription(res, bodyRes, http.StatusInternalServerError)
			return
		}
	} else { //param is empty
		bodyRes.Response = "id needs to be filled"
		responses.UserSubscription(res, bodyRes, http.StatusBadRequest)
	}
}

//Get all users subscriptions
func GetAllUserSubscriptions(res http.ResponseWriter, req *http.Request) {
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
