package routes

import (
	"net/http"
	"strconv"

	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/api/router/responses"
	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/db"
)

type ServerToken struct {
	ServerID int `json:"serverID"`
}

// Creates an API token for a server
func CreateServerToken(res http.ResponseWriter, req *http.Request) {
	var bodyReq ServerToken
	bodyRes := responses.StandardResponse{}
	bearerToken := req.Header.Get("Bearer")

	err := ParseRequest(req, &bodyReq)
	if err != nil {
		combinedLogger.Error("Parsing request " + err.Error())
		bodyRes.Response = "Error parsing request"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	if bodyReq.ServerID <= 0 {
		bodyRes.Response = "serverID must be filled"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	adminPerms := []int{db.SERVER_TOKEN_ADD}
	_, validAdminErr := db.ValidatePerms(bearerToken, adminPerms)

	if validAdminErr != nil {
		bodyRes.Response = "user does not have permission or an error occurred"
		responses.Standard(res, bodyRes, http.StatusForbidden)
		return
	} else if validAdminErr == nil { //If request is from admin with perms
		adminSubErr := db.CreateServerToken(bodyReq.ServerID)
		if adminSubErr != nil {
			bodyRes.Response = adminSubErr.Error()
			responses.Standard(res, bodyRes, http.StatusBadRequest)
			return
		}
		bodyRes.Response = "added server token successfully"
		responses.Standard(res, bodyRes, http.StatusAccepted)
		return
	} else {
		bodyRes.Response = "an error occurred"
		responses.Standard(res, bodyRes, http.StatusInternalServerError)
		return
	}
}

func DeleteServerToken(res http.ResponseWriter, req *http.Request) {
	bodyRes := responses.StandardResponse{}
	bearerToken := req.Header.Get("Bearer")

	queryVars := req.URL.Query()

	serverTokenID := queryVars.Get("id")

	if serverTokenID == "" {
		bodyRes.Response = "id needs to be filled"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	serverTokenIDInt, convErr := strconv.Atoi(serverTokenID)
	if convErr != nil {
		bodyRes.Response = "id could not be converted to integer"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	if serverTokenIDInt <= 0 {
		bodyRes.Response = "serverTokenID must be filled"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	adminPerms := []int{db.SERVER_TOKEN_DELETE}
	_, validAdminErr := db.ValidatePerms(bearerToken, adminPerms)

	if validAdminErr != nil {
		bodyRes.Response = "user does not have permission or an error occurred"
		responses.Standard(res, bodyRes, http.StatusForbidden)
		return
	} else if validAdminErr == nil { //If request is from admin with perms
		adminSubErr := db.DeleteServerToken(serverTokenIDInt)
		if adminSubErr != nil {
			bodyRes.Response = adminSubErr.Error()
			responses.Standard(res, bodyRes, http.StatusBadRequest)
			return
		}
		bodyRes.Response = "deleted server token successfully"
		responses.Standard(res, bodyRes, http.StatusAccepted)
		return
	} else {
		bodyRes.Response = "an error occurred"
		responses.Standard(res, bodyRes, http.StatusInternalServerError)
		return
	}
}
