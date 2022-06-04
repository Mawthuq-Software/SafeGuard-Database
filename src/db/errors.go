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
	ErrUserGroupNotFound = errors.New("user's groups were not found")
	ErrCreatingUserGroup = errors.New("could not create new user group")
	//Group
	ErrGroupNotFound = errors.New("group was not found")
	//GroupPolicy
	ErrGroupPolicyNotFound = errors.New("group policies were not found")
	//Authentication
	ErrAuthNotFound  = errors.New("authentication was not found")
	ErrCreatingAuth  = errors.New("error when creating new authentication")
	ErrIncorrectPass = errors.New("password was incorrect")

	//Hashing
	ErrHashing = errors.New("error when hashing")

	//Token
	ErrCreatingToken = errors.New("error creating token")
)

// func ErrorCheck(error) error {
// 	return errors.New(""
// }
