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
	user.HandleFunc("/add", routes.AddUser).Methods("POST")
	user.HandleFunc("/login", routes.LoginWithUsername).Methods("POST")
	user.HandleFunc("/change-password", routes.ChangeUserPassword).Methods("POST")

	userSubscription := user.PathPrefix("/subscription").Subrouter()
	// userSubscription.HandleFunc("/edit", routes.EditingUserSubscription).Methods("POST")
	userSubscription.HandleFunc("/get", routes.GetUserSubscription).Methods("POST")
	userSubscription.HandleFunc("/get-all", routes.GetAllUserSubscriptions).Methods("POST")
	userSubscription.HandleFunc("/create", routes.CreateUserSubscription).Methods("POST")

	subscription := router.PathPrefix("/subscription").Subrouter()
	subscription.HandleFunc("/", routes.CreateSubscription).Methods("POST")
	subscription.HandleFunc("/", routes.ReadSubscription).Methods("GET")
	subscription.HandleFunc("/update-numberofkeys", routes.UpdateSubscriptionNumKeys).Methods("PUT")
	subscription.HandleFunc("/update-bandwidth", routes.UpdateSubscriptionBandwidth).Methods("PUT")
	subscription.HandleFunc("/", routes.DeleteSubscription).Methods("DELETE")

	key := router.PathPrefix("/key").Subrouter()
	key.HandleFunc("/add", routes.AddKey).Methods("POST")
	key.HandleFunc("/delete", routes.DeleteKey).Methods("POST")
	key.HandleFunc("/toggle-usage", routes.EnableDisableKey).Methods("POST")
	key.HandleFunc("/get-all", routes.GetAllKeys).Methods("POST")

	token := router.PathPrefix("/token").Subrouter()
	token.HandleFunc("/validate", routes.Validate).Methods("POST")

	router.MethodNotAllowedHandler = http.HandlerFunc(setCorsHeader) //if method is not found allow OPTIONS
	return router
}
