package model

import "github.com/google/uuid"

type NotificationSettings struct {
	UserID                 uuid.UUID
	Email                  string
	Language               string
	IsEmailNotificationsOn bool
	WhiteLabelID           uuid.UUID
}
