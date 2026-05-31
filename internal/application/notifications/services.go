package notifications

import "fmt"

func CreateNotification(userId uint, title, notificationtype , message string)(Notification, error){
var Notification Notification

Notification.UserID = userId
Notification.Title = title
Notification.Type = notificationtype
Notification.Message = message
Notification.Status = "unread"

createdNotification, err := Notification.CreateNotification()

if err != nil {
	return Notification, fmt.Errorf("failed to create notification: %v", err)
}

return *createdNotification, nil


}