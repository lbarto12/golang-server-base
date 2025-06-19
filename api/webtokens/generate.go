package webtokens

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey []byte

func Init() error {
	secretKey, exists := os.LookupEnv("JWT_SECRET_KEY")
	if !exists {
		return errors.New("JWT_SECRET_KEY environment variable not set")
	}
	jwtKey = []byte(secretKey)
	return nil
}

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
