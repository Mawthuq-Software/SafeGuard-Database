/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"strconv"

	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/db"
	"github.com/spf13/cobra"
)

// serverConfigurationCmd represents the serverConfiguration command
var serverConfigurationCmd = &cobra.Command{
	Use:     "configuration",
	Aliases: []string{"conf"},
	Short:   "A command to manage server configurations.",
	Long:    `This command allows you to link and unlink server configurations.`,
}

var serverConfigurationAddCmd = &cobra.Command{
	Use:     "add SERVER_ID CONFIG_ID",
	Aliases: []string{"a"},
	Short:   "A command to link server configurations.",
	Long:    `This command allows you to link new server configurations in the database.`,
	Example: `server configuration add 12 92`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New(`"server configuration add" requires at least 2 arguments`)
		}

		serverID, err := strconv.Atoi(args[0])
		if err != nil || serverID < 1 {
			return errors.New("ServerID is not a proper integer")
		}

		configID, err := strconv.Atoi(args[1])
		if err != nil || configID < 1 {
			return errors.New("ConfigID is not a proper integer")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		db.DBStart(false)
		serverID, _ := strconv.Atoi(args[0])
		configID, _ := strconv.Atoi(args[1])

		err := db.CreateServerConfig(serverID, configID)
		if err != nil {
			combinedLogger.Warn("An error occurred when linking server configuration: " + err.Error())
			return
		}

		combinedLogger.Info("Linked server configuration successfully")
	},
}

var serverConfigurationReadCmd = &cobra.Command{
	Use:     "read SERVER_CONFIG_ID",
	Aliases: []string{"r"},
	Short:   "A command to read server configurations.",
	Long:    `This command allows you read server configurations in the database.`,
	Example: `server configuration read 2`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New(`"server configuration read" requires at least 1 argument`)
		}

		serverConfigID, err := strconv.Atoi(args[0])
		if err != nil || serverConfigID < 1 {
			return errors.New("ServerConfigID is not a proper integer")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		db.DBStart(false)
		serverConfigID, _ := strconv.Atoi(args[0])

		serverConfig, err := db.ReadServerConfig(serverConfigID)
		if err != nil {
			combinedLogger.Warn("An error occurred when reading server configuration: " + err.Error())
			return
		}

		combinedLogger.Sugar().Infoln("Server Config ID", serverConfig.ID)
		combinedLogger.Sugar().Infoln("Server ID", serverConfig.ServerID)
		combinedLogger.Sugar().Infoln("Config ID", serverConfig.ConfigID)

		combinedLogger.Info("Read server configuration successfully")
	},
}

var serverConfigurationReadServerIDCmd = &cobra.Command{
	Use:     "serverID SERVER_ID",
	Aliases: []string{"s"},
	Short:   "A command to read server configurations using the serverID.",
	Long:    `This command allows you read server configurations in the database using server ID.`,
	Example: `server configuration read serverID 2`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New(`"server configuration read serverID" requires at least 1 argument`)
		}

		serverID, err := strconv.Atoi(args[0])
		if err != nil || serverID < 1 {
			return errors.New("ServerID is not a proper integer")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		db.DBStart(false)
		serverID, _ := strconv.Atoi(args[0])

		serverConfig, err := db.ReadServerConfigFromServerID(serverID)
		if err != nil {
			combinedLogger.Warn("An error occurred when reading server configuration: " + err.Error())
			return
		}

		combinedLogger.Sugar().Infoln("Server Config ID", serverConfig.ID)
		combinedLogger.Sugar().Infoln("Server ID", serverConfig.ServerID)
		combinedLogger.Sugar().Infoln("Config ID", serverConfig.ConfigID)

		combinedLogger.Info("Read server configuration successfully")
	},
}

var serverConfigurationDeleteCmd = &cobra.Command{
	Use:     "delete SERVER_CONFIG_ID",
	Aliases: []string{"d"},
	Short:   "A command to delete server configurations.",
	Long:    `This command allows you delete server configurations in the database.`,
	Example: `server configuration redeletead 2`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New(`"server configuration read" requires at least 1 argument`)
		}

		serverConfigID, err := strconv.Atoi(args[0])
		if err != nil || serverConfigID < 1 {
			return errors.New("ServerConfigID is not a proper integer")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		db.DBStart(false)
		serverConfigID, _ := strconv.Atoi(args[0])

		err := db.DeleteServerConfig(serverConfigID)
		if err != nil {
			combinedLogger.Warn("An error occurred when reading server configuration: " + err.Error())
			return
		}

		combinedLogger.Info("Deleted server configuration successfully")
	},
}

func init() {
	serverCmd.AddCommand(serverConfigurationCmd)

	serverConfigurationCmd.AddCommand(serverConfigurationAddCmd)

	serverConfigurationCmd.AddCommand(serverConfigurationReadCmd)
	serverConfigurationReadCmd.AddCommand(serverConfigurationReadServerIDCmd)

	serverConfigurationCmd.AddCommand(serverConfigurationDeleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverConfigurationCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverConfigurationCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
