package repository

import (
	"context"

	"code.emcdtech.com/b2b/swap/model"
)

type OrderFee interface {
	UpdateAll(ctx context.Context, fee map[string]*model.Fee) error
	GetFee(ctx context.Context, symbol string) (*model.Fee, error)
}
