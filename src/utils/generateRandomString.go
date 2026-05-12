package utils

import (
	"crypto/rand"
	"encoding/base64"
)

// creates a random base64 URL-encoded string
func GenerateRandomString(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	// Use URLEncoding to ensure it is safe to put in a URL query parameter
	return base64.URLEncoding.EncodeToString(b), nil
}
