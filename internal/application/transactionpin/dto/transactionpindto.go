package dto

type SetTransactionPinRequest struct {
	Pin string `json:"pin"`
}

type SetTransactionPinResponse struct {
	Message string `json:"message"`
}
