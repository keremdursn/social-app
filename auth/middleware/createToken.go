package middleware

import (
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // 72 saat geçerli olacak
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySigningKey)
}