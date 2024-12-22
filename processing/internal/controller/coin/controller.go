package coin

import (
	"context"

	"code.emcdtech.com/b2b/processing/internal/service"
	"code.emcdtech.com/b2b/processing/protocol/coinpb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Controller struct {
	coinpb.UnimplementedCoinsServiceServer
	service service.CoinService
}

func NewController(service service.CoinService) *Controller {
	return &Controller{
		service: service,
	}
}

func (c *Controller) GetCoins(_ context.Context, _ *emptypb.Empty) (*coinpb.GetCoinsResponse, error) {
	coins := c.service.GetCoins()
	response := &coinpb.GetCoinsResponse{
		Coins: make([]*coinpb.Coin, len(coins)),
	}

	for i, coin := range coins {
		networks := make([]*coinpb.Network, len(coin.Networks))
		for j, network := range coin.Networks {
			networks[j] = &coinpb.Network{
				Id:    network.ID,
				Title: network.Title,
			}
		}

		response.Coins[i] = &coinpb.Coin{
			Id:          coin.ID,
			Title:       coin.Title,
			Description: coin.Description,
			MediaUrl:    coin.MediaURL,
			Networks:    networks,
		}
	}

	return response, nil
}
