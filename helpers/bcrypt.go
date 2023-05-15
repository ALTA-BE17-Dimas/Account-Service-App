package helpers

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// DefaultSalt is the default number of rounds of hashing used to hash the password.
const DefaultSalt = 8

// HashPass returns the hashed password for the given plain text password.
func HashPass(p string) string {
	password := []byte(strings.TrimSpace(p))
	hash, _ := bcrypt.GenerateFromPassword(password, DefaultSalt)

	return string(hash)
}

// ComparePass compares the given hashed password with the given plain text password.
// It returns an error if the comparison fails.
func ComparePass(h, p []byte) error {
	hash, pass := []byte(h), []byte(p)

	err := bcrypt.CompareHashAndPassword(hash, pass)

	return err
}
