package wallet

import (
	"encoding/json"
	"fmt"
	"net/http"

	"payme/internal/api/middleware"
	"payme/internal/application/wallet/dto"
	"payme/pkg/utils"

	"github.com/gorilla/mux"
)

type ChargeWithTokenRequest struct {
	Token  string `json:"token"`
	Amount int    `json:"amount"`
	Email  string `json:"email"`
	TxRef  string `json:"tx_ref"`
}

type WalletController struct {
	service WalletService
}

func NewWalletController(service WalletService) *WalletController {
	return &WalletController{
		service: service,
	}
}

func (h *WalletController) GetWalletBalance(w http.ResponseWriter, r *http.Request) {
	// GetBalance logic here
}

func (h *WalletController) InitiateWalletFunding(w http.ResponseWriter, r *http.Request) {
	var userCard dto.ChargeRequest
	utils.ParseBody(r, &userCard)

	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	txref, result, err := h.service.InitiateCardWalletFunding(r.Context(), userCard, userID)
	response := map[string]interface{}{
		"tx_ref":      txref,
		"flutterwave": result,
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(result)
	resp, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.JSON(w,http.StatusOK,resp)
}

func (h *WalletController) AuthorizeCardFunding(w http.ResponseWriter, r *http.Request) {
	// The user only needs to send their PIN and the tx_ref from step 1.
	// The full card payload is retrieved automatically from the database.
	var req struct {
		TxRef string `json:"tx_ref"`
		Pin   string `json:"pin"`
	}
	utils.ParseBody(r, &req)

	if req.TxRef == "" || req.Pin == "" {
		http.Error(w, "tx_ref and pin are required", http.StatusBadRequest)
		return
	}

	resp, err := h.service.AuthorizeCardFundingService(r.Context(), req.TxRef, req.Pin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

		utils.JSON(w,http.StatusOK,resp)
}

func (h *WalletController) ValidateWalletFunding(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FlwRef string `json:"flw_ref"`
		Otp    string `json:"otp"`
	}

	utils.ParseBody(r, &req)

	if req.FlwRef == "" || req.Otp == "" {
		http.Error(w, "flw_ref and otp are required", http.StatusBadRequest)
		return
	}

	result, err := h.service.ValidateCardCharge(r.Context(), req.FlwRef, req.Otp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Just return result — DO NOT CREDIT WALLET HERE
	utils.JSON(w,http.StatusOK,map[string]interface{}{
		"message": "Payment processing, awaiting confirmation",
		"data":    result,
		"status":"pending",
	})
}

func (h *WalletController) VerifyCardCharge(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		http.Error(w, "transaction id is required", http.StatusBadRequest)
		return
	}

	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "user not authenticated", http.StatusUnauthorized)
		return
	}

	result, err := h.service.VerifyCard(r.Context(), id, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, _ := json.Marshal(result)
	utils.JSON(w,http.StatusOK,resp)
}
