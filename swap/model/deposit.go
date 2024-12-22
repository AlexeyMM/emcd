package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type DepositStatus uint8

const (
	DepositSuccess DepositStatus = iota + 1
	DepositPending
	DepositFailed
)

type DepositType uint8

const (
	DepositNormal = iota + 1
	DepositAbnormal
)

type Deposit struct {
	TxID        string
	SwapID      uuid.UUID
	Coin        string
	Amount      decimal.Decimal
	Fee         decimal.Decimal
	Status      DepositStatus
	UpdatedAt   time.Time
	DepositType DepositType
}

type DepositFilter struct {
	SwapID    *uuid.UUID
	UpdatedAt *time.Time // Фильтрует депозиты, обновленные после указанной даты
}

type Deposits []*Deposit
