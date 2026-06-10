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
	Status      string `json:"status"` // active, inactive
	isVerified  bool   `json:"is_verified"`
	

}

type EmailVerification struct {
	gorm.Model
	Email      string
	FullName   string
	Password   string
	OTP        string
	ExpiresAt  time.Time
}