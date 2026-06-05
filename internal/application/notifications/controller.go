package notifications

import (
	"encoding/json"
	"net/http"
	"payme/internal/api/middleware"
)

type NotificationController struct {
	service NotificationService
}

func NewNotificationController(service NotificationService) *NotificationController {
	return &NotificationController{
		service: service,
	}
}

// GetNotifications returns all notifications for the authenticated user,
// ordered from newest to oldest.
func (h *NotificationController) GetNotifications(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	notifs, err := h.service.GetNotifications(r.Context(), userID)
	if err != nil {
		http.Error(w, "Could not fetch notifications", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifs)
}

// MarkNotificationRead marks a single notification as read for the authenticated user.
func (h *NotificationController) MarkNotificationRead(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var body struct {
		NotificationID uint `json:"notification_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.NotificationID == 0 {
		http.Error(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}

	rowsAffected, err := h.service.MarkNotificationRead(r.Context(), body.NotificationID, userID)
	if err != nil {
		http.Error(w, "Could not update notification", http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Notification not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Notification marked as read",
	})
}
