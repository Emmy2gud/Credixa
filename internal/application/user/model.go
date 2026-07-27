package user

import (
	"gorm.io/gorm"
	"time"
)

type OTPPurpose string

const (
	PurposeSignup      OTPPurpose = "signup"
	PurposePinReset    OTPPurpose = "pin_reset"
	PurposePasswordReset OTPPurpose = "password_reset"
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
	Token string `json:"token"`

}


//this is temporary
type OtpVerification struct {
	gorm.Model
	Email      string  `gorm:"unique" json:"email"`
	FullName   string `json:"full_name"`
	Token      string `json:"token"`
	Password   string `json:"password"` //password or pin
	Purpose    OTPPurpose `json:"purpose"`
	OTP        string `json:"otp"` //hashed otp 
    OTPRequestCount int `json:"otp_request_count"`
    FirstOTPRequestAt time.Time `json:"first_otp_request_at"`
	ExpiresAt  time.Time `json:"expires_at"`
}