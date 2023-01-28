package db

import (
	"errors"

	"gorm.io/gorm"
)

func CreateServerToken(serverID int) (err error) {
	findServer, serverErr := ReadServer(serverID)
	if serverErr != nil {
		return serverErr
	}

	_, err = ReadServerTokenFromServerID(serverID)
	if err != ErrServerTokenNotFound {
		return ErrServerTokenExists
	}

	_, tokenID, err := CreateToken(findServer.Name)
	if err != nil {
		return
	}

	token, err := ReadTokensFromName(findServer.Name)
	if (err != ErrTokenNotFound || len(token) > 0) && err != nil {
		return
	}

	err = createServerTokenLink(serverID, tokenID)
	return
}

func createServerTokenLink(serverID int, tokenID int) (err error) {
	db := DBSystem

	newServerToken := ServerTokens{ServerID: serverID, TokenID: tokenID}
	newTokenInfo := db.Create(&newServerToken)
	if newTokenInfo.Error != nil {
		return ErrServerTokenAddingLink
	}
	return
}

func ReadServerToken(serverTokenID int) (serverToken ServerTokens, err error) {
	db := DBSystem

	findInfo := db.Where("id = ?", serverTokenID).First(&serverToken)
	if errors.Is(findInfo.Error, gorm.ErrRecordNotFound) {
		err = ErrServerTokenNotFound
		return
	} else if findInfo.Error != nil {
		err = ErrServerTokenSearch
		return
	}
	return
}

func ReadServerTokenFromServerID(serverID int) (serverToken ServerTokens, err error) {
	db := DBSystem

	findInfo := db.Where("server_id = ?", serverID).First(&serverToken)
	if errors.Is(findInfo.Error, gorm.ErrRecordNotFound) {
		err = ErrServerTokenNotFound
		return
	} else if findInfo.Error != nil {
		err = ErrServerTokenSearch
		return
	}
	return
}

func ReadServerTokenFromTokenID(tokenID int) (serverToken ServerTokens, err error) {
	db := DBSystem

	findInfo := db.Where("token_id = ?", tokenID).First(&serverToken)
	if errors.Is(findInfo.Error, gorm.ErrRecordNotFound) {
		err = ErrServerTokenNotFound
		return
	} else if findInfo.Error != nil {
		err = ErrServerTokenSearch
		return
	}
	return
}

func DeleteServerToken(tokenID int) (err error) {
	serverToken, err := ReadServerTokenFromTokenID(tokenID)
	if err != nil {
		return
	}
	err = deleteServerTokenLink(serverToken)
	if err != nil {
		return
	}
	//final err so no need to check if its nil
	err = DeleteToken(tokenID)
	return
}

// func DeleteServerTokenFromName(name string) (err error) {
// }

func deleteServerTokenLink(serverToken ServerTokens) (err error) {
	db := DBSystem

	deleteInfo := db.Delete(&serverToken)
	if deleteInfo.Error != nil {
		return ErrDeletingServerToken
	}
	return
}
