package helpers

import (
	"math/rand"
	"time"
)

func GenerateID() string {
	// Generate a random number or use a specific algorithm to generate a unique ID
	// In this example, we generate a random 6-digit alphanumeric ID
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const idLength = 6

	rand.Seed(time.Now().UnixNano())

	id := make([]byte, idLength)
	for i := range id {
		id[i] = charset[rand.Intn(len(charset))]
	}

	return string(id)
}
