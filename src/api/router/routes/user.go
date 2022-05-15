package routes

import "net/http"

type User struct {
	Username    string `json:"username"`
	UserID      int    `json:"userID"`
	Group       string `json:"group"`
	AccessToken string `json:"accessToken"`
}

func AddUser(res http.ResponseWriter, req *http.Request) {

}
