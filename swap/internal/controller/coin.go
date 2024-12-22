package controller

import (
	"context"
	"fmt"

	"code.emcdtech.com/b2b/swap/internal/controller/mapping"
	"code.emcdtech.com/b2b/swap/internal/service"
	"code.emcdtech.com/b2b/swap/protocol/swapCoin"
)

type Coin struct {
	srv service.Coin
	swapCoin.UnimplementedSwapCoinServiceServer
}

func NewCoin(srv service.Coin) *Coin {
	return &Coin{
		srv: srv,
	}
}

func (c *Coin) GetAll(ctx context.Context, req *swapCoin.GetAllRequest) (*swapCoin.GetAllResponse, error) {
	coins, err := c.srv.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("getAll: %w", err)
	}
	return mapping.MapCoinsToProtoGetAllResponse(coins), nil
}
