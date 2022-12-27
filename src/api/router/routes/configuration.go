package routes

import (
	"net/http"
	"strconv"

	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/api/router/responses"
	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/db"
)

type Configuration struct {
	ConfigurationID int    `json:"configurationID"`
	Name            string `json:"name"`
	DNS             string `json:"dns"`
	Mask            int    `json:"mask"`
	NumberOfKeys    int    `json:"numberOfKeys"`
}

func CreateConfiguration(res http.ResponseWriter, req *http.Request) {
	var bodyReq Configuration
	bodyRes := responses.StandardResponse{}

	err := ParseRequest(req, &bodyReq)
	if err != nil {
		combinedLogger.Error("Parsing request " + err.Error())
		bodyRes.Response = "Error parsing request"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	if bodyReq.Name == "" {
		bodyRes.Response = "name must be filled"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	} else if bodyReq.DNS == "" {
		bodyRes.Response = "dns must be filled"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	} else if bodyReq.Mask < 16 || bodyReq.Mask > 24 {
		bodyRes.Response = "mask must be >= 16 or <= 24"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	} else if bodyReq.NumberOfKeys < 1 || bodyReq.NumberOfKeys > 5000 {
		bodyRes.Response = "numberOfKeys must be > 0 and < 5000"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	// Check perms
	bearerToken := req.Header.Get("Bearer")

	adminPerms := []int{db.CONFIGURATION_ADD}
	_, validAdminErr := db.ValidatePerms(bearerToken, adminPerms)

	if validAdminErr != nil {
		bodyRes.Response = "user does not have permission or an error occurred"
		responses.Standard(res, bodyRes, http.StatusForbidden)
		return
	}

	createConf := db.CreateConfiguration(bodyReq.Name, bodyReq.DNS, bodyReq.Mask, bodyReq.NumberOfKeys)
	if createConf != nil {
		bodyRes.Response = createConf.Error()
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}
	bodyRes.Response = "successfully created configuration"
	responses.Standard(res, bodyRes, http.StatusAccepted)
}

func ReadConfiguration(res http.ResponseWriter, req *http.Request) {
	bodyRes := responses.ConfigurationResponse{}

	queryVars := req.URL.Query()

	confIDStr := queryVars.Get("id")
	if confIDStr == "" {
		bodyRes.Response = "id needs to be filled"
		responses.Configuration(res, bodyRes, http.StatusBadRequest)
		return
	}

	// Check perms
	bearerToken := req.Header.Get("Bearer")

	adminPerms := []int{db.CONFIGURATION_READ}
	_, validAdminErr := db.ValidatePerms(bearerToken, adminPerms)

	if validAdminErr != nil {
		bodyRes.Response = "user does not have permission or an error occurred"
		responses.Configuration(res, bodyRes, http.StatusForbidden)
		return
	}
	// Convert to int
	confID, convErr := strconv.Atoi(confIDStr)
	if convErr != nil {
		bodyRes.Response = "id needs to be converted to int"
		responses.Configuration(res, bodyRes, http.StatusBadRequest)
		return
	}
	// Read conf
	conf, confErr := db.ReadConfiguration(confID)
	if confErr != nil {
		bodyRes.Response = confErr.Error()
		responses.Configuration(res, bodyRes, http.StatusBadRequest)
		return
	}

	bodyRes.Response = "successfully pulled configuration"
	bodyRes.Configuration = conf
	responses.Configuration(res, bodyRes, http.StatusAccepted)
}

func ReadAllConfigurations(res http.ResponseWriter, req *http.Request) {
	bodyRes := responses.DumpConfigurationResponse{}

	// Check perms
	bearerToken := req.Header.Get("Bearer")

	adminPerms := []int{db.CONFIGURATION_READ, db.PERSONAL_KEYS_ADD}
	_, validAdminErr := db.ValidatePerms(bearerToken, adminPerms)

	if validAdminErr != nil {
		bodyRes.Response = "user does not have permission or an error occurred"
		responses.DumpConfigurations(res, bodyRes, http.StatusForbidden)
		return
	}

	// Read conf
	conf, confErr := db.ReadAllConfigurations()
	if confErr != nil {
		bodyRes.Response = confErr.Error()
		responses.DumpConfigurations(res, bodyRes, http.StatusBadRequest)
		return
	}

	bodyRes.Response = "successfully pulled configurations"
	bodyRes.Configurations = conf
	responses.DumpConfigurations(res, bodyRes, http.StatusAccepted)
}

func UpdateConfiguration(res http.ResponseWriter, req *http.Request) {
	bodyRes := responses.StandardResponse{}

	//check perms
	bearerToken := req.Header.Get("Bearer")

	userPerms := []int{db.CONFIGURATION_MODIFY}
	_, validAdminErr := db.ValidatePerms(bearerToken, userPerms)

	if validAdminErr != nil {
		bodyRes.Response = "user does not have permission or an error occurred"
		responses.Standard(res, bodyRes, http.StatusForbidden)
		return
	}

	queryVars := req.URL.Query()

	configID := queryVars.Get("id")
	configName := queryVars.Get("name")
	configDNS := queryVars.Get("dns")
	configMask := queryVars.Get("mask")

	if configID == "" || configID == "0" {
		bodyRes.Response = "id needs to be filled"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	convID, convErr := strconv.Atoi(configID)
	if convErr != nil {
		bodyRes.Response = "unable to convert id to int"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	convInt, convErr := strconv.Atoi(configMask)
	if convErr != nil {
		bodyRes.Response = "unable to convert the mask to an int"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	confErr := db.CheckConfig(configDNS, convInt)
	if confErr != nil {
		bodyRes.Response = confErr.Error()
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	conf, readErr := db.ReadConfiguration(convID)
	if readErr != nil {
		bodyRes.Response = readErr.Error()
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	if configName != "" {
		conf.Name = configName
	}
	if configDNS != "" {
		conf.DNS = configDNS
	}
	if convInt != -1 {
		conf.Mask = convInt
	}

	updateErr := db.UpdateConfiguration(conf)
	if updateErr != nil {
		bodyRes.Response = updateErr.Error()
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}
	bodyRes.Response = "successfully updated config"
	responses.Standard(res, bodyRes, http.StatusAccepted)
}

func DeleteConfiguration(res http.ResponseWriter, req *http.Request) {
	bodyRes := responses.StandardResponse{}

	//check perms
	bearerToken := req.Header.Get("Bearer")

	userPerms := []int{db.CONFIGURATION_DELETE}
	_, validAdminErr := db.ValidatePerms(bearerToken, userPerms)

	if validAdminErr != nil {
		bodyRes.Response = "user does not have permission or an error occurred"
		responses.Standard(res, bodyRes, http.StatusForbidden)
		return
	}

	queryVars := req.URL.Query()

	configID := queryVars.Get("id")
	if configID == "" || configID == "0" {
		bodyRes.Response = "id needs to be filled"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	convID, convErr := strconv.Atoi(configID)
	if convErr != nil {
		bodyRes.Response = "unable to convert id to int"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	delErr := db.DeleteConfiguration(convID)
	if delErr != nil {
		bodyRes.Response = delErr.Error()
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}
	bodyRes.Response = "successfully delete configuration"
	responses.Standard(res, bodyRes, http.StatusAccepted)
}
