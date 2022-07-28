package routes

import (
	"net/http"

	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/api/router/responses"
	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/db"
)

type Key struct {
	ID       int `json:"keyID"`
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

type KeyDelete struct {
	Key
}

type KeyGetInfo struct {
	Key
	KeyPublic
}

func AddKey(res http.ResponseWriter, req *http.Request) {
	var bodyReq KeyAdd
	bodyRes := responses.StandardResponse{}

	bearerToken := req.Header.Get("Bearer")
	perms := []int{db.PERSONAL_KEYS_ADD, db.KEYS_ADD_ALL}
	username, validErr := db.ValidateUserPerms(bearerToken, perms)
	if validErr != nil {
		bodyRes.Response = "user does not have permission or an error occurred"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	err := ParseRequest(req, &bodyReq)
	if err != nil {
		combinedLogger.Error("Parsing request " + err.Error())
		bodyRes.Response = "Error parsing request"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}
	if bodyReq.ServerID <= 0 {
		bodyRes.Response = "serverID cannot be a null value"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}
	if bodyReq.PublicKey == "" {
		bodyRes.Response = "client publicKey cannot be empty"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}
	if bodyReq.PresharedKey == "" {
		bodyRes.Response = "client presharedKey cannot be empty"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	user, err := db.FindUserFromUsername(username)
	if err != nil {
		bodyRes.Response = err.Error()
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}
	err = db.AddUserKey(user.ID, bodyReq.ServerID, bodyReq.PublicKey, bodyReq.PresharedKey)
	if err != nil {
		bodyRes.Response = err.Error()
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}
	bodyRes.Response = "added key successfully"
	responses.Standard(res, bodyRes, http.StatusBadRequest)
}

func DeleteKey(res http.ResponseWriter, req *http.Request) {
	var bodyReq KeyDelete
	bodyRes := responses.StandardResponse{}

	bearerToken := req.Header.Get("Bearer")
	perms := []int{db.PERSONAL_KEYS_DELETE, db.KEYS_DELETE_ALL}

	_, validErr := db.ValidateUserPerms(bearerToken, perms)
	if validErr != nil {
		bodyRes.Response = "user does not have permission or an error occurred"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	if bodyReq.ID <= 0 {
		bodyRes.Response = "client keyID cannot be empty"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}
	err := db.DeleteUserKey(bodyReq.ID)
	if err != nil {
		bodyRes.Response = err.Error()
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}
	bodyRes.Response = "deleted key successfully"
	responses.Standard(res, bodyRes, http.StatusBadRequest)
}

func EnableDisableKey(res http.ResponseWriter, req *http.Request) {
	var bodyReq Key
	bodyRes := responses.StandardResponse{}
	bearerToken := req.Header.Get("Bearer")
	perms := []int{db.PERSONAL_KEYS_MODIFY, db.KEYS_MODIFY_ALL}

	_, validErr := db.ValidateUserPerms(bearerToken, perms)
	if validErr != nil {
		bodyRes.Response = "user does not have permission or an error occurred"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
	}

	if bodyReq.ID <= 0 {
		bodyRes.Response = "client keyID cannot be empty"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	err := db.ToggleKey(bodyReq.ID)
	if err != nil {
		bodyRes.Response = err.Error()
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}
	bodyRes.Response = "toggled key successfully"
	responses.Standard(res, bodyRes, http.StatusBadRequest)
}

func GetAllKeys(res http.ResponseWriter, req *http.Request) {
	bodyRes := responses.StandardResponse{}
	bearerToken := req.Header.Get("Bearer")
	perms := []int{db.KEYS_VIEW_ALL}

	_, validErr := db.ValidateUserPerms(bearerToken, perms)
	if validErr != nil {
		bodyRes.Response = "user does not have permission or an error occurred"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
	}

	//ADD LOGIC
}
