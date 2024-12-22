package model

import "github.com/google/uuid"

type NotificationSettings struct {
	UserID                 uuid.UUID
	Email                  string
	Language               string
	IsTgNotificationsOn    bool
	IsEmailNotificationsOn bool
	TgID                   int64
	WhiteLabelID           uuid.UUID
	IsPushNotificationsOn  bool
}

type ChangeableNotificationSettings struct {
	UserID                 uuid.UUID
	IsTgNotificationsOn    bool
	IsEmailNotificationsOn bool
	TgID                   int64
	IsPushNotificationsOn  bool
}
