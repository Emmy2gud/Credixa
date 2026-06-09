package transfer

import (
	"encoding/json"
	"net/http"
	"payme/internal/api/middleware"
	"payme/internal/application/transfer/dto"
	"payme/pkg/utils"
)

type TransferController struct {
	service TransferService
}

func NewTransferController(service TransferService) *TransferController {
	return &TransferController{
		service: service,
	}
}

func (h *TransferController) ResolveBankDetails(w http.ResponseWriter, r *http.Request) {
	type Input struct {
		AccountNumber string `json:"account_number"`
		AccountBank   string `json:"account_bank"`
	}
	var input Input
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	body, err := h.service.ResolveBankDetails(r.Context(), input.AccountNumber, input.AccountBank)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func (h *TransferController) InitializeFunding(w http.ResponseWriter, r *http.Request) {
	var input struct {
		AccountNumber string `json:"account_number"`
		AccountBank   string `json:"account_bank"`
		Amount        int64  `json:"amount"`
		Narration     string `json:"narration"`
		IdempotencyKey string `json:"idempotency_key"`
	}

	utils.ParseBody(r, &input)

	userID, ok := middleware.GetUserID(r)

	if !ok {
		http.Error(w,"unauthorized",http.StatusUnauthorized)
		return
	}
	req := dto.CreateTransferRequest{
		UserID:        userID,
		AccountNumber: input.AccountNumber,
		AccountBank:   input.AccountBank,
		Amount:        input.Amount,
		Narration:     input.Narration,
		IdempotencyKey: input.IdempotencyKey,
	}
	body, err := h.service.InitializeFunding(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.JSON(w, http.StatusOK, body)
}

func (h *TransferController) VerifyFunding(w http.ResponseWriter, r *http.Request) {
	err := h.service.VerifyFunding(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// VerifyFunding logic here
}
