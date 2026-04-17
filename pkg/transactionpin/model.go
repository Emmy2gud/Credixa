package transactionpin

import "gorm.io/gorm"

type TransactionPin struct {
	gorm.Model
	UserID uint `json:"user_id"`
	Pin    string `json:"pin"`
}
