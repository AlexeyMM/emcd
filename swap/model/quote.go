package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Quote struct {
	ID           string
	Rate         decimal.Decimal
	FromCoin     string
	FromCoinType string
	ToCoin       string
	ToCoinType   string
	FromAmount   decimal.Decimal
	ToAmount     decimal.Decimal
	ExpiredTime  time.Time
}
