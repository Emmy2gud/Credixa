package splits
import (
	"time"

	"gorm.io/gorm"
)

type SplitBill struct {
	gorm.Model
	ID  uint64 `json:"id"`
    CreatorID uint64 `json:"creator_id"`
    Title string `json:"title"`
    Description string `json:"description"`
    TotalAmount int64 `json:"total_amount"`
    SplitType string `json:"split_type"` // equal, exact_amount, percentage
    Status string `json:"status"`
    ParticipantsCount int `json:"participants_count"`

}

type SplitBillParticipants struct {
	gorm.Model
	ID        uint64 `json:"id"`
	SplitBillID uint64 `json:"split_bill_id"`
	UserID      uint64 `json:"user_id"`
	Amount      int64 `json:"amount"`
	Percentage  uint64 `json:"percentage"`
	Status      string `json:"status"`
	RespondedAt *time.Time `json:"responded_at"`
	PaidAt      *time.Time `json:"paid_at"`
}


