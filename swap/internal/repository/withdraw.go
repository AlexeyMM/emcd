package repository

import (
	"context"

	"code.emcdtech.com/b2b/swap/model"
	transactor "code.emcdtech.com/emcd/sdk/pg"
)

type Withdraw interface {
	Add(ctx context.Context, withdraw *model.Withdraw) error
	Find(ctx context.Context, filter *model.WithdrawFilter) (model.Withdraws, error)
	FindOne(ctx context.Context, filter *model.WithdrawFilter) (*model.Withdraw, error)
	Update(ctx context.Context, withdraw *model.Withdraw, filter *model.WithdrawFilter, partial *model.WithdrawPartial) error
	transactor.PgxTransactor
}
