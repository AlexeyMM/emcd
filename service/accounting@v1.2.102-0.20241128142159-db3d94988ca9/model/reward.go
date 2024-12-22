package model

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type RefType = string

const (
	Normal RefType = "normal"
	Fee    RefType = "fee"
	WlFee  RefType = "wlFee"
	RefFee RefType = "referral"
)

type ReferralCalculation struct {
	UserId uuid.UUID
	Type   RefType
	Amount decimal.Decimal
}

type UserIncome struct {
	UserID uuid.UUID
	Type   TransactionType
	Amount decimal.Decimal
}
