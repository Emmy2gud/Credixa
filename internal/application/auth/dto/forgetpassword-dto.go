package dto
type ForgotPasswordRequest struct {
	Email       string `json:"email"`
	FrontendUrl string `json:"frontend_url"`
}

type ForgotPasswordResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type ResetPasswordRequest struct {
	Password string `json:"password"`
}

type ResetPasswordResponse struct {
	Message string `json:"message"`
}