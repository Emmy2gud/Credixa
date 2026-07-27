package transfer

import (
	"net/http"
	"payme/internal/api/middleware"
	"payme/internal/application/transfer/dto"
	"payme/pkg/utils"
	"strconv"

	"github.com/gorilla/mux"
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

func (h *TransferController) GetUserTransfers (w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	page := vars["page"]
	limit := vars["limit"]
	//coverting limit from string to int
	pageparam, _ := strconv.Atoi(page)

	limitparam, _ := strconv.Atoi(limit)
	offset := (pageparam - 1) * limitparam

	notifs, err := h.service.GetUserTransfers(r.Context(), userID,pageparam, offset, limitparam)
	if err != nil {
		http.Error(w, "Could not fetch notifications", http.StatusInternalServerError)
		return
	}
    utils.JSON(w, http.StatusOK, notifs)
}