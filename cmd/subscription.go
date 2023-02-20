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

// subscriptionCmd represents the subscription command
var subscriptionCmd = &cobra.Command{
	Use:     "subscription",
	Aliases: []string{"sub"},
	Short:   "A command to manage subscriptions.",
	Long:    "A command to manage subscriptions in the database.",
}

var subscriptionAddCmd = &cobra.Command{
	Use:     "add NAME NUMBER_OF_KEYS TOTAL_BANDWIDTH_IN_MB",
	Aliases: []string{"a"},
	Short:   "A command to add subscriptions.",
	Long:    `This command allows you to add new subscriptions to the database.`,
	Example: `subscription add basic_subscription 2 10000`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 3 {
			return errors.New(`"subscription add" requires at least 3 arguments`)
		}
		numKeys, err := strconv.Atoi(args[1])
		if err != nil || numKeys < 1 {
			return errors.New("NumberOfKeys is not a proper integer")
		}

		totalBW, err := strconv.Atoi(args[2])
		if err != nil || totalBW < 1 {
			return errors.New("TotalBandwidth is not a proper integer")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		db.DBStart(false)
		name := args[0]
		numKeys, _ := strconv.Atoi(args[1])
		totalBW, _ := strconv.Atoi(args[2])

		err := db.CreateSubscription(name, numKeys, totalBW)
		if err != nil {
			combinedLogger.Warn("An error occurred when adding subscription: " + err.Error())
			return
		}

		combinedLogger.Info("Added subscription successfully")
	},
}

var subscriptionReadCmd = &cobra.Command{
	Use:     "read SUBSCRIPTION_ID",
	Aliases: []string{"r"},
	Short:   "A command to read subscriptions.",
	Long:    `This command allows you to read subscriptions in the database by id.`,
	Example: `subscription read 1`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New(`"subscription read" requires at least 1 argument`)
		}
		subscriptionID, err := strconv.Atoi(args[0])
		if err != nil || subscriptionID < 1 {
			return errors.New("SubscriptionID is not valid")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		db.DBStart(false)
		subscriptionID, _ := strconv.Atoi(args[0])

		subs, err := db.ReadSubscription(subscriptionID)
		if err != nil {
			combinedLogger.Warn("An error occurred when reading subscription: " + err.Error())
			return
		}
		combinedLogger.Sugar().Infoln("Subscription ID: ", subs.ID)
		combinedLogger.Sugar().Infoln("Name: ", subs.Name)
		combinedLogger.Sugar().Infoln("Number of Keys: ", subs.NumberOfKeys)
		combinedLogger.Sugar().Infoln("Total Bandwidth: ", subs.TotalBandwidth)
	},
}

var subscriptionReadNameCmd = &cobra.Command{
	Use:     "name NAME",
	Aliases: []string{"n"},
	Short:   "A command to read subscriptions by the name.",
	Long:    `This command allows you to read subscriptions in the database by name.`,
	Example: `subscription read name basic_subscription`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New(`"subscription read name" requires at least 1 argument`)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		db.DBStart(false)
		name := args[0]

		subs, err := db.ReadSubscriptionByName(name)
		if err != nil {
			combinedLogger.Warn("An error occurred when reading subscription: " + err.Error())
			return
		}
		combinedLogger.Sugar().Infoln("Subscription ID: ", subs.ID)
		combinedLogger.Sugar().Infoln("Name: ", subs.Name)
		combinedLogger.Sugar().Infoln("Number of Keys: ", subs.NumberOfKeys)
		combinedLogger.Sugar().Infoln("Total Bandwidth: ", subs.TotalBandwidth)
	},
}

var subscriptionUpdateCmd = &cobra.Command{
	Use:     "update SUBSCRIPTION_ID",
	Aliases: []string{"u"},
	Short:   "A command to update subscriptions by the id.",
	Long:    `This command allows you to update subscriptions in the database by id.`,
	Example: `subscription update 1`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New(`"subscription update" requires at least 1 argument`)
		}
		subscriptionID, err := strconv.Atoi(args[0])
		if err != nil || subscriptionID < 1 {
			return errors.New("SubscriptionID is not valid")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		db.DBStart(false)
		subscriptionID, _ := strconv.Atoi(args[0])

		subscription, err := db.ReadSubscription(subscriptionID)
		if err != nil {
			combinedLogger.Warn("An error occurred when updating subscription: " + err.Error())
			return
		}

		name, _ := cmd.Flags().GetString("name")
		if name != "" {
			subscription.Name = name
		} else {
			combinedLogger.Warn("Name is not set, not updating")
		}

		numOfKeys, _ := cmd.Flags().GetInt("number_of_keys")
		if numOfKeys > 1 {
			subscription.NumberOfKeys = numOfKeys
		} else if numOfKeys < 1 {
			combinedLogger.Warn("NumberOfKeys is set to less than 1, not updating.")
		}

		totalBW, _ := cmd.Flags().GetInt("total_bandwidth")
		if totalBW > 1 {
			subscription.TotalBandwidth = totalBW
		} else if totalBW < 1 {
			combinedLogger.Warn("TotalBandwidth is set to less than 1, not updating.")
		}

		err = db.UpdateSubscription(subscription)
		if err != nil {
			combinedLogger.Warn("An error occurred when updating subscription: " + err.Error())
			return
		}
	},
}

var subscriptionDeleteCmd = &cobra.Command{
	Use:     "delete SUBSCRIPTION_ID",
	Aliases: []string{"d"},
	Short:   "A command to delete subscriptions.",
	Long:    `This command allows you to delete subscriptions in the database by id.`,
	Example: `subscription delete 10`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New(`"subscription delete" requires at least 1 argument`)
		}
		subscriptionID, err := strconv.Atoi(args[0])
		if err != nil || subscriptionID < 1 {
			return errors.New("SubscriptionID is not valid")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		db.DBStart(false)
		subscriptionID, _ := strconv.Atoi(args[0])

		_, err := db.ReadSubscription(subscriptionID)
		if err != nil {
			combinedLogger.Warn("An error occurred when deleting subscription: " + err.Error())
			return
		}

		userSubs, err := db.ReadUserSubscriptionWithSubscriptionID(subscriptionID)
		if err != nil && err != db.ErrUserSubscriptionsNotFound {
			combinedLogger.Warn("Error when reading subscriptions")
			return
		}

		err = db.DeleteSubscription(subscriptionID)
		if err != nil {
			combinedLogger.Warn("Error when deleting subscription: " + err.Error())
			return
		}

		if len(userSubs) > 0 {
			combinedLogger.Info("User's with subscription exist, remove them first")
			return
		}

		combinedLogger.Info("Deleted subscription")
	},
}

func init() {
	rootCmd.AddCommand(subscriptionCmd)
	subscriptionCmd.AddCommand(subscriptionAddCmd)

	subscriptionCmd.AddCommand(subscriptionReadCmd)
	subscriptionReadCmd.AddCommand(subscriptionReadNameCmd)

	subscriptionCmd.AddCommand(subscriptionUpdateCmd)
	subscriptionUpdateCmd.Flags().StringP("name", "n", "", `The name of the subscription to be updated`)
	subscriptionUpdateCmd.Flags().IntP("number_of_keys", "k", -1, `The number of keys of the subscription to be updated`)
	subscriptionUpdateCmd.Flags().IntP("total_bandwidth", "t", -1, `The total bandwidth of the subscription to be updated`)

	subscriptionCmd.AddCommand(subscriptionDeleteCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// subscriptionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// subscriptionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
