package api

import (
	"net"
	"net/http"

	"github.com/spf13/viper"
	"gitlab.com/mawthuq-software/wireguard-manager-authenticator/src/api/router"
	"gitlab.com/mawthuq-software/wireguard-manager-authenticator/src/logger"
)

var combinedLogger = logger.GetCombinedLogger()

func API() {
	newRouter := router.NewRouter()

	serverDev := viper.GetBool("SERVER.SECURITY")
	if !serverDev {
		port := viper.GetString("SERVER.PORT")
		combinedLogger.Info("HTTP about to listen on " + port)

		resolve, _ := net.ResolveTCPAddr("tcp4", "0.0.0.0:"+port)
		resolveTCP, _ := net.ListenTCP("tcp4", resolve)

		errServer := http.Serve(resolveTCP, newRouter)
		combinedLogger.Fatal("Startup of API server " + errServer.Error())
	} else {
		port := viper.GetString("SERVER.PORT")
		fullchainCert := viper.GetString("SERVER.CERT.FULLCHAIN")
		privKeyCert := viper.GetString("SERVER.CERT.PK")

		combinedLogger.Info("HTTPS about to listen on " + port)

		resolve, _ := net.ResolveTCPAddr("tcp4", "0.0.0.0:"+port)
		resolveTCP, _ := net.ListenTCP("tcp4", resolve)

		errServer := http.ServeTLS(resolveTCP, newRouter, fullchainCert, privKeyCert)
		combinedLogger.Fatal("Startup of API server " + errServer.Error())
	}
}
