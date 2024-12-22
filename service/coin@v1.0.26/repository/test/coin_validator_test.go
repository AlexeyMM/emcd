package test

import (
	coinPb "code.emcdtech.com/emcd/service/coin/protocol/coin"
	coinMockPb "code.emcdtech.com/emcd/service/coin/protocol/coin/mocks"
	"code.emcdtech.com/emcd/service/coin/repository"
	"context"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

const cacheUpdateInterval = time.Millisecond

func TestCoinValidator_Serve(t *testing.T) {
	t.Run("normal work flow", func(t *testing.T) {
		coinMock := coinMockPb.NewMockCoinServiceClient(t)
		repo := repository.NewCoinValidatorRepository(coinMock)

		ctx, cancel := context.WithCancel(context.Background())

		wg := new(sync.WaitGroup)

		btcIdLegacy := int32(1)
		usdtIdLegacy := int32(10)
		invalidIdLegacy := int32(999)

		btcCode := "btc"
		usdtCode := "usdt"
		invalidCode := "invalid"

		req := &coinPb.GetCoinsRequest{
			Limit:  9999999,
			Offset: 0,
			WlId:   nil,
		}

		coinBtc := &coinPb.Coin{
			Id:                    btcCode,
			IsActive:              true,
			Title:                 "",
			Description:           "",
			SortPriorityMining:    0,
			SortPriorityWallet:    0,
			MediaUrl:              "",
			IsWithdrawalsDisabled: false,
			Networks:              nil,
			LegacyCoinId:          btcIdLegacy,
			MiningRewardType:      "",
		}

		coinUsdt := &coinPb.Coin{
			Id:                    usdtCode,
			IsActive:              true,
			Title:                 "",
			Description:           "",
			SortPriorityMining:    0,
			SortPriorityWallet:    0,
			MediaUrl:              "",
			IsWithdrawalsDisabled: false,
			Networks:              nil,
			LegacyCoinId:          usdtIdLegacy,
			MiningRewardType:      "",
		}

		coinsResponse := &coinPb.GetCoinsResponse{
			Coins:      []*coinPb.Coin{coinBtc, coinUsdt},
			TotalCount: 0,
		}

		coinMock.EXPECT().
			GetCoins(ctx, req).
			Return(coinsResponse, nil)

		wg.Add(1)
		go repo.Serve(ctx, wg, cacheUpdateInterval)

		time.Sleep(100 * time.Millisecond)

		var id int32
		var code string
		var ok bool

		id, ok = repo.GetIdByCode(btcCode)
		require.Equal(t, id, btcIdLegacy)
		require.True(t, ok)

		id, ok = repo.GetIdByCode(invalidCode)
		require.Empty(t, id)
		require.False(t, ok)

		code, ok = repo.GetCodeById(btcIdLegacy)
		require.Equal(t, code, btcCode)
		require.True(t, ok)

		code, ok = repo.GetCodeById(invalidIdLegacy)
		require.Empty(t, code)
		require.False(t, ok)

		require.True(t, repo.IsValidCode(btcCode))
		require.True(t, repo.IsValidIdLegacy(btcIdLegacy))
		require.False(t, repo.IsValidCode(invalidCode))
		require.False(t, repo.IsValidIdLegacy(invalidIdLegacy))

		require.Len(t, repo.GetCodes(), len(coinsResponse.Coins))
		require.Len(t, repo.GetIdsLegacy(), len(coinsResponse.Coins))

		cancel()
		wg.Wait()

	})

	t.Run("error mock return", func(t *testing.T) {
		coinMock := coinMockPb.NewMockCoinServiceClient(t)
		repo := repository.NewCoinValidatorRepository(coinMock)

		ctx, cancel := context.WithCancel(context.Background())

		wg := new(sync.WaitGroup)

		btcIdLegacy := int32(1)
		invalidIdLegacy := int32(999)

		btcCode := "btc"
		invalidCode := "invalid"

		req := &coinPb.GetCoinsRequest{
			Limit:  9999999,
			Offset: 0,
			WlId:   nil,
		}

		retError := newMockError()

		coinMock.EXPECT().
			GetCoins(ctx, req).
			Return(nil, retError)

		wg.Add(1)
		go repo.Serve(ctx, wg, cacheUpdateInterval)

		time.Sleep(10 * time.Millisecond)

		var id int32
		var code string
		var ok bool

		id, ok = repo.GetIdByCode(btcCode)
		require.Empty(t, id, btcIdLegacy)
		require.False(t, ok)

		id, ok = repo.GetIdByCode(invalidCode)
		require.Empty(t, id)
		require.False(t, ok)

		code, ok = repo.GetCodeById(btcIdLegacy)
		require.Empty(t, code, btcCode)
		require.False(t, ok)

		code, ok = repo.GetCodeById(invalidIdLegacy)
		require.Empty(t, code)
		require.False(t, ok)

		require.False(t, repo.IsValidCode(btcCode))
		require.False(t, repo.IsValidIdLegacy(btcIdLegacy))
		require.False(t, repo.IsValidCode(invalidCode))
		require.False(t, repo.IsValidIdLegacy(invalidIdLegacy))

		cancel()
		wg.Wait()

	})
}
