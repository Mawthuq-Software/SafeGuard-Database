package routes

import (
	"log"
	"net/http"

	"gitlab.com/mawthuq-software/wireguard-manager-authenticator/src/api/router/responses"
	"gitlab.com/mawthuq-software/wireguard-manager-authenticator/src/db"
)

type User struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	UserID      int    `json:"userID"`
	Group       string `json:"group"`
	AccessToken string `json:"accessToken"`
}

func AddUser(res http.ResponseWriter, req *http.Request) {
	var incomingJSON User
	response := responses.StandardResponse{}

	err := parseRequest(req, &incomingJSON)
	if err != nil {
		log.Println("Error - Parsing request", err)
		response.Completed = false
		response.Response = "Error parsing request"
		responses.Standard(res, response, http.StatusBadRequest)
		return
	}

	if incomingJSON.Username == "" {
		response.Completed = false
		response.Response = "Username cannot be blank"
		responses.Standard(res, response, http.StatusBadRequest)
		return
	} else if incomingJSON.Password == "" {
		response.Completed = false
		response.Response = "Password cannot be blank"
		responses.Standard(res, response, http.StatusBadRequest)
		return
	} else if incomingJSON.Email == "" {
		response.Completed = false
		response.Response = "Email cannot be blank"
		responses.Standard(res, response, http.StatusBadRequest)
		return
	}

	dbRes := db.AddUser(incomingJSON.Username, incomingJSON.Password, incomingJSON.Email)
	response.Completed = dbRes.Proccessed
	response.Response = dbRes.Response
	if !dbRes.Proccessed {
		responses.Standard(res, response, http.StatusBadRequest)
		return
	}
	responses.Standard(res, response, http.StatusAccepted)
}
