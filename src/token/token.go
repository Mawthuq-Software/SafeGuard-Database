package token

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateNew(toTokenize string, expiryTime time.Time) (string, error) {
	const SecretKey = "ThisIsMySecretKey"
	numericDate := jwt.NewNumericDate(expiryTime)

	claims := jwt.NewWithClaims(
		jwt.SigningMethodHS512,
		jwt.RegisteredClaims{
			Issuer:    toTokenize,
			ExpiresAt: numericDate,
		})

	token, err := claims.SignedString([]byte(SecretKey))
	return token, err
}
