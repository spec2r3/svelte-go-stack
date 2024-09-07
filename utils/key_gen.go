package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// Function to generate a random 16-character hexadecimal API key
func GenerateKey() (string, error) {
	bytes := make([]byte, 8) // 8 bytes * 2 hex digits per byte = 16 hex characters
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
