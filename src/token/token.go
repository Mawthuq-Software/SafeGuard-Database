package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

func GenerateUser(userID string, expiryTime time.Time) (string, error) {
	issuer := viper.GetString("TOKEN.STRING")
	secretKey := viper.GetString("TOKEN.SECRETKEY")
	numericDate := jwt.NewNumericDate(expiryTime)

	claims := jwt.MapClaims{
		"UserID": userID,
		"RegisteredClaims": jwt.RegisteredClaims{
			Issuer:    issuer,
			ExpiresAt: numericDate,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	signedToken, err := token.SignedString([]byte(secretKey))
	return signedToken, err
}

func ValidateUser(tokenStr string) (string, error) {
	secretKey := viper.GetString("TOKEN.SECRETKEY")
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secretKey), nil
	})

	claims, valid := token.Claims.(jwt.MapClaims)
	if valid && token.Valid {
		return claims["UserID"].(string), nil
	} else {
		return "", err
	}
}
