package models

import "time"

type History struct {
	ID          string
	PhoneNumber string
	Amount      float64
	CreatedAt   time.Time
}
