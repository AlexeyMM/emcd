package model

import (
	"sync"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"
)

// OrderBook инициализация только c NewOrderBook()
type OrderBook struct {
	Mu   sync.RWMutex
	Bids *treemap.Map // Цены отсортированы от наибольшей к наименьшей
	Asks *treemap.Map // Цены отсортированы от наименьшей к наибольшей
}

func NewOrderBook() *OrderBook {
	bidComparator := func(a, b interface{}) int {
		return -utils.Float64Comparator(a, b)
	}

	return &OrderBook{
		Bids: treemap.NewWith(bidComparator),
		Asks: treemap.NewWith(utils.Float64Comparator),
	}
}

type OrderBookUpdateMessage struct {
	Symbol     string
	Bids       [][2]string
	Asks       [][2]string
	IsSnapshot bool
}
