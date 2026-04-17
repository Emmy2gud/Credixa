package savings

import (
	"time"

	"gorm.io/gorm"
)


type PersonalSaving struct {
	gorm.Model	
	ID            string    `json:"id"`
	UserID        uint    `json:"user_id"`
	WalletID      uint    `json:"wallet_id"`
	TargetAmount  float64   `json:"target_amount"`
	CurrentAmount float64   `json:"current_amount"`
	Purpose       string    `json:"purpose"`
	Deadline      time.Time `json:"deadline"`
	Status        string    `json:"status"`
	
}

type GroupSaving struct {
	gorm.Model	
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	TargetAmount  float64   `json:"target_amount"`
	CurrentAmount float64   `json:"current_amount"`
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