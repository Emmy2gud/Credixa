package notifications

import (
	"payme/internal/config"

	"gorm.io/gorm"
)

type Notification struct {
	gorm.Model
	ID        uint       `json:"id"`
	UserID    uint       `json:"user_id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	IsRead    bool      `json:"is_read"`
	Status    string    `json:"status"`// unread, read
	Type    string    `json:"type"`//splits,bills-payment,wallet-fund,transfers

}

func (n *Notification) CreateNotification() (*Notification, error) {
	err := config.GetDB().Create(&n).Error
	if err != nil {
		return nil, err
	}
	return n, nil
}


func GetAllNotifications(userID uint,limit int, offset int)([]Notification , int64) {
	var notifications []Notification
	var totalRows int64


	config.GetDB().Model(&Notification{}).Where("user_id = ?", userID).Count(&totalRows)

	config.GetDB().Limit(limit).Offset(offset).Where("user_id = ?", userID).Find(&notifications)
	return notifications,totalRows
}