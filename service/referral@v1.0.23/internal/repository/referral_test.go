package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"

	pgTx "code.emcdtech.com/emcd/sdk/pg"

	"code.emcdtech.com/emcd/service/referral/internal/model"
	"code.emcdtech.com/emcd/service/referral/internal/repository"
)

func TestReferral_Create(t *testing.T) {
	ctx := context.Background()

	repo := repository.NewReferral(pgTx.NewPgxTransactor(db), db)

	t.Run("Empty_WhitelabelID_AND_ReferralID", func(t *testing.T) {
		in := &model.Referral{
			UserID:        uuid.New(),
			Product:       "p2p",
			Coin:          "btc",
			Fee:           decimal.NewFromFloat(15),
			WhitelabelFee: decimal.NewFromFloat(17),
			ReferralFee:   decimal.NewFromFloat(12),
		}

		err := repo.Create(ctx, in)
		require.NoError(t, err)

		defer func(userID uuid.UUID, product, coin string, wlID uuid.UUID) {
			err = repo.Delete(ctx, in.UserID, in.Product, in.Coin)
			require.NoError(t, err)
		}(in.UserID, in.Product, in.Coin, in.WhitelabelID)

		out, err := repo.Get(ctx, in.UserID, in.Product, in.Coin)
		require.NoError(t, err)
		require.Equal(t, in.UserID, out.UserID)
		require.Equal(t, in.Product, out.Product)
		require.Equal(t, in.Coin, out.Coin)
		require.Equal(t, in.WhitelabelID, out.WhitelabelID)
		require.Equal(t, in.Fee, out.Fee)
		require.Equal(t, in.WhitelabelFee, out.WhitelabelFee)
		require.Equal(t, in.ReferralFee, out.ReferralFee)
		require.Equal(t, in.ReferralID, out.ReferralID)
	})

	t.Run("Empty_ReferralID", func(t *testing.T) {
		in := &model.Referral{
			UserID:        uuid.New(),
			Product:       "p2p",
			Coin:          "btc",
			WhitelabelID:  uuid.New(),
			Fee:           decimal.NewFromFloat(15),
			WhitelabelFee: decimal.NewFromFloat(17),
			ReferralFee:   decimal.NewFromFloat(12),
		}

		err := repo.Create(ctx, in)
		require.NoError(t, err)

		defer func(userID uuid.UUID, product, coin string, wlID uuid.UUID) {
			err = repo.Delete(ctx, in.UserID, in.Product, in.Coin)
			require.NoError(t, err)
		}(in.UserID, in.Product, in.Coin, in.WhitelabelID)

		out, err := repo.Get(ctx, in.UserID, in.Product, in.Coin)
		require.NoError(t, err)
		require.Equal(t, in.UserID, out.UserID)
		require.Equal(t, in.Product, out.Product)
		require.Equal(t, in.Coin, out.Coin)
		require.Equal(t, in.WhitelabelID, out.WhitelabelID)
		require.Equal(t, in.Fee, out.Fee)
		require.Equal(t, in.WhitelabelFee, out.WhitelabelFee)
		require.Equal(t, in.ReferralFee, out.ReferralFee)
		require.Equal(t, in.ReferralID, out.ReferralID)
	})

	t.Run("Full", func(t *testing.T) {
		in := &model.Referral{
			UserID:        uuid.New(),
			Product:       "p2p",
			Coin:          "btc",
			WhitelabelID:  uuid.New(),
			Fee:           decimal.NewFromFloat(15),
			WhitelabelFee: decimal.NewFromFloat(17),
			ReferralFee:   decimal.NewFromFloat(12),
			ReferralID:    uuid.New(),
		}

		err := repo.Create(ctx, in)
		require.NoError(t, err)

		defer func(userID uuid.UUID, product, coin string, wlID uuid.UUID) {
			err = repo.Delete(ctx, in.UserID, in.Product, in.Coin)
			require.NoError(t, err)
		}(in.UserID, in.Product, in.Coin, in.WhitelabelID)

		out, err := repo.Get(ctx, in.UserID, in.Product, in.Coin)
		require.NoError(t, err)
		require.Equal(t, in.UserID, out.UserID)
		require.Equal(t, in.Product, out.Product)
		require.Equal(t, in.Coin, out.Coin)
		require.Equal(t, in.WhitelabelID, out.WhitelabelID)
		require.Equal(t, in.Fee, out.Fee)
		require.Equal(t, in.WhitelabelFee, out.WhitelabelFee)
		require.Equal(t, in.ReferralFee, out.ReferralFee)
		require.Equal(t, in.ReferralID, out.ReferralID)
	})
}

func TestReferral_Update(t *testing.T) {
	ctx := context.Background()

	repo := repository.NewReferral(pgTx.NewPgxTransactor(db), db)

	t.Run("Empty_WhitelabelID_AND_ReferralID", func(t *testing.T) {
		in := &model.Referral{
			UserID:        uuid.New(),
			Product:       "p2p",
			Coin:          "btc",
			Fee:           decimal.NewFromFloat(15),
			WhitelabelFee: decimal.NewFromFloat(17),
			ReferralFee:   decimal.NewFromFloat(12),
		}

		err := repo.Create(ctx, in)
		require.NoError(t, err)

		defer func(userID uuid.UUID, product, coin string, wlID uuid.UUID) {
			err = repo.Delete(ctx, in.UserID, in.Product, in.Coin)
			require.NoError(t, err)
		}(in.UserID, in.Product, in.Coin, in.WhitelabelID)

		out, err := repo.Get(ctx, in.UserID, in.Product, in.Coin)
		require.NoError(t, err)
		require.Equal(t, in.UserID, out.UserID)
		require.Equal(t, in.Product, out.Product)
		require.Equal(t, in.Coin, out.Coin)
		require.Equal(t, in.WhitelabelID, out.WhitelabelID)
		require.Equal(t, in.Fee, out.Fee)
		require.Equal(t, in.WhitelabelFee, out.WhitelabelFee)
		require.Equal(t, in.ReferralFee, out.ReferralFee)
		require.Equal(t, in.ReferralID, out.ReferralID)

		in.Fee = decimal.NewFromFloat(34)
		in.ReferralFee = decimal.NewFromFloat(36)

		err = repo.Update(ctx, in)
		require.NoError(t, err)

		out, err = repo.Get(ctx, in.UserID, in.Product, in.Coin)
		require.NoError(t, err)
		require.Equal(t, in.UserID, out.UserID)
		require.Equal(t, in.Product, out.Product)
		require.Equal(t, in.Coin, out.Coin)
		require.Equal(t, in.WhitelabelID, out.WhitelabelID)
		require.Equal(t, in.Fee, out.Fee)
		require.Equal(t, in.ReferralFee, out.ReferralFee)
		require.Equal(t, in.ReferralID, out.ReferralID)
	})

	t.Run("Empty_ReferralID", func(t *testing.T) {
		in := &model.Referral{
			UserID:        uuid.New(),
			Product:       "p2p",
			Coin:          "btc",
			WhitelabelID:  uuid.New(),
			Fee:           decimal.NewFromFloat(15),
			WhitelabelFee: decimal.NewFromFloat(17),
			ReferralFee:   decimal.NewFromFloat(12),
		}

		err := repo.Create(ctx, in)
		require.NoError(t, err)

		defer func(userID uuid.UUID, product, coin string, wlID uuid.UUID) {
			err = repo.Delete(ctx, in.UserID, in.Product, in.Coin)
			require.NoError(t, err)
		}(in.UserID, in.Product, in.Coin, in.WhitelabelID)

		out, err := repo.Get(ctx, in.UserID, in.Product, in.Coin)
		require.NoError(t, err)
		require.Equal(t, in.UserID, out.UserID)
		require.Equal(t, in.Product, out.Product)
		require.Equal(t, in.Coin, out.Coin)
		require.Equal(t, in.WhitelabelID, out.WhitelabelID)
		require.Equal(t, in.Fee, out.Fee)
		require.Equal(t, in.WhitelabelFee, out.WhitelabelFee)
		require.Equal(t, in.ReferralFee, out.ReferralFee)
		require.Equal(t, in.ReferralID, out.ReferralID)

		in.Fee = decimal.NewFromFloat(34)
		in.ReferralFee = decimal.NewFromFloat(36)
		in.ReferralID = uuid.New()

		err = repo.Update(ctx, in)
		require.NoError(t, err)

		out, err = repo.Get(ctx, in.UserID, in.Product, in.Coin)
		require.NoError(t, err)
		require.Equal(t, in.UserID, out.UserID)
		require.Equal(t, in.Product, out.Product)
		require.Equal(t, in.Coin, out.Coin)
		require.Equal(t, in.WhitelabelID, out.WhitelabelID)
		require.Equal(t, in.Fee, out.Fee)
		require.Equal(t, in.ReferralFee, out.ReferralFee)
		require.Equal(t, in.ReferralID, out.ReferralID)
	})

	t.Run("Full", func(t *testing.T) {
		in := &model.Referral{
			UserID:        uuid.New(),
			Product:       "p2p",
			Coin:          "btc",
			WhitelabelID:  uuid.New(),
			Fee:           decimal.NewFromFloat(15),
			WhitelabelFee: decimal.NewFromFloat(17),
			ReferralFee:   decimal.NewFromFloat(12),
			ReferralID:    uuid.New(),
		}

		err := repo.Create(ctx, in)
		require.NoError(t, err)

		defer func(userID uuid.UUID, product, coin string, wlID uuid.UUID) {
			err = repo.Delete(ctx, in.UserID, in.Product, in.Coin)
			require.NoError(t, err)
		}(in.UserID, in.Product, in.Coin, in.WhitelabelID)

		out, err := repo.Get(ctx, in.UserID, in.Product, in.Coin)
		require.NoError(t, err)
		require.Equal(t, in.UserID, out.UserID)
		require.Equal(t, in.Product, out.Product)
		require.Equal(t, in.Coin, out.Coin)
		require.Equal(t, in.WhitelabelID, out.WhitelabelID)
		require.Equal(t, in.Fee, out.Fee)
		require.Equal(t, in.WhitelabelFee, out.WhitelabelFee)
		require.Equal(t, in.ReferralFee, out.ReferralFee)
		require.Equal(t, in.ReferralID, out.ReferralID)

		in.Fee = decimal.NewFromFloat(34)
		in.ReferralFee = decimal.NewFromFloat(36)
		in.ReferralID = uuid.New()

		err = repo.Update(ctx, in)
		require.NoError(t, err)

		out, err = repo.Get(ctx, in.UserID, in.Product, in.Coin)
		require.NoError(t, err)
		require.Equal(t, in.UserID, out.UserID)
		require.Equal(t, in.Product, out.Product)
		require.Equal(t, in.Coin, out.Coin)
		require.Equal(t, in.WhitelabelID, out.WhitelabelID)
		require.Equal(t, in.Fee, out.Fee)
		require.Equal(t, in.ReferralFee, out.ReferralFee)
		require.Equal(t, in.ReferralID, out.ReferralID)
	})
}

func TestReferral_Delete(t *testing.T) {
	ctx := context.Background()

	repo := repository.NewReferral(pgTx.NewPgxTransactor(db), db)

	t.Run("Empty_WhitelabelID_AND_ReferralID", func(t *testing.T) {
		in := &model.Referral{
			UserID:        uuid.New(),
			Product:       "p2p",
			Coin:          "btc",
			Fee:           decimal.NewFromFloat(15),
			WhitelabelFee: decimal.NewFromFloat(17),
			ReferralFee:   decimal.NewFromFloat(12),
		}

		err := repo.Create(ctx, in)
		require.NoError(t, err)

		defer func(userID uuid.UUID, product, coin string, wlID uuid.UUID) {
			err = repo.Delete(ctx, in.UserID, in.Product, in.Coin)
			require.NoError(t, err)
		}(in.UserID, in.Product, in.Coin, in.WhitelabelID)

		out, err := repo.Get(ctx, in.UserID, in.Product, in.Coin)
		require.NoError(t, err)
		require.Equal(t, in.UserID, out.UserID)
		require.Equal(t, in.Product, out.Product)
		require.Equal(t, in.Coin, out.Coin)
		require.Equal(t, in.WhitelabelID, out.WhitelabelID)
		require.Equal(t, in.Fee, out.Fee)
		require.Equal(t, in.WhitelabelFee, out.WhitelabelFee)
		require.Equal(t, in.ReferralFee, out.ReferralFee)
		require.Equal(t, in.ReferralID, out.ReferralID)

		err = repo.Delete(ctx, in.UserID, in.Product, in.Coin)
		require.NoError(t, err)

		out, err = repo.Get(ctx, in.UserID, in.Product, in.Coin)
		require.Error(t, err)
	})

	t.Run("Empty_ReferralID", func(t *testing.T) {
		in := &model.Referral{
			UserID:        uuid.New(),
			Product:       "p2p",
			Coin:          "btc",
			WhitelabelID:  uuid.New(),
			Fee:           decimal.NewFromFloat(15),
			WhitelabelFee: decimal.NewFromFloat(17),
			ReferralFee:   decimal.NewFromFloat(12),
		}

		err := repo.Create(ctx, in)
		require.NoError(t, err)

		defer func(userID uuid.UUID, product, coin string, wlID uuid.UUID) {
			err = repo.Delete(ctx, in.UserID, in.Product, in.Coin)
			require.NoError(t, err)
		}(in.UserID, in.Product, in.Coin, in.WhitelabelID)

		out, err := repo.Get(ctx, in.UserID, in.Product, in.Coin)
		require.NoError(t, err)
		require.Equal(t, in.UserID, out.UserID)
		require.Equal(t, in.Product, out.Product)
		require.Equal(t, in.Coin, out.Coin)
		require.Equal(t, in.WhitelabelID, out.WhitelabelID)
		require.Equal(t, in.Fee, out.Fee)
		require.Equal(t, in.WhitelabelFee, out.WhitelabelFee)
		require.Equal(t, in.ReferralFee, out.ReferralFee)
		require.Equal(t, in.ReferralID, out.ReferralID)

		err = repo.Delete(ctx, in.UserID, in.Product, in.Coin)
		require.NoError(t, err)

		out, err = repo.Get(ctx, in.UserID, in.Product, in.Coin)
		require.Error(t, err)
	})

	t.Run("Full", func(t *testing.T) {
		in := &model.Referral{
			UserID:        uuid.New(),
			Product:       "p2p",
			Coin:          "btc",
			WhitelabelID:  uuid.New(),
			Fee:           decimal.NewFromFloat(15),
			WhitelabelFee: decimal.NewFromFloat(17),
			ReferralFee:   decimal.NewFromFloat(12),
			ReferralID:    uuid.New(),
		}

		err := repo.Create(ctx, in)
		require.NoError(t, err)

		defer func(userID uuid.UUID, product, coin string, wlID uuid.UUID) {
			err = repo.Delete(ctx, in.UserID, in.Product, in.Coin)
			require.NoError(t, err)
		}(in.UserID, in.Product, in.Coin, in.WhitelabelID)

		out, err := repo.Get(ctx, in.UserID, in.Product, in.Coin)
		require.NoError(t, err)
		require.Equal(t, in.UserID, out.UserID)
		require.Equal(t, in.Product, out.Product)
		require.Equal(t, in.Coin, out.Coin)
		require.Equal(t, in.WhitelabelID, out.WhitelabelID)
		require.Equal(t, in.Fee, out.Fee)
		require.Equal(t, in.WhitelabelFee, out.WhitelabelFee)
		require.Equal(t, in.ReferralFee, out.ReferralFee)
		require.Equal(t, in.ReferralID, out.ReferralID)

		err = repo.Delete(ctx, in.UserID, in.Product, in.Coin)
		require.NoError(t, err)

		out, err = repo.Get(ctx, in.UserID, in.Product, in.Coin)
		require.Error(t, err)
	})
}

func TestReferral_Get(t *testing.T) {
	ctx := context.Background()

	repo := repository.NewReferral(pgTx.NewPgxTransactor(db), db)

	t.Run("Empty_WhitelabelID_AND_ReferralID", func(t *testing.T) {
		in := &model.Referral{
			UserID:        uuid.New(),
			Product:       "p2p",
			Coin:          "btc",
			Fee:           decimal.NewFromFloat(15),
			WhitelabelFee: decimal.NewFromFloat(17),
			ReferralFee:   decimal.NewFromFloat(12),
		}

		err := repo.Create(ctx, in)
		require.NoError(t, err)

		defer func(userID uuid.UUID, product, coin string, wlID uuid.UUID) {
			err = repo.Delete(ctx, in.UserID, in.Product, in.Coin)
			require.NoError(t, err)
		}(in.UserID, in.Product, in.Coin, in.WhitelabelID)

		out, err := repo.Get(ctx, in.UserID, in.Product, in.Coin)
		require.NoError(t, err)
		require.Equal(t, in.UserID, out.UserID)
		require.Equal(t, in.Product, out.Product)
		require.Equal(t, in.Coin, out.Coin)
		require.Equal(t, in.WhitelabelID, out.WhitelabelID)
		require.Equal(t, in.Fee, out.Fee)
		require.Equal(t, in.WhitelabelFee, out.WhitelabelFee)
		require.Equal(t, in.ReferralFee, out.ReferralFee)
		require.Equal(t, in.ReferralID, out.ReferralID)
	})

	t.Run("Empty_ReferralID", func(t *testing.T) {
		in := &model.Referral{
			UserID:        uuid.New(),
			Product:       "p2p",
			Coin:          "btc",
			WhitelabelID:  uuid.New(),
			Fee:           decimal.NewFromFloat(15),
			WhitelabelFee: decimal.NewFromFloat(17),
			ReferralFee:   decimal.NewFromFloat(12),
		}

		err := repo.Create(ctx, in)
		require.NoError(t, err)

		defer func(userID uuid.UUID, product, coin string, wlID uuid.UUID) {
			err = repo.Delete(ctx, in.UserID, in.Product, in.Coin)
			require.NoError(t, err)
		}(in.UserID, in.Product, in.Coin, in.WhitelabelID)

		out, err := repo.Get(ctx, in.UserID, in.Product, in.Coin)
		require.NoError(t, err)
		require.Equal(t, in.UserID, out.UserID)
		require.Equal(t, in.Product, out.Product)
		require.Equal(t, in.Coin, out.Coin)
		require.Equal(t, in.WhitelabelID, out.WhitelabelID)
		require.Equal(t, in.Fee, out.Fee)
		require.Equal(t, in.WhitelabelFee, out.WhitelabelFee)
		require.Equal(t, in.ReferralFee, out.ReferralFee)
		require.Equal(t, in.ReferralID, out.ReferralID)
	})

	t.Run("Full", func(t *testing.T) {
		in := &model.Referral{
			UserID:        uuid.New(),
			Product:       "p2p",
			Coin:          "btc",
			WhitelabelID:  uuid.New(),
			Fee:           decimal.NewFromFloat(15),
			WhitelabelFee: decimal.NewFromFloat(17),
			ReferralFee:   decimal.NewFromFloat(12),
			ReferralID:    uuid.New(),
		}

		err := repo.Create(ctx, in)
		require.NoError(t, err)

		defer func(userID uuid.UUID, product, coin string, wlID uuid.UUID) {
			err = repo.Delete(ctx, in.UserID, in.Product, in.Coin)
			require.NoError(t, err)
		}(in.UserID, in.Product, in.Coin, in.WhitelabelID)

		out, err := repo.Get(ctx, in.UserID, in.Product, in.Coin)
		require.NoError(t, err)
		require.Equal(t, in.UserID, out.UserID)
		require.Equal(t, in.Product, out.Product)
		require.Equal(t, in.Coin, out.Coin)
		require.Equal(t, in.WhitelabelID, out.WhitelabelID)
		require.Equal(t, in.Fee, out.Fee)
		require.Equal(t, in.WhitelabelFee, out.WhitelabelFee)
		require.Equal(t, in.ReferralFee, out.ReferralFee)
		require.Equal(t, in.ReferralID, out.ReferralID)
	})
}

func TestReferral_List(t *testing.T) {
	ctx := context.Background()

	repo := repository.NewReferral(pgTx.NewPgxTransactor(db), db)

	t.Run("Empty_WhitelabelID_AND_ReferralID", func(t *testing.T) {
		in := &model.Referral{
			UserID:        uuid.New(),
			Product:       "p2p",
			Coin:          "btc",
			Fee:           decimal.NewFromFloat(15),
			WhitelabelFee: decimal.NewFromFloat(17),
			ReferralFee:   decimal.NewFromFloat(12),
		}

		err := repo.Create(ctx, in)
		require.NoError(t, err)

		defer func(userID uuid.UUID, product, coin string, wlID uuid.UUID) {
			err = repo.Delete(ctx, in.UserID, in.Product, in.Coin)
			require.NoError(t, err)
		}(in.UserID, in.Product, in.Coin, in.WhitelabelID)

		in.Product = "mining"
		in.Coin = "usdt"

		err = repo.Create(ctx, in)
		require.NoError(t, err)

		defer func(userID uuid.UUID, product, coin string, wlID uuid.UUID) {
			err = repo.Delete(ctx, in.UserID, in.Product, in.Coin)
			require.NoError(t, err)
		}(in.UserID, in.Product, in.Coin, in.WhitelabelID)

		list, _, err := repo.List(ctx, in.UserID, 0, 2)
		require.NoError(t, err)
		require.Equal(t, 2, len(list))
	})

	t.Run("Empty_ReferralID", func(t *testing.T) {
		in := &model.Referral{
			UserID:        uuid.New(),
			Product:       "p2p",
			Coin:          "btc",
			WhitelabelID:  uuid.New(),
			Fee:           decimal.NewFromFloat(15),
			WhitelabelFee: decimal.NewFromFloat(17),
			ReferralFee:   decimal.NewFromFloat(12),
		}

		err := repo.Create(ctx, in)
		require.NoError(t, err)

		defer func(userID uuid.UUID, product, coin string, wlID uuid.UUID) {
			err = repo.Delete(ctx, in.UserID, in.Product, in.Coin)
			require.NoError(t, err)
		}(in.UserID, in.Product, in.Coin, in.WhitelabelID)

		in.Product = "mining"
		in.Coin = "usdt"

		err = repo.Create(ctx, in)
		require.NoError(t, err)

		defer func(userID uuid.UUID, product, coin string, wlID uuid.UUID) {
			err = repo.Delete(ctx, in.UserID, in.Product, in.Coin)
			require.NoError(t, err)
		}(in.UserID, in.Product, in.Coin, in.WhitelabelID)

		list, _, err := repo.List(ctx, in.UserID, 0, 2)
		require.NoError(t, err)
		require.Equal(t, 2, len(list))
	})

	t.Run("Full", func(t *testing.T) {
		in := &model.Referral{
			UserID:        uuid.New(),
			Product:       "p2p",
			Coin:          "btc",
			WhitelabelID:  uuid.New(),
			Fee:           decimal.NewFromFloat(15),
			WhitelabelFee: decimal.NewFromFloat(17),
			ReferralFee:   decimal.NewFromFloat(12),
			ReferralID:    uuid.New(),
		}

		err := repo.Create(ctx, in)
		require.NoError(t, err)

		defer func(userID uuid.UUID, product, coin string, wlID uuid.UUID) {
			err = repo.Delete(ctx, in.UserID, in.Product, in.Coin)
			require.NoError(t, err)
		}(in.UserID, in.Product, in.Coin, in.WhitelabelID)

		in.Product = "mining"
		in.Coin = "usdt"

		err = repo.Create(ctx, in)
		require.NoError(t, err)

		defer func(userID uuid.UUID, product, coin string, wlID uuid.UUID) {
			err = repo.Delete(ctx, in.UserID, in.Product, in.Coin)
			require.NoError(t, err)
		}(in.UserID, in.Product, in.Coin, in.WhitelabelID)

		list, _, err := repo.List(ctx, in.UserID, 0, 2)
		require.NoError(t, err)
		require.Equal(t, 2, len(list))
	})
}

func TestReferral_History(t *testing.T) {
	ctx := context.Background()

	repo := repository.NewReferral(pgTx.NewPgxTransactor(db), db)

	t.Run("Empty_WhitelabelID_AND_ReferralID", func(t *testing.T) {
		in := &model.Referral{
			UserID:        uuid.New(),
			Product:       "p2p",
			Coin:          "btc",
			Fee:           decimal.NewFromFloat(15),
			WhitelabelFee: decimal.NewFromFloat(17),
			ReferralFee:   decimal.NewFromFloat(12),
		}

		err := repo.Create(ctx, in)
		require.NoError(t, err)

		defer func(userID uuid.UUID, product, coin string, wlID uuid.UUID) {
			err = repo.Delete(ctx, in.UserID, in.Product, in.Coin)
			require.NoError(t, err)
		}(in.UserID, in.Product, in.Coin, in.WhitelabelID)

		in.Fee = decimal.NewFromFloat(34)
		in.ReferralFee = decimal.NewFromFloat(36)
		in.WhitelabelFee = decimal.NewFromFloat(45)

		err = repo.Update(ctx, in)
		require.NoError(t, err)

		in.Fee = decimal.NewFromFloat(23)
		in.ReferralFee = decimal.NewFromFloat(16)
		in.WhitelabelFee = decimal.NewFromFloat(34)

		err = repo.Update(ctx, in)
		require.NoError(t, err)

		out, err := repo.History(ctx, in.UserID, in.Product, in.Coin)
		require.NoError(t, err)
		require.NotEqual(t, 0, len(out))
	})

	t.Run("Empty_ReferralID", func(t *testing.T) {
		in := &model.Referral{
			UserID:        uuid.New(),
			Product:       "p2p",
			Coin:          "btc",
			WhitelabelID:  uuid.New(),
			Fee:           decimal.NewFromFloat(15),
			WhitelabelFee: decimal.NewFromFloat(34),
			ReferralFee:   decimal.NewFromFloat(12),
		}

		err := repo.Create(ctx, in)
		require.NoError(t, err)

		defer func(userID uuid.UUID, product, coin string, wlID uuid.UUID) {
			err = repo.Delete(ctx, in.UserID, in.Product, in.Coin)
			require.NoError(t, err)
		}(in.UserID, in.Product, in.Coin, in.WhitelabelID)

		in.Fee = decimal.NewFromFloat(34)
		in.ReferralFee = decimal.NewFromFloat(36)
		in.WhitelabelFee = decimal.NewFromFloat(45)

		err = repo.Update(ctx, in)
		require.NoError(t, err)

		in.Fee = decimal.NewFromFloat(23)
		in.ReferralFee = decimal.NewFromFloat(16)
		in.WhitelabelFee = decimal.NewFromFloat(48)

		err = repo.Update(ctx, in)
		require.NoError(t, err)

		out, err := repo.History(ctx, in.UserID, in.Product, in.Coin)
		require.NoError(t, err)
		require.NotEqual(t, 0, len(out))
	})

	t.Run("Full", func(t *testing.T) {
		in := &model.Referral{
			UserID:        uuid.New(),
			Product:       "p2p",
			Coin:          "btc",
			WhitelabelID:  uuid.New(),
			Fee:           decimal.NewFromFloat(15),
			WhitelabelFee: decimal.NewFromFloat(23),
			ReferralFee:   decimal.NewFromFloat(12),
			ReferralID:    uuid.New(),
		}

		err := repo.Create(ctx, in)
		require.NoError(t, err)

		defer func(userID uuid.UUID, product, coin string, wlID uuid.UUID) {
			err = repo.Delete(ctx, in.UserID, in.Product, in.Coin)
			require.NoError(t, err)
		}(in.UserID, in.Product, in.Coin, in.WhitelabelID)

		in.Fee = decimal.NewFromFloat(34)
		in.ReferralFee = decimal.NewFromFloat(36)
		in.WhitelabelFee = decimal.NewFromFloat(45)

		err = repo.Update(ctx, in)
		require.NoError(t, err)

		in.Fee = decimal.NewFromFloat(23)
		in.ReferralFee = decimal.NewFromFloat(16)
		in.WhitelabelFee = decimal.NewFromFloat(62)

		err = repo.Update(ctx, in)
		require.NoError(t, err)

		out, err := repo.History(ctx, in.UserID, in.Product, in.Coin)
		require.NoError(t, err)
		require.NotEqual(t, 0, len(out))
	})
}

func TestReferral_CreateMultiple(t *testing.T) {
	ctx := context.Background()
	defer truncateReferrals(ctx)

	exp1 := model.Referral{
		UserID:        uuid.New(),
		Product:       "p2p",
		Coin:          "btc",
		WhitelabelID:  uuid.New(),
		Fee:           decimal.NewFromFloat(15),
		WhitelabelFee: decimal.NewFromFloat(23),
		ReferralFee:   decimal.NewFromFloat(12),
		ReferralID:    uuid.New(),
		CreatedAt:     time.Now().UTC().Truncate(time.Millisecond),
	}
	exp2 := exp1
	exp2.Coin = "ltc"

	exp3 := exp1
	exp3.Product = "mining"
	exp3.Coin = "kas"

	repo := repository.NewReferral(pgTx.NewPgxTransactor(db), db)
	err := repo.CreateMultiple(ctx, []*model.Referral{&exp1, &exp2, &exp3})
	require.NoError(t, err)

	for _, exp := range []*model.Referral{&exp1, &exp2, &exp3} {
		actual, err := repo.Get(ctx, exp.UserID, exp.Product, exp.Coin)
		require.NoError(t, err)
		require.Equal(t, exp, actual)
	}
}

func TestReferral_UpdateWithMultiplierOne(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewReferral(pgTx.NewPgxTransactor(db), db)
	defer truncateReferrals(ctx)

	userID := uuid.New()
	fee := decimal.NewFromInt(15)
	err := repo.Create(ctx, &model.Referral{
		Fee:       fee,
		UserID:    userID,
		Product:   "mining",
		Coin:      "BTC",
		CreatedAt: time.Now().UTC(),
	})
	require.NoError(t, err)

	multiplier := decimal.NewFromFloat32(0.15)
	err = repo.UpdateWithMultiplier(ctx, userID, "mining", []string{"BTC"}, multiplier)
	require.NoError(t, err)

	actual, err := repo.Get(ctx, userID, "mining", "BTC")
	require.NoError(t, err)
	require.Equal(t, fee.Mul(multiplier).String(), actual.Fee.String())
}

func TestReferral_UpdateWithMultiplierMany(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewReferral(pgTx.NewPgxTransactor(db), db)
	defer truncateReferrals(ctx)

	fees := map[string]decimal.Decimal{"BTC": decimal.NewFromFloat32(15), "LTC": decimal.NewFromFloat32(10), "DOGE": decimal.NewFromFloat32(2)}
	coins := []string{"BTC", "LTC", "DOGE"}
	expected := &model.Referral{
		UserID:        uuid.New(),
		Product:       "mining",
		Coin:          "BTC",
		CreatedAt:     time.Now().UTC(),
		WhitelabelFee: decimal.NewFromFloat(44),
		ReferralID:    uuid.New(),
		ReferralFee:   decimal.NewFromFloat(11),
	}
	for i := 0; i < len(coins); i++ {
		expected.Fee = fees[coins[i]]
		expected.Coin = coins[i]
		err := repo.Create(ctx, expected)
		require.NoError(t, err)
	}
	multiplier := decimal.NewFromFloat32(0.6)

	multMap := map[string]decimal.Decimal{"BTC": multiplier, "LTC": multiplier, "DOGE": decimal.NewFromInt(1)}
	err := repo.UpdateWithMultiplier(ctx, expected.UserID, "mining", []string{"BTC", "LTC"}, multiplier)
	require.NoError(t, err)

	for i := 0; i < len(coins); i++ {
		actual, err := repo.Get(ctx, expected.UserID, "mining", coins[i])
		require.NoError(t, err)
		require.Equal(t, fees[coins[i]].Mul(multMap[coins[i]]).String(), actual.Fee.String())
	}
}

func TestReferral_UpdateFee(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewReferral(pgTx.NewPgxTransactor(db), db)
	defer truncateReferrals(ctx)

	fees := map[string]decimal.Decimal{"BTC": decimal.NewFromFloat32(15), "LTC": decimal.NewFromFloat32(10), "DOGE": decimal.NewFromFloat32(2)}
	coins := []string{"BTC", "LTC", "DOGE"}
	expected := &model.Referral{
		UserID:        uuid.New(),
		Product:       "mining",
		Coin:          "BTC",
		CreatedAt:     time.Now().UTC(),
		WhitelabelFee: decimal.NewFromFloat(44),
		ReferralID:    uuid.New(),
		ReferralFee:   decimal.NewFromFloat(11),
	}
	for i := 0; i < len(coins); i++ {
		expected.Fee = fees[coins[i]]
		expected.Coin = coins[i]
		err := repo.Create(ctx, expected)
		require.NoError(t, err)
	}
	newFees := map[string]decimal.Decimal{"LTC": decimal.NewFromFloat(8.885), "DOGE": decimal.NewFromFloat(32.222)}

	err := repo.UpdateFee(ctx, expected.UserID, expected.Product, newFees)
	require.NoError(t, err)

	expectedFees := map[string]decimal.Decimal{"LTC": decimal.NewFromFloat(8.885), "DOGE": decimal.NewFromFloat(32.222), "BTC": decimal.NewFromFloat32(15)}
	for coin, fee := range expectedFees {
		actual, err := repo.Get(ctx, expected.UserID, "mining", coin)
		require.NoError(t, err)
		require.Equal(t, fee.String(), actual.Fee.String())
	}
}

func TestReferral_UpdateWithPromoCodeByCoin(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewReferral(pgTx.NewPgxTransactor(db), db)
	defer truncateReferrals(ctx)

	fees := map[string]decimal.Decimal{"BTC": decimal.NewFromFloat32(15), "LTC": decimal.NewFromFloat32(10), "DOGE": decimal.NewFromFloat32(2)}
	coins := []string{"BTC", "LTC", "DOGE"}
	expected := &model.Referral{
		UserID:        uuid.New(),
		Product:       "mining",
		Coin:          "BTC",
		CreatedAt:     time.Now().UTC().Truncate(time.Second),
		WhitelabelFee: decimal.NewFromFloat(44),
		ReferralID:    uuid.New(),
		ReferralFee:   decimal.NewFromFloat(11),
	}
	for i := 0; i < len(coins); i++ {
		expected.Fee = fees[coins[i]]
		expected.Coin = coins[i]
		err := repo.Create(ctx, expected)
		require.NoError(t, err)
	}
	feeMult := map[string]decimal.Decimal{"LTC": decimal.NewFromFloat(2), "DOGE": decimal.NewFromFloat(2.2)}
	refFeeMult := map[string]decimal.Decimal{"LTC": decimal.NewFromFloat(4), "DOGE": decimal.NewFromFloat(4.4)}
	now := time.Now().UTC().Truncate(time.Millisecond)
	actionID := uuid.New()
	for coin := range feeMult {
		err := repo.UpdateWithPromoCodeByCoin(ctx, &model.CoinMultiplier{
			Coin:             coin,
			UserID:           expected.UserID,
			CreatedAt:        now,
			Product:          "mining",
			ActionID:         actionID,
			FeeMultiplier:    feeMult[coin],
			RefFeeMultiplier: refFeeMult[coin],
		})
		require.NoError(t, err)
	}
	for coin := range feeMult {
		actual, err := repo.Get(ctx, expected.UserID, "mining", coin)
		require.NoError(t, err)
		require.Equal(t, fees[coin].Mul(feeMult[coin]).String(), actual.Fee.String(), coin)
		require.Equal(t, expected.ReferralFee.Mul(refFeeMult[coin]).String(), actual.ReferralFee.String())
	}
}

func TestReferral_UpdateWithPromoCodeByCoinConflict(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewReferral(pgTx.NewPgxTransactor(db), db)
	defer truncateReferrals(ctx)

	expected := &model.Referral{
		UserID:        uuid.New(),
		Product:       "mining",
		Coin:          "BTC",
		CreatedAt:     time.Now().UTC().Truncate(time.Second),
		WhitelabelFee: decimal.NewFromFloat(44),
		ReferralID:    uuid.New(),
		ReferralFee:   decimal.NewFromFloat(11),
		Fee:           decimal.NewFromFloat(1),
	}

	err := repo.Create(ctx, expected)
	require.NoError(t, err)

	actionID := uuid.New()
	err = repo.UpdateWithPromoCodeByCoin(ctx, &model.CoinMultiplier{
		UserID:   expected.UserID,
		Coin:     expected.Coin,
		Product:  "mining",
		ActionID: actionID,
	})
	require.NoError(t, err)
	err = repo.UpdateWithPromoCodeByCoin(ctx, &model.CoinMultiplier{
		UserID:   expected.UserID,
		Coin:     expected.Coin,
		Product:  "mining",
		ActionID: actionID,
	})
	require.Error(t, err)
}

func truncateReferrals(ctx context.Context) {
	_, err := db.Exec(ctx, "TRUNCATE TABLE referral.referrals")
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
}
