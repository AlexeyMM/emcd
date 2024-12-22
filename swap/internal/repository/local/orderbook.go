package local

import (
	"fmt"
	"strconv"
	"sync"

	"code.emcdtech.com/b2b/swap/internal/business_error"
	"code.emcdtech.com/b2b/swap/model"
	"github.com/emirpasic/gods/maps/treemap"
)

type OrderBookStore struct {
	mu         sync.RWMutex
	orderBooks map[string]*model.OrderBook
}

func NewOrderBookStore() *OrderBookStore {
	return &OrderBookStore{
		orderBooks: make(map[string]*model.OrderBook),
	}
}

func (o *OrderBookStore) Init(symbols []*model.Symbol) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	for i := range symbols {
		o.orderBooks[symbols[i].Title] = model.NewOrderBook()
	}

	return nil
}

func (o *OrderBookStore) AddSnapshot(symbol string, bids, asks [][2]string) error {
	o.mu.Lock()

	ob := model.NewOrderBook()
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	o.orderBooks[symbol] = ob

	o.mu.Unlock()

	if len(bids) > 0 {
		err := o.update(ob.Bids, bids)
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}
	}

	if len(asks) > 0 {
		err := o.update(ob.Asks, asks)
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}
	}

	return nil
}

func (o *OrderBookStore) AddDelta(symbol string, bids, asks [][2]string) error {
	o.mu.Lock()
	ob, exists := o.orderBooks[symbol]
	if !exists {
		o.mu.Unlock()
		return businessError.OrderBookNotFoundErr
	}

	ob.Mu.Lock()
	defer ob.Mu.Unlock()

	o.mu.Unlock()

	if len(bids) > 0 {
		err := o.update(ob.Bids, bids)
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}
	}

	if len(asks) > 0 {
		err := o.update(ob.Asks, asks)
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}
	}

	return nil
}

func (o *OrderBookStore) update(tm *treemap.Map, delta [][2]string) error {
	for _, level := range delta {
		price, err := strconv.ParseFloat(level[0], 64)
		if err != nil {
			return fmt.Errorf("invalid price in update orderbook: %w", err)
		}
		size, err := strconv.ParseFloat(level[1], 64)
		if err != nil {
			return fmt.Errorf("invalid size in update orderbook: %w", err)
		}

		if size == 0 {
			tm.Remove(price)
		} else {
			tm.Put(price, size)
		}
	}
	return nil
}

func (o *OrderBookStore) GetBidTopLevels(symbol string, n int) ([][2]float64, error) {
	o.mu.RLock()
	ob, exists := o.orderBooks[symbol]
	if !exists {
		o.mu.RUnlock()
		return nil, businessError.OrderBookNotFoundErr
	}

	ob.Mu.RLock()
	defer ob.Mu.RUnlock()

	o.mu.RUnlock()

	var bids [][2]float64

	bidsIterator := ob.Bids.Iterator()
	for bidsIterator.Begin(); bidsIterator.Next() && len(bids) < n; {
		price := bidsIterator.Key().(float64)
		size := bidsIterator.Value().(float64)
		bids = append(bids, [2]float64{price, size})
	}

	return bids, nil
}

func (o *OrderBookStore) GetAskTopLevels(symbol string, n int) ([][2]float64, error) {
	o.mu.RLock()
	ob, exists := o.orderBooks[symbol]
	if !exists {
		o.mu.RUnlock()
		return nil, businessError.OrderBookNotFoundErr
	}

	ob.Mu.RLock()
	defer ob.Mu.RUnlock()

	o.mu.RUnlock()

	var asks [][2]float64

	asksIterator := ob.Asks.Iterator()
	for asksIterator.Begin(); asksIterator.Next() && len(asks) < n; {
		price := asksIterator.Key().(float64)
		size := asksIterator.Value().(float64)
		asks = append(asks, [2]float64{price, size})
	}

	return asks, nil
}

func (o *OrderBookStore) IsExist(symbol string) bool {
	o.mu.RLock()
	_, exists := o.orderBooks[symbol]
	o.mu.RUnlock()
	return exists
}

func (o *OrderBookStore) Len() int {
	o.mu.Lock()
	defer o.mu.Unlock()

	return len(o.orderBooks)
}
