package main

import (
	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/api"
	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/config"
	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/db"
	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/logger"
)

func main() {
	combinedLogger := logger.GetCombinedLogger()
	combinedLogger.Info("Firing up Wireguard Manager Authenticator")
	config.LoadConfig()
	db.DBStart()
	api.API()
}
