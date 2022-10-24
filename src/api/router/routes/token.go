package routes

import (
	"net/http"

	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/api/router/responses"
	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/db"
	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/token"
)

type Token struct {
	Name    string `json:"name"`
	TokenID int    `json:"tokenID"`
}

// Creates an API token
func CreateToken(res http.ResponseWriter, req *http.Request) {
	var bodyReq Token
	bodyRes := responses.TokenResponse{}
	bearerToken := req.Header.Get("Bearer")

	err := ParseRequest(req, &bodyReq)
	if err != nil {
		combinedLogger.Error("Parsing request " + err.Error())
		bodyRes.Response = "Error parsing request"
		responses.Token(res, bodyRes, http.StatusBadRequest)
		return
	}

	if bodyReq.Name == "" {
		bodyRes.Response = "name must be filled"
		responses.Token(res, bodyRes, http.StatusBadRequest)
		return
	}

	adminPerms := []int{db.TOKEN_ADD}
	_, validAdminErr := db.ValidatePerms(bearerToken, adminPerms)

	if validAdminErr != nil {
		bodyRes.Response = "user does not have permission or an error occurred"
		responses.Token(res, bodyRes, http.StatusForbidden)
		return
	} else if validAdminErr == nil { //If request is from admin with perms
		uuid, adminSubErr := db.CreateToken(bodyReq.Name)
		if adminSubErr != nil {
			bodyRes.Response = adminSubErr.Error()
			responses.Token(res, bodyRes, http.StatusBadRequest)
			return
		}
		bodyRes.Response = "added successfully"
		bodyRes.Token = uuid
		responses.Token(res, bodyRes, http.StatusAccepted)
		return
	} else {
		bodyRes.Response = "an error occurred"
		responses.Token(res, bodyRes, http.StatusInternalServerError)
		return
	}
}

func DeleteToken(res http.ResponseWriter, req *http.Request) {
	var bodyReq Token
	bodyRes := responses.StandardResponse{}
	bearerToken := req.Header.Get("Bearer")

	err := ParseRequest(req, &bodyReq)
	if err != nil {
		combinedLogger.Error("Parsing request " + err.Error())
		bodyRes.Response = "Error parsing request"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	if bodyReq.TokenID <= 0 {
		bodyRes.Response = "tokenID must be filled"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}
	// , db.PERSONAL_KEYS_ADD
	adminPerms := []int{db.TOKEN_DELETE}
	_, validAdminErr := db.ValidatePerms(bearerToken, adminPerms)

	if validAdminErr != nil {
		bodyRes.Response = "user does not have permission or an error occurred"
		responses.Standard(res, bodyRes, http.StatusForbidden)
		return
	} else if validAdminErr == nil { //If request is from admin with perms
		adminSubErr := db.DeleteToken(bodyReq.TokenID)
		if adminSubErr != nil {
			bodyRes.Response = adminSubErr.Error()
			responses.Standard(res, bodyRes, http.StatusBadRequest)
			return
		}
		bodyRes.Response = "deleted token successfully"
		responses.Standard(res, bodyRes, http.StatusAccepted)
		return
	} else {
		bodyRes.Response = "an error occurred"
		responses.Standard(res, bodyRes, http.StatusInternalServerError)
		return
	}
}

// Validates a user's JWT Bearer token
func Validate(res http.ResponseWriter, req *http.Request) {
	bodyRes := responses.StandardResponse{}
	bearerToken := req.Header.Get("Bearer")
	if bearerToken == "" {
		bodyRes.Response = "bearer token cannot be blank"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}
	_, tokenErr := token.ValidateUser(bearerToken)
	if tokenErr != nil {
		bodyRes.Response = tokenErr.Error() //use the direct error text here
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}
	bodyRes.Response = "token verified"
	responses.Standard(res, bodyRes, http.StatusOK)
}
