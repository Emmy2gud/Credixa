package services

import (
	"fmt"
	"net/http"
	"payme/pkg/config"
	"payme/pkg/transaction"
	"payme/pkg/wallet"
	"strconv"
)

func HandleFunding(payload map[string]interface{}, w http.ResponseWriter, r *http.Request, source map[string]interface{}) {


	// 2. Map fields (supports both snake_case and camelCase)
	event, _ := payload["event"].(string)
	if event == "" {
		event, _ = payload["event.type"].(string)
	}

	txRef, ok := source["tx_ref"].(string)
	if !ok {
		txRef, _ = source["txRef"].(string)
	}

	status, ok := source["status"].(string)
	if !ok {
		status, _ = source["status_code"].(string)
	}

	var amount float64
	switch v := source["amount"].(type) {
	case float64:
		amount = v
	case int:
		amount = float64(v)
	case string:
		if val, err := strconv.ParseFloat(v, 64); err == nil {
			amount = val
		}
	}

	fmt.Printf("Webhook parsed: event=%s, status=%s, txRef=%s, amount=%f\n", event, status, txRef, amount)

	if txRef == "" {
		fmt.Println("Webhook Error: Missing tx_ref")
		w.WriteHeader(http.StatusOK)
		return
	}

	// 3. Only process successful payments
	if status != "successful" {
		fmt.Printf("Webhook Info: Status is %s, skipping wallet update\n", status)
		w.WriteHeader(http.StatusOK)
		return
	}

	var session transaction.Transaction

	err := config.DB.Where("reference = ? AND status = ?", txRef, "pending").First(&session).Error
	if err != nil {
		fmt.Printf("Webhook DB Error: Could not find pending transaction with reference %s: %v\n", txRef, err)
		w.WriteHeader(http.StatusOK)
		return
	}

	// 5. Credit wallet (IDEMPOTENT)
	fmt.Printf("Crediting wallet for user %d with amount %f\n", session.UserID, session.Amount)
	if err := wallet.UpdateWalletBalance(uint(session.UserID), session.Amount); err != nil {
		fmt.Printf("Webhook Error: Failed to update wallet balance: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 6. Mark as successful
	session.Status = "successful"
	if err := config.DB.Save(&session).Error; err != nil {
		fmt.Printf("Webhook Error: Failed to save transaction status: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Printf("✅ Success: Wallet credited and transaction updated for ref: %s\n", txRef)
	w.WriteHeader(http.StatusOK)
}