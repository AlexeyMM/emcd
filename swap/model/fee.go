package model

import "github.com/shopspring/decimal"

type Fee struct {
	Symbol   string
	MakerFee decimal.Decimal
	TakerFee decimal.Decimal
}
