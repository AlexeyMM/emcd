package repository

import (
	"context"

	"code.emcdtech.com/b2b/swap/model"
)

type User interface {
	Add(ctx context.Context, user *model.User) error
	FindOne(ctx context.Context, filter *model.UserFilter) (*model.User, error)
}
