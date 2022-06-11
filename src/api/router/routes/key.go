package routes

import (
	"net/http"

	"gitlab.com/mawthuq-software/wireguard-manager-authenticator/src/db"
)

func AddKey(res http.ResponseWriter, req *http.Request) {
	perms := []int{db.PERSONAL_KEYS_ADD, db.KEYS_ADD_ALL}
	ValidatePerms(res, req, perms)
}

func DeleteKey(res http.ResponseWriter, req *http.Request) {
	perms := []int{db.PERSONAL_KEYS_DELETE, db.KEYS_DELETE_ALL}
	ValidatePerms(res, req, perms)
}

func EnableDisableKey(res http.ResponseWriter, req *http.Request) {
	perms := []int{db.PERSONAL_KEYS_MODIFY, db.KEYS_MODIFY_ALL}
	ValidatePerms(res, req, perms)
}

func GetAllKeys(res http.ResponseWriter, req *http.Request) {
	perms := []int{db.KEYS_VIEW_ALL}
	ValidatePerms(res, req, perms)
}
