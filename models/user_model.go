package models

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
}
