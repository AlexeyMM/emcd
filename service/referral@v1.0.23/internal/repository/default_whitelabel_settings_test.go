package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"

	"code.emcdtech.com/emcd/service/referral/internal/model"
	"code.emcdtech.com/emcd/service/referral/internal/repository"
)

func TestDefaultWhitelabelSettings_Create(t *testing.T) {
	ctx := context.Background()

	repo := repository.NewDefaultWhitelabelSettings(db)

	in := &model.DefaultWhitelabelSettingsV2{
		WhitelabelID:  uuid.New(),
		Product:       "p2p",
		Coin:          "btc",
		Fee:           decimal.NewFromFloat(15),
		ReferralFee:   decimal.NewFromFloat(12),
		WhiteLabelFee: decimal.NewFromFloat(8.888),
	}

	err := repo.Create(ctx, in)
	require.NoError(t, err)

	defer func(wlID uuid.UUID, product string, coin string) {
		err = repo.Delete(ctx, in.Product, in.Coin, in.WhitelabelID)
		require.NoError(t, err)
	}(in.WhitelabelID, in.Product, in.Coin)

	out, err := repo.GetV2ByCoin(ctx, in.Product, in.Coin, in.WhitelabelID)
	require.NoError(t, err)
	verifyEqualWlSettingsV2(t, in, out)
}

func TestDefaultWhitelabelSettings_Update(t *testing.T) {
	ctx := context.Background()

	repo := repository.NewDefaultWhitelabelSettings(db)

	in := &model.DefaultWhitelabelSettingsV2{
		WhitelabelID:  uuid.New(),
		Product:       "p2p",
		Coin:          "btc",
		Fee:           decimal.NewFromFloat(15),
		ReferralFee:   decimal.NewFromFloat(12),
		WhiteLabelFee: decimal.NewFromFloat(22.222),
	}

	err := repo.Create(ctx, in)
	require.NoError(t, err)

	defer func(wlID uuid.UUID, product string, coin string) {
		err = repo.Delete(ctx, in.Product, in.Coin, in.WhitelabelID)
		require.NoError(t, err)
	}(in.WhitelabelID, in.Product, in.Coin)

	out, err := repo.GetV2ByCoin(ctx, in.Product, in.Coin, in.WhitelabelID)
	require.NoError(t, err)
	verifyEqualWlSettingsV2(t, in, out)

	in.Fee = decimal.NewFromFloat(34)
	in.ReferralFee = decimal.NewFromFloat(36)

	err = repo.Update(ctx, in)
	require.NoError(t, err)

	out, err = repo.GetV2ByCoin(ctx, in.Product, in.Coin, in.WhitelabelID)
	require.NoError(t, err)
	verifyEqualWlSettingsV2(t, in, out)
}

func TestDefaultWhitelabelSettings_Delete(t *testing.T) {
	ctx := context.Background()

	repo := repository.NewDefaultWhitelabelSettings(db)

	in := &model.DefaultWhitelabelSettingsV2{
		WhitelabelID:  uuid.New(),
		Product:       "p2p",
		Coin:          "btc",
		Fee:           decimal.NewFromFloat(15),
		ReferralFee:   decimal.NewFromFloat(12),
		WhiteLabelFee: decimal.NewFromFloat(5.555),
	}

	err := repo.Create(ctx, in)
	require.NoError(t, err)

	defer func(wlID uuid.UUID, product string, coin string) {
		err = repo.Delete(ctx, in.Product, in.Coin, in.WhitelabelID)
		require.NoError(t, err)
	}(in.WhitelabelID, in.Product, in.Coin)

	out, err := repo.GetV2ByCoin(ctx, in.Product, in.Coin, in.WhitelabelID)
	require.NoError(t, err)
	verifyEqualWlSettingsV2(t, in, out)

	err = repo.Delete(ctx, in.Product, in.Coin, in.WhitelabelID)
	require.NoError(t, err)

	out, err = repo.GetV2ByCoin(ctx, in.Product, in.Coin, in.WhitelabelID)
	require.NoError(t, err)
	require.Nil(t, out)
}

func TestDefaultWhitelabelSettings_GetAll(t *testing.T) {
	ctx := context.Background()

	repo := repository.NewDefaultWhitelabelSettings(db)

	in := &model.DefaultWhitelabelSettingsV2{
		WhitelabelID:  uuid.New(),
		Product:       "p2p",
		Coin:          "btc",
		Fee:           decimal.NewFromFloat(15),
		ReferralFee:   decimal.NewFromFloat(12),
		WhiteLabelFee: decimal.NewFromFloat(5.533),
	}

	err := repo.Create(ctx, in)
	require.NoError(t, err)

	defer func(wlID uuid.UUID, product string, coin string) {
		err = repo.Delete(ctx, in.Product, in.Coin, in.WhitelabelID)
		require.NoError(t, err)
	}(in.WhitelabelID, in.Product, in.Coin)

	in.Coin = "ltc"

	err = repo.Create(ctx, in)
	require.NoError(t, err)

	defer func(wlID uuid.UUID, product string, coin string) {
		err = repo.Delete(ctx, in.Product, in.Coin, in.WhitelabelID)
		require.NoError(t, err)
	}(in.WhitelabelID, in.Product, in.Coin)

	in.Product = "mining"

	err = repo.Create(ctx, in)
	require.NoError(t, err)

	defer func(wlID uuid.UUID, product string, coin string) {
		err = repo.Delete(ctx, in.Product, in.Coin, in.WhitelabelID)
		require.NoError(t, err)
	}(in.WhitelabelID, in.Product, in.Coin)

	list, _, err := repo.GetAllWithFilters(ctx, 0, 3, make(map[string]string))
	require.NoError(t, err)
	require.Equal(t, 3, len(list))
}

func TestDefaultWhitelabelSettings_GetAllWithoutPagination(t *testing.T) {
	ctx := context.Background()

	repo := repository.NewDefaultWhitelabelSettings(db)
	truncateWlSettings(ctx)

	defer truncateWlSettings(ctx)
	exp1 := model.DefaultWhitelabelSettingsV2{
		WhitelabelID:  uuid.New(),
		Product:       "p2p",
		Coin:          "btc",
		Fee:           decimal.NewFromFloat(15),
		ReferralFee:   decimal.NewFromFloat(12),
		WhiteLabelFee: decimal.NewFromFloat(5.533),
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

	actual, err := repo.GetAllWithoutPaginationWithFilters(ctx, make(map[string]string))
	require.NoError(t, err)

	for i, exp := range []*model.DefaultWhitelabelSettingsV2{&exp1, &exp2, &exp3} {
		verifyEqualWlSettingsV2(t, exp, actual[i])
	}
}

func TestDefaultWhitelabelSettings_GetAllWithoutPaginationWithFilters(t *testing.T) {
	ctx := context.Background()

	repo := repository.NewDefaultWhitelabelSettings(db)
	truncateWlSettings(ctx)

	defer truncateWlSettings(ctx)
	exp1 := model.DefaultWhitelabelSettingsV2{
		WhitelabelID:  uuid.New(),
		Product:       "p2p",
		Coin:          "btc",
		Fee:           decimal.NewFromFloat(15),
		ReferralFee:   decimal.NewFromFloat(12),
		WhiteLabelFee: decimal.NewFromFloat(5.533),
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

	actual, err := repo.GetAllWithoutPaginationWithFilters(ctx, map[string]string{
		"coin": "btc",
	})
	require.NoError(t, err)
	require.Len(t, actual, 2)
	verifyEqualWlSettingsV2(t, &exp1, actual[0])
	verifyEqualWlSettingsV2(t, &exp3, actual[1])
}

func TestDefaultWhitelabelSettings_GetAllWithoutPaginationWithAllFilters(t *testing.T) {
	ctx := context.Background()

	repo := repository.NewDefaultWhitelabelSettings(db)
	truncateWlSettings(ctx)

	defer truncateWlSettings(ctx)
	exp1 := model.DefaultWhitelabelSettingsV2{
		WhitelabelID:  uuid.New(),
		Product:       "p2p",
		Coin:          "btc",
		Fee:           decimal.NewFromFloat(15),
		ReferralFee:   decimal.NewFromFloat(12),
		WhiteLabelFee: decimal.NewFromFloat(5.533),
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

	actual, err := repo.GetAllWithoutPaginationWithFilters(ctx, map[string]string{
		"coin":    "btc",
		"product": "p2p",
	})
	require.NoError(t, err)
	require.Len(t, actual, 1)
	verifyEqualWlSettingsV2(t, &exp1, actual[0])
}

func verifyEqualWlSettingsV2(t *testing.T, in, out *model.DefaultWhitelabelSettingsV2) {
	require.Equal(t, in.WhitelabelID, out.WhitelabelID)
	require.Equal(t, in.Product, out.Product)
	require.Equal(t, in.Coin, out.Coin)
	require.Equal(t, in.Fee, out.Fee)
	require.Equal(t, in.ReferralFee, out.ReferralFee)
	require.Equal(t, time.Now().UTC().Truncate(time.Second), out.CreatedAt.Truncate(time.Second))
	require.Equal(t, in.WhiteLabelFee, out.WhiteLabelFee)
}

func truncateWlSettings(ctx context.Context) {
	_, err := db.Exec(ctx, `TRUNCATE TABLE referral.default_whitelabel_settings`)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
}

func TestDefaultWhitelabelSettings_GetV2(t *testing.T) {
	ctx := context.Background()

	defer truncateWlSettings(ctx)
	wlID := uuid.New()
	input := []*model.DefaultWhitelabelSettingsV2{
		{
			Product:       "min",
			Coin:          "BTC",
			WhitelabelID:  wlID,
			Fee:           decimal.NewFromFloat(0.3432),
			WhiteLabelFee: decimal.NewFromFloat(0.22),
			CreatedAt:     time.Now().UTC(),
			ReferralFee:   decimal.NewFromFloat(0.2),
		},
		{
			Product:       "min",
			Coin:          "LTC",
			WhitelabelID:  wlID,
			Fee:           decimal.NewFromFloat(0.32),
			WhiteLabelFee: decimal.NewFromFloat(0.12),
			CreatedAt:     time.Now().UTC(),
			ReferralFee:   decimal.NewFromFloat(0.22),
		},
		{
			Product:       "min",
			Coin:          "BTC",
			WhitelabelID:  uuid.New(),
			Fee:           decimal.NewFromFloat(0.3432),
			WhiteLabelFee: decimal.NewFromFloat(0.22),
			CreatedAt:     time.Now().UTC(),
			ReferralFee:   decimal.NewFromFloat(0.2),
		},
	}

	repo := repository.NewDefaultWhitelabelSettings(db)
	createWhitelabelSettingsV2(ctx, db, input)
	actual, err := repo.GetV2(ctx, wlID)
	require.NoError(t, err)
	require.Len(t, actual, 2)
	for i := 0; i < 2; i++ {
		verifyEqualWlSettingsV2(t, input[i], actual[i])
	}
}

func createWhitelabelSettingsV2(ctx context.Context, pool *pgxpool.Pool, ds []*model.DefaultWhitelabelSettingsV2) {
	query := `INSERT INTO referral.default_whitelabel_settings (product, whitelabel_id, coin, fee, referral_fee, created_at, whitelabel_fee)
				VALUES ($1, $2, $3, $4, $5, $6, $7)`
	for i := range ds {
		_, err := pool.Exec(ctx, query, ds[i].Product, ds[i].WhitelabelID, ds[i].Coin, ds[i].Fee, ds[i].ReferralFee, ds[i].CreatedAt, ds[i].WhiteLabelFee)
		if err != nil {
			log.Fatal().Msgf("createWhitelabelSettingsV2: %v", err)
		}
	}
}

func TestDefaultWhitelabelSettings_GetV2ByCoin(t *testing.T) {
	ctx := context.Background()

	defer truncateWlSettings(ctx)
	wlID := uuid.New()
	input := []*model.DefaultWhitelabelSettingsV2{
		{
			Product:       "min",
			Coin:          "BTC",
			WhitelabelID:  wlID,
			Fee:           decimal.NewFromFloat(0.3432),
			WhiteLabelFee: decimal.NewFromFloat(0.22),
			CreatedAt:     time.Now().UTC(),
			ReferralFee:   decimal.NewFromFloat(0.2),
		},
		{
			Product:       "min",
			Coin:          "LTC",
			WhitelabelID:  wlID,
			Fee:           decimal.NewFromFloat(0.32),
			WhiteLabelFee: decimal.NewFromFloat(0.12),
			CreatedAt:     time.Now().UTC(),
			ReferralFee:   decimal.NewFromFloat(0.22),
		},
		{
			Product:       "min",
			Coin:          "BTC",
			WhitelabelID:  uuid.New(),
			Fee:           decimal.NewFromFloat(0.3432),
			WhiteLabelFee: decimal.NewFromFloat(0.22),
			CreatedAt:     time.Now().UTC(),
			ReferralFee:   decimal.NewFromFloat(0.2),
		},
	}

	repo := repository.NewDefaultWhitelabelSettings(db)
	createWhitelabelSettingsV2(ctx, db, input)
	actual, err := repo.GetV2ByCoin(ctx, "min", "LTC", wlID)
	require.NoError(t, err)
	verifyEqualWlSettingsV2(t, input[1], actual)
}
