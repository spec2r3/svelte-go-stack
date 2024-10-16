package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// Function to generate a random 16-character hexadecimal API key
func GenerateKey() (string, error) {
	bytes := make([]byte, 8)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
