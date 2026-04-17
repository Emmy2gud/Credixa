package bill_payments

import "gorm.io/gorm"


type BillPayment struct {
	gorm.Model	
	ID        string    `json:"id"`
	UserID    uint    `json:"user_id"`
	WalletID  uint    `json:"wallet_id"`
	BillType  string    `json:"bill_type"` // airtime, data, electricity, cable
	Provider  string    `json:"provider"`  // mtn, glo, dstv, etc.
	Amount    float64   `json:"amount"`
	Token string `json:"token"`
	Reference string    `json:"reference"`
	Status    string    `json:"status"`   // pending, successful, failed
	Metadata  string    `json:"metadata"` // e.g., JSON string for meter number, phone number

}
