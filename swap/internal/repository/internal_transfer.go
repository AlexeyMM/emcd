package repository

import (
	"context"

	"code.emcdtech.com/b2b/swap/model"
)

type Transfer interface {
	Add(ctx context.Context, transfer *model.InternalTransfer) error
	Find(ctx context.Context, filter *model.InternalTransferFilter) (model.InternalTransfers, error)
	FindOne(ctx context.Context, filter *model.InternalTransferFilter) (*model.InternalTransfer, error)
	Update(ctx context.Context, internalTransfer *model.InternalTransfer, filter *model.InternalTransferFilter,
		partial *model.InternalTransferPartial) error
}
