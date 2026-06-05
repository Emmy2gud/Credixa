package transfer

import (
	"encoding/json"
	"net/http"
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
	type Input struct {
		AccountNumber string  `json:"account_number"`
		AccountBank   string  `json:"account_bank"`
		Amount        float64 `json:"amount"`
		Narration     string  `json:"narration"`
	}
	var input Input
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	body, err := h.service.InitializeFunding(r.Context(), input.AccountNumber, input.AccountBank, input.Amount, input.Narration)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func (h *TransferController) VerifyFunding(w http.ResponseWriter, r *http.Request) {
	err := h.service.VerifyFunding(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// VerifyFunding logic here
}
