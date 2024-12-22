package coin

import (
	"context"
	"fmt"

	"code.emcdtech.com/b2b/processing/internal/client"
	"code.emcdtech.com/b2b/processing/model"
	coinspb "code.emcdtech.com/emcd/service/coin/protocol/coin"
)

const (
	getCoinsLimit = 1000
)

type Coin struct {
	coinServiceClient coinspb.CoinServiceClient
}

func NewCoinClient(coinServiceClient coinspb.CoinServiceClient) client.CoinClient {
	return &Coin{coinServiceClient: coinServiceClient}
}

func (c *Coin) GetCoins(ctx context.Context) ([]*model.Coin, error) {
	coins, err := c.coinServiceClient.GetCoins(ctx, &coinspb.GetCoinsRequest{
		Limit:  getCoinsLimit,
		Offset: 0,
	})
	if err != nil {
		return nil, fmt.Errorf("getCoins: %w", err)
	}

	result := make([]*model.Coin, len(coins.GetCoins()))

	for i, coin := range coins.GetCoins() {
		coinNetworks := make([]*model.Network, len(coin.GetNetworks()))

		for j, network := range coin.GetNetworks() {
			coinNetworks[j] = &model.Network{
				ID:     network.GetNetworkId(),
				CoinID: network.GetCoinId(),
				Title:  network.GetTitle(),
			}
		}

		result[i] = &model.Coin{
			ID:          coin.GetId(),
			Title:       coin.GetTitle(),
			Description: coin.GetDescription(),
			MediaURL:    coin.GetMediaUrl(),
			Networks:    coinNetworks,
		}
	}

	return result, nil
}
