package responses

import (
	"encoding/json"
	"net/http"
)

type StandardResponse struct {
	Response  string `json:"response"`
	Completed bool   `json:"completed"`
}

func Standard(res http.ResponseWriter, resStruct StandardResponse, httpStatusCode int) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(httpStatusCode)
	jsonResp, _ := json.Marshal(resStruct)
	res.Write(jsonResp)
}
