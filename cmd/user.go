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

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "A command to manage users.",
	Long:  `A command to manage users in the database.`,
}

var userAddCmd = &cobra.Command{
	Use:     "add USERNAME PASSWORD EMAIL",
	Aliases: []string{"a"},
	Short:   "A command to add a user.",
	Long:    `This command allows you to add new users to the database.`,
	Example: `user add TheClown01 @jackandTheBeanStalk1 theclown01@randommail.com`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 3 {
			return errors.New(`"user add" requires at least 3 arguments`)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		db.DBStart(false)
		username := args[0]
		password := args[1]
		email := args[2]

		err := db.CreateUser(username, password, email)
		if err != nil {
			combinedLogger.Warn("An error occurred when adding user: " + err.Error())
			return
		}

		combinedLogger.Info("Added user successfully")
	},
}

var userReadCmd = &cobra.Command{
	Use:     "read USER_ID",
	Aliases: []string{"r"},
	Short:   "A command to read a user.",
	Long:    `This command allows you to read users in the database.`,
	Example: `user read 1`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New(`"user read" requires at least 1 argument`)
		}
		userID, err := strconv.Atoi(args[0])
		if err != nil || userID < 1 {
			return errors.New("UserID is not valid")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		db.DBStart(false)
		userID, _ := strconv.Atoi(args[0])
		user, err := db.ReadUser(userID)
		if err != nil {
			combinedLogger.Warn("An error occurred when reading user: " + err.Error())
			return
		}
		userKeys, err := db.ReadUserKeys(userID)
		if err != nil {
			combinedLogger.Warn("An error occurred when reading user: " + err.Error())
			return
		}

		combinedLogger.Info("USER INFO")
		combinedLogger.Sugar().Infoln(`AuthID: `, user.AuthID, "\n")

		for i := 0; i < len(userKeys); i++ {
			keyID := userKeys[i].KeyID
			key, err := db.ReadKey(keyID)
			if err != nil {
				combinedLogger.Sugar().Warnf("Error reading key of ID: ", keyID, ""+err.Error())
				continue
			}
			combinedLogger.Sugar().Infoln(`Key ID:`, keyID)
			combinedLogger.Sugar().Infoln(`Server ID:`, key.ServerID)
			combinedLogger.Sugar().Infoln(`Enabled:`, key.Enabled)
			combinedLogger.Sugar().Infoln(`Private IPv4:`, key.PrivateIPv4)
			combinedLogger.Sugar().Infoln(`Private IPv6:`, key.PrivateIPv6)
			combinedLogger.Sugar().Infoln(`Public Key:`, key.PublicKey)
			combinedLogger.Sugar().Infoln(`Used Bandwidth:`, key.UsedBandwidth)
			combinedLogger.Sugar().Infoln(`Total Bandwidth:`, key.TotalBandwidth, "\n")

		}

	},
}

var userDeleteCmd = &cobra.Command{
	Use:     "delete USER_ID",
	Aliases: []string{"d"},
	Short:   "A command to delete a user.",
	Long:    `This command allows you to delete users in the database.`,
	Example: `user delete 1`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New(`"user delete" requires at least 1 argument`)
		}
		userID, err := strconv.Atoi(args[0])
		if err != nil || userID < 1 {
			return errors.New("UserID is not valid")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		db.DBStart(false)
		userID, _ := strconv.Atoi(args[0])
		err := db.DeleteUser(userID)
		if err != nil {
			combinedLogger.Warn("An error occurred when delete user: " + err.Error())
			return
		}
		combinedLogger.Info("Deleted user successfully")
	},
}

func init() {
	rootCmd.AddCommand(userCmd)

	userCmd.AddCommand(userAddCmd)

	userCmd.AddCommand(userReadCmd)

	userCmd.AddCommand(userDeleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// userCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// userCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
