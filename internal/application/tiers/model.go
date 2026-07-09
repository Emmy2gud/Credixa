package tiers

import (
	"time"

	"gorm.io/gorm"
)

type KYC struct {
	gorm.Model

	UserID uint `json:"user_id"`

	Reference string `json:"reference"`

	Tier            uint   `json:"tier"`
	NINVerified     bool   `json:"nin_verified"`
	AddressVerified bool   `json:"address_verified"`
	Status          string // pending, approved, rejected

	VerifiedBy uint

	VerifiedAt *time.Time

	RejectionReason string
}
