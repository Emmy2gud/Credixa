
package accounts

import (
	"gorm.io/gorm"
)

type VirtualAccount struct {
	gorm.Model
	WalletID      uint   `json:"wallet_id"`
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
	BankName      string `json:"bank_name"`
	Provider      string `json:"provider"`
	Status        string `json:"status"` // active, inactive
}
