package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type InternalTransferStatus uint8

const (
	ItsUnknown InternalTransferStatus = iota
	ItsSuccess
	ItsPending
	ItsFailed
)

type InternalTransfer struct {
	ID              uuid.UUID
	Coin            string
	Amount          decimal.Decimal
	FromAccountID   int64
	ToAccountID     int64
	FromAccountType string
	ToAccountType   string
	Status          InternalTransferStatus
	UpdatedAt       time.Time
}

type InternalTransferFilter struct {
	ID            *uuid.UUID
	FromAccountID *int64
	IsLast        *bool
}

type InternalTransferPartial struct {
	Status *InternalTransferStatus
}

func (it *InternalTransfer) Update(filter *InternalTransferPartial) {
	if filter.Status != nil {
		it.Status = *filter.Status
	}
}

type InternalTransfers []*InternalTransfer

type AddressData struct {
	Address string
	Tag     string
}
