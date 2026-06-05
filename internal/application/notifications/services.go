package notifications

import (
	"context"
	"fmt"
	"gorm.io/gorm"
)

type NotificationService interface {
	GetNotifications(ctx context.Context, userID uint) ([]Notification, error)
	MarkNotificationRead(ctx context.Context, notificationID uint, userID uint) (int64, error)
	CreateNotification(ctx context.Context, userID uint, title, notificationType, message string) (Notification, error)
}

type notificationService struct {
	db *gorm.DB
}

func NewNotificationService(db *gorm.DB) NotificationService {
	return &notificationService{
		db: db,
	}
}

func (s *notificationService) GetNotifications(ctx context.Context, userID uint) ([]Notification, error) {
	var notifs []Notification
	if err := s.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("id DESC").
		Find(&notifs).Error; err != nil {
		return nil, err
	}
	return notifs, nil
}

func (s *notificationService) MarkNotificationRead(ctx context.Context, notificationID uint, userID uint) (int64, error) {
	result := s.db.WithContext(ctx).
		Model(&Notification{}).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Update("is_read", true)

	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (s *notificationService) CreateNotification(ctx context.Context, userID uint, title, notificationType, message string) (Notification, error) {
	notif := Notification{
		UserID:  userID,
		Title:   title,
		Type:    notificationType,
		Message: message,
		Status:  "unread",
	}

	if err := s.db.WithContext(ctx).Create(&notif).Error; err != nil {
		return notif, err
	}
	return notif, nil
}

// Keep package-level for Splits package compatibility.
func CreateNotification(userId uint, title, notificationtype , message string) (Notification, error) {
	var notif Notification
	notif.UserID = userId
	notif.Title = title
	notif.Type = notificationtype
	notif.Message = message
	notif.Status = "unread"

	createdNotification, err := notif.CreateNotification()
	if err != nil {
		return notif, fmt.Errorf("failed to create notification: %v", err)
	}

	return *createdNotification, nil
}