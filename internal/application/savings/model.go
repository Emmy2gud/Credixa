package savings

import (
	"time"

	"gorm.io/gorm"
)


type PersonalSaving struct {
    gorm.Model
    UserID        uint      `json:"user_id" gorm:"not null;index"`
    WalletID      uint      `json:"wallet_id" gorm:"not null"`
    TargetAmount  uint64   `json:"target_amount" gorm:"type:decimal(15,2);not null"`
    CurrentAmount uint64   `json:"current_amount" gorm:"type:decimal(15,2);default:0"`
    Purpose       string    `json:"purpose" gorm:"not null"` // e.g. "Rent", "Laptop"
    Status        string    `json:"status" gorm:"default:'active'"` // active | completed | cancelled
    AutoSave      bool      `json:"auto_save" gorm:"default:false"`
    AutoSaveFrequency string `json:"auto_save_frequency,omitempty"` // daily | weekly | monthly
    AutoSaveAmount    uint64 `json:"auto_save_amount,omitempty" gorm:"type:decimal(15,2)"`
	LastAutoSaveDate  time.Time `json:"last_auto_save_date"`
}

type GroupSaving struct {
	gorm.Model	
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	TargetAmount  uint64   `json:"target_amount"`
	CurrentAmount uint64   `json:"current_amount"`
	Status        string    `json:"status"`
	
}

type GroupMember struct {
	GroupID      uint  `json:"group_id"`
	UserID       uint  `json:"user_id"`
	Contribution float64 `json:"contribution"`
	Role string `json:"role"`//admin, member
}




type GroupContribution struct {
	ID string `json:"id"`
	GroupID uint `json:"group_id"`
	UserID uint `json:"user_id"`
	Amount float64 `json:"amount"`
	Reference string `json:"reference"`
	
}