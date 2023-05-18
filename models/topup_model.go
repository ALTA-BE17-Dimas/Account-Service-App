package models
import "time"

type TopUpHistory struct {
	BaseModel
	UserID 		uint
	Amount 		float64
	CreatedAt 	time.Time
}
