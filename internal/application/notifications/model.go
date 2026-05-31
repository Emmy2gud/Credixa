package notifications

import (


	"payme/internal/config"
)

type Notification struct {
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
