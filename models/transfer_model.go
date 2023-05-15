package models

type TransferHistory struct {
	BaseModel
	UserIDSender    uint
	UserIDRecipient uint
	Amount          float64
}
