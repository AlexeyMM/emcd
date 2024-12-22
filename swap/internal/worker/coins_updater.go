package worker

import (
	"context"
	"time"

	"code.emcdtech.com/b2b/swap/internal/service"
	"code.emcdtech.com/emcd/sdk/log"
)

const (
	coinsUpdateInterval      = 2 * time.Hour
	ctxTimeoutForCoinsUpdate = 30 * time.Second
)

type CoinsUpdater struct {
	coinSrv service.Coin
}

func NewCoinsUpdater(coinSrv service.Coin) *CoinsUpdater {
	return &CoinsUpdater{
		coinSrv: coinSrv,
	}
}

func (c *CoinsUpdater) Run(ctx context.Context) error {
	t := time.NewTicker(coinsUpdateInterval)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
			newCtx, cancel := context.WithTimeout(ctx, ctxTimeoutForCoinsUpdate)

			err := c.coinSrv.SyncWithAPI(newCtx)
			if err != nil {
				log.Error(ctx, "coinsUpdater: syncWithAPI: %w", err)
				cancel()
				continue
			}
			cancel()
		}
	}
}
