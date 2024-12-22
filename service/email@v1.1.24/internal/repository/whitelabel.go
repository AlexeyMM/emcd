package repository

import (
	"context"

	"github.com/google/uuid"

	"code.emcdtech.com/emcd/service/email/internal/model"
)

//go:generate mockgen -source=$GOFILE -destination=./mocks/whitelabel.go -package=mockstore

type Whitelabel interface {
	GetVersion(ctx context.Context, wlID uuid.UUID) (int, error)
	GetWlFullDomain(ctx context.Context, wlOwnerUserID int32) (string, error)
	GetWlByID(ctx context.Context, wlID uuid.UUID) (*model.Whitelabel, error)
}
