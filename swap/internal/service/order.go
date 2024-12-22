package service

import (
	"context"
	"fmt"

	"code.emcdtech.com/b2b/swap/internal/repository"
	"code.emcdtech.com/b2b/swap/model"
)

type Order interface {
	Find(ctx context.Context, filter *model.OrderFilter) (model.Orders, error)
	FindOne(ctx context.Context, filter *model.OrderFilter) (*model.Order, error)
	Update(ctx context.Context, order *model.Order, filter *model.OrderFilter, partial *model.OrderPartial) error
}

type order struct {
	repo repository.Order
}

func NewOrder(repo repository.Order) *order {
	return &order{
		repo: repo,
	}
}

func (o *order) Find(ctx context.Context, filter *model.OrderFilter) (model.Orders, error) {
	orders, err := o.repo.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("find: %w", err)
	}
	return orders, nil
}

func (o *order) FindOne(ctx context.Context, filter *model.OrderFilter) (*model.Order, error) {
	ord, err := o.repo.FindOne(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("findOne: %w", err)
	}
	return ord, nil
}

func (o *order) Update(ctx context.Context, order *model.Order, filter *model.OrderFilter, partial *model.OrderPartial) error {
	err := o.repo.Update(ctx, order, filter, partial)
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}
	return nil
}
