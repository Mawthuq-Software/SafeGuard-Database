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
	user.HandleFunc("/register", routes.AddUser).Methods("POST")                   // DONE
	user.HandleFunc("/login", routes.LoginWithUsername).Methods("POST")            // DONE
	user.HandleFunc("/change-password", routes.ChangeUserPassword).Methods("POST") // DONE

	userSubscription := user.PathPrefix("/subscription").Subrouter()
	userSubscription.HandleFunc("/", routes.CreateUserSubscription).Methods("POST")              // DONE
	userSubscription.HandleFunc("/", routes.ReadUserSubscription).Methods("GET")                 // DONE
	userSubscription.HandleFunc("/userid", routes.ReadUserSubscriptionFromUserID).Methods("GET") // DONE
	userSubscription.HandleFunc("/all", routes.ReadAllUserSubscriptions).Methods("GET")          // DONE
	userSubscription.HandleFunc("/", routes.UpdateUserSubscription).Methods("PUT")               // DONE
	userSubscription.HandleFunc("/", routes.DeleteUserSubscription).Methods("DELETE")            // DONE

	subscription := router.PathPrefix("/subscription").Subrouter()
	subscription.HandleFunc("/", routes.CreateSubscription).Methods("POST")     // DONE
	subscription.HandleFunc("/", routes.ReadSubscription).Methods("GET")        // DONE
	subscription.HandleFunc("/all", routes.ReadAllSubscriptions).Methods("GET") // DONE
	subscription.HandleFunc("/", routes.UpdateSubscription).Methods("PUT")      // DONE
	subscription.HandleFunc("/", routes.DeleteSubscription).Methods("DELETE")   // DONE

	key := router.PathPrefix("/key").Subrouter()
	key.HandleFunc("/", routes.AddKey).Methods("POST")                       // DONE
	key.HandleFunc("/", routes.DeleteKey).Methods("DELETE")                  // DONE
	key.HandleFunc("/", routes.GetKeys).Methods("GET")                       // DONE
	key.HandleFunc("/get-all", routes.GetAllKeys).Methods("GET")             // DONE
	key.HandleFunc("/toggle-usage", routes.EnableDisableKey).Methods("POST") // NOT DONE NEEDS FURTHER TESTING

	token := router.PathPrefix("/token").Subrouter()
	token.HandleFunc("/", routes.CreateToken).Methods("POST")
	token.HandleFunc("/", routes.DeleteToken).Methods("DELETE")
	token.HandleFunc("/validate", routes.Validate).Methods("POST") // DONE

	server := router.PathPrefix("/server").Subrouter()
	server.HandleFunc("/", routes.CreateServer).Methods("POST")     // DONE
	server.HandleFunc("/", routes.ReadServer).Methods("GET")        // DONE
	server.HandleFunc("/all", routes.ReadAllServers).Methods("GET") // DONE
	server.HandleFunc("/", routes.UpdateServer).Methods("PUT")      // DONE
	server.HandleFunc("/", routes.DeleteServer).Methods("DELETE")   // DONE

	serverToken := server.PathPrefix("/token").Subrouter()
	serverToken.HandleFunc("/", routes.CreateServerToken).Methods("POST")
	serverToken.HandleFunc("/", routes.DeleteServerToken).Methods("DELETE")

	configuration := router.PathPrefix("/configuration").Subrouter()
	configuration.HandleFunc("/", routes.CreateConfiguration).Methods("POST")     // DONE
	configuration.HandleFunc("/", routes.ReadConfiguration).Methods("GET")        // DONE
	configuration.HandleFunc("/all", routes.ReadAllConfigurations).Methods("GET") // DONE
	configuration.HandleFunc("/", routes.UpdateConfiguration).Methods("PUT")      // DONE
	configuration.HandleFunc("/", routes.DeleteConfiguration).Methods("DELETE")   // DONE

	serverConfiguration := router.PathPrefix("/server-configuration").Subrouter()
	serverConfiguration.HandleFunc("/", routes.CreateServerConfiguration).Methods("POST")    // DONE
	serverConfiguration.HandleFunc("/", routes.ReadServerConfiguration).Methods("GET")       // DONE
	serverConfiguration.HandleFunc("/all", routes.ReadAllServerConfiguration).Methods("GET") // DONE
	serverConfiguration.HandleFunc("/", routes.UpdateServerConfiguration).Methods("PUT")     // DONE
	serverConfiguration.HandleFunc("/", routes.DeleteServerConfiguration).Methods("DELETE")  // DONE

	wgInstance := router.PathPrefix("/wireguard-instance").Subrouter()
	wgInstance.HandleFunc("/", routes.CreateWireguardInstance).Methods("POST")
	wgInstance.HandleFunc("/", routes.ReadWireguardInstance).Methods("GET")

	router.MethodNotAllowedHandler = http.HandlerFunc(setCorsHeader) //if method is not found allow OPTIONS
	return router
}
