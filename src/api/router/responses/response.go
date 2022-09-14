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
	Keys []db.VPNKeys `json:"keys"`
}

type UserSubscriptionResponse struct {
	StandardResponse
	UserSubscription db.UserSubscriptions `json:"userSubscription"`
}

type DumpUserSubscriptionResponse struct {
	StandardResponse
	UserSubscription []db.UserSubscriptions `json:"userSubscription"`
}

type SubscriptionResponse struct {
	StandardResponse
	Subscription db.Subscriptions `json:"subscription"`
}

type DumpSubscriptionResponse struct {
	StandardResponse
	Subscriptions []db.Subscriptions `json:"subscriptions"`
}

type ServerResponse struct {
	StandardResponse
	Server db.Servers `json:"server"`
}

type DumpServerResponse struct {
	StandardResponse
	Servers []db.Servers `json:"servers"`
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

func Subscription(res http.ResponseWriter, resStruct SubscriptionResponse, httpStatusCode int) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(httpStatusCode)
	jsonResp, _ := json.Marshal(resStruct)
	res.Write(jsonResp)
}

func DumpSubscriptions(res http.ResponseWriter, resStruct DumpSubscriptionResponse, httpStatusCode int) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(httpStatusCode)
	jsonResp, _ := json.Marshal(resStruct)
	res.Write(jsonResp)
}

func Server(res http.ResponseWriter, resStruct ServerResponse, httpStatusCode int) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(httpStatusCode)
	jsonResp, _ := json.Marshal(resStruct)
	res.Write(jsonResp)
}

func DumpServers(res http.ResponseWriter, resStruct DumpServerResponse, httpStatusCode int) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(httpStatusCode)
	jsonResp, _ := json.Marshal(resStruct)
	res.Write(jsonResp)
}
