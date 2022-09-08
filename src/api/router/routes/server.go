package routes

import (
	"net/http"

	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/api/router/responses"
	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/db"
)

type Server struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Region    string `json:"region"`
	Country   string `json:"country"`
	IPAddress string `json:"ipAddress"`
}

func CreateServer(res http.ResponseWriter, req *http.Request) {
	var bodyReq Server
	bodyRes := responses.StandardResponse{}

	//check perms
	bearerToken := req.Header.Get("Bearer")

	adminPerms := []int{db.SERVER_ADD_ALL}
	_, validAdminErr := db.ValidatePerms(bearerToken, adminPerms)

	if validAdminErr != nil {
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

	if bodyReq.Name == "" {
		bodyRes.Response = "name cannot be blank"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	} else if bodyReq.Region == "" {
		bodyRes.Response = "region cannot be blank"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	} else if bodyReq.Country == "" {
		bodyRes.Response = "country cannot be blank"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	} else if bodyReq.IPAddress == "" {
		bodyRes.Response = "ipAddress cannot be blank"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	queryError := db.CreateServer(bodyReq.Name, bodyReq.Region, bodyReq.Country, bodyReq.IPAddress)
	if queryError != nil {
		bodyRes.Response = queryError.Error()
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}
	bodyRes.Response = "server added successfully"
	responses.Standard(res, bodyRes, http.StatusAccepted)
}

func ReadServerFromServerName(res http.ResponseWriter, req *http.Request) {
	bodyRes := responses.ServerResponse{}
	queryVars := req.URL.Query()

	serverName := queryVars.Get("name")

	//check perms
	bearerToken := req.Header.Get("Bearer")

	userPerms := []int{db.SERVER_VIEW}
	_, validUserErr := db.ValidatePerms(bearerToken, userPerms)

	if validUserErr != nil {
		bodyRes.Response = "user does not have permission or an error occurred"
		responses.Server(res, bodyRes, http.StatusForbidden)
		return
	}

	server, serverErr := db.ReadServerFromServerName(serverName)
	if serverErr != nil {
		bodyRes.Response = serverErr.Error()
		responses.Server(res, bodyRes, http.StatusBadRequest)
		return
	}
	bodyRes.Server = server
	bodyRes.Response = "pulled server successfully"
	responses.Server(res, bodyRes, http.StatusAccepted)
}
