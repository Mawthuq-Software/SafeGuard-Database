package db

import "errors"

var (
	//General
	ErrQuery = errors.New("error in database query")

	//User
	ErrUserNotFound = errors.New("user was not found")
	ErrUserExists   = errors.New("user already exists")
	ErrCreatingUser = errors.New("error creating new user")
	//Authentication
	ErrAuthNotFound  = errors.New("authentication was not found")
	ErrCreatingAuth  = errors.New("error when creating new authentication")
	ErrIncorrectPass = errors.New("password was incorrect")
	//Group
	ErrGroupNotFound = errors.New("group was not found")

	//Hashing
	ErrHashing = errors.New("error when hashing")

	//Token
	ErrCreatingToken = errors.New("error creating token")
)

// func ErrorCheck(error) error {
// 	return errors.New(""
// }
