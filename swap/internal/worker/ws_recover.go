package worker

import (
	"context"
	"time"

	"code.emcdtech.com/b2b/swap/internal/client"
	"code.emcdtech.com/b2b/swap/model"
	"code.emcdtech.com/emcd/sdk/log"
)

var recoveryInterval = time.Minute

type WsRecover struct {
	subscriber          client.Subscriber
	brokenOrderBookWsCh chan []*model.Symbol
}

func NewWsRecover(subscriber client.Subscriber, brokenOrderBookWsCh chan []*model.Symbol) *WsRecover {
	return &WsRecover{
		subscriber:          subscriber,
		brokenOrderBookWsCh: brokenOrderBookWsCh,
	}
}

func (r *WsRecover) Run(ctx context.Context) error {
	log.Debug(ctx, "wsRecover: start")
	defer log.Debug(ctx, "wsRecover: end")

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case symbols := <-r.brokenOrderBookWsCh:
			go r.recoverOrderBookWs(ctx, symbols)
		}
	}
}

func (r *WsRecover) recoverOrderBookWs(ctx context.Context, symbols []*model.Symbol) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(recoveryInterval):
			err := r.subscriber.SubscribeOnOrderbooks(ctx, symbols)
			if err != nil {
				log.Error(ctx, "wsRecover: recoverOrderBookWs: subscribeOnOrderbooks: %w", err)
				continue
			}
			log.Debug(ctx, "wsRecover: orderbook ws successfully recovered: %+v", symbols)
			return
		}
	}
}
