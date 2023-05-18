package models

import "time"

type TransferHistory struct {
	ID              uint
	UserIDSender    string
	UserIDRecipient string
	Amount          float64
	CreatedAt       time.Time
}
