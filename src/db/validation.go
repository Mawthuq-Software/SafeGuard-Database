package db

import "github.com/Mawthuq-Software/Wireguard-Central-Node/src/token"

func ValidateUserPerms(bearerToken string, perms []int) (username string, err error) {
	if bearerToken == "" {
		return "", ErrBearerEmpty
	}
	username, tokenErr := token.ValidateUser(bearerToken)
	if tokenErr != nil {
		err = ErrBearerInvalid
		return
	}
	user, userErr := FindAuthFromUsername(username)
	if userErr != nil {
		err = userErr
		return
	}
	userHasPerm := false
	for i := 0; i < len(perms); i++ {
		hasPerms, permErr := CheckPermission(user.ID, perms[i])
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
