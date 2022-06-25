package routes

import (
	"net/http"

	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/api/router/responses"
	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/db"
)

type Key struct {
	ID       int `json:"id"`
	ServerID int `json:"serverID"`
}

type KeyPreshared struct {
	PresharedKey string `json:"presharedKey"`
}

type KeyPublic struct {
	PublicKey string `json:"publicKey"`
}

type KeyAdd struct {
	Key
	KeyPreshared
	KeyPublic
}

type KeyGetInfo struct {
	Key
	KeyPublic
}

func AddKey(res http.ResponseWriter, req *http.Request) {
	bodyRes := responses.StandardResponse{}
	bearerToken := req.Header.Get("Bearer")
	perms := []int{db.PERSONAL_KEYS_ADD, db.KEYS_ADD_ALL}
	validErr := db.ValidateUserPerms(bearerToken, perms)
	if validErr != nil {
		bodyRes.Response = "user does not have permission or an error occurred"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
	}
	//ADD LOGIC
}

func DeleteKey(res http.ResponseWriter, req *http.Request) {
	bodyRes := responses.StandardResponse{}
	bearerToken := req.Header.Get("Bearer")
	perms := []int{db.PERSONAL_KEYS_DELETE, db.KEYS_DELETE_ALL}

	validErr := db.ValidateUserPerms(bearerToken, perms)
	if validErr != nil {
		bodyRes.Response = "user does not have permission or an error occurred"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
	}
	//ADD LOGIC
}

func EnableDisableKey(res http.ResponseWriter, req *http.Request) {
	bodyRes := responses.StandardResponse{}
	bearerToken := req.Header.Get("Bearer")
	perms := []int{db.PERSONAL_KEYS_MODIFY, db.KEYS_MODIFY_ALL}

	validErr := db.ValidateUserPerms(bearerToken, perms)
	if validErr != nil {
		bodyRes.Response = "user does not have permission or an error occurred"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
	}
	//ADD LOGIC
}

func GetAllKeys(res http.ResponseWriter, req *http.Request) {
	bodyRes := responses.StandardResponse{}
	bearerToken := req.Header.Get("Bearer")
	perms := []int{db.KEYS_VIEW_ALL}

	validErr := db.ValidateUserPerms(bearerToken, perms)
	if validErr != nil {
		bodyRes.Response = "user does not have permission or an error occurred"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
	}
	//ADD LOGIC
}
