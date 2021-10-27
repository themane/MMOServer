package models

import controllerModels "github.com/themane/MMOServer/controllers/models"

type Notification struct {
	Id           string                        `json:"_id" bson:"_id"`
	Notification controllerModels.Notification `json:"notification" bson:"notification"`
	Sent         bool                          `json:"sent" bson:"sent"`
	UserId       string                        `json:"user_id" bson:"user_id"`
}

type NotificationRepository interface {
	FindByUserId(userId string) ([]Notification, error)
	AddNotification(notification controllerModels.Notification, userId string) (*Notification, error)
	MarkSent(Id string) error
}
