package routes

import (
	"net/http"

	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/api/router/responses"
	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/db"
)

func EditingSubscription(res http.ResponseWriter, req *http.Request) {
	bodyRes := responses.StandardResponse{}
	bearerToken := req.Header.Get("Bearer")
	perms := []int{db.PERSONAL_SUBSCRIPTION_MODIFY, db.SUBSCRIPTION_MODIFY_ALL}

	_, validErr := db.ValidateUserPerms(bearerToken, perms)
	if validErr != nil {
		bodyRes.Response = "user does not have permission or an error occurred"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
	}
	//ADD LOGIC
}

func GetSubscription(res http.ResponseWriter, req *http.Request) {
	bodyRes := responses.StandardResponse{}
	bearerToken := req.Header.Get("Bearer")
	perms := []int{db.PERSONAL_SUBSCRIPTION_VIEW, db.SUBSCRIPTION_VIEW_ALL}

	_, validErr := db.ValidateUserPerms(bearerToken, perms)
	if validErr != nil {
		bodyRes.Response = "user does not have permission or an error occurred"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
	}
	//ADD LOGIC
}

func GetAllSubscriptions(res http.ResponseWriter, req *http.Request) {
	bodyRes := responses.StandardResponse{}
	bearerToken := req.Header.Get("Bearer")
	perms := []int{db.SUBSCRIPTION_VIEW_ALL}

	_, validErr := db.ValidateUserPerms(bearerToken, perms)
	if validErr != nil {
		bodyRes.Response = "user does not have permission or an error occurred"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
	}
	//ADD LOGIC
}
