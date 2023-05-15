package models

type TopUpHistory struct {
	BaseModel
	UserID uint
	Amount float64
}
