package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func PassHash(pwd string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), 10)
	return string(bytes), err
}

func CheckPassHash(pwd string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
	return (err == nil)
}
