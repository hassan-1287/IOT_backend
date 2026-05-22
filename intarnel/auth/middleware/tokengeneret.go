package middleware

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var Secretjwt = []byte("hassan")

func GenerateToken(userID string) (string, error) {

	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(Secretjwt)
}
