package webtokensapi

import (
	"errors"
	"os"
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
