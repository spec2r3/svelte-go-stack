package models

import (
	"crypto/rand"
	"encoding/hex"
)

// Function to generate a random 23-character API key
func GenerateKey() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes)[:23], nil
}
