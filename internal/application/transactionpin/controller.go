package transactionpin

import (
	"encoding/json"
	"net/http"
	"payme/internal/api/middleware"
	"payme/internal/application/transactionpin/dto"
	"payme/pkg/utils"
)
type transactionPinController struct {
	service TransactionPinService
}

func NewTransactionPinController(service TransactionPinService) *transactionPinController {
	return &transactionPinController{
		service: service,
	}
}

func(h *transactionPinController) CreateTransactionPin(w http.ResponseWriter, r *http.Request) {
	var tp dto.SetTransactionPinRequest
	utils.ParseBody(r, &tp)
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}
   resp,err :=h.service.SetTransactionPin(r.Context(),tp,userID)
   if err != nil {
	   	http.Error(w, err.Error(), http.StatusUnauthorized)
		return
   }


	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)

}

func (h *transactionPinController) VerifyTransactionPin(w http.ResponseWriter, r *http.Request) {
	// VerifyTransactionPin logic here
}

func (h *transactionPinController) UpdateTransactionPin(w http.ResponseWriter, r *http.Request) {
	// UpdateTransactionPin logic here
}

func DeleteTransactionPin(w http.ResponseWriter, r *http.Request) {
	// DeleteTransactionPin logic here
}
