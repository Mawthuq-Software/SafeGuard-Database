package db

import (
	"log"
	"os"
	"strconv"

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
	errMigrate := db.AutoMigrate(&Policies{}, &Groups{}, &GroupPolicies{}, &Authentications{}, &Users{}, &Servers{}) //Migrate tables to sqlite
	if errMigrate != nil {
		log.Fatal("Error - Migrating database", errMigrate)
	} else {
		log.Println("Info - Successfully migrated db")
	}
	startupCreation()
}

const (
	//Standard User Permissions
	PERSONAL_KEYS_VIEW      int = 1
	PERSONAL_KEYS_ADD       int = 2
	PERSONAL_KEYS_DELETE    int = 3
	PERSONAL_KEYS_MODIFY    int = 4
	PERSONAL_PASSWORD_RESET int = 5
	PERSONAL_LOGIN          int = 6

	//Admin User Permissions
	KEYS_VIEW_ALL      int = 7
	KEYS_ADD_ALL       int = 8
	KEYS_ADD_OVERRIDE  int = 12
	KEYS_DELETE_ALL    int = 9
	KEYS_MODIFY_ALL    int = 10
	PASSWORD_RESET_ALL int = 11
)

func startupCreation() {

	standardVPNPerms := []int{PERSONAL_KEYS_VIEW, PERSONAL_KEYS_ADD}
	AddPolicy("STANDARD_USER_VPN", standardVPNPerms)

	advancedUserVPNPerms := []int{PERSONAL_KEYS_MODIFY, PERSONAL_KEYS_DELETE}
	AddPolicy("ADVANCED_USER_VPN", advancedUserVPNPerms)

	userSettingPerms := []int{PERSONAL_PASSWORD_RESET, PERSONAL_LOGIN}
	AddPolicy("STANDARD_USER_SETTINGS", userSettingPerms)

	adminPerms := []int{KEYS_VIEW_ALL, KEYS_ADD_ALL, KEYS_ADD_OVERRIDE, KEYS_DELETE_ALL, KEYS_MODIFY_ALL, PASSWORD_RESET_ALL}
	AddPolicy("ADMIN_USER", adminPerms)

	AddGroup("User")

	userPolicies := []string{"STANDARD_USER_VPN", "STANDARD_USER_SETTINGS"}
	AddGroupPolicies("User", userPolicies)

	adminPolicies := []string{"STANDARD_USER_VPN", "STANDARD_USER_SETTINGS", "ADVANCED_USER_VPN", "ADMIN_USER"}
	AddGroupPolicies("Admin", adminPolicies)
}

func AddPolicy(policyName string, perms []int) error {
	db := DBSystem

	totalPerms := ""
	for i := 0; i < len(perms); i++ {
		perm := strconv.Itoa(perms[i])
		totalPerms += perm + ";"
	}

	newPerms := Policies{Name: policyName, Permissions: totalPerms}
	result := db.Create(&newPerms)
	if result.Error != nil {
		log.Println("Warning - Adding policy"+policyName+"to db", result.Error)
		return result.Error
	}
	return nil
}

func AddGroup(groupName string) error {
	db := DBSystem

	newGroup := Groups{Name: groupName}
	resultGroup := db.Create(&newGroup)
	if resultGroup.Error != nil {
		log.Println("Warning - Adding group"+groupName+"to db", resultGroup.Error)
		return resultGroup.Error
	}
	return nil
}

func AddGroupPolicies(groupName string, policyNames []string) error {
	db := DBSystem
	var findGroup Groups

	resFindGroup := db.Where("name = ?", groupName).First(&findGroup)
	if resFindGroup.Error != nil {
		log.Println("Error - Finding group ", resFindGroup.Error)
		return resFindGroup.Error
	}
	// if errors.Is(resFindGroup.Error, gorm.ErrRecordNotFound) {
	// 	processed.Proccessed = false
	// 	processed.Response = "Group was not found on the server"
	// 	return processed
	// }

	for i := 0; i < len(policyNames); i++ {
		var findPolicy Policies
		resFindPol := db.Where("name = ?", policyNames[i]).First(&findPolicy)
		if resFindPol.Error != nil {
			log.Println("Warning - Finding policy"+groupName+"to db", resFindPol.Error)
			return resFindPol.Error
		}

		groupPolicy := GroupPolicies{GroupID: findGroup.ID, PolicyID: findPolicy.ID}
		createGroupPolicy := db.Create(&groupPolicy)
		if createGroupPolicy.Error != nil {
			log.Println("Error - Creating group"+groupName+" policy "+policyNames[i]+"to db", resFindPol.Error)
			return createGroupPolicy.Error
		}
	}
	return nil
}
