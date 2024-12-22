package service

import (
	"context"
	"fmt"
	"sort"

	"code.emcdtech.com/b2b/swap/internal/client"
	"code.emcdtech.com/b2b/swap/internal/repository"
	"code.emcdtech.com/b2b/swap/model"
)

type Coin interface {
	SyncWithAPI(ctx context.Context) error
	GetAll(ctx context.Context) ([]*model.Coin, error)
}

type coin struct {
	repo     repository.Coin
	marketCl client.Market
	coinCli  client.Coin
}

func NewCoin(repo repository.Coin, marketCl client.Market, coinCli client.Coin) *coin {
	return &coin{
		repo:     repo,
		marketCl: marketCl,
		coinCli:  coinCli,
	}
}

func (c *coin) SyncWithAPI(ctx context.Context) error {
	coins, err := c.marketCl.GetCoinInfo(ctx)
	if err != nil {
		return fmt.Errorf("getCoinInfo: %w", err)
	}

	// Удалить SUI и SUIA из списка
	for i := 0; i < len(coins); {
		if coins[i].Title == "SUI" || coins[i].Title == "SUIA" {
			coins = append(coins[:i], coins[i+1:]...)
			// Не увеличиваем i, чтобы не пропустить сдвинувшийся элемент
		} else {
			i++
		}
	}

	list, err := c.marketCl.GetConvertCoinList(ctx)
	if err != nil {
		return fmt.Errorf("getConvertCoinList: %w", err)
	}

	coinsInfo, err := c.coinCli.GetCoinsInfo(ctx)
	if err != nil {
		return fmt.Errorf("getCoinsInfo: %w", err)
	}

	for i := range coins {
		accuracy, ok := list[coins[i].Title]
		if !ok {
			//log.Warn(ctx, "accuracy for coin: %s not found", coins[i].Title)
		} else {
			coins[i].Accuracy = accuracy
		}

		info, ok := coinsInfo[coins[i].Title]
		if !ok {
			coinInfo := model.CoinInfo{
				Rating: 9999, // монеты без рейтинга должны отображаться в конце списка
			}
			info = &coinInfo
		}
		coins[i].Info = info
	}

	sort.Slice(coins, func(i, j int) bool {
		if coins[i].Info.Rating == coins[j].Info.Rating {
			return coins[i].Title < coins[j].Title
		}
		return coins[i].Info.Rating < coins[j].Info.Rating
	})

	err = c.repo.UpdateAll(ctx, coins)
	if err != nil {
		return fmt.Errorf("updateAll: %w", err)
	}

	return nil
}

func (c *coin) GetAll(ctx context.Context) ([]*model.Coin, error) {
	allCoins, err := c.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("getAll: %w", err)
	}
	return allCoins, nil
}
