package routes

import (
	"net/http"

	"gitlab.com/mawthuq-software/wireguard-manager-authenticator/src/api/router/responses"
	"gitlab.com/mawthuq-software/wireguard-manager-authenticator/src/db"
	"gitlab.com/mawthuq-software/wireguard-manager-authenticator/src/token"
)

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

func ValidatePerms(res http.ResponseWriter, req *http.Request, perms []int) {
	bodyRes := responses.StandardResponse{}
	bearerToken := req.Header.Get("Bearer")
	if bearerToken == "" {
		bodyRes.Response = "bearer token cannot be blank"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}
	username, tokenErr := token.ValidateUser(bearerToken)
	if tokenErr != nil {
		bodyRes.Response = "token is not valid"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}
	user, userErr := db.FindAuthFromUsername(username)
	if userErr != nil {
		bodyRes.Response = "could not find user"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}
	userHasPerm := false
	for i := 0; i < len(perms); i++ {
		hasPerms, permErr := db.CheckPermission(user.ID, perms[i])
		if permErr != nil {
			bodyRes.Response = "error when checking permissions"
			responses.Standard(res, bodyRes, http.StatusBadRequest)
			return
		}
		if hasPerms {
			userHasPerm = true
			break
		}
	}
	if !userHasPerm {
		bodyRes.Response = "user does not have permission"
		responses.Standard(res, bodyRes, http.StatusBadRequest)
		return
	}

	bodyRes.Response = "entity has access"
	responses.Standard(res, bodyRes, http.StatusAccepted)
}
