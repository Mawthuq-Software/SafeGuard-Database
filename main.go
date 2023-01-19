package main

import (
	"github.com/Mawthuq-Software/Wireguard-Central-Node/cmd"
	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/config"
	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/logger"
)

func main() {
	combinedLogger := logger.GetCombinedLogger()
	combinedLogger.Info("Firing up Safeguard Central Node")
	config.LoadConfig()
	cmd.Execute()
}
