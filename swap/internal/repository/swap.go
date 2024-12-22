package repository

import (
	"context"

	transactor "code.emcdtech.com/emcd/sdk/pg"

	"code.emcdtech.com/b2b/swap/model"
)

type Swap interface {
	transactor.PgxTransactor

	Add(ctx context.Context, swap *model.Swap) error
	Find(ctx context.Context, filter *model.SwapFilter) (model.Swaps, error)
	FindOne(ctx context.Context, filter *model.SwapFilter) (*model.Swap, error)
	Update(ctx context.Context, swap *model.Swap, filter *model.SwapFilter, partial *model.SwapPartial) error
	CountSwapsByStatus(ctx context.Context) (map[model.Status]int, error)
	CountTotalWithFilter(ctx context.Context, filter *model.SwapFilter) (int, error)
}
