package repository

import (
	"context"

	"github.com/google/uuid"

	"code.emcdtech.com/emcd/service/email/internal/model"
)

//go:generate mockgen -source=$GOFILE -destination=./mocks/provider_settings.go -package=mockstore

type ProvideSettings interface {
	Create(ctx context.Context, setting model.Setting) error
	Get(ctx context.Context, whitelabelID uuid.UUID) (model.Setting, error)
	Update(ctx context.Context, setting model.Setting) error
	List(ctx context.Context, pagination Pagination) ([]model.Setting, error)
}
