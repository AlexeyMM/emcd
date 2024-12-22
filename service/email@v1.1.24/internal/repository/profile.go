package repository

import (
	"context"

	"github.com/google/uuid"

	"code.emcdtech.com/emcd/service/email/internal/model"
)

//go:generate mockgen -source=$GOFILE -destination=./mocks/profile.go -package=mockstore

type Profile interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) (model.Profile, error)
	GetNotificationSettings(ctx context.Context, userID uuid.UUID) (model.NotificationSettings, error)
}
