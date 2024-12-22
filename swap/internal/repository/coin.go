package repository

import (
	"context"

	"code.emcdtech.com/b2b/swap/model"
)

type Coin interface {
	UpdateAll(ctx context.Context, coins []*model.Coin) error
	Get(ctx context.Context, coin string) (*model.Coin, error)
	GetNetwork(ctx context.Context, coin, network string) (*model.Network, error)
	GetAnyNetwork(ctx context.Context, coin string) (*model.Network, error)
	GetAccuracyForWithdrawAndDeposit(ctx context.Context, coin string, network string) (int, error)
	GetWithdrawFee(ctx context.Context, coin, network string) (*model.WithdrawFee, error)
	GetAll(ctx context.Context) ([]*model.Coin, error)
}
