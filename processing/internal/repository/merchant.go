package repository

import (
	"context"

	"code.emcdtech.com/b2b/processing/model"
	"github.com/google/uuid"
)

type MerchantAdmin interface {
	SaveMerchant(ctx context.Context, merchant *model.Merchant) error
}

type Merchant interface {
	Get(ctx context.Context, id uuid.UUID) (*model.Merchant, error)
}
