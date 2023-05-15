package models

import "time"

type User struct {
	BaseModel
	FullName             string
	SingleIdentityNumber string
	BirthDate            time.Time
	Address              string
	Email                string
	PhoneNumber          string
	Password             string
	Balance              float64
	DeletedAt            time.Time
}
