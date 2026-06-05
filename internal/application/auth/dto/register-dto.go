package dto

type SignUpRequest struct {
    ID       string `json:"id"`
    Email    string `json:"email"    binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
    FullName string `json:"full_name" binding:"required"`
    Role     string `json:"role"`
}

type SignUpResponse struct {
    Message string `json:"message"`
    UserID  uint   `json:"user_id"`
	Token   string `json:"access_token"`
}