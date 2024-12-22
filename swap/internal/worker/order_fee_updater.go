package worker

import (
	"context"
	"time"

	"code.emcdtech.com/b2b/swap/internal/service"
	"code.emcdtech.com/emcd/sdk/log"
)

const (
	orderFeeUpdateInterval      = 24 * time.Hour
	ctxTimeoutForOrderFeeUpdate = 30 * time.Second
)

type OrderFeeUpdater struct {
	srv service.OrderFee
}

func NewOrderFeeUpdater(srv service.OrderFee) *OrderFeeUpdater {
	return &OrderFeeUpdater{
		srv: srv,
	}
}

func (o *OrderFeeUpdater) Run(ctx context.Context) error {
	t := time.NewTicker(orderFeeUpdateInterval)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
			log.Debug(ctx, "orderFeeUpdater: uploadFee execute")

			newCtx, cancel := context.WithTimeout(ctx, ctxTimeoutForOrderFeeUpdate)

			err := o.srv.SyncWithAPI(newCtx)
			if err != nil {
				log.Error(ctx, "orderFeeUpdater: syncWithAPI: %s", err.Error())
				cancel()
				continue
			}
			cancel()
		}
	}
}
