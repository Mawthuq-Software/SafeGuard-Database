package routes

import (
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
	var bodyReq User
	bodyRes := responses.StandardResponse{}

	err := ParseRequest(req, &bodyReq)
	if err != nil {
		combinedLogger.Error("Parsing request " + err.Error())
		bodyRes.Response = "Error parsing request"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	if bodyReq.Username == "" {
		bodyRes.Response = "username cannot be blank"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	} else if bodyReq.Password == "" {
		bodyRes.Response = "password cannot be blank"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	} else if bodyReq.Email == "" {
		bodyRes.Response = "email cannot be blank"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}
	queryError := db.AddUser(bodyReq.Username, bodyReq.Password, bodyReq.Email)
	if queryError != nil {
		bodyRes.Response = queryError.Error()
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}
	bodyRes.Response = "user added"
	responses.Standard(res, bodyRes, http.StatusAccepted)
}

func LoginWithUsername(res http.ResponseWriter, req *http.Request) {
	var bodyReq User
	bodyRes := responses.TokenResponse{}

	err := ParseRequest(req, &bodyReq)
	if err != nil {
		combinedLogger.Error("Parsing request " + err.Error())
		bodyRes.Response = "Error parsing request"
		responses.Token(res, bodyRes, http.StatusBadRequest)
		return
	}

	if bodyReq.Username == "" {
		bodyRes.Response = "username cannot be blank"
		responses.Token(res, bodyRes, http.StatusBadRequest)
		return
	} else if bodyReq.Password == "" {
		bodyRes.Response = "password cannot be blank"
		responses.Token(res, bodyRes, http.StatusBadRequest)
		return
	}
	token, dbErr := db.LoginWithUsername(bodyReq.Username, bodyReq.Password)
	if dbErr != nil {
		bodyRes.Response = "could not login user"
		responses.Token(res, bodyRes, http.StatusBadRequest)
		return
	}
	bodyRes.Token = token
	bodyRes.Response = "successfully created token"
	responses.Token(res, bodyRes, http.StatusAccepted)
}
