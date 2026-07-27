package transfer

import (
	"payme/internal/config"

	"gorm.io/gorm"
)


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

func GetAllUserTransfers(userID uint,limit int, offset int)([]Transfer , int64) {
	var transfers []Transfer
	var totalRows int64


	config.DB.Model(&Transfer{}).Where("sender_id = ?", userID).Count(&totalRows)

	config.DB.Limit(limit).Offset(offset).Where("user_id = ?", userID).Find(&transfers)
	return transfers,totalRows
}