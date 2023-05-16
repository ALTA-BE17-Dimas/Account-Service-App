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

func ValidateEmail(email string) (bool, error) {
	if strings.TrimSpace(email) == "" {
		return false, fmt.Errorf("email cannot be empty")
	}

	// Regular expression pattern for email validation
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// Create a regular expression object
	regex := regexp.MustCompile(pattern)

	// Use the MatchString function to check if the email matches the pattern
	if !regex.MatchString(email) {
		return false, fmt.Errorf("invalid email format")
	}

	return true, nil
}

func ValidatePhoneNumber(phoneNumber string) (bool, error) {
	if strings.TrimSpace(phoneNumber) == "" {
		return false, fmt.Errorf("phone number cannot be empty")
	}

	containLowerCase := strings.ContainsAny(phoneNumber, "abcdefghijklmnopqrstuvwxyz")
	containUpperCase := strings.ContainsAny(phoneNumber, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")

	if containLowerCase || containUpperCase {
		return false, fmt.Errorf("phone number cannot contain letters")
	}

	pattern := `[^a-zA-Z0-9]`
	regex := regexp.MustCompile(pattern)

	if regex.MatchString(phoneNumber) {
		return false, fmt.Errorf("phone number cannot contain special characters")
	}

	return true, nil
}

func ValidatePassword(password string) (bool, error) {
	if strings.TrimSpace(password) == "" {
		return false, fmt.Errorf("password cannot be empty")
	}

	if len(password) <= 7 {
		return false, fmt.Errorf("password must be at least 8 characters long")
	}

	if !strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz") {
		return false, fmt.Errorf("password must contain at least one lowercase letter")
	}

	if !strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		return false, fmt.Errorf("password must contain at least one uppercase letter")
	}

	// Check for at least one special character
	pattern := `[^a-zA-Z0-9]`
	regex := regexp.MustCompile(pattern)

	if !regex.MatchString(password) {
		return false, fmt.Errorf("password must contain at least one special character")
	}

	return true, nil
}
