package models

import "time"

type User struct {
	ID             string
	FullName       string
	IdentityNumber string
	BirthDate      string
	Address        string
	Email          string
	PhoneNumber    string
	Password       string
	Balance        float64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
