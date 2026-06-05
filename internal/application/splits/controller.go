package splits

import (
	"encoding/json"
	"net/http"
	"payme/internal/api/middleware"
	"payme/pkg/utils"
	"strconv"

	"github.com/gorilla/mux"
)

type SplitController struct {
	service SplitService
}

func NewSplitController(service SplitService) *SplitController {
	return &SplitController{service: service}
}

func (h *SplitController) CreateSplitBill(w http.ResponseWriter, r *http.Request) {
	var input CreateSplitBillInput
	utils.ParseBody(r, &input)

	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	splitID, err := h.service.CreateSplitBill(r.Context(), input, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Split bill created successfully",
		"split_id": splitID,
	})
}

func (h *SplitController) GetSplitBills(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	splitBills, err := h.service.GetSplitBills(r.Context(), userID)
	if err != nil {
		http.Error(w, "Could not fetch split bills", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(splitBills)
}

func (h *SplitController) AcceptSplitBill(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	splitID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid split bill ID", http.StatusBadRequest)
		return
	}

	if err := h.service.AcceptSplitBill(r.Context(), splitID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Split bill accepted successfully",
	})
}

func (h *SplitController) DeclineSplitBill(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	splitID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid split bill ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeclineSplitBill(r.Context(), splitID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Split bill declined successfully",
	})
}
