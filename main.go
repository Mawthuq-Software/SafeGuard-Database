package main

import (
	"gitlab.com/mawthuq-software/wireguard-manager-authenticator/src/api"
	"gitlab.com/mawthuq-software/wireguard-manager-authenticator/src/config"
	"gitlab.com/mawthuq-software/wireguard-manager-authenticator/src/db"
	"gitlab.com/mawthuq-software/wireguard-manager-authenticator/src/logger"
)

func main() {
	combinedLogger := logger.GetCombinedLogger()
	combinedLogger.Info("Firing up Wireguard Manager Authenticator")
	config.LoadConfig()
	db.DBStart()
	api.API()
}
