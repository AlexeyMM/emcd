package service_test

import (
	"context"
	"fmt"
	"testing"

	wlPb "code.emcdtech.com/emcd/service/whitelabel/protocol/whitelabel"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"code.emcdtech.com/emcd/service/coin/internal/config"
	repository "code.emcdtech.com/emcd/service/coin/internal/repository/mocks"
	srvMock "code.emcdtech.com/emcd/service/coin/internal/repository/mocks"
	"code.emcdtech.com/emcd/service/coin/internal/service"
	"code.emcdtech.com/emcd/service/coin/model"
)

func TestCoin_GetCoins(t *testing.T) {
	ctx := context.Background()

	coinRepoMock := repository.MockCoin{}
	wlMock := srvMock.MockWhiteLabel{}

	svc := service.NewCoin(&coinRepoMock, nil, &wlMock, nil, config.Config{})

	t.Run("Success_empty_wl_id", func(t *testing.T) {
		limit := int32(100)
		offset := int32(0)

		coinRepoMock.On("GetCoinsNetworks", ctx).Return([]*model.CoinNetwork{
			{
				CoinID: "btc",
				Title:  "bitcoin",
			},
			{
				CoinID: "ltc",
				Title:  "lightcoin",
			},
			{
				CoinID: "kas",
				Title:  "kaspa",
			},
		}, nil).Once()

		coinRepoMock.On("GetCoins", ctx, limit, offset).
			Return([]*model.Coin{
				{
					ID:       "btc",
					IsActive: true,
				},
				{
					ID:       "ltc",
					IsActive: true,
				},
				{
					ID:       "kas",
					IsActive: true,
				},
			}, int32(3), nil).Once()

		res, count, err := svc.GetCoins(ctx, nil, limit, offset)
		require.NoError(t, err)
		require.True(t, len(res) == 3)
		require.True(t, count == 3)
	})

	t.Run("Success_empty_with_id", func(t *testing.T) {
		wlID := uuid.NewString()
		limit := int32(100)
		offset := int32(0)

		coinRepoMock.On("GetCoinsNetworks", ctx).Return([]*model.CoinNetwork{
			{
				CoinID: "btc",
				Title:  "bitcoin",
			},
			{
				CoinID: "ltc",
				Title:  "lightcoin",
			},
			{
				CoinID: "kas",
				Title:  "kaspa",
			},
		}, nil).Once()

		coinRepoMock.On("GetCoins", ctx, limit, offset).
			Return([]*model.Coin{
				{
					ID:       "btc",
					IsActive: true,
				},
				{
					ID:       "ltc",
					IsActive: true,
				},
				{
					ID:       "kas",
					IsActive: true,
				},
			}, int32(3), nil).Once()

		wlMock.On("GetCoins", ctx, &wlPb.GetCoinsRequest{
			WlId: wlID,
		}).Return(&wlPb.GetCoinsResponse{
			Coins: []*wlPb.Coin{
				{
					CoinId: "ltc",
				},
				{
					CoinId: "kas",
				},
			},
		}, nil)

		res, count, err := svc.GetCoins(ctx, &wlID, limit, offset)
		require.NoError(t, err)
		require.True(t, len(res) == 2)
		require.True(t, count == 2)
		require.True(t, res[0].ID == "ltc")
		require.True(t, res[1].ID == "kas")
	})

	t.Run("Failed_wl_id_not_found", func(t *testing.T) {
		wlID := uuid.NewString()
		limit := int32(100)
		offset := int32(0)

		coinRepoMock.On("GetCoinsNetworks", ctx).Return([]*model.CoinNetwork{
			{
				CoinID: "btc",
				Title:  "bitcoin",
			},
			{
				CoinID: "ltc",
				Title:  "lightcoin",
			},
			{
				CoinID: "kas",
				Title:  "kaspa",
			},
		}, nil).Once()

		coinRepoMock.On("GetCoins", ctx, limit, offset).
			Return([]*model.Coin{
				{
					ID:       "btc",
					IsActive: true,
				},
				{
					ID:       "ltc",
					IsActive: true,
				},
				{
					ID:       "kas",
					IsActive: true,
				},
			}, int32(3), nil).Once()

		wlMock.On("GetCoins", ctx, &wlPb.GetCoinsRequest{
			WlId: wlID,
		}).Return(nil, fmt.Errorf("not found"))

		res, count, err := svc.GetCoins(ctx, &wlID, limit, offset)
		require.Error(t, err)
		require.Nil(t, res)
		require.True(t, count == 0)
	})

	t.Run("Failed_wrong_uuid_format", func(t *testing.T) {
		wlID := "uuid.NewString()"
		limit := int32(100)
		offset := int32(0)

		coinRepoMock.On("GetCoinsNetworks", ctx).Return([]*model.CoinNetwork{
			{
				CoinID: "btc",
				Title:  "bitcoin",
			},
			{
				CoinID: "ltc",
				Title:  "lightcoin",
			},
			{
				CoinID: "kas",
				Title:  "kaspa",
			},
		}, nil).Once()

		coinRepoMock.On("GetCoins", ctx, limit, offset).
			Return([]*model.Coin{
				{
					ID:       "btc",
					IsActive: true,
				},
				{
					ID:       "ltc",
					IsActive: true,
				},
				{
					ID:       "kas",
					IsActive: true,
				},
			}, int32(3), nil).Once()

		res, count, err := svc.GetCoins(ctx, &wlID, limit, offset)
		require.Error(t, err)
		require.Nil(t, res)
		require.True(t, count == 0)
	})
}
