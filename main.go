package main

import (
	"fmt"

	"gitlab.com/mawthuq-software/wireguard-manager-authenticator/src/config"
	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/api"
	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/db"
	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/logger"
)

func main() {
	fmt.Println("WG MANAGER AND API STARTING UP")

	fmt.Println("Env file loading - 1/4")
	config.LoadConfig()
	fmt.Println("Logger starting up - 2/4")
	logger.LoggerSetup()

	fmt.Println("Starting database - 3/4")
	db.DBStart()

	fmt.Println("Starting API - 4/4")
	api.API()
}
