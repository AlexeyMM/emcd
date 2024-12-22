package local

import (
	"context"
	"sync"

	"code.emcdtech.com/b2b/swap/internal/business_error"
	"code.emcdtech.com/b2b/swap/model"
)

type Coins struct {
	mu sync.RWMutex
	// The same data in the map and slice
	Map   map[string]*model.Coin
	Slice []*model.Coin
}

// NewCoins needs UpdateAll
func NewCoins() *Coins {
	return &Coins{}
}

func (c *Coins) UpdateAll(ctx context.Context, coins []*model.Coin) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	m := make(map[string]*model.Coin, len(coins))
	for i := range coins {
		m[coins[i].Title] = coins[i]
	}

	c.Map = m
	c.Slice = coins

	return nil
}

func (c *Coins) Get(ctx context.Context, coin string) (*model.Coin, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	myCoin, ok := c.Map[coin]
	if !ok {
		return nil, businessError.CoinNotFoundErr
	}
	return myCoin, nil
}

func (c *Coins) GetNetwork(ctx context.Context, coin, network string) (*model.Network, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	myCoin, ok := c.Map[coin]
	if !ok {
		return nil, businessError.CoinNotFoundErr
	}
	for i := range myCoin.Networks {
		if myCoin.Networks[i].Title == network {
			return myCoin.Networks[i], nil
		}
	}
	return nil, businessError.NetworkNotFountErr
}

func (c *Coins) GetAnyNetwork(ctx context.Context, coin string) (*model.Network, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	myCoin, ok := c.Map[coin]
	if !ok {
		return nil, businessError.CoinNotFoundErr
	}
	for i := range myCoin.Networks {
		return myCoin.Networks[i], nil
	}
	return nil, businessError.NetworkNotFountErr
}

func (c *Coins) GetAccuracyForWithdrawAndDeposit(ctx context.Context, coin string, network string) (int, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	myCoin, ok := c.Map[coin]
	if !ok {
		return 0, nil
	}
	for i := range myCoin.Networks {
		if myCoin.Networks[i].Title == network {
			return myCoin.Networks[i].AccuracyWithdrawAndDeposit, nil
		}
	}
	return 0, nil
}

func (c *Coins) GetWithdrawFee(ctx context.Context, coin, network string) (*model.WithdrawFee, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	myCoin, ok := c.Map[coin]
	if !ok {
		return &model.WithdrawFee{}, nil
	}
	for i := range myCoin.Networks {
		if myCoin.Networks[i].Title == network {
			return myCoin.Networks[i].WithdrawFee, nil
		}
	}

	return &model.WithdrawFee{}, nil
}

func (c *Coins) GetAll(ctx context.Context) ([]*model.Coin, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.Slice, nil
}
