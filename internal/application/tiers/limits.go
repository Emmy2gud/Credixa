package tiers

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type TierLimits struct {
	Tier                 uint  `json:"tier"`
	SingleTransferLimit  int64 `json:"single_transfer_limit"`
	DailyTransferLimit   int64 `json:"daily_transfer_limit"`
	SingleReceivingLimit int64 `json:"single_receiving_limit"`
	DailyReceivingLimit  int64 `json:"daily_receiving_limit"`
}

var DefaultTierLimits = map[uint]TierLimits{
	1: {
		Tier:                 1,
		SingleTransferLimit:  5000000,   // 50,000 NGN (in kobo/cents)
		DailyTransferLimit:   30000000,  // 300,000 NGN
		SingleReceivingLimit: 10000000,  // 100,000 NGN
		DailyReceivingLimit:  50000000,  // 500,000 NGN
	},
	2: {
		Tier:                 2,
		SingleTransferLimit:  20000000,  // 200,000 NGN
		DailyTransferLimit:   500000000, // 5,000,000 NGN
		SingleReceivingLimit: 500000000, // 5,000,000 NGN
		DailyReceivingLimit:  1000000000,// 10,000,000 NGN
	},
	3: {
		Tier:                 3,
		SingleTransferLimit:  500000000, // 5,000,000 NGN
		DailyTransferLimit:   5000000000,// 50,000,000 NGN
		SingleReceivingLimit: 10000000000,// 100,000,000 NGN
		DailyReceivingLimit:  50000000000,// Unlimited/500,000,000 NGN
	},
}

func ValidateTierLimits(db *gorm.DB, userID uint, tier uint, amount int64, isDebit bool) error {
	limits, exists := DefaultTierLimits[tier]
	if !exists {
		limits = DefaultTierLimits[1] // fallback to Tier 1
	}

	if isDebit {
		if amount > limits.SingleTransferLimit {
			return fmt.Errorf("transaction amount exceeds your Tier %d single transfer limit of %d", tier, limits.SingleTransferLimit)
		}

		// Calculate today's total debits
		startOfDay := time.Now().Truncate(24 * time.Hour)
		var dailyTotal int64
		db.Table("transactions").
			Where("user_id = ? AND type = ? AND status = ? AND created_at >= ?", userID, "debit", "completed", startOfDay).
			Select("COALESCE(SUM(amount), 0)").Scan(&dailyTotal)

		if dailyTotal+amount > limits.DailyTransferLimit {
			return fmt.Errorf("transaction exceeds your Tier %d daily limit of %d", tier, limits.DailyTransferLimit)
		}
	} else {
		if amount > limits.SingleReceivingLimit {
			return fmt.Errorf("transaction amount exceeds your Tier %d single receiving limit of %d", tier, limits.SingleReceivingLimit)
		}

		startOfDay := time.Now().Truncate(24 * time.Hour)
		var dailyTotal int64
		db.Table("transactions").
			Where("user_id = ? AND type = ? AND status = ? AND created_at >= ?", userID, "credit", "completed", startOfDay).
			Select("COALESCE(SUM(amount), 0)").Scan(&dailyTotal)

		if dailyTotal+amount > limits.DailyReceivingLimit {
			return fmt.Errorf("transaction exceeds your Tier %d daily receiving limit of %d", tier, limits.DailyReceivingLimit)
		}
	}

	return nil
}
