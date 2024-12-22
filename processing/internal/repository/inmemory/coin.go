package inmemory

import (
	"sync"

	"code.emcdtech.com/b2b/processing/model"
)

type Coin struct {
	coins    []*model.Coin
	coinsMap map[string]*model.Coin
	sync.RWMutex
}

func NewCoin() *Coin {
	return &Coin{
		coins:    []*model.Coin{},
		coinsMap: make(map[string]*model.Coin),
	}
}

func (c *Coin) SetCoins(coins []*model.Coin) {
	c.Lock()
	defer c.Unlock()
	c.coins = coins
	c.coinsMap = make(map[string]*model.Coin, len(coins))

	for _, coin := range coins {
		c.coinsMap[coin.ID] = coin
	}
}

func (c *Coin) GetCoins() []*model.Coin {
	c.RLock()
	defer c.RUnlock()

	return c.coins
}

func (c *Coin) GetCoin(id string) (*model.Coin, error) {
	c.RLock()
	defer c.RUnlock()

	coin, ok := c.coinsMap[id]
	if !ok {
		return nil, &model.Error{
			Code:    model.ErrorCodeNoSuchCoin,
			Message: "No such coin with id: " + id,
		}
	}

	return coin, nil
}

func (c *Coin) GetNetwork(coinID, networkID string) (*model.Network, error) {
	c.RLock()
	defer c.RUnlock()

	coin, err := c.GetCoin(coinID)
	if err != nil {
		return nil, err
	}

	for _, n := range coin.Networks {
		if n.ID == networkID {
			return n, nil
		}
	}

	return nil, &model.Error{
		Code:    model.ErrorCodeNoSuchNetwork,
		Message: "No such network with id: " + networkID,
	}
}
