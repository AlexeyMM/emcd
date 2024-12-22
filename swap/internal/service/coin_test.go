package service

import (
	"context"
	"testing"

	"code.emcdtech.com/b2b/swap/mocks/internal_/client"
	"code.emcdtech.com/b2b/swap/mocks/internal_/repository"
	"code.emcdtech.com/b2b/swap/model"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCoin_SyncWithAPI(t *testing.T) {
	ctx := context.Background()

	// Подготовка данных для теста
	mockCoins := []*model.Coin{
		{Title: "BTC"},
		{Title: "ETH"},
		{Title: "ADA"}, // без рейтинга
		{Title: "BCH"}, // без рейтинга
	}

	mockConvertCoinList := map[string]int{
		"BTC": 8,
		"ETH": 18,
	}

	mockCoinsInfo := map[string]*model.CoinInfo{
		"BTC": {Rating: 1, IconURL: "btc.png"},
		"ETH": {Rating: 2, IconURL: "eth.png"},
	}

	updatedCoins := []*model.Coin{
		{
			Title:    "BTC",
			Accuracy: 8,
			Info:     &model.CoinInfo{Rating: 1, IconURL: "btc.png"},
		},
		{
			Title:    "ETH",
			Accuracy: 18,
			Info:     &model.CoinInfo{Rating: 2, IconURL: "eth.png"},
		},

		// Без рейтинга сортируем по алфавиту
		{
			Title:    "ADA",
			Accuracy: 0,
			Info:     &model.CoinInfo{Rating: 9999, IconURL: ""},
		},
		{
			Title:    "BCH",
			Accuracy: 0,
			Info:     &model.CoinInfo{Rating: 9999, IconURL: ""},
		},
	}

	// Моки для репозитория
	coinRep := repository.NewMockCoin(t)
	coinRep.On("UpdateAll", ctx, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		coins := args.Get(1).([]*model.Coin)
		require.Equal(t, updatedCoins, coins)
	})
	coinRep.On("GetAll", ctx).Return(updatedCoins, nil)

	// Моки для клиента marketCl
	marketCli := client.NewMockMarket(t)
	marketCli.On("GetCoinInfo", ctx).Return(mockCoins, nil)
	marketCli.On("GetConvertCoinList", ctx).Return(mockConvertCoinList, nil)

	// Моки для клиента coinCli
	coinCli := client.NewMockCoin(t)
	coinCli.On("GetCoinsInfo", ctx).Return(mockCoinsInfo, nil)

	// Создание сервиса
	coinSrv := NewCoin(coinRep, marketCli, coinCli)

	// Выполнение метода SyncWithAPI
	err := coinSrv.SyncWithAPI(ctx)
	require.NoError(t, err)

	// Проверка результата через GetAll
	actualCoins, err := coinSrv.GetAll(ctx)
	require.NoError(t, err)
	require.Equal(t, updatedCoins, actualCoins)

	// Проверка вызовов моков
	coinRep.AssertExpectations(t)
	marketCli.AssertExpectations(t)
	coinCli.AssertExpectations(t)
}
