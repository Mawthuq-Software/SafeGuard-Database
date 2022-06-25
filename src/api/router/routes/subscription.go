package routes

import (
	"net/http"

	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/db"
)

func EditingSubscription(res http.ResponseWriter, req *http.Request) {
	perms := []int{db.PERSONAL_SUBSCRIPTION_MODIFY, db.SUBSCRIPTION_MODIFY_ALL}
	ValidatePerms(res, req, perms)
}

func GetSubscription(res http.ResponseWriter, req *http.Request) {
	perms := []int{db.PERSONAL_SUBSCRIPTION_VIEW, db.SUBSCRIPTION_VIEW_ALL}
	ValidatePerms(res, req, perms)
}

func GetAllSubscriptions(res http.ResponseWriter, req *http.Request) {
	perms := []int{db.SUBSCRIPTION_VIEW_ALL}
	ValidatePerms(res, req, perms)
}
