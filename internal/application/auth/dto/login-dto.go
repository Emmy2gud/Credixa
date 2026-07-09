package dto

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"access_token"`
	FullName string `json:"full_name"`
	Email    string `json:"email" `
}