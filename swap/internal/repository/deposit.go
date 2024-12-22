package repository

import (
	"context"

	"code.emcdtech.com/b2b/swap/model"
	transactor "code.emcdtech.com/emcd/sdk/pg"
)

type Deposit interface {
	Add(ctx context.Context, deposit *model.Deposit) error
	Find(ctx context.Context, filter *model.DepositFilter) (model.Deposits, error)
	FindOne(ctx context.Context, filter *model.DepositFilter) (*model.Deposit, error)
	transactor.PgxTransactor
}
