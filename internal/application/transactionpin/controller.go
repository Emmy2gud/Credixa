package transactionpin

import (
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

func (h *transactionPinController) CreateTransactionPin(w http.ResponseWriter, r *http.Request) {
	var tp dto.SetTransactionPinRequest
	utils.ParseBody(r, &tp)
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}
	resp, err := h.service.SetTransactionPin(r.Context(), tp, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.JSON(w, http.StatusOK, resp)

}

func (h *transactionPinController) VerifyTransactionPin(w http.ResponseWriter, r *http.Request) {
	// VerifyTransactionPin logic here
}

func (h *transactionPinController) UpdateTransactionPin(w http.ResponseWriter, r *http.Request) {
	var tp dto.UpdateTransactionPin
	utils.ParseBody(r, &tp)
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}
	resp, err := h.service.UpdateTransactionPin(r.Context(), tp, userID)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, err)
		return
	}

	utils.JSON(w, http.StatusOK, resp)
}

func (h *transactionPinController) DeleteTransactionPin(w http.ResponseWriter, r *http.Request) {
	// DeleteTransactionPin logic here
}
func (h *transactionPinController) ForgotTransactionPin(w http.ResponseWriter, r *http.Request) {
	var tp dto.SetTransactionPinRequest
	utils.ParseBody(r, &tp)
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}
	resp, err := h.service.ForgotTransactionPin(r.Context(), tp, userID)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, err)
		return
	}

	utils.JSON(w, http.StatusOK, resp)
}

func (h *transactionPinController) ResetTransactionPin(w http.ResponseWriter, r *http.Request) {
 	var tp dto.VerifyTransactionPinOTPRequest
	utils.ParseBody(r, &tp)
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}
	resp, err := h.service.ResetTransactionPin(r.Context(), tp, userID)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, err)
		return
	}

	utils.JSON(w, http.StatusOK, resp)
}
