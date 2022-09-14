package routes

import (
	"net"
	"net/http"
	"strconv"

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

func ReadServer(res http.ResponseWriter, req *http.Request) {
	bodyRes := responses.ServerResponse{}
	queryVars := req.URL.Query()

	serverID := queryVars.Get("id")
	serverName := queryVars.Get("name")

	if serverID == "" && serverName == "" {
		bodyRes.Response = "id or name needs to be filled"
		responses.Server(res, bodyRes, http.StatusBadRequest)
		return
	}

	//check perms
	bearerToken := req.Header.Get("Bearer")

	userPerms := []int{db.SERVER_VIEW}
	_, validUserErr := db.ValidatePerms(bearerToken, userPerms)

	if validUserErr != nil {
		bodyRes.Response = "user does not have permission or an error occurred"
		responses.Server(res, bodyRes, http.StatusForbidden)
		return
	}

	if serverID != "" {
		serverIDInt, convErr := strconv.Atoi(serverID)
		if convErr != nil {
			bodyRes.Response = "could not convert id into int"
			responses.Server(res, bodyRes, http.StatusBadRequest)
			return
		}

		server, serverErr := db.ReadServer(serverIDInt)
		if serverErr != nil {
			bodyRes.Response = serverErr.Error()
			responses.Server(res, bodyRes, http.StatusBadRequest)
			return
		}
		bodyRes.Server = server
		bodyRes.Response = "pulled server successfully"
		responses.Server(res, bodyRes, http.StatusAccepted)
	} else if serverName != "" {
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
}

func ReadAllServers(res http.ResponseWriter, req *http.Request) {
	bodyRes := responses.DumpServerResponse{}
	//check perms
	bearerToken := req.Header.Get("Bearer")

	userPerms := []int{db.SERVER_VIEW}
	_, validUserErr := db.ValidatePerms(bearerToken, userPerms)

	if validUserErr != nil {
		bodyRes.Response = "user does not have permission or an error occurred"
		responses.DumpServers(res, bodyRes, http.StatusForbidden)
		return
	}

	servers, serverErr := db.ReadAllServers()
	if serverErr != nil {
		bodyRes.Response = serverErr.Error()
		responses.DumpServers(res, bodyRes, http.StatusBadRequest)
		return
	}

	bodyRes.Servers = servers
	bodyRes.Response = "pulled servers successfully"
	responses.DumpServers(res, bodyRes, http.StatusAccepted)
}

func UpdateServer(res http.ResponseWriter, req *http.Request) {
	bodyRes := responses.StandardResponse{}

	//check perms
	bearerToken := req.Header.Get("Bearer")

	userPerms := []int{db.SERVER_MODIFY_ALL}
	_, validAdminErr := db.ValidatePerms(bearerToken, userPerms)

	if validAdminErr != nil {
		bodyRes.Response = "user does not have permission or an error occurred"
		responses.Standard(res, bodyRes, http.StatusForbidden)
		return
	}

	queryVars := req.URL.Query()

	serverID := queryVars.Get("id")
	serverName := queryVars.Get("name")
	serverRegion := queryVars.Get("region")
	serverCountry := queryVars.Get("country")
	serverIPAddress := queryVars.Get("ipAddress")

	serverIDInt, serverConvErr := strconv.Atoi(serverID)
	if serverConvErr != nil {
		bodyRes.Response = "could not convert id to int"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	server, serverErr := db.ReadServer(serverIDInt)
	if serverErr != nil {
		bodyRes.Response = serverErr.Error()
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	if serverName != "" {
		server.Name = serverName
	}
	if serverRegion != "" {
		server.Region = serverRegion
	}
	if serverCountry != "" {
		server.Country = serverCountry
	}
	if serverIPAddress != "" {
		ip := net.ParseIP(serverIPAddress)
		if ip == nil {
			bodyRes.Response = "could not parse ipAddress"
			responses.Standard(res, bodyRes, http.StatusBadRequest)
			return
		}
		server.IPAddress = serverIPAddress
	}
	updateErr := db.UpdateServer(server)
	if updateErr != nil {
		bodyRes.Response = updateErr.Error()
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}
	bodyRes.Response = "updated server successfully"
	responses.Standard(res, bodyRes, http.StatusAccepted)
}

func DeleteServer(res http.ResponseWriter, req *http.Request) {
	bodyRes := responses.StandardResponse{}

	//check perms
	bearerToken := req.Header.Get("Bearer")

	userPerms := []int{db.SERVER_MODIFY_ALL}
	_, validAdminErr := db.ValidatePerms(bearerToken, userPerms)

	if validAdminErr != nil {
		bodyRes.Response = "user does not have permission or an error occurred"
		responses.Standard(res, bodyRes, http.StatusForbidden)
		return
	}

	queryVars := req.URL.Query()

	serverID := queryVars.Get("id")
	if serverID == "" {
		bodyRes.Response = "id needs to be filled"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	serverIDInt, serverConvErr := strconv.Atoi(serverID)
	if serverConvErr != nil {
		bodyRes.Response = "could not convert id to int"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	updateErr := db.DeleteServer(serverIDInt)
	if updateErr != nil {
		bodyRes.Response = updateErr.Error()
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}
	bodyRes.Response = "deleted server successfully"
	responses.Standard(res, bodyRes, http.StatusAccepted)
}
