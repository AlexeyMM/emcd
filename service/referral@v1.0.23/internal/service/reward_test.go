package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"

	"code.emcdtech.com/emcd/service/referral/internal/model"
	"code.emcdtech.com/emcd/service/referral/internal/repository"
	"code.emcdtech.com/emcd/service/referral/internal/service"
)

func TestReferral_Calculate(t *testing.T) {
	ctx := context.Background()

	referralRps := repository.NewMockReferral(t)

	// d.Repository.Profile,
	//	d.Repository.DefaultSettings,
	//	d.Repository.DefaultWhitelabelSettings,

	svc := service.NewReward(referralRps, nil, nil, nil)
	t.Run("Empty_WhitelabelID_AND_ReferralID", func(t *testing.T) {
		in := &model.Referral{
			UserID:  uuid.New(),
			Product: "p2p",
			Coin:    "btc",
			Fee:     decimal.NewFromFloat(10),
		}

		amount := decimal.NewFromFloat(100)

		referralRps.On("Get", ctx, in.UserID, in.Product, in.Coin).Return(in, nil)

		// normalAmount(90) + systemFee(10)
		txs, err := svc.Calculate(ctx, in.UserID, in.Product, in.Coin, amount)
		require.NoError(t, err)
		require.Equal(t, 2, len(txs))

		require.Equal(t, in.UserID, txs[0].UserID)
		require.Equal(t, float64(90), txs[0].Amount.InexactFloat64())
		require.Equal(t, "normal", txs[0].Type)
		require.Equal(t, in.UserID, txs[1].UserID)
		require.Equal(t, float64(10), txs[1].Amount.InexactFloat64())
		require.Equal(t, "fee", txs[1].Type)
	})

	t.Run("Empty_ReferralID", func(t *testing.T) {
		in := &model.Referral{
			UserID:        uuid.New(),
			Product:       "p2p",
			Coin:          "btc",
			Fee:           decimal.NewFromFloat(10),
			WhitelabelFee: decimal.NewFromFloat(20),
			WhitelabelID:  uuid.New(),
		}

		amount := decimal.NewFromFloat(100)

		referralRps.On("Get", ctx, in.UserID, in.Product, in.Coin).Return(in, nil)

		// normalAmount(70) + systemFee(10) + whitelabelFee(20)
		txs, err := svc.Calculate(ctx, in.UserID, in.Product, in.Coin, amount)
		require.NoError(t, err)
		require.Equal(t, 3, len(txs))

		require.Equal(t, in.UserID, txs[0].UserID)
		require.Equal(t, float64(70), txs[0].Amount.InexactFloat64())
		require.Equal(t, "normal", txs[0].Type)
		require.Equal(t, in.UserID, txs[1].UserID)
		require.Equal(t, float64(10), txs[1].Amount.InexactFloat64())
		require.Equal(t, "fee", txs[1].Type)
		require.Equal(t, in.WhitelabelID, txs[2].UserID)
		require.Equal(t, float64(20), txs[2].Amount.InexactFloat64())
		require.Equal(t, "wlFee", txs[2].Type)
	})

	t.Run("Empty_WhitelabelID", func(t *testing.T) {
		in := &model.Referral{
			UserID:      uuid.New(),
			Product:     "p2p",
			Coin:        "btc",
			Fee:         decimal.NewFromFloat(10),
			ReferralFee: decimal.NewFromFloat(20),
			ReferralID:  uuid.New(),
		}

		amount := decimal.NewFromFloat(100)

		referralRps.On("Get", ctx, in.UserID, in.Product, in.Coin).Return(in, nil)

		// normalAmount(90) + systemFee(8) + referralFee(2)
		txs, err := svc.Calculate(ctx, in.UserID, in.Product, in.Coin, amount)
		require.NoError(t, err)
		require.Equal(t, 3, len(txs))

		require.Equal(t, in.UserID, txs[0].UserID)
		require.Equal(t, float64(90), txs[0].Amount.InexactFloat64())
		require.Equal(t, "normal", txs[0].Type)
		require.Equal(t, in.UserID, txs[1].UserID)
		require.Equal(t, float64(8), txs[1].Amount.InexactFloat64())
		require.Equal(t, "fee", txs[1].Type)
		require.Equal(t, in.ReferralID, txs[2].UserID)
		require.Equal(t, float64(2), txs[2].Amount.InexactFloat64())
		require.Equal(t, "referral", txs[2].Type)
	})

	t.Run("Full", func(t *testing.T) {
		in := &model.Referral{
			UserID:        uuid.New(),
			Product:       "p2p",
			Coin:          "btc",
			WhitelabelID:  uuid.New(),
			Fee:           decimal.NewFromFloat(10),
			WhitelabelFee: decimal.NewFromFloat(20),
			ReferralFee:   decimal.NewFromFloat(20),
			ReferralID:    uuid.New(),
		}

		amount := decimal.NewFromFloat(100)

		referralRps.On("Get", ctx, in.UserID, in.Product, in.Coin).Return(in, nil)

		// normalAmount(70) + systemFee(10) + whitelabelFee(16) + referralFee(4)
		txs, err := svc.Calculate(ctx, in.UserID, in.Product, in.Coin, amount)
		require.NoError(t, err)
		require.Equal(t, 4, len(txs))

		require.Equal(t, in.UserID, txs[0].UserID)
		require.Equal(t, float64(70), txs[0].Amount.InexactFloat64())
		require.Equal(t, "normal", txs[0].Type)
		require.Equal(t, in.UserID, txs[1].UserID)
		require.Equal(t, float64(10), txs[1].Amount.InexactFloat64())
		require.Equal(t, "fee", txs[1].Type)
		require.Equal(t, in.WhitelabelID, txs[2].UserID)
		require.Equal(t, float64(16), txs[2].Amount.InexactFloat64())
		require.Equal(t, "wlFee", txs[2].Type)
		require.Equal(t, in.ReferralID, txs[3].UserID)
		require.Equal(t, float64(4), txs[3].Amount.InexactFloat64())
		require.Equal(t, "referral", txs[3].Type)
	})
}
