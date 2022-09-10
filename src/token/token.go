package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

type customClaims struct {
	UserID int `json:"userID"`
	jwt.RegisteredClaims
}

func GenerateUser(userID int, expiryTime time.Time) (string, error) {
	issuer := viper.GetString("TOKEN.STRING")
	secretKey := viper.GetString("TOKEN.SECRETKEY")
	numericDate := jwt.NewNumericDate(expiryTime)

	claims := customClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			ExpiresAt: numericDate,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	signedToken, err := token.SignedString([]byte(secretKey))
	return signedToken, err
}

func ValidateUser(tokenStr string) (userID int, err error) {
	secretKey := viper.GetString("TOKEN.SECRETKEY")
	token, err := jwt.ParseWithClaims(tokenStr, &customClaims{},
		func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		})
	if err != nil {
		return
	}
	claims, valid := token.Claims.(*customClaims)
	//NEED TO PARSE TOKEN!!
	if valid && token.Valid {
		userIDTokenInterface := claims.UserID
		if userIDTokenInterface == 0 {
			return 0, fmt.Errorf("token is not valid")
		}
		return userIDTokenInterface, nil
	} else {
		return 0, fmt.Errorf("token is not valid")
	}
}
