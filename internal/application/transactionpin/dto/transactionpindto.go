package dto

type SetTransactionPinRequest struct {
	Pin string `json:"pin"`
}

type UpdateTransactionPin struct{
	NewPin string `json:"new_pin"`
	OldPin string `json:"old_pin"`
}
type SetTransactionPinResponse struct {
	Message string `json:"message"`
}

type VerifyTransactionPinOTPRequest struct {
    OTP      string `json:"otp"`
    Pin string `json:"pin"`
}