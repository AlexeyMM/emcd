package repository

import (
	"context"

	"github.com/google/uuid"

	"code.emcdtech.com/emcd/service/email/internal/model"
)

//go:generate mockgen -source=$GOFILE -destination=./mocks/templates.go -package=mockstore

type Template interface {
	Create(ctx context.Context, template model.Template) error
	Get(ctx context.Context, whitelabelID uuid.UUID, language string, _type model.CodeTemplate) (model.Template, error)
	Update(ctx context.Context, template model.Template) error
	List(ctx context.Context, pagination Pagination) ([]model.Template, error)
}
