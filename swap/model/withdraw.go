package model

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type WithdrawStatus uint8

const (
	WsUnknown WithdrawStatus = iota
	WsSuccess
	WsPending
	WsFailed
	WsBlockchainConfirmed
)

type Withdraw struct {
	InternalID         uuid.UUID
	ID                 int64
	SwapID             uuid.UUID
	HashID             string
	Coin               string
	Network            string
	Address            string
	Tag                string
	Amount             decimal.Decimal
	IncludeFeeInAmount bool
	Status             WithdrawStatus
	ExplorerLink       string
}

type WithdrawFilter struct {
	ID     *int64
	SwapID *uuid.UUID
}

type WithdrawPartial struct {
	ID     *int64
	Amount *decimal.Decimal
	Status *WithdrawStatus
	HashID *string
}

func (w *Withdraw) Update(partial *WithdrawPartial) {
	if partial.ID != nil {
		w.ID = *partial.ID
	}
	if partial.Amount != nil {
		w.Amount = *partial.Amount
	}
	if partial.Status != nil {
		w.Status = *partial.Status
	}
	if partial.HashID != nil {
		w.HashID = *partial.HashID
	}
}

type Withdraws []*Withdraw
