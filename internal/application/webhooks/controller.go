package webhooks

import (
	"net/http"
)

type WebhookController struct {
	service WebhookService
}

func NewWebhookController(service WebhookService) *WebhookController {
	return &WebhookController{service: service}
}

func (h *WebhookController) FlutterwaveWebhook(w http.ResponseWriter, r *http.Request) {
	if err := h.service.ProcessFlutterwaveWebhook(r.Context(), w, r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
