package withdrawal

import (
	"time"

	"gorm.io/gorm"
)

type Withdrawal struct {
	gorm.Model	
	ID            string    `json:"id"`
	UserID        uint    `json:"user_id"`
	WalletID      uint    `json:"wallet_id"`
	TransactionID uint    `json:"transaction_id"`
	Amount        int64     `json:"amount"` // Kobo/Cents
	Fee           int64     `json:"fee"`
	BankName      string    `json:"bank_name"`
	AccountNumber string    `json:"account_number"`
	AccountName   string    `json:"account_name"`
	BankCode      string    `json:"bank_code"`
	Status        string    `json:"status"` // pending, processing, success, failed
	Reference     string    `json:"reference"`
	CreatedAt     time.Time `json:"created_at"`
}
