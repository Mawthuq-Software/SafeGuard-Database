package responses

import (
	"encoding/json"
	"net/http"
)

type StandardResponse struct {
	Response string `json:"response"`
}

type TokenResponse struct {
	StandardResponse
	Token string `json:"token"`
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
