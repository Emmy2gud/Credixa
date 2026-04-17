package pendingcard

import (
	"encoding/json"

	"gorm.io/gorm"
)

type PendingCard struct {
	gorm.Model
	UserID  uint            `json:"user_id"`
	Payload json.RawMessage `json:"payload"`
	TxRef   string          `json:"tx_ref"`
}
