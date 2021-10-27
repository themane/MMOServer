package models

import "github.com/themane/MMOServer/models"

type NotificationData struct {
	Id           string              `json:"_id" bson:"_id"`
	Notification models.Notification `json:"notification" bson:"notification"`
	Sent         bool                `json:"sent" bson:"sent"`
	UserId       string              `json:"user_id" bson:"user_id"`
}

type NotificationRepository interface {
	FindByUserId(userId string) ([]NotificationData, error)
	AddNotification(notification models.Notification, userId string) (*NotificationData, error)
	MarkSent(Id string) error
}
