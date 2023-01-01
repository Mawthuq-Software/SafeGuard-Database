package routes

import (
	"net/http"
	"strconv"

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

//Adds a personal user key
func AddKey(res http.ResponseWriter, req *http.Request) {
	var bodyReq KeyAdd
	bodyRes := responses.StandardResponse{}

	bearerToken := req.Header.Get("Bearer")

	userPerms := []int{db.PERSONAL_KEYS_ADD, db.KEYS_ADD_ALL} //check perms
	userID, userValidErr := db.ValidatePerms(bearerToken, userPerms)

	adminPerms := []int{db.KEYS_MODIFY_ALL}
	_, adminValidErr := db.ValidatePerms(bearerToken, adminPerms)

	if userValidErr != nil && adminValidErr != nil {
		bodyRes.Response = "user does not have permission or an error occurred"
		responses.Standard(res, bodyRes, http.StatusForbidden)
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
	} else if userValidErr == nil {
		user, err := db.ReadUser(userID)
		if err != nil {
			bodyRes.Response = err.Error()
			responses.Standard(res, bodyRes, http.StatusBadRequest)
			return
		}
		if user.ID != userID { //verify user is the actual user
			bodyRes.Response = "user does not have permission or an error occurred"
			responses.Standard(res, bodyRes, http.StatusForbidden)
			return
		}
	}

	// We assume based on logic user is either admin or a user with permission
	err = db.CreateKeyAndLink(userID, bodyReq.ServerID, bodyReq.PublicKey, bodyReq.PresharedKey)
	if err != nil {
		bodyRes.Response = err.Error()
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}
	bodyRes.Response = "added key successfully"
	responses.Standard(res, bodyRes, http.StatusAccepted)
}

//Deletes a personal user key
func DeleteKey(res http.ResponseWriter, req *http.Request) {
	bodyRes := responses.StandardResponse{}

	bearerToken := req.Header.Get("Bearer")
	perms := []int{db.PERSONAL_KEYS_DELETE, db.KEYS_DELETE_ALL}

	_, validErr := db.ValidatePerms(bearerToken, perms)
	if validErr != nil {
		bodyRes.Response = "user does not have permission or an error occurred"
		responses.Standard(res, bodyRes, http.StatusForbidden)
		return
	}

	queryVars := req.URL.Query()

	keyID := queryVars.Get("id")
	if keyID == "" {
		bodyRes.Response = "id needs to be filled"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	keyIDInt, convErr := strconv.Atoi(keyID)
	if convErr != nil {
		bodyRes.Response = "id could not be converted to integer"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}
	err := db.DeleteKeyAndLink(keyIDInt)
	if err != nil {
		bodyRes.Response = err.Error()
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}
	bodyRes.Response = "deleted key successfully"
	responses.Standard(res, bodyRes, http.StatusAccepted)
}

//Toggles a key usability
func EnableDisableKey(res http.ResponseWriter, req *http.Request) {
	var bodyReq Key
	bodyRes := responses.StandardResponse{}
	bearerToken := req.Header.Get("Bearer")
	perms := []int{db.PERSONAL_KEYS_MODIFY, db.KEYS_MODIFY_ALL}

	_, validErr := db.ValidatePerms(bearerToken, perms)
	if validErr != nil {
		bodyRes.Response = "user does not have permission or an error occurred"
		responses.Standard(res, bodyRes, http.StatusForbidden)
		return
	}

	err := ParseRequest(req, &bodyReq)
	if err != nil {
		combinedLogger.Error("Parsing request " + err.Error())
		bodyRes.Response = "Error parsing request"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	if bodyReq.ID <= 0 {
		bodyRes.Response = "client keyID cannot be empty"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	err = db.ToggleKey(bodyReq.ID)
	if err != nil {
		bodyRes.Response = err.Error()
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}
	bodyRes.Response = "toggled key successfully"
	responses.Standard(res, bodyRes, http.StatusAccepted)
}

func GetKeys(res http.ResponseWriter, req *http.Request) {
	bodyRes := responses.UserKeysResponse{}
	bearerToken := req.Header.Get("Bearer")
	userPerms := []int{db.PERSONAL_KEYS_VIEW}
	adminPerms := []int{db.KEYS_VIEW_ALL}

	userID, userValidErr := db.ValidatePerms(bearerToken, userPerms)
	_, adminValidErr := db.ValidatePerms(bearerToken, adminPerms)

	queryVars := req.URL.Query()

	userIDQuery := queryVars.Get("userID")
	if userIDQuery == "" {
		bodyRes.Response = "userID needs to be filled"
		responses.UserKeys(res, bodyRes, http.StatusBadRequest)
		return
	}

	userIDQueryInt, convErr := strconv.Atoi(userIDQuery)
	if convErr != nil {
		bodyRes.Response = "id could not be converted to integer"
		responses.UserKeys(res, bodyRes, http.StatusBadRequest)
		return
	}

	if userValidErr != nil && adminValidErr != nil {
		bodyRes.Response = "user does not have permission or an error occurred"
		responses.UserKeys(res, bodyRes, http.StatusForbidden)
		return
	} else if userValidErr == nil {
		if userID != userIDQueryInt {
			bodyRes.Response = "user does not have permission or an error occurred"
			responses.UserKeys(res, bodyRes, http.StatusForbidden)
			return
		}
	}

	userKeys, err := db.ReadUserKeys(userID)
	if err != nil {
		bodyRes.Response = err.Error()
		responses.UserKeys(res, bodyRes, http.StatusBadRequest)
		return
	}
	bodyRes.UserKeys = userKeys
	bodyRes.Response = "got all keys successfully"
	responses.UserKeys(res, bodyRes, http.StatusAccepted)
}

//Gets all keys from the database
func GetAllKeys(res http.ResponseWriter, req *http.Request) {
	bodyRes := responses.AllKeyResponse{}
	bearerToken := req.Header.Get("Bearer")
	perms := []int{db.KEYS_VIEW_ALL}

	_, validErr := db.ValidatePerms(bearerToken, perms)
	if validErr != nil {
		bodyRes.Response = "user does not have permission or an error occurred"
		responses.AllKeys(res, bodyRes, http.StatusForbidden)
		return
	}

	keys, err := db.ReadAllKeys()
	if err != nil {
		bodyRes.Response = err.Error()
		responses.AllKeys(res, bodyRes, http.StatusBadRequest)
		return
	}
	bodyRes.Keys = keys
	bodyRes.Response = "got all keys successfully"
	responses.AllKeys(res, bodyRes, http.StatusAccepted)
}
