package notifications

import (
	"encoding/json"
	"net/http"
	"payme/pkg/config"
	"payme/pkg/middleware"
)

// GetNotifications returns all notifications for the authenticated user,
// ordered from newest to oldest.
func GetNotifications(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var notifs []Notification
	if err := config.GetDB().
		Where("user_id = ?", userID).
		Order("id DESC").
		Find(&notifs).Error; err != nil {
		http.Error(w, "Could not fetch notifications", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifs)
}

// MarkNotificationRead marks a single notification as read for the authenticated user.
func MarkNotificationRead(w http.ResponseWriter, r *http.Request) {
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

	result := config.GetDB().
		Model(&Notification{}).
		Where("id = ? AND user_id = ?", body.NotificationID, userID).
		Update("is_read", true)

	if result.Error != nil {
		http.Error(w, "Could not update notification", http.StatusInternalServerError)
		return
	}
	if result.RowsAffected == 0 {
		http.Error(w, "Notification not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Notification marked as read",
	})
}