package wallet

import (

	"gorm.io/gorm"
)


type Wallet struct {
	gorm.Model	
	UserID   uint    `json:"user_id"`

	Balance  int64 `json:"balance"`
	Currency string  `json:"currency"`
	Status   string  `json:"status"`
}

type SavingsWallet struct {
	gorm.Model	
	UserID   uint    `json:"user_id"`
	Balance  int64 `json:"balance"`
	Currency string  `json:"currency"`
	Status   string  `json:"status"`
	
}

type WalletLedger struct{
  gorm.Model
  WalletID       uint `json:"wallet_id"`
  UserID         uint `json:"user_id"`
  TransactionID  uint `json:"transaction_id"`
  Amount          int64 `json:"amount"`
  BalanceBefore  int64 `json:"balance_before"`
  BalanceAfter   int64 `json:"balance_after"`
  Description     string `json:"description"`
  Status          string `json:"status"`
 EntryType      string `json:"entry_type"`
}
// type FundingSession struct {
// 	gorm.Model	
// 	UserID         uint      `json:"user_id"`
// 	Amount         float64   `json:"amount"`
// 	PaymentGateway string    `json:"payment_gateway"`
// 	Status         string    `json:"status"` // pending, successful, failed
// 	Reference      string    `json:"reference" gorm:"unique"`
// 	CreatedAt      time.Time `json:"created_at"`
// }
