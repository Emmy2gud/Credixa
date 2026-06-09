package transfer

import "gorm.io/gorm"


type Transfer struct {
	gorm.Model
	ID string `json:"id"`
	SenderID uint `json:"sender_id"`
	ReceiverID *uint `json:"receiver_id"`
	Amount int64 `json:"amount"`
	Reference string `json:"reference"`
	Status string `json:"status"`//pending,completed,failed
	TransferType string `json:"transfer_type"`//internal,external
	AccountNumber string `json:"account_number"`
	BankCode *string `json:"bank_code"`
	AccountName string `json:"account_name"`
	Narration string `json:"narration"`
	
	
	
}

