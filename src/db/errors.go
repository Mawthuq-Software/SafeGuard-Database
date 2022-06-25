package db

import "errors"

var (
	//General
	ErrQuery = errors.New("error in database query")

	//User
	ErrUserNotFound = errors.New("user was not found")
	ErrUserExists   = errors.New("user already exists")
	ErrCreatingUser = errors.New("error creating new user")
	//User Group
	ErrUserGroupNotFound = errors.New("user's groups was not found")
	ErrCreatingUserGroup = errors.New("could not create new user group")
	//Group
	ErrGroupNotFound = errors.New("group was not found")
	//GroupPolicy
	ErrGroupPolicyNotFound = errors.New("group policies was not found")
	//Policy
	ErrPolicyNotFound = errors.New("policy was nto found")
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
)

// func ErrorCheck(error) error {
// 	return errors.New(""
// }
