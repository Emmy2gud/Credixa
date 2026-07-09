package user

import (
	"gorm.io/gorm"
	"time"
)

// User represents a customer user in the system
type User struct {
	gorm.Model
	FullName    string `json:"full_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password" `
	Role        string `json:"role"` // user, admin
	ProfilePicture string `json:"profile_picture"`
	IsDeleted   bool   `json:"is_deleted"`
	Address     string `json:"address"`
	DateOfBirth string `json:"date_of_birth"`
    NextOfKinName string `json:"next_of_kin_name"`
	NextOfKinRelationship string `json:"next_of_kin_relationship"`
	Status      string `json:"status"` // active, inactive
	IsVerified  bool   `json:"is_verified"`
	FailedAttempts int `json:"failed_attempts"`
    LockedUntil    *time.Time `json:"locked_until"`
	KYCStatus string `json:"kyc_status"` // pending, verified, rejected
	Tier uint `json:"tier"` // 1,2,3

}

type EmailVerification struct {
	gorm.Model
	Email      string  
	FullName   string
	Password   string
	OTP        string
    OTPRequestCount int
    FirstOTPRequestAt time.Time
	ExpiresAt  time.Time
}