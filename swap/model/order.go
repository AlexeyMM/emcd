package model

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Direction uint8

const (
	Sell = iota + 1
	Buy
)

func (d Direction) String() string {
	switch d {
	case Buy:
		return "buy"
	case Sell:
		return "sell"
	}
	return "unknown"
}

type OrderStatus uint8

const (
	OrderUnknown OrderStatus = iota
	OrderCreated
	OrderPending
	OrderFilled
	OrderPartiallyFilled
	OrderFailed
)

type Order struct {
	ID         uuid.UUID
	SwapID     uuid.UUID
	AccountID  int64
	Category   string
	Symbol     string
	Direction  Direction
	AmountFrom decimal.Decimal
	AmountTo   decimal.Decimal
	Status     OrderStatus
	IsFirst    bool
}

type OrderFilter struct {
	ID        *uuid.UUID
	AccountID *int64
	IsFirst   *bool
	LtStatus  *OrderStatus // Для фильтрации ордеров со статусом меньше указанного
}

type OrderPartial struct {
	AmountFrom *decimal.Decimal
	AmountTo   *decimal.Decimal
	Status     *OrderStatus
}

func (o *Order) Update(partial *OrderPartial) {
	if partial.AmountFrom != nil {
		o.AmountFrom = *partial.AmountFrom
	}
	if partial.AmountTo != nil {
		o.AmountTo = *partial.AmountTo
	}
	if partial.Status != nil {
		o.Status = *partial.Status
	}
}

type Orders []*Order

func (os Orders) FindFirst() (*Order, error) {
	for i := range os {
		if os[i].IsFirst {
			return os[i], nil
		}
	}
	return nil, fmt.Errorf("order not found")
}

func (os Orders) FindSecond() (*Order, error) {
	for i := range os {
		if !os[i].IsFirst {
			return os[i], nil
		}
	}
	return nil, fmt.Errorf("order not found")
}
