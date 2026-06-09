package services

import (
	"fmt"
	"net/http"
	"payme/internal/application/transaction"
	"payme/internal/application/transfer"
	"payme/internal/application/wallet"
	"payme/internal/config"

	"gorm.io/gorm"
)
func HandleTransfer(payload map[string]interface{}, w http.ResponseWriter, source map[string]interface{},db *gorm.DB) {
   var transfer transfer.Transfer
	txRef, ok := source["tx_ref"].(string)
	if !ok {
		txRef, _ = source["txRef"].(string)
	}
	if txRef == "" {
		fmt.Println("Webhook Error: Missing tx_ref")
		w.WriteHeader(http.StatusOK)
		return
	}

	status, _ := source["status"].(string)

	// Find the pending transfer record
	var t transaction.Transaction
	if err := config.DB.Where("reference = ? AND status = ?", txRef, "pending").First(&t).Error; err != nil {
		fmt.Printf("Webhook: no pending transaction for ref %s: %v\n", txRef, err)
		w.WriteHeader(http.StatusOK)
		return
	}

	switch status {
	case "successful":
	
		if err := config.DB.Model(&transaction.Transaction{}).
			Where("reference = ?", txRef).
			Update("status", "completed").Error; err != nil {
				http.Error(w, "failed to update transaction", http.StatusInternalServerError)
				return
			}
		db.Model(&transfer).
			Where("reference = ?", txRef).
			Update("status", "completed")

		fmt.Printf(" Transfer completed for ref: %s\n", txRef)

	case "failed":
	
		if err := wallet.UpdateWalletBalance(uint(t.UserID), t.Amount); err != nil {
			http.Error(w, "failed to refund wallet", http.StatusInternalServerError)
			return
		}
		db.Model(&transaction.Transaction{}).
			Where("reference = ?", txRef).
			Update("status", "failed")
		db.Model(&transfer).
			Where("reference = ?", txRef).
			Update("status", "failed")

		fmt.Printf(" Transfer failed, wallet refunded for ref: %s\n", txRef)

	default:
		fmt.Printf("Webhook Info: unhandled transfer status %s for ref %s\n", status, txRef)
	}

	w.WriteHeader(http.StatusOK)
}