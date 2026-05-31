package user

import (
	"gorm.io/gorm"
)

// User represents a customer user in the system
type User struct {
	gorm.Model
	FullName string `json:"full_name"`
	Email    string `json:"email" `
	Password string `json:"password" `
	Role     string `json:"role"` // user, admin

}
