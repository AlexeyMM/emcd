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

type Merchant struct {
	transactor.PgxTransactor
}

func NewMerchant(pool *pgxpool.Pool) *Merchant {
	return &Merchant{PgxTransactor: transactor.NewPgxTransactor(pool)}
}

func (m *Merchant) Get(ctx context.Context, id uuid.UUID) (*model.Merchant, error) {
	q := sqlc.New(m.Runner(ctx))

	merchantRow, err := q.GetMerchantWithTariff(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &model.Error{
				Code:    model.ErrorCodeNoSuchMerchant,
				Message: fmt.Sprintf("no such merchant: %s", id),
				Inner:   err,
			}
		}

		return nil, fmt.Errorf("getMerchantWithTariff: %w", err)
	}

	merchant := &model.Merchant{
		ID: merchantRow.ID,
		Tariff: &model.Tariff{
			UpperFee: merchantRow.UpperFee,
			LowerFee: merchantRow.LowerFee,
			MinPay:   merchantRow.MinPay,
			MaxPay:   merchantRow.MaxPay,
		},
	}

	return merchant, nil
}
