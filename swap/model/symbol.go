package model

import "github.com/shopspring/decimal"

type Symbol struct {
	Title          string
	BaseCoin       string
	QuoteCoin      string
	BasePrecision  decimal.Decimal
	QuotePrecision decimal.Decimal
	MinOrderQty    decimal.Decimal
	MaxOrderQty    decimal.Decimal
	MinOrderAmt    decimal.Decimal
	MaxOrderAmt    decimal.Decimal
	Accuracy       *Accuracy
}

type Accuracy struct {
	BaseAccuracy  int32
	QuoteAccuracy int32
}
