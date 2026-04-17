package transactionpin

import (
	"encoding/json"
	"net/http"
	"payme/pkg/config"
	"payme/pkg/middleware"
	"payme/pkg/utils"
	
	
)

func CreateTransactionPin(w http.ResponseWriter, r *http.Request) {
	var tp TransactionPin
	utils.ParseBody(r, &tp)
	// 1️⃣ Get user ID from context
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}
	//pin length is 4 
	if len(tp.Pin) != 4 {
		http.Error(w, "Pin must be 4 digits", http.StatusBadRequest)
		return
	}
	pin, err := utils.HashPassword(tp.Pin)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	tp.UserID = userID

	tp.Pin = pin

	config.DB.Create(&tp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Transaction pin created successfully",
	})

}

func VerifyTransactionPin(w http.ResponseWriter, r *http.Request) {
	// VerifyTransactionPin logic here
}

func UpdateTransactionPin(w http.ResponseWriter, r *http.Request) {
	// UpdateTransactionPin logic here
}

func DeleteTransactionPin(w http.ResponseWriter, r *http.Request) {
	// DeleteTransactionPin logic here
}
