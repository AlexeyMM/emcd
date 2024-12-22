package coin

import (
	"context"
	"fmt"
	"strings"

	"code.emcdtech.com/emcd/service/coin/protocol/coin"
)

type Client struct {
	coinClient coin.CoinServiceClient
}

func NewClient(coinClient coin.CoinServiceClient) *Client {
	return &Client{
		coinClient: coinClient,
	}
}

func (c *Client) GetCoins(ctx context.Context) (map[int32]string, error) {
	resp, err := c.coinClient.GetCoins(ctx, &coin.GetCoinsRequest{Limit: 100, Offset: 0})
	if err != nil {
		return nil, fmt.Errorf("get coins: %w", err)
	}
	coins := make(map[int32]string, len(resp.Coins))
	for _, c2 := range resp.GetCoins() {
		coins[c2.LegacyCoinId] = strings.ToUpper(c2.GetId())
	}
	return coins, nil
}
