package db

import (
	"os"

	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/logger"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DBSystem *gorm.DB
var combinedLogger = logger.GetCombinedLogger()

func DBStart() {
	combinedLogger.Info("Database connection starting")

	errCreateDir := os.MkdirAll("/opt/wgManagerAuth/db", 0755) //create dir if not exist
	if errCreateDir != nil {
		combinedLogger.Fatal("Creating db directory " + errCreateDir.Error())
	}

	user := viper.GetString("DB.USER")
	password := viper.GetString("DB.PASSWORD")
	dbIP := viper.GetString("DB.IP")
	port := viper.GetString("DB.PORT")
	database := viper.GetString("DB.DATABASE")

	dsn := user + ":" + password + "@tcp(" + dbIP + ":" + port + ")/" + database + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		combinedLogger.Fatal("Failed to connect database")
	}

	DBSystem = db //set global variable up

	// Migrate the schema
	errMigrate := db.AutoMigrate(
		&Policies{}, &Groups{}, &GroupPolicies{},
		&Authentications{}, &UserGroups{}, &Users{},
		&UserKeys{}, &Keys{}, &Servers{},
		&KeyIPv4{}, &PublicIPv4{}, &PrivateIPv4{},
		&IPv4Interfaces{}, &KeyIPv6{}, &PublicIPv6{},
		&PrivateIPv6{}, &IPv6Interfaces{}, &WireguardInterfaces{},
		&ServerInterfaces{}, &ServerTokens{}, &Tokens{},
		&ServerConfigurations{}, &Configurations{}, &Subscriptions{},
		&UserSubscriptions{}) //Migrate tables to sqlite
	if errMigrate != nil {
		combinedLogger.Fatal("Migrating database " + errMigrate.Error())
	} else {
		combinedLogger.Info("Successfully migrated db")
	}
	startupCreation()
}
