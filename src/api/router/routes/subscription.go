package routes

import (
	"net/http"
	"strconv"

	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/api/router/responses"
	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/db"
)

type Subscription struct {
	SubscriptionID int    `json:"subscriptionID"`
	Name           string `json:"name"`
	NumberOfKeys   int    `json:"numberOfKeys"`
	TotalBandwidth int    `json:"totalBandwidth"`
}

func CreateSubscription(res http.ResponseWriter, req *http.Request) {
	var bodyReq Subscription
	bodyRes := responses.StandardResponse{}

	err := ParseRequest(req, &bodyReq)
	if err != nil {
		combinedLogger.Error("Parsing request " + err.Error())
		bodyRes.Response = "Error parsing request"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	if bodyReq.NumberOfKeys < 0 {
		bodyRes.Response = "numberOfKeys must be >= 0"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	} else if bodyReq.TotalBandwidth < 0 {
		bodyRes.Response = "subscriptionID must be >= 0"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	//check perms
	bearerToken := req.Header.Get("Bearer")

	adminPerms := []int{db.SUBSCRIPTION_MODIFY_ALL}
	_, validAdminErr := db.ValidatePerms(bearerToken, adminPerms)

	if validAdminErr != nil {
		bodyRes.Response = "user does not have permission or an error occurred"
		responses.Standard(res, bodyRes, http.StatusForbidden)
		return
	} else if validAdminErr == nil { //If request is from admin with perms
		adminSubErr := db.CreateSubscription(bodyReq.Name, bodyReq.NumberOfKeys, bodyReq.TotalBandwidth)
		if adminSubErr != nil {
			bodyRes.Response = adminSubErr.Error()
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

func ReadSubscription(res http.ResponseWriter, req *http.Request) {
	bodyRes := responses.SubscriptionResponse{}
	queryVars := req.URL.Query()

	subscriptionIDStr := queryVars.Get("id")
	subscriptionName := queryVars.Get("name")

	// No need to check permission

	if subscriptionIDStr != "" {
		subscriptionID, convErr := strconv.Atoi(subscriptionIDStr)
		if convErr != nil {
			bodyRes.Response = "could not convert id to integer"
			responses.Subscription(res, bodyRes, http.StatusBadRequest)
			return
		}
		sub, subErr := db.ReadSubscription(subscriptionID)
		if subErr != nil {
			bodyRes.Response = subErr.Error()
			responses.Subscription(res, bodyRes, http.StatusBadRequest)
			return
		}
		bodyRes.Subscription = sub
		bodyRes.Response = "pulled subscription successfully"
	} else if subscriptionName != "" {
		sub, subErr := db.ReadSubscriptionByName(subscriptionName)
		if subErr != nil {
			bodyRes.Response = subErr.Error()
			responses.Subscription(res, bodyRes, http.StatusBadRequest)
			return
		}
		bodyRes.Subscription = sub
		bodyRes.Response = "pulled subscription successfully"
		responses.Subscription(res, bodyRes, http.StatusBadRequest)
		return
	} else {
		bodyRes.Response = "id or name query must be filled"
		responses.Subscription(res, bodyRes, http.StatusBadRequest)
		return
	}
}

func ReadAllSubscriptions(res http.ResponseWriter, req *http.Request) {
	bodyRes := responses.DumpSubscriptionResponse{}
	// No need to check permission

	subs, subErr := db.ReadAllSubscriptions()
	if subErr != nil {
		bodyRes.Response = subErr.Error()
		responses.DumpSubscriptions(res, bodyRes, http.StatusBadRequest)
		return
	}
	bodyRes.Subscriptions = subs
	bodyRes.Response = "pulled subscriptions successfully"
	responses.DumpSubscriptions(res, bodyRes, http.StatusAccepted)
}

func UpdateSubscription(res http.ResponseWriter, req *http.Request) {
	bodyRes := responses.StandardResponse{}
	queryVars := req.URL.Query()

	subscriptionName := queryVars.Get("name")
	totalBandwidth := queryVars.Get("totalBandwidth")
	numberOfKeys := queryVars.Get("numberOfKeys")

	if subscriptionName == "" {
		bodyRes.Response = "name must be filled"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	} else if totalBandwidth == "" && numberOfKeys == "" {
		bodyRes.Response = "totalBandwidth or numberOfKeys needs to be filled"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	//check perms
	bearerToken := req.Header.Get("Bearer")

	adminPerms := []int{db.SUBSCRIPTION_MODIFY_ALL}
	_, validAdminErr := db.ValidatePerms(bearerToken, adminPerms)

	if validAdminErr != nil {
		bodyRes.Response = "user does not have permission or an error occurred"
		responses.Standard(res, bodyRes, http.StatusForbidden)
		return
	}

	numKeys := -1
	totalBW := -1

	if numberOfKeys != "" {
		numKeysInt, convErr := strconv.Atoi(numberOfKeys)
		if convErr != nil {
			bodyRes.Response = "numberOfKeys could not be converted to integer"
			responses.Standard(res, bodyRes, http.StatusBadRequest)
			return
		}
		if numKeysInt < 0 {
			bodyRes.Response = "numberOfKeys needs to be >= 0"
			responses.Standard(res, bodyRes, http.StatusBadRequest)
			return
		}
		numKeys = numKeysInt
	}
	if totalBandwidth != "" {
		totalBWInt, convErr := strconv.Atoi(totalBandwidth)
		if convErr != nil {
			bodyRes.Response = "totalBandwidth could not be converted to integer"
			responses.Standard(res, bodyRes, http.StatusBadRequest)
			return
		}
		if totalBWInt < 0 {
			bodyRes.Response = "totalBandwidth needs to be >= 0"
			responses.Standard(res, bodyRes, http.StatusBadRequest)
			return
		}
		totalBW = totalBWInt
	}
	updateErr := db.UpdateSubscription(subscriptionName, numKeys, totalBW)
	if updateErr != nil {
		bodyRes.Response = updateErr.Error()
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}
	bodyRes.Response = "Updated subscription successfully"
	responses.Standard(res, bodyRes, http.StatusBadRequest)
}

func DeleteSubscription(res http.ResponseWriter, req *http.Request) {
	queryVars := req.URL.Query()

	subscriptionIDStr := queryVars.Get("id")
	subscriptionName := queryVars.Get("name")
	bodyRes := responses.StandardResponse{}

	//check perms
	bearerToken := req.Header.Get("Bearer")

	adminPerms := []int{db.SUBSCRIPTION_MODIFY_ALL}
	_, validAdminErr := db.ValidatePerms(bearerToken, adminPerms)

	if validAdminErr != nil {
		bodyRes.Response = "user does not have permission or an error occurred"
		responses.Standard(res, bodyRes, http.StatusForbidden)
		return
	}

	if subscriptionIDStr != "" {
		subscriptionID, convErr := strconv.Atoi(subscriptionIDStr)
		if convErr != nil {
			bodyRes.Response = "could not convert id to integer"
			responses.Standard(res, bodyRes, http.StatusBadRequest)
			return
		}
		subErr := db.DeleteSubscription(subscriptionID)
		if subErr != nil {
			bodyRes.Response = subErr.Error()
			responses.Standard(res, bodyRes, http.StatusBadRequest)
			return
		}
		bodyRes.Response = "deleted subscription successfully"
	} else if subscriptionName != "" {
		subErr := db.DeleteSubscriptionByName(subscriptionName)
		if subErr != nil {
			bodyRes.Response = subErr.Error()
			responses.Standard(res, bodyRes, http.StatusBadRequest)
			return
		}
		bodyRes.Response = "deleted subscription successfully"
		responses.Standard(res, bodyRes, http.StatusAccepted)
		return
	} else {
		bodyRes.Response = "id or name query must be filled"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}
}
