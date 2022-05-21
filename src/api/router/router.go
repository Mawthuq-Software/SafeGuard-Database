package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/mawthuq-software/wireguard-manager-authenticator/src/api/router/routes"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter() //Router for routes
	router.Use(setHeader)     //need to allow CORS and OPTIONS

	user := router.PathPrefix("/user").Subrouter()
	user.HandleFunc("/add", routes.AddUser).Methods("POST")

	router.MethodNotAllowedHandler = http.HandlerFunc(setCorsHeader) //if method is not found allow OPTIONS
	return router
}
