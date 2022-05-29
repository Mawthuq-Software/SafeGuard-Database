package routes

import (
	"net/http"

	"gitlab.com/mawthuq-software/wireguard-manager-authenticator/src/api/router/responses"
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
