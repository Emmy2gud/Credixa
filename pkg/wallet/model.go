package wallet

import (
	"time"

	"gorm.io/gorm"
)


type Wallet struct {
	gorm.Model	
	UserID   uint    `json:"user_id"`
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
	Status   string  `json:"status"`
}

type SavingsWallet struct {
	gorm.Model	
	UserID   uint    `json:"user_id"`
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
	Status   string  `json:"status"`
}

type FundingSession struct {
	gorm.Model	
	UserID         uint      `json:"user_id"`
	Amount         float64   `json:"amount"`
	PaymentGateway string    `json:"payment_gateway"`
	Status         string    `json:"status"` // pending, successful, failed
	Reference      string    `json:"reference" gorm:"unique"`
	CreatedAt      time.Time `json:"created_at"`
}
