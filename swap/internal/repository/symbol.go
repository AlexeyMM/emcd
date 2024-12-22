package repository

import (
	"context"

	"code.emcdtech.com/b2b/swap/model"
)

type Symbol interface {
	UpdateAll(ctx context.Context, symbols map[string]*model.Symbol) error
	Get(ctx context.Context, title string) (*model.Symbol, error)
	GetAll(ctx context.Context) ([]*model.Symbol, error)
	GetAccuracy(ctx context.Context, symbol string) (*model.Accuracy, error)
}
