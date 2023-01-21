/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/db"
	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/logger"
	"github.com/spf13/cobra"
)

var combinedLogger = logger.GetCombinedLogger()

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A command to manage servers.",
	Long:  `This command allows you to run different commands to manage servers.`,
}

var serverAddCmd = &cobra.Command{
	Use:     "add NAME REGION COUNTRY IP_ADDRESS",
	Aliases: []string{"a"},
	Short:   "A command to add servers.",
	Long:    `This command allows you to add new servers to the database.`,
	Example: `server add Server01 NYC USA 1.1.1.9`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 4 {
			return errors.New(`"server add" requires at least 4 arguments`)
		}
		ipAddressStr := args[3]
		if net.ParseIP(ipAddressStr) == nil {
			return db.ErrServerIPInvalid
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		db.DBStart(false)
		name := args[0]
		region := args[1]
		country := args[2]
		ipAddress := args[3]

		err := db.CreateServer(name, region, country, ipAddress)
		if err != nil {
			combinedLogger.Warn("An error occurred when adding server: " + err.Error())
			return
		}

		combinedLogger.Info("Added server successfully")
	},
}

var serverReadCmd = &cobra.Command{
	Use:     "read SERVER_ID",
	Aliases: []string{"r"},

	Short:   "A command to read a server.",
	Long:    `This command allows you to read a server from the database.`,
	Example: `server read 1`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New(`"server read" requires at least 1 argument`)
		}
		serverID, err := strconv.Atoi(args[0])
		if err != nil || serverID < 1 {
			return errors.New("ServerID argument is not a valid integer")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		db.DBStart(false)

		serverID, _ := strconv.Atoi(args[0])

		server, err := db.ReadServer(serverID)
		if err != nil {
			combinedLogger.Warn("An error occurred when reading server: " + err.Error())
			return
		}

		fmt.Println("SERVER INFO")
		fmt.Println("ServerID: ", server.ID)
		fmt.Println("Server Name: " + server.Name)
		fmt.Println("Server Region: " + server.Region)
		fmt.Println("Server Country: " + server.Country)
		fmt.Println("Server IP Address: " + server.IPAddress)
		fmt.Println("Server Last Online: ", server.LastOnline)
	},
}

var serverReadNameCmd = &cobra.Command{
	Use:     "name SERVER_NAME",
	Aliases: []string{"n"},
	Short:   "A command to read a server using the name.",
	Long:    `This command allows you to read a server using the server name.`,
	Example: `server read name Server1`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New(`"server read name" requires at least 1 argument`)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		db.DBStart(false)

		serverName := args[0]

		server, err := db.ReadServerFromServerName(serverName)
		if err != nil {
			combinedLogger.Warn("An error occurred when reading server: " + err.Error())
			return
		}

		fmt.Println("SERVER INFO")
		fmt.Println("ServerID: ", server.ID)
		fmt.Println("Server Name: " + server.Name)
		fmt.Println("Server Region: " + server.Region)
		fmt.Println("Server Country: " + server.Country)
		fmt.Println("Server IP Address: " + server.IPAddress)
		fmt.Println("Server Last Online: ", server.LastOnline)
	},
}

var serverUpdateCmd = &cobra.Command{
	Use:     "update SERVER_ID",
	Aliases: []string{"u"},
	Short:   "A command to update a server.",
	Long:    `This command allows you to update a server in the database.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New(`"server update" requires at least 1 argument`)
		}
		serverID, err := strconv.Atoi(args[0])
		if err != nil || serverID < 1 {
			return errors.New("ServerID argument is not a valid integer")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		db.DBStart(false)
		serverID, _ := strconv.Atoi(args[0])

		server, err := db.ReadServer(serverID)
		if err != nil {
			combinedLogger.Warn("An error occurred when updating server: " + err.Error())
			return
		}

		name, _ := cmd.Flags().GetString("name")
		if name != "" {
			server.Name = name
		}

		region, _ := cmd.Flags().GetString("region")
		if region != "" {
			server.Region = region
		}

		country, _ := cmd.Flags().GetString("country")
		if country != "" {
			server.Country = country
		}

		ipAddress, _ := cmd.Flags().GetString("ip_address")
		if ipAddress != "" {
			server.IPAddress = ipAddress
		}

		err = db.UpdateServer(server)
		if err != nil {
			combinedLogger.Warn("An error occurred when updating server: " + err.Error())
			return
		}
	},
}

var serverDeleteCmd = &cobra.Command{
	Use:     "delete SERVER_ID",
	Aliases: []string{"d"},
	Short:   "A command to delete a server.",
	Long:    `This command allows you to delete a server in the database.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New(`"server delete" requires at least 1 argument`)
		}
		serverID, err := strconv.Atoi(args[0])
		if err != nil || serverID < 1 {
			return errors.New("ServerID argument is not a valid integer")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		db.DBStart(false)
		serverID, _ := strconv.Atoi(args[0])

		err := db.DeleteServer(serverID)
		if err != nil {
			combinedLogger.Warn("An error occurred when deleting server: " + err.Error())
			return
		}
		combinedLogger.Info("Deleted server successfully")
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.AddCommand(serverAddCmd)

	serverCmd.AddCommand(serverReadCmd)
	serverReadCmd.AddCommand(serverReadNameCmd)

	serverCmd.AddCommand(serverUpdateCmd)
	serverUpdateCmd.Flags().StringP("name", "n", "", `The name of the server to be updated`)
	serverUpdateCmd.Flags().StringP("region", "r", "", `The region of the server to be updated`)
	serverUpdateCmd.Flags().StringP("country", "c", "", `The country of the server to be updated`)
	serverUpdateCmd.Flags().StringP("ip_address", "i", "", `The ip address of the server to be updated`)

	serverCmd.AddCommand(serverDeleteCmd)

	serverCmd.AddCommand(tokenCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
