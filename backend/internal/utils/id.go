package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateID generates a random hex string of given length (in bytes)
func GenerateID(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

// GenerateToken generates a cryptographic token (24 bytes = 48 hex chars)
func GenerateToken() string {
	return GenerateID(24)
}

// GenerateShortID generates a short ID (8 bytes = 16 hex chars)
func GenerateShortID() string {
	return GenerateID(8)
}
