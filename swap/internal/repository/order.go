package repository

import (
	"context"

	"code.emcdtech.com/b2b/swap/model"
)

type Order interface {
	Add(ctx context.Context, order *model.Order) error
	Find(ctx context.Context, filter *model.OrderFilter) (model.Orders, error)
	FindOne(ctx context.Context, filter *model.OrderFilter) (*model.Order, error)
	Update(ctx context.Context, order *model.Order, filter *model.OrderFilter, partial *model.OrderPartial) error
}
