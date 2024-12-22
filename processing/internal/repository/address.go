package repository

import (
	"context"

	"code.emcdtech.com/b2b/processing/model"
	transactor "code.emcdtech.com/emcd/sdk/pg"
	"github.com/google/uuid"
)

type DepositAddressPool interface {
	transactor.PgxTransactor
	// OccupyAddress finds available address and marks it not available
	OccupyAddress(ctx context.Context, merchantID uuid.UUID, networkID string) (*model.Address, error)
	Save(ctx context.Context, address *model.Address) error
}
