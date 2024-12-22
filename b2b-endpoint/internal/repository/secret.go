package repository

import (
	"context"

	"code.emcdtech.com/b2b/endpoint/internal/model"
)

type Secret interface {
	Add(ctx context.Context, secret *model.Secret) error
	Find(ctx context.Context, filter *model.SecretFilter) ([]*model.Secret, error)
	FindOne(ctx context.Context, filter *model.SecretFilter) (*model.Secret, error)
	Update(ctx context.Context, secret *model.Secret, filter *model.SecretFilter, partial *model.SecretPartial) error
}
