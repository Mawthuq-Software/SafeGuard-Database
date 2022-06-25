package api

import (
	"net"
	"net/http"

	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/api/router"
	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/logger"
	"github.com/spf13/viper"
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
