package models

import "time"

type User struct {
	BaseModel
	FullName       string
	IdentityNumber string
	BirthDate      string
	Address        string
	Email          string
	PhoneNumber    string
	Password       string
	Balance        float64
	UpdatedAt      time.Time
}
