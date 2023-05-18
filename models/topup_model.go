package models

import "time"

type TopUpHistory struct {
	ID        uint
	UserID    string
	Amount    float64
	CreatedAt time.Time
}
