package client

import (
	"context"

	"code.emcdtech.com/b2b/processing/model"
)

type CoinClient interface {
	GetCoins(ctx context.Context) ([]*model.Coin, error)
}
