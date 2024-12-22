package client

import (
	"context"
	"fmt"
	"strings"

	"code.emcdtech.com/emcd/service/coin/protocol/coin"

	"code.emcdtech.com/b2b/swap/model"
)

type Coin interface {
	GetCoinsInfo(ctx context.Context) (map[string]*model.CoinInfo, error)
}

type CoinImp struct {
	cli coin.CoinServiceClient
}

func NewCoin(cli coin.CoinServiceClient) *CoinImp {
	return &CoinImp{
		cli: cli,
	}
}

func (c *CoinImp) GetCoinsInfo(ctx context.Context) (map[string]*model.CoinInfo, error) {
	resp, err := c.cli.GetCoins(ctx, &coin.GetCoinsRequest{
		Limit:  1000,
		Offset: 0,
	})
	if err != nil {
		return nil, fmt.Errorf("getCoins: %w", err)
	}

	m := make(map[string]*model.CoinInfo, len(resp.Coins))
	for _, co := range resp.Coins {
		m[strings.ToUpper(co.Id)] = &model.CoinInfo{
			Rating:  int(co.SortPrioritySwap),
			IconURL: co.MediaUrl,
		}
	}

	return m, nil
}
