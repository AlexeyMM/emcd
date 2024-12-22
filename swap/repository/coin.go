package repository

import (
	"context"
	"fmt"

	"code.emcdtech.com/b2b/swap/internal/controller/mapping"
	"code.emcdtech.com/b2b/swap/model"
	"code.emcdtech.com/b2b/swap/protocol/swapCoin"
	"code.emcdtech.com/emcd/sdk/log"
)

type Coin interface {
	GetAll(ctx context.Context) ([]*model.Coin, error)
}

type coin struct {
	handler swapCoin.SwapCoinServiceClient
}

func NewCoin(handler swapCoin.SwapCoinServiceClient) *coin {
	return &coin{
		handler,
	}
}

func (c *coin) GetAll(ctx context.Context) ([]*model.Coin, error) {
	resp, err := c.handler.GetAll(ctx, &swapCoin.GetAllRequest{})
	if err != nil {
		log.Error(ctx, "coin: getAll: %s", err.Error())
		return nil, fmt.Errorf("getAll: %w", err)
	}
	return mapping.MapProtoGetAllResponseToCoins(resp), nil
}
