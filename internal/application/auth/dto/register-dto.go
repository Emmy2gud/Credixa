package dto

type SignUpRequest struct {
    ID       string `json:"id"`
    Email    string `json:"email"    binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
	PhoneNumber string `json:"phone_number" binding:"required"`
    FullName string `json:"full_name" binding:"required"`
    Role     string `json:"role"`
    Token  string `json:"token"`
}

type SignUpResponse struct {
    Message string `json:"message"`
    UserID  uint   `json:"user_id"`
	// Token   string `json:"access_token"`
    FullName string `json:"full_name"`
	Email    string `json:"email" `
}

type VerifyOTPRequest struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
      FullName string `json:"full_name"`
}