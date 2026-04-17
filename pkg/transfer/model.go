package transfer

import "gorm.io/gorm"


type Transfer struct {
	gorm.Model
	ID string `json:"id"`
	SenderID uint `json:"sender_id"`
	ReceiverID uint `json:"receiver_id"`
	Amount float64 `json:"amount"`
	Reference string `json:"reference"`
	Status string `json:"status"`//pending,completed,failed
	
}