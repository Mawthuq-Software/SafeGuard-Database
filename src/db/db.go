package db

import (
	"log"
	"os"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DBSystem *gorm.DB

func DBStart() {
	log.Println("Info - Database connection starting")
	errCreateDir := os.MkdirAll("/opt/wgManagerAuth/db", 0755) //create dir if not exist
	if errCreateDir != nil {
		log.Fatal("Error - Creating db directory", errCreateDir)
	}

	user := viper.GetString("DB.USER")
	password := viper.GetString("DB.PASSWORD")
	dbIP := viper.GetString("DB.IP")
	port := viper.GetString("DB.PORT")
	database := viper.GetString("DB.DATABASE")

	dsn := user + ":" + password + "@tcp(" + dbIP + ":" + port + ")/" + database + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Error - Failed to connect database")
	}

	DBSystem = db //set global variable up

	// Migrate the schema
	errMigrate := db.AutoMigrate(&Policies{}, &Groups{}, &GroupPolicies{}, &Authentications{}, &Users{}, &Servers{}, &UserGroups{}) //Migrate tables to sqlite
	if errMigrate != nil {
		log.Fatal("Error - Migrating database", errMigrate)
	} else {
		log.Println("Info - Successfully migrated db")
	}
	startupCreation()
}
