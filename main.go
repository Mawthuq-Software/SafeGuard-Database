package main

import (
	"fmt"

	"gitlab.com/mawthuq-software/wireguard-manager-authenticator/src/api"
	"gitlab.com/mawthuq-software/wireguard-manager-authenticator/src/config"
	"gitlab.com/mawthuq-software/wireguard-manager-authenticator/src/db"
	"gitlab.com/mawthuq-software/wireguard-manager-authenticator/src/logger"
)

func main() {
	fmt.Println("WG MANAGER AND API STARTING UP")

	fmt.Println("Logger starting up - 1/4")
	logger.LoggerSetup()

	fmt.Println("Env file loading - 2/4")
	config.LoadConfig()

	fmt.Println("Starting database - 3/4")
	db.DBStart()

	fmt.Println("Starting API - 4/4")
	api.API()
}
