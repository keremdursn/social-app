package middleware

import (
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateToken(Mail string) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = Mail
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token 24 saat geçerli
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySigningKey)
}
