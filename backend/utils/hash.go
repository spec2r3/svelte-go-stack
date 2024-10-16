package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func HashKey(APIKey string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(APIKey), 14)
	return string(bytes), err
}

func CheckPassword(password, hashPassword string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	return err == nil
}

func CheckAPIKey(APIKey, hashKey string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(APIKey), []byte(hashKey))
	return err == nil
}
