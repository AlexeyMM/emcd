package model

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Tariff struct {
	UpperFee decimal.Decimal
	LowerFee decimal.Decimal
	MinPay   decimal.Decimal
	MaxPay   decimal.Decimal
}

type Merchant struct {
	ID     uuid.UUID
	Tariff *Tariff
}
