package routes

import (
	"fmt"
	"net/http"

	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/api/router/responses"
	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/db"
)

type User struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	UserID      int    `json:"userID"`
	Group       string `json:"group"`
	AccessToken string `json:"accessToken"`
}

type UserPasswordChange struct {
	User
	NewPassword string `json:"newPassword"`
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
		fmt.Println(dbErr)
		bodyRes.Response = "could not login user"
		responses.Token(res, bodyRes, http.StatusBadRequest)
		return
	}
	bodyRes.Token = token
	bodyRes.Response = "successfully created token"
	responses.Token(res, bodyRes, http.StatusAccepted)
}

func ChangeUserPassword(res http.ResponseWriter, req *http.Request) {
	var bodyReq UserPasswordChange
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
	} else if bodyReq.NewPassword == "" {
		bodyRes.Response = "new password cannot be blank"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	passErr := db.ChangeUserPassword(bodyReq.Username, bodyReq.Password, bodyReq.NewPassword)
	if passErr != nil {
		bodyRes.Response = "there was an issue changing the password"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}
	bodyRes.Response = "password changed successfully"
	responses.Standard(res, bodyRes, http.StatusAccepted)
}
