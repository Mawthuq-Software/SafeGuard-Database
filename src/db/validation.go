package db

import "github.com/Mawthuq-Software/Wireguard-Central-Node/src/token"

func ValidateUserPerms(bearerToken string, perms []int) error {
	if bearerToken == "" {
		return ErrBearerEmpty
	}
	username, tokenErr := token.ValidateUser(bearerToken)
	if tokenErr != nil {
		return ErrBearerInvalid
	}
	user, userErr := FindAuthFromUsername(username)
	if userErr != nil {
		return userErr
	}
	userHasPerm := false
	for i := 0; i < len(perms); i++ {
		hasPerms, permErr := CheckPermission(user.ID, perms[i])
		if permErr != nil {
			return permErr
		}
		if hasPerms {
			userHasPerm = true
			break
		}
	}
	if !userHasPerm {
		return ErrInvalidUserPermission
	}

	return nil
}
