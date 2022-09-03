package responses

import (
	"encoding/json"
	"net/http"

	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/db"
)

type StandardResponse struct {
	Response string `json:"response"`
}

type TokenResponse struct {
	StandardResponse
	Token string `json:"token"`
}

type AllKeyResponse struct {
	StandardResponse
	Keys []db.Keys `json:"keys"`
}

type UserSubscriptionResponse struct {
	StandardResponse
	UserSubscription db.UserSubscriptions `json:"userSubscription"`
}

type DumpUserSubscriptionResponse struct {
	StandardResponse
	UserSubscription []db.UserSubscriptions `json:"userSubscription"`
}

func Standard(res http.ResponseWriter, resStruct StandardResponse, httpStatusCode int) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(httpStatusCode)
	jsonResp, _ := json.Marshal(resStruct)
	res.Write(jsonResp)
}

func Token(res http.ResponseWriter, resStruct TokenResponse, httpStatusCode int) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(httpStatusCode)
	jsonResp, _ := json.Marshal(resStruct)
	res.Write(jsonResp)
}

func AllKeys(res http.ResponseWriter, resStruct AllKeyResponse, httpStatusCode int) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(httpStatusCode)
	jsonResp, _ := json.Marshal(resStruct)
	res.Write(jsonResp)
}

func UserSubscription(res http.ResponseWriter, resStruct UserSubscriptionResponse, httpStatusCode int) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(httpStatusCode)
	jsonResp, _ := json.Marshal(resStruct)
	res.Write(jsonResp)
}

func DumpUserSubscription(res http.ResponseWriter, resStruct DumpUserSubscriptionResponse, httpStatusCode int) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(httpStatusCode)
	jsonResp, _ := json.Marshal(resStruct)
	res.Write(jsonResp)
}
