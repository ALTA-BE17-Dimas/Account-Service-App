package helpers

import "regexp"

func ValidateEmail(email string) bool {
	// Regular expression pattern for email validation
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// Create a regular expression object
	regex := regexp.MustCompile(pattern)

	// Use the MatchString function to check if the email matches the pattern
	return regex.MatchString(email)
}
