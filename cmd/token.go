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

var tokenReadCmd = &cobra.Command{
	Use:     "read TOKEN_ID",
	Aliases: []string{"r"},
	Short:   "A command to read server tokens.",
	Long:    `This command allows you to read server tokens in the database by id.`,
	Example: `server token read 1`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New(`"server token read" requires at least 1 argument`)
		}
		tokenID, err := strconv.Atoi(args[0])
		if err != nil || tokenID < 1 {
			return errors.New("TokenID is not valid")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		db.DBStart(false)
		tokenID, _ := strconv.Atoi(args[0])

		token, err := db.ReadToken(tokenID)
		if err != nil {
			combinedLogger.Warn("An error occurred when reading token: " + err.Error())
			return
		}

		combinedLogger.Sugar().Infoln("Token ID: ", token.ID)
		combinedLogger.Sugar().Infoln("Name: ", token.Name)
		combinedLogger.Sugar().Infoln("Hashed token is redacted")
	},
}

var tokenReadNameCmd = &cobra.Command{
	Use:     "name TOKEN_NAME",
	Aliases: []string{"n"},
	Short:   "A command to read server tokens by name.",
	Long:    `This command allows you to read server tokens in the database by name.`,
	Example: `server token read name tokenEU`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New(`"server token read name" requires at least 1 argument`)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		db.DBStart(false)
		tokenName := args[0]

		tokens, err := db.ReadTokensFromName(tokenName)
		if err != nil {
			combinedLogger.Warn("An error occurred when reading tokens: " + err.Error())
			return
		}
		if len(tokens) < 1 {
			combinedLogger.Info("No tokens were found with that name")
			return
		}

		for i := 0; i < len(tokens); i++ {
			token := tokens[i]
			combinedLogger.Sugar().Infoln("Token ID: ", token.ID)
			combinedLogger.Sugar().Infoln("Name: ", token.Name)
			combinedLogger.Sugar().Infoln("Hashed token is redacted\n")
		}
	},
}

var tokenDeleteCmd = &cobra.Command{
	Use:     "delete SERVER_ID",
	Aliases: []string{"d"},
	Short:   "A command to delete server tokens by server ID.",
	Long:    `This command allows you to delete server tokens in the database by server id.`,
	Example: `server token delete 2`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New(`"server token delete" requires at least 1 argument`)
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
		serverToken, err := db.ReadServerTokenFromServerID(serverID)
		if err != nil {
			combinedLogger.Warn("An error occurred when deleting tokens: " + err.Error())
			return
		}

		err = db.DeleteServerToken(serverToken.ID)
		if err != nil {
			combinedLogger.Warn("An error occurred when deleting tokens: " + err.Error())
			return
		}
	},
}

func init() {
	tokenCmd.AddCommand(tokenAddCmd)

	tokenCmd.AddCommand(tokenReadCmd)
	tokenReadCmd.AddCommand(tokenReadNameCmd)

	tokenCmd.AddCommand(tokenDeleteCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tokenCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tokenCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
