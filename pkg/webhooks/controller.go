package webhooks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"payme/pkg/webhooks/services"
)

func FlutterwaveWebhook(w http.ResponseWriter, r *http.Request) {
	// 1. Verify signature
	signature := r.Header.Get("verif-hash")
	secret := os.Getenv("FLW_SECRET_HASH")
	var payload map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		fmt.Printf("Webhook JSON decode error: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// fmt.Printf("--- Incoming Flutterwave Webhook ---\nPayload: %+v\n", payload)
	pretty, _ := json.MarshalIndent(payload, "", "  ")
	fmt.Println("Webhook JSON:\n", string(pretty))

	// 1. Extract transaction data
	// Some Flutterwave webhooks use a nested 'data' object, others use a flat root structure.
	var source map[string]interface{}
	//check if the payload has a data field
	if data, ok := payload["data"].(map[string]interface{}); ok {
		source = data
	} else {
		source = payload
	}

	if signature == "" || signature != secret {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	paymentType := source["event.type"].(string)
	fmt.Println("payment type is", paymentType)
	switch paymentType {
	case "CARD_TRANSACTION":
		services.HandleFunding(payload, w, r, source)
	case "bank_transfer":
		// services.HandleTransfer(payload, w, r,source)
		fmt.Println("bank transfer")
	case "ussd":
		// services.HandleBillapay(payload, w, r,source)
		fmt.Println("ussd")
	case "transfer":
		// services.HandleWithdrawal(payload, w, r,source)
		fmt.Println("transfer")
	default:
		http.Error(w, "unknown payment type", http.StatusBadRequest)
		return
	}
}
