package pg

import (
	"context"
	"fmt"

	"code.emcdtech.com/b2b/processing/internal/repository/pg/sqlc"
	"code.emcdtech.com/b2b/processing/model"
	transactor "code.emcdtech.com/emcd/sdk/pg"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MerchantAdmin struct {
	transactor.PgxTransactor
}

func NewMerchantAdmin(pool *pgxpool.Pool) *MerchantAdmin {
	return &MerchantAdmin{PgxTransactor: transactor.NewPgxTransactor(pool)}
}

func (r *MerchantAdmin) SaveMerchant(ctx context.Context, merchant *model.Merchant) error {
	return r.WithinTransaction(ctx, func(ctx context.Context) error {
		q := sqlc.New(r.Runner(ctx))

		if err := q.SaveMerchantID(ctx, merchant.ID); err != nil {
			return fmt.Errorf("saveMerchantID: %w", err)
		}

		err := q.SaveMerchantTariff(ctx, &sqlc.SaveMerchantTariffParams{
			MerchantID: merchant.ID,
			UpperFee:   merchant.Tariff.UpperFee,
			LowerFee:   merchant.Tariff.LowerFee,
			MinPay:     merchant.Tariff.MinPay,
			MaxPay:     merchant.Tariff.MaxPay,
		})
		if err != nil {
			return fmt.Errorf("saveMerchantTariff: %w", err)
		}

		return nil
	})
}
