package utils

import (
	// "crypto/rand"
	// "encoding/base64"
	// "log"
	// "golang.org/x/crypto/bcrypt"

	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(email string) (string, error) {
	jwtSecret := []byte(os.Getenv("JWTSECRET"))
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}

func GetGmailToken(email string) ([]byte, error) {
	fmt.Printf("gmail token getter reached")
	return []byte{}, nil
}
