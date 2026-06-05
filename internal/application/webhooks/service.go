package webhooks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"payme/internal/application/webhooks/services"
)

type WebhookService interface {
	ProcessFlutterwaveWebhook(ctx context.Context, w http.ResponseWriter, r *http.Request) error
}

type webhookService struct{}

func NewWebhookService() WebhookService {
	return &webhookService{}
}

func (s *webhookService) ProcessFlutterwaveWebhook(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	signature := r.Header.Get("verif-hash")
	secret := os.Getenv("FLW_SECRET_HASH")
	var payload map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		fmt.Printf("Webhook JSON decode error: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	var source map[string]interface{}
	if data, ok := payload["data"].(map[string]interface{}); ok {
		source = data
	} else {
		source = payload
	}

	if signature == "" || signature != secret {
		w.WriteHeader(http.StatusUnauthorized)
		return nil
	}

	paymentType, _ := source["event.type"].(string)
	switch paymentType {
	case "CARD_TRANSACTION":
		services.HandleFunding(payload, w, r, source)
	case "bank_transfer":
		services.HandleTransfer(payload, w, r, source)
	case "ussd":
		fmt.Println("ussd")
	case "transfer":
		fmt.Println("transfer")
	default:
		http.Error(w, "unknown payment type", http.StatusBadRequest)
		return fmt.Errorf("unknown payment type")
	}

	return nil
}
