package db

import (
	"github.com/Mawthuq-Software/Wireguard-Central-Node/src/token"
)

func ValidatePerms(bearerToken string, perms []int) (userIDInt int, err error) {
	if bearerToken == "" {
		err = ErrBearerEmpty
		return
	}
	userID, tokenErr := token.ValidateUser(bearerToken) //validate bearer token
	if tokenErr != nil {
		err = ErrBearerInvalid
		return
	}
	user, userErr := FindAuthFromUserID(userID) //find user info from bearer
	if userErr != nil {
		err = userErr
		return
	}
	userHasPerm := false
	for i := 0; i < len(perms); i++ {
		hasPerms, permErr := CheckPermission(user.ID, perms[i]) //check to see if matches incoming perms
		if permErr != nil {
			err = permErr
			return
		}
		if hasPerms {
			userHasPerm = true
			break
		}
	}
	if !userHasPerm {
		err = ErrInvalidUserPermission
		return
	}
	return
}
