package repository

import (
	"code.emcdtech.com/b2b/processing/model"
)

type Coin interface {
	SetCoins([]*model.Coin)
	GetCoins() []*model.Coin
	GetNetwork(coinID string, networkID string) (*model.Network, error)
}
