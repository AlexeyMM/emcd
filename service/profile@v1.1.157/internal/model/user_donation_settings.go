package model

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type UserDonationSettings struct {
	UserID       uuid.UUID
	IsDonationOn bool
	Percent      decimal.Decimal
}
