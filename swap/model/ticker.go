package model

import "github.com/shopspring/decimal"

type Ticker struct {
	Symbol string
	Ask    decimal.Decimal
	Bid    decimal.Decimal
}
