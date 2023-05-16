package helpers

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

func ValidateDate(dateStr string) (bool, time.Time, error) {
	// parsing and checking if date is in DD-MM-YYYY format
	birthDate, err := time.Parse("02-01-2006", dateStr)
	if err != nil {
		return false, time.Time{}, fmt.Errorf("failed to parse birth date: %s", err.Error())
	}

	return true, birthDate, nil
}

func ValidateEmail(email string) bool {
	// Regular expression pattern for email validation
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// Create a regular expression object
	regex := regexp.MustCompile(pattern)

	// Use the MatchString function to check if the email matches the pattern
	return regex.MatchString(email)
}

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
