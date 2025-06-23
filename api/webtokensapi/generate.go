package webtokensapi

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(subject string, expires time.Time) (string, error) {
	claims := jwt.MapClaims{
		"sub": subject,
		"exp": expires.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func VerifyToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
}
