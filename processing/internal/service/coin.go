package service

import (
	"context"

	"code.emcdtech.com/b2b/processing/model"
)

type CoinService interface {
	GetCoins() []*model.Coin
	UpdateCoins(ctx context.Context) error
}
