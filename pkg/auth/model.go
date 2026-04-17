package auth




type RegisterDTO struct {
	FullName string `json:"full_name" binding:"required" gorm:"column:full_name"`
	PhoneNumber string `json:"phone_number" binding:"required" gorm:"column:phone_number"`
	Email    string `json:"email" binding:"required,email" gorm:"column:email"`
	Password string `json:"password" binding:"required,min=6" gorm:"column:password_hash"`
}

type LoginDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type TransactioPin struct {
	UserID string `json:"user_id" binding:"required"`
	Pin    string `json:"pin" binding:"required"`
	PinAttempts int `json:"pin_attempts"`



}