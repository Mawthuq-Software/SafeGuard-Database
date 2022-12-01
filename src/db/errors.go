package db

import "errors"

var (
	//General
	ErrQuery = errors.New("error in database query")

	//User
	ErrUserNotFound = errors.New("user was not found")
	ErrUserExists   = errors.New("user already exists")
	ErrCreatingUser = errors.New("error creating new user")
	//UserGroup
	ErrUserGroupNotFound = errors.New("user's groups was not found")
	ErrCreatingUserGroup = errors.New("could not create new user group")
	//UserSubscription
	ErrUserSubscriptionNotFound   = errors.New("user does not have subscription")
	ErrUserSubscriptionsNotFound  = errors.New("userSubscriptions not found or table empty")
	ErrUserSubscriptionValidation = errors.New("userID does not match userSubscription userID")
	ErrUserSubscriptionExists     = errors.New("user already has a subscription")
	ErrCreatingUserSubscription   = errors.New("error creating a user subscription")
	ErrUsersSubscriptionExists    = errors.New("user(s) with existing user subscriptions")
	//Subscription
	ErrSubscriptionNotFound      = errors.New("subscription was not found")
	ErrSubscriptionExists        = errors.New("subscription exists")
	ErrSubscriptionExpired       = errors.New("subscription has expired")
	ErrSubscriptionKeyMaxed      = errors.New("user has maxed out all available keys")
	ErrSubscriptionUserSubExists = errors.New("user subscriptions with subscription ID exist, delete them first")
	//Group
	ErrGroupNotFound = errors.New("group was not found")
	//GroupPolicy
	ErrGroupPolicyNotFound = errors.New("group policies was not found")
	//Policy
	ErrPolicyNotFound = errors.New("policy was not found")
	//Authentication
	ErrAuthNotFound   = errors.New("authentication was not found")
	ErrCreatingAuth   = errors.New("error when creating new authentication")
	ErrIncorrectPass  = errors.New("password was incorrect")
	ErrSavingPassword = errors.New("error when saving user password")
	//Hashing
	ErrHashing = errors.New("error when hashing")
	//Token or Bearer
	ErrCreatingToken = errors.New("error creating token")
	ErrBearerEmpty   = errors.New("bearer token is empty")
	ErrBearerInvalid = errors.New("bearer token is invalid")
	//Permissions
	ErrMissingPermission     = errors.New("entity does not have permission to perform this function")
	ErrInvalidUserPermission = errors.New("user does not have valid permission")
	//Server
	ErrServerNotFound  = errors.New("server not found")
	ErrServerIPInvalid = errors.New("ip address is not valid")
	ErrServerSaving    = errors.New("an error occurred when saving")
	ErrDeletingServer  = errors.New("an error occurred when deleting server")
	ErrServerKeysExist = errors.New("keys with this server exist, delete them first")
	ErrServerNameExist = errors.New("server with same name already exists")
	//Key
	ErrCreatingKey               = errors.New("issue creating key in database")
	ErrDeletingKey               = errors.New("error deleting key from database")
	ErrPublicKeyIncorrectForm    = errors.New("public key is not in the correct form")
	ErrPresharedKeyIncorrectForm = errors.New("preshared key is not in the correct form")
	ErrKeyNotFound               = errors.New("key(s) was not found")
	//UserKey
	ErrUserKeyNotFound = errors.New("user's keys were not found")
	ErrUserKeyLink     = errors.New("error creating user key link")
	ErrDeletingUserKey = errors.New("error deleting user key link")
	//Configuration
	ErrConfMaskInvalid = errors.New("mask is invalid")
	ErrConfDNSInvalid  = errors.New("DNS is invalid")
	ErrCreatingConf    = errors.New("an error occurred when creating new configuration")
	ErrConfNotFound    = errors.New("configuration(s) was not found")
	//ServerConfiguration
	ErrServerConfNotFound      = errors.New("server configuration(s) not found")
	ErrServerConfAlreadyExists = errors.New("server configuration already exists on server")
	//Token
	ErrGeneratingToken = errors.New("an error occurred when creating a new token")
	ErrAddingToken     = errors.New("an error occurred when adding a new token to the database")
	ErrTokenNotFound   = errors.New("token was not found")
	//ServerToken
	ErrServerTokenAddingLink = errors.New("an error occurred when creating a server token link")
	ErrServerTokenNotFound   = errors.New("server token does not exist")
	ErrServerTokenSearch     = errors.New("an error occurred when searching for the server token")
	ErrDeletingServerToken   = errors.New("an error occurred when deleting the server token link")
	ErrServerTokenExists     = errors.New("server already has an existing token")
	//WireguardInterface
	ErrCreatingWireguardInterface = errors.New("an error occurred when creating a new wireguard interface")
	ErrFindingWireguardInterface  = errors.New("an error occurred when searching for the wireguard interface")
	ErrDeletingWireguardInterface = errors.New("an error occurred when deleting the wireguard interface")
	//ServerInterface
	ErrCreatingServerInterface    = errors.New("an error occurred when creating a new server interface link")
	ErrFindingServerInterface     = errors.New("an error occurred when searching for the server interface")
	ErrKeysExistOnServerInterface = errors.New("keys exist on the server interface, delete them first")
	ErrDeletingServerInterface    = errors.New("an error occurred when deleting the server interface")
)

// func ErrorCheck(error) error {
// 	return errors.New(""
// }
