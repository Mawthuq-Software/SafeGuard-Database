package router

import (
	"net/http"

	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/api/router/routes"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter() //Router for routes
	router.Use(setHeader)     //need to allow CORS and OPTIONS

	user := router.PathPrefix("/user").Subrouter()
	user.HandleFunc("/register", routes.AddUser).Methods("POST")
	user.HandleFunc("/login", routes.LoginWithUsername).Methods("POST")
	user.HandleFunc("/change-password", routes.ChangeUserPassword).Methods("POST")

	userSubscription := user.PathPrefix("/subscription").Subrouter()
	userSubscription.HandleFunc("/", routes.CreateUserSubscription).Methods("POST")
	userSubscription.HandleFunc("/", routes.ReadUserSubscription).Methods("GET")
	userSubscription.HandleFunc("/all", routes.ReadAllUserSubscriptions).Methods("GET")
	userSubscription.HandleFunc("/", routes.UpdateUserSubscription).Methods("PUT")
	userSubscription.HandleFunc("/", routes.DeleteUserSubscription).Methods("DELETE")

	subscription := router.PathPrefix("/subscription").Subrouter()
	subscription.HandleFunc("/", routes.CreateSubscription).Methods("POST")
	subscription.HandleFunc("/", routes.ReadSubscription).Methods("GET")
	subscription.HandleFunc("/all", routes.ReadAllSubscriptions).Methods("GET")
	subscription.HandleFunc("/", routes.UpdateSubscription).Methods("PUT")
	subscription.HandleFunc("/", routes.DeleteSubscription).Methods("DELETE")

	key := router.PathPrefix("/key").Subrouter()
	key.HandleFunc("/add", routes.AddKey).Methods("POST")
	key.HandleFunc("/delete", routes.DeleteKey).Methods("POST")
	key.HandleFunc("/toggle-usage", routes.EnableDisableKey).Methods("POST")
	key.HandleFunc("/get-all", routes.GetAllKeys).Methods("POST")

	token := router.PathPrefix("/token").Subrouter()
	token.HandleFunc("/validate", routes.Validate).Methods("POST")

	server := router.PathPrefix("/server").Subrouter()
	server.HandleFunc("/", routes.CreateServer).Methods("POST")
	server.HandleFunc("/", routes.ReadServer).Methods("GET")
	server.HandleFunc("/all", routes.ReadAllServers).Methods("GET")
	server.HandleFunc("/", routes.UpdateServer).Methods("PUT")
	server.HandleFunc("/", routes.DeleteServer).Methods("DELETE")

	router.MethodNotAllowedHandler = http.HandlerFunc(setCorsHeader) //if method is not found allow OPTIONS
	return router
}
