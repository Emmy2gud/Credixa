package beneficiary

import "gorm.io/gorm"

type Beneficiary struct {
	gorm.Model
	UserID        uint   `gorm:"index;not null" json:"user_id"`
	AccountName   string `gorm:"not null" json:"account_name"`
	AccountNumber string `gorm:"not null" json:"account_number"`
	BankCode      string `gorm:"not null" json:"bank_code"`
	BankName      string `json:"bank_name"`
}
