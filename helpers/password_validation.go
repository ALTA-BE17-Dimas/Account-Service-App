package helpers

import (
	"regexp"
	"strings"
)

func ValidatePassword(password string) bool {
	// Check password length
	if len(password) <= 7 {
		return false
	}

	// Check for at least one lowercase letter
	if !strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz") {
		return false
	}

	// Check for at least one uppercase letter
	if !strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		return false
	}

	// Check for at least one special character
	pattern := `[^a-zA-Z0-9]`
	regex := regexp.MustCompile(pattern)

	return regex.MatchString(password)
}
