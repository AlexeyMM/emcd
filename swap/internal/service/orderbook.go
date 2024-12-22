package service

import (
	"context"
	"fmt"

	"code.emcdtech.com/b2b/swap/internal/client"
	"code.emcdtech.com/b2b/swap/internal/repository"
	"code.emcdtech.com/b2b/swap/model"
)

const SymbolsPerConnection = 5

type OrderBook interface {
	Subscribe(ctx context.Context, symbols []*model.Symbol) error
	Update(ctx context.Context, symbol string, bids, asks [][2]string, isSnapshot bool) error
}

type orderBook struct {
	orderBookRepo repository.OrderBook
	subscriberCli client.Subscriber
}

func NewOrderBook(orderBookRepo repository.OrderBook, subscriberCli client.Subscriber) *orderBook {
	return &orderBook{
		orderBookRepo: orderBookRepo,
		subscriberCli: subscriberCli,
	}
}

func (o *orderBook) Subscribe(ctx context.Context, symbols []*model.Symbol) error {
	err := o.orderBookRepo.Init(symbols)
	if err != nil {
		return fmt.Errorf("init: %w", err)
	}

	var subSymbols []*model.Symbol
	for i := range symbols {
		subSymbols = append(subSymbols, symbols[i])
		if (i+1)%SymbolsPerConnection == 0 {
			err = o.subscriberCli.SubscribeOnOrderbooks(ctx, subSymbols)
			if err != nil {
				return fmt.Errorf("subscribeOnOrders 1: %w", err)
			}
			subSymbols = []*model.Symbol{}
		}
	}

	if len(subSymbols) > 0 {
		err = o.subscriberCli.SubscribeOnOrderbooks(ctx, subSymbols)
		if err != nil {
			return fmt.Errorf("subscribeOnOrders 2: %w", err)
		}
	}

	return nil
}

func (o *orderBook) Update(ctx context.Context, symbol string, bids, asks [][2]string, isSnapshot bool) error {
	if isSnapshot {
		err := o.orderBookRepo.AddSnapshot(symbol, bids, asks)
		if err != nil {
			return fmt.Errorf("addSnapshot: %w", err)
		}
	} else {
		err := o.orderBookRepo.AddDelta(symbol, bids, asks)
		if err != nil {
			return fmt.Errorf("addDelta: %w", err)
		}
	}
	return nil
}
