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

// configurationCmd represents the configuration command
var configurationCmd = &cobra.Command{
	Use:     "configuration",
	Aliases: []string{"conf"},
	Short:   "Manages configurations of servers.",
	Long:    "A command which allows you to manage configurations of different servers.",
}

var configurationAddCmd = &cobra.Command{
	Use:     "add NAME DNS SUBNET_MASK NUM_OF_KEYS",
	Aliases: []string{"a"},
	Short:   "A command to add configurations.",
	Long:    `This command allows you to add new configurations to the database.`,
	Example: `configuration add Conf01 1.2.3.4 23 10`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 4 {
			return errors.New(`"configuration add" requires at least 4 arguments`)
		}

		dnsStr := args[1]
		subnet, _ := strconv.Atoi(args[2])
		err := db.CheckConfig(dnsStr, subnet)
		if err != nil {
			return err
		}

		numOfKeys, err := strconv.Atoi(args[3])
		if err != nil || numOfKeys < 1 {
			return errors.New("NumberOfKeys is not a proper integer")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		db.DBStart(false)
		name := args[0]
		dnsStr := args[1]
		subnet, _ := strconv.Atoi(args[2])
		numKeys, _ := strconv.Atoi(args[3])

		err := db.CreateConfiguration(name, dnsStr, subnet, numKeys)
		if err != nil {
			combinedLogger.Warn("An error occurred when adding configuration: " + err.Error())
			return
		}

		combinedLogger.Info("Added configuration successfully")
	},
}

var configurationReadCmd = &cobra.Command{
	Use:     "read CONF_ID",
	Aliases: []string{"r"},
	Short:   "A command to read a configuration.",
	Long:    `This command allows you to read configurations in the database.`,
	Example: `configuration read 1`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New(`"configuration read" requires at least 1 argument`)
		}
		confID, err := strconv.Atoi(args[0])
		if err != nil || confID < 1 {
			return errors.New("ConfigurationID is not valid")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		db.DBStart(false)
		confID, _ := strconv.Atoi(args[0])
		conf, err := db.ReadConfiguration(confID)
		if err != nil {
			combinedLogger.Warn("An error occurred when reading configuration: " + err.Error())
			return
		}

		combinedLogger.Info("CONFIGURATION INFO")
		combinedLogger.Sugar().Infoln(`Name:`, conf.Name)
		combinedLogger.Sugar().Infoln(`DNS:`, conf.DNS)
		combinedLogger.Sugar().Infoln(`Mask:`, conf.Mask)
		combinedLogger.Sugar().Infoln(`Number Of Keys:`, conf.NumberOfKeys, "\n")
	},
}

var configurationReadNameCmd = &cobra.Command{
	Use:     "name NAME",
	Aliases: []string{"n"},
	Short:   "A command to read a configuration by name.",
	Long:    `This command allows you to read configurations in the database by name.`,
	Example: `configuration read name Conf01`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New(`"configuration read name" requires at least 1 argument`)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		db.DBStart(false)
		name := args[0]
		conf, err := db.ReadConfigurationFromName(name)
		if err != nil {
			combinedLogger.Warn("An error occurred when reading configuration: " + err.Error())
			return
		}

		combinedLogger.Info("CONFIGURATION INFO")
		combinedLogger.Sugar().Infoln(`Name:`, conf.Name)
		combinedLogger.Sugar().Infoln(`DNS:`, conf.DNS)
		combinedLogger.Sugar().Infoln(`Mask:`, conf.Mask)
		combinedLogger.Sugar().Infoln(`Number Of Keys:`, conf.NumberOfKeys, "\n")
	},
}

var configurationReadAllCmd = &cobra.Command{
	Use:     "all",
	Aliases: []string{"a"},
	Short:   "A command to read all configurations.",
	Long:    `This command allows you to read all configurations in the database.`,
	Example: `configuration read all`,
	Run: func(cmd *cobra.Command, args []string) {
		db.DBStart(false)
		confs, err := db.ReadAllConfigurations()
		if err != nil {
			combinedLogger.Warn("An error occurred when reading configurations: " + err.Error())
			return
		}

		for i := 0; i < len(confs); i++ {
			combinedLogger.Info("CONFIGURATION INFO")
			combinedLogger.Sugar().Infoln(`ConfigID:`, confs[i].ID)
			combinedLogger.Sugar().Infoln(`Name:`, confs[i].Name)
			combinedLogger.Sugar().Infoln(`DNS:`, confs[i].DNS)
			combinedLogger.Sugar().Infoln(`Mask:`, confs[i].Mask)
			combinedLogger.Sugar().Infoln(`Number Of Keys:`, confs[i].NumberOfKeys, "\n")
		}
		if len(confs) == 0 {
			combinedLogger.Sugar().Info("No configurations found")
		}
	},
}

var configurationDeleteCmd = &cobra.Command{
	Use:     "delete CONF_ID",
	Aliases: []string{"d"},
	Short:   "A command to delete configurations.",
	Long:    `This command allows you to delete configurations in the database.`,
	Example: `configuration delete Conf01`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New(`"configuration read" requires at least 1 argument`)
		}
		confID, err := strconv.Atoi(args[0])
		if err != nil || confID < 1 {
			return errors.New("ConfigurationID is not valid")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		db.DBStart(false)
		confID, _ := strconv.Atoi(args[0])

		err := db.DeleteConfiguration(confID)
		if err != nil {
			combinedLogger.Warn("An error occurred when deleting configuration: " + err.Error())
			return
		}

		combinedLogger.Info("Deleted configuration successfully")
	},
}

func init() {
	rootCmd.AddCommand(configurationCmd)

	configurationCmd.AddCommand(configurationAddCmd)

	configurationCmd.AddCommand(configurationReadCmd)

	configurationReadCmd.AddCommand(configurationReadNameCmd)
	configurationReadCmd.AddCommand(configurationReadAllCmd)

	configurationCmd.AddCommand(configurationDeleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configurationCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configurationCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
