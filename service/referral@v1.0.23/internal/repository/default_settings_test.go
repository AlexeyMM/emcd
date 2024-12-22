package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"

	"code.emcdtech.com/emcd/service/referral/internal/model"
	"code.emcdtech.com/emcd/service/referral/internal/repository"
)

func TestDefaultSettings_Create(t *testing.T) {
	ctx := context.Background()

	repo := repository.NewDefaultSettings(db)

	in := &model.DefaultSettings{
		Product:     "p2p",
		Coin:        "btc",
		Fee:         decimal.NewFromFloat(15),
		ReferralFee: decimal.NewFromFloat(12),
	}

	err := repo.Create(ctx, in)
	require.NoError(t, err)

	defer func(product string, coin string) {
		err = repo.Delete(ctx, in.Product, in.Coin)
		require.NoError(t, err)
	}(in.Product, in.Coin)

	out, err := repo.Get(ctx, in.Product, in.Coin)
	require.NoError(t, err)
	require.Equal(t, in.Product, out.Product)
	require.Equal(t, in.Coin, out.Coin)
	require.Equal(t, in.Fee, out.Fee)
	require.Equal(t, in.ReferralFee, out.ReferralFee)
}

func TestDefaultSettings_Update(t *testing.T) {
	ctx := context.Background()

	repo := repository.NewDefaultSettings(db)

	in := &model.DefaultSettings{
		Product:     "p2p",
		Coin:        "btc",
		Fee:         decimal.NewFromFloat(15),
		ReferralFee: decimal.NewFromFloat(12),
	}

	err := repo.Create(ctx, in)
	require.NoError(t, err)

	defer func(product string, coin string) {
		err = repo.Delete(ctx, in.Product, in.Coin)
		require.NoError(t, err)
	}(in.Product, in.Coin)

	out, err := repo.Get(ctx, in.Product, in.Coin)
	require.NoError(t, err)
	require.Equal(t, in.Product, out.Product)
	require.Equal(t, in.Coin, out.Coin)
	require.Equal(t, in.Fee, out.Fee)
	require.Equal(t, in.ReferralFee, out.ReferralFee)

	in.Fee = decimal.NewFromFloat(34)
	in.ReferralFee = decimal.NewFromFloat(36)

	err = repo.Update(ctx, in)
	require.NoError(t, err)

	out, err = repo.Get(ctx, in.Product, in.Coin)
	require.NoError(t, err)
	require.Equal(t, in.Product, out.Product)
	require.Equal(t, in.Coin, out.Coin)
	require.Equal(t, in.Fee, out.Fee)
	require.Equal(t, in.ReferralFee, out.ReferralFee)
}

func TestDefaultSettings_Delete(t *testing.T) {
	ctx := context.Background()

	repo := repository.NewDefaultSettings(db)

	in := &model.DefaultSettings{
		Product:     "p2p",
		Coin:        "btc",
		Fee:         decimal.NewFromFloat(15),
		ReferralFee: decimal.NewFromFloat(12),
	}

	err := repo.Create(ctx, in)
	require.NoError(t, err)

	defer func(product string, coin string) {
		err = repo.Delete(ctx, in.Product, in.Coin)
		require.NoError(t, err)
	}(in.Product, in.Coin)

	out, err := repo.Get(ctx, in.Product, in.Coin)
	require.NoError(t, err)
	require.Equal(t, in.Product, out.Product)
	require.Equal(t, in.Coin, out.Coin)
	require.Equal(t, in.Fee, out.Fee)
	require.Equal(t, in.ReferralFee, out.ReferralFee)

	err = repo.Delete(ctx, in.Product, in.Coin)
	require.NoError(t, err)

	_, err = repo.Get(ctx, in.Product, in.Coin)
	require.Error(t, err)
}

func TestDefaultSettings_Get(t *testing.T) {
	ctx := context.Background()

	repo := repository.NewDefaultSettings(db)

	in := &model.DefaultSettings{
		Product:     "p2p",
		Coin:        "btc",
		Fee:         decimal.NewFromFloat(15),
		ReferralFee: decimal.NewFromFloat(12),
	}

	err := repo.Create(ctx, in)
	require.NoError(t, err)

	defer func(product string, coin string) {
		err = repo.Delete(ctx, in.Product, in.Coin)
		require.NoError(t, err)
	}(in.Product, in.Coin)

	out, err := repo.Get(ctx, in.Product, in.Coin)
	require.NoError(t, err)
	require.Equal(t, in.Product, out.Product)
	require.Equal(t, in.Coin, out.Coin)
	require.Equal(t, in.Fee, out.Fee)
	require.Equal(t, in.ReferralFee, out.ReferralFee)
}

func TestDefaultSettings_GetAll(t *testing.T) {
	ctx := context.Background()

	repo := repository.NewDefaultSettings(db)

	in := &model.DefaultSettings{
		Product:     "p2p",
		Coin:        "btc",
		Fee:         decimal.NewFromFloat(15),
		ReferralFee: decimal.NewFromFloat(12),
	}

	err := repo.Create(ctx, in)
	require.NoError(t, err)

	defer func(product string, coin string) {
		err = repo.Delete(ctx, in.Product, in.Coin)
		require.NoError(t, err)
		require.NoError(t, err)
	}(in.Product, in.Coin)

	in.Coin = "ltc"

	err = repo.Create(ctx, in)
	require.NoError(t, err)

	defer func(product string, coin string) {
		err = repo.Delete(ctx, in.Product, in.Coin)
		require.NoError(t, err)
	}(in.Product, in.Coin)

	in.Product = "mining"

	err = repo.Create(ctx, in)
	require.NoError(t, err)

	defer func(product string, coin string) {
		err = repo.Delete(ctx, in.Product, in.Coin)
		require.NoError(t, err)
	}(in.Product, in.Coin)

	list, _, err := repo.GetAll(ctx, 0, 3)
	require.NoError(t, err)
	require.Equal(t, 3, len(list))
}

func TestDefaultSettings_GetAllWithoutPagination(t *testing.T) {
	ctx := context.Background()

	truncateSettings(ctx)
	repo := repository.NewDefaultSettings(db)
	defer truncateSettings(ctx)
	exp1 := model.DefaultSettings{
		Product:     "p2p",
		Coin:        "btc",
		Fee:         decimal.NewFromFloat(15),
		ReferralFee: decimal.NewFromFloat(12),
	}

	err := repo.Create(ctx, &exp1)
	require.NoError(t, err)

	exp2 := exp1
	exp2.Coin = "ltc"

	err = repo.Create(ctx, &exp2)
	require.NoError(t, err)

	exp3 := exp1
	exp3.Product = "mining"

	err = repo.Create(ctx, &exp3)
	require.NoError(t, err)

	actual, err := repo.GetAllWithoutPagination(ctx)
	require.NoError(t, err)

	for i, exp := range []*model.DefaultSettings{&exp1, &exp2, &exp3} {
		verifyEqualSettings(t, exp, actual[i])
	}
}

func verifyEqualSettings(t *testing.T, in, out *model.DefaultSettings) {
	require.Equal(t, in.Product, out.Product)
	require.Equal(t, in.Coin, out.Coin)
	require.Equal(t, in.Fee, out.Fee)
	require.Equal(t, in.ReferralFee, out.ReferralFee)
	require.Equal(t, time.Now().UTC().Truncate(time.Second), out.CreatedAt.Truncate(time.Second))
}

func truncateSettings(ctx context.Context) {
	_, err := db.Exec(ctx, `TRUNCATE TABLE referral.default_settings`)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
}
