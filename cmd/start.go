/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/api"
	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/db"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "This command starts the Safeguard Central Node.",
	Long:  `This command starts the Safeguard Central Node which connects to the database and runs the API server.`,
	Run: func(cmd *cobra.Command, args []string) {
		db.DBStart(true)
		api.API()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
