package repository

import (
	"context"

	"code.emcdtech.com/emcd/service/coin/model"
)

//go:generate mockery --name=Coin --structname=MockCoin --outpkg=mocks --output ./mocks --filename $GOFILE

type Coin interface {
	GetCoinFromLegacyID(ctx context.Context, legacyCoinID int32) (*model.Coin, error)
	GetCoin(ctx context.Context, coinID string) (*model.Coin, error)
	GetCoins(ctx context.Context, limit, offset int32) ([]*model.Coin, int32, error)
	GetCoinsNetworks(ctx context.Context) ([]*model.CoinNetwork, error)
	GetCoinNetwork(ctx context.Context, coinID, networkID string) (*model.CoinNetwork, error)
}
