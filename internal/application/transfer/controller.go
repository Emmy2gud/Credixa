package transfer

import (

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

	var input dto.ResolveBankDetailsRequest
	utils.ParseBody(r, &input)
	body, err := h.service.ResolveBankDetails(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.JSON(w,http.StatusOK,body)
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
