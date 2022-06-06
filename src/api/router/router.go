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
	user.HandleFunc("/login", routes.LoginWithUsername).Methods("POST")

	key := router.PathPrefix("/key").Subrouter()
	key.HandleFunc("/add", routes.AddKey).Methods("GET")

	token := router.PathPrefix("/token").Subrouter()
	token.HandleFunc("/validate", routes.Validate).Methods("POST")

	router.MethodNotAllowedHandler = http.HandlerFunc(setCorsHeader) //if method is not found allow OPTIONS
	return router
}
