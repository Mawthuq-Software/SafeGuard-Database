package routes

import (
	"net/http"
	"strconv"

	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/api/router/responses"
	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/db"
)

type ServerConfiguration struct {
	ServerID int `json:"serverID"`
	ConfigID int `json:"configID"`
}

func CreateServerConfiguration(res http.ResponseWriter, req *http.Request) {
	var bodyReq ServerConfiguration
	bodyRes := responses.StandardResponse{}

	err := ParseRequest(req, &bodyReq)
	if err != nil {
		combinedLogger.Error("Parsing request " + err.Error())
		bodyRes.Response = "Error parsing request"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	if bodyReq.ServerID == 0 {
		bodyRes.Response = "serverID must be filled"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	} else if bodyReq.ConfigID == 0 {
		bodyRes.Response = "configID must be filled"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	// Check perms
	bearerToken := req.Header.Get("Bearer")

	adminPerms := []int{db.SERVER_CONFIGURATION_ADD}
	_, validAdminErr := db.ValidatePerms(bearerToken, adminPerms)

	if validAdminErr != nil {
		bodyRes.Response = "user does not have permission or an error occurred"
		responses.Standard(res, bodyRes, http.StatusForbidden)
		return
	}

	createServerConf := db.CreateServerConfig(bodyReq.ServerID, bodyReq.ConfigID)
	if createServerConf != nil {
		bodyRes.Response = createServerConf.Error()
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}
	bodyRes.Response = "successfully created server configuration"
	responses.Standard(res, bodyRes, http.StatusAccepted)
}

func ReadServerConfiguration(res http.ResponseWriter, req *http.Request) {
	bodyRes := responses.ServerConfigurationResponse{}
	queryVars := req.URL.Query()

	serverConfigID := queryVars.Get("id")

	//check perms
	bearerToken := req.Header.Get("Bearer")

	adminPerms := []int{db.SERVER_CONFIGURATION_READ}
	_, validAdminErr := db.ValidatePerms(bearerToken, adminPerms)

	if validAdminErr != nil {
		bodyRes.Response = "user does not have permission or an error occurred"
		responses.ServerConfiguration(res, bodyRes, http.StatusForbidden)
		return
	}

	if serverConfigID == "" {
		bodyRes.Response = "id needs to be filled"
		responses.ServerConfiguration(res, bodyRes, http.StatusBadRequest)
		return
	}

	serverConfigIDInt, serverConfConvErr := strconv.Atoi(serverConfigID)
	if serverConfConvErr != nil {
		bodyRes.Response = "could not convert id to int"
		responses.ServerConfiguration(res, bodyRes, http.StatusBadRequest)
		return
	}

	serverConfigs, serverConfErr := db.ReadServerConfig(serverConfigIDInt)
	if serverConfErr != nil {
		bodyRes.Response = serverConfErr.Error()
		responses.ServerConfiguration(res, bodyRes, http.StatusBadRequest)
		return
	}

	bodyRes.Response = "successfully pulled server configurations"
	bodyRes.ServerConfiguration = serverConfigs
	responses.ServerConfiguration(res, bodyRes, http.StatusAccepted)
}

func ReadAllServerConfiguration(res http.ResponseWriter, req *http.Request) {
	bodyRes := responses.DumpServerConfigurationResponse{}

	//check perms
	bearerToken := req.Header.Get("Bearer")

	adminPerms := []int{db.SERVER_CONFIGURATION_READ}
	_, validAdminErr := db.ValidatePerms(bearerToken, adminPerms)

	if validAdminErr != nil {
		bodyRes.Response = "user does not have permission or an error occurred"
		responses.DumpServerConfigurations(res, bodyRes, http.StatusForbidden)
		return
	}

	serverConfigs, serverConfErr := db.ReadAllServerConfigs()
	if serverConfErr != nil {
		bodyRes.Response = serverConfErr.Error()
		responses.DumpServerConfigurations(res, bodyRes, http.StatusBadRequest)
		return
	}

	bodyRes.Response = "successfully pulled server configurations"
	bodyRes.ServerConfigurations = serverConfigs
	responses.DumpServerConfigurations(res, bodyRes, http.StatusAccepted)
}

func UpdateServerConfiguration(res http.ResponseWriter, req *http.Request) {
	bodyRes := responses.StandardResponse{}
	queryVars := req.URL.Query()

	serverID := queryVars.Get("serverID")
	configID := queryVars.Get("configID")

	//check perms
	bearerToken := req.Header.Get("Bearer")

	adminPerms := []int{db.SERVER_CONFIGURATION_MODIFY}
	_, validAdminErr := db.ValidatePerms(bearerToken, adminPerms)

	if validAdminErr != nil {
		bodyRes.Response = "user does not have permission or an error occurred"
		responses.Standard(res, bodyRes, http.StatusForbidden)
		return
	}

	if serverID == "" {
		bodyRes.Response = "serverID needs to be filled"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	if configID == "" {
		bodyRes.Response = "configID needs to be filled"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	serverIDInt, serverConfConvErr := strconv.Atoi(serverID)
	if serverConfConvErr != nil {
		bodyRes.Response = "could not convert serverID to int"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	configIDInt, confConvErr := strconv.Atoi(configID)
	if confConvErr != nil {
		bodyRes.Response = "could not convert configID to int"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	serverConf, serverConfErr := db.ReadServerConfigFromServerID(serverIDInt)
	if serverConfErr != nil {
		bodyRes.Response = "server configuration doesn't exist or an error occurred."
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	_, confErr := db.ReadConfiguration(configIDInt)
	if confErr != nil {
		bodyRes.Response = "configuration doesn't exist or an error occurred."
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	serverConf.ConfigID = configIDInt
	updateErr := db.UpdateServerConfig(serverConf)

	if updateErr != nil {
		bodyRes.Response = updateErr.Error()
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	bodyRes.Response = "successfully updated server configurations"
	responses.Standard(res, bodyRes, http.StatusAccepted)
}

func DeleteServerConfiguration(res http.ResponseWriter, req *http.Request) {
	bodyRes := responses.StandardResponse{}
	queryVars := req.URL.Query()

	serverID := queryVars.Get("serverID")

	//check perms
	bearerToken := req.Header.Get("Bearer")

	adminPerms := []int{db.SERVER_CONFIGURATION_DELETE}
	_, validAdminErr := db.ValidatePerms(bearerToken, adminPerms)

	if validAdminErr != nil {
		bodyRes.Response = "user does not have permission or an error occurred"
		responses.Standard(res, bodyRes, http.StatusForbidden)
		return
	}

	if serverID == "" {
		bodyRes.Response = "serverID needs to be filled"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	serverIDInt, serverConfConvErr := strconv.Atoi(serverID)
	if serverConfConvErr != nil {
		bodyRes.Response = "could not convert serverID to int"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	serverConf, confErr := db.ReadServerConfigFromServerID(serverIDInt)
	if confErr != nil {
		bodyRes.Response = "server configuration doesn't exist or an error occurred."
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	delErr := db.DeleteServerConfig(serverConf.ID)

	if delErr != nil {
		bodyRes.Response = delErr.Error()
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	bodyRes.Response = "successfully deleted server configurations"
	responses.Standard(res, bodyRes, http.StatusAccepted)
}
