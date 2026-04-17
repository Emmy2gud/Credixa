package transaction

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	UserID      uint    `json:"user_id"`
	WalletID    uint    `json:"wallet_id"`
	Type        string  `json:"type"`     //credit or debit
	Category    string  `json:"category"` //funding, transfer, airtime, electricity, withdrawal, savings
	Amount      float64 `json:"amount"`
	Fee         float64 `json:"fee"`
	Reference   string  `json:"reference"`
	Status      string  `json:"status"` //pending, successful, failed
	Description string  `json:"description"`
}