package pg

import (
	"context"
	"errors"
	"fmt"

	"code.emcdtech.com/b2b/processing/internal/repository/pg/sqlc"
	"code.emcdtech.com/b2b/processing/model"
	transactor "code.emcdtech.com/emcd/sdk/pg"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DepositAddressPool struct {
	transactor.PgxTransactor
}

func NewDepositAddressPool(pool *pgxpool.Pool) *DepositAddressPool {
	return &DepositAddressPool{PgxTransactor: transactor.NewPgxTransactor(pool)}
}

func (p *DepositAddressPool) OccupyAddress(
	ctx context.Context,
	merchantID uuid.UUID,
	networkID string,
) (*model.Address, error) {
	q := sqlc.New(p.Runner(ctx))

	address, err := q.OccupyDepositAddress(ctx, &sqlc.OccupyDepositAddressParams{
		NetworkID:  networkID,
		MerchantID: merchantID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &model.Error{Code: model.ErrorCodeNoAvailableAddress}
		}

		return nil, fmt.Errorf("occupyDepositAddress: %w", err)
	}

	return &model.Address{
		Address:    address,
		NetworkID:  networkID,
		MerchantID: merchantID,
		Available:  true, // lets agree to call it available until transaction is committed
	}, nil
}

func (p *DepositAddressPool) Save(ctx context.Context, address *model.Address) error {
	q := sqlc.New(p.Runner(ctx))

	err := q.SaveDepositAddress(ctx, &sqlc.SaveDepositAddressParams{
		Address:    address.Address,
		NetworkID:  address.NetworkID,
		MerchantID: address.MerchantID,
		Available:  address.Available,
	})
	if err != nil {
		return fmt.Errorf("saveDepositAddress: %w", err)
	}

	return nil
}
