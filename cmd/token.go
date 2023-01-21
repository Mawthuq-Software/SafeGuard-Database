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

// tokenCmd represents the token command
var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "A command to manage tokens",
	Long:  "This command allows you to manage server tokens in the database",
}

var tokenAddCmd = &cobra.Command{
	Use:     "add SERVER_ID",
	Aliases: []string{"a"},
	Short:   "A command to add a server token.",
	Long:    `This command allows you to add a new server token for a server.`,
	Example: `server token add 12`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New(`"server token add" requires at least 1 argument`)
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
		err := db.CreateServerToken(serverID)
		if err != nil {
			combinedLogger.Warn("An error occurred when adding server token: " + err.Error())
			return
		}

		combinedLogger.Info("Added server token successfully")
	},
}

func init() {
	tokenCmd.AddCommand(tokenAddCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tokenCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tokenCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
