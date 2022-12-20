package routes

import (
	"net/http"

	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/api/router/responses"
	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/db"
)

type WireguardInstance struct {
	ServerID    int    `json:"serverID"`
	ListenPort  int    `json:"listenPort"`
	PublicKey   string `json:"publicKey"`
	IPV4Address string `json:"ipv4Address"`
	IPV6Address string `json:"ipv6Address"`
}

func CreateWireguardInstance(res http.ResponseWriter, req *http.Request) {
	var bodyReq WireguardInstance
	bodyRes := responses.StandardResponse{}

	err := ParseRequest(req, &bodyReq)
	if err != nil {
		combinedLogger.Error("Parsing request " + err.Error())
		bodyRes.Response = "Error parsing request"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	if bodyReq.ServerID <= 0 {
		bodyRes.Response = "serverID must be >= 0"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	} else if bodyReq.ListenPort <= 0 {
		bodyRes.Response = "listenPort must be >= 0"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	} else if bodyReq.PublicKey == "" {
		bodyRes.Response = "publicKey must be filled"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	} else if bodyReq.IPV4Address == "" {
		bodyRes.Response = "ipv4Address must be filled"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	} else if bodyReq.IPV6Address == "" {
		bodyRes.Response = "ipv6Address must be filled"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	//check perms
	bearerToken := req.Header.Get("Bearer")

	adminPerms := []int{db.PERSONAL_KEYS_ADD}
	_, validAdminErr := db.ValidatePerms(bearerToken, adminPerms)

	if validAdminErr != nil {
		bodyRes.Response = "user does not have permission or an error occurred"
		responses.Standard(res, bodyRes, http.StatusForbidden)
		return
	} else if validAdminErr == nil { //If request is from admin with perms
		adminSubErr := db.CreateServerWireguardInterface(bodyReq.ServerID, bodyReq.ListenPort, bodyReq.PublicKey, bodyReq.IPV4Address, bodyReq.IPV6Address)
		if adminSubErr != nil {
			bodyRes.Response = adminSubErr.Error()
			responses.Standard(res, bodyRes, http.StatusBadRequest)
			return
		}
		bodyRes.Response = "created wireguard instance successfully"
		responses.Standard(res, bodyRes, http.StatusAccepted)
		return
	} else {
		bodyRes.Response = "an error occurred"
		responses.Standard(res, bodyRes, http.StatusInternalServerError)
		return
	}
}
