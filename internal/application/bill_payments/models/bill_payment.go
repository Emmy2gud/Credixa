package models

import "gorm.io/gorm"

type BillPayment struct {
	gorm.Model
	ID            string `json:"id"`
	UserID        uint   `json:"user_id"`
	TransactionID uint   `json:"transaction_id"`
	WalletID      uint   `json:"wallet_id"`
	BillType      string `json:"bill_type"`
	Provider      string `json:"provider"`
	Amount        uint64 `json:"amount"`
	Token         string `json:"token"`
	Reference     string `json:"reference"`
	Status        string `json:"status"`
	Metadata      string `json:"metadata"`
}
