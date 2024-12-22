package repository

import (
	"context"

	"github.com/google/uuid"

	"code.emcdtech.com/b2b/swap/model"
)

type SwapStatusHistory interface {
	Add(ctx context.Context, swapID uuid.UUID, item *model.SwapStatusHistoryItem) error
	Find(ctx context.Context, filter *model.SwapStatusHistoryFilter) ([]*model.SwapStatusHistoryItem, error)
}
