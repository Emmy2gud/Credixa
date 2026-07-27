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
type WebhookTransferService interface {
	
	HandleFunding(payload map[string]interface{}, w http.ResponseWriter, r *http.Request, source map[string]interface{})
}
type webHookTransferService struct {
	db *gorm.DB
}

func NewWebhookTransferService() WebhookTransferService {
	return &webHookTransferService{
		db: config.DB,
	}
}

func(wh *webHookTransferService) HandleFunding(payload map[string]interface{}, w http.ResponseWriter, r *http.Request, source map[string]interface{}) {
	
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
  	err := wh.db.Transaction(func(tx *gorm.DB) error {
		if err := wh.db.Model(&transaction.Transaction{}).
			Where("reference = ?", txRef).
			Update("status", "completed").Error; err != nil {
			
			return fmt.Errorf("failed to update transaction")
		}
		wh.db.Model(&transfer).
			Where("reference = ?", txRef).
			Update("status", "completed")

		fmt.Printf(" Transfer completed for ref: %s\n", txRef)
		return nil
	})
	if err != nil {
		http.Error(w, "failed to update transaction", http.StatusInternalServerError)
		return
	}	
	case "failed":
 	err := wh.db.Transaction(func(tx *gorm.DB) error {
		if err := wallet.UpdateWalletBalance(tx,t.WalletID, 0, t.Amount, txRef, "transfer", "Debit", "failed", "transfer"); err != nil {
			return fmt.Errorf("failed to refund wallet")
		}

		wh.db.Model(&transfer).
			Where("reference = ?", txRef).
			Update("status", "failed")
		return nil
	})
	if err != nil {
		http.Error(w, "failed to update transaction", http.StatusInternalServerError)
		return
	}
		fmt.Printf(" Transfer failed, wallet refunded for ref: %s\n", txRef)

	default:
		fmt.Printf("Webhook Info: unhandled transfer status %s for ref %s\n", status, txRef)
	}

	w.WriteHeader(http.StatusOK)
}
