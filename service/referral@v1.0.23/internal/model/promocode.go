package model

import (
	"time"

	"github.com/google/uuid"
)

type PromoCode struct {
	ID                      int32
	Code                    string
	ValidDaysAmount         int32
	HasNoLimit              bool
	FeePercent              float64
	ReferralEnabled         bool
	IsActive                bool
	IsDisposable            bool
	RefID                   int32
	CoinIDs                 []int32
	IsSummable              bool
	IsOnlyForRegistration   bool
	IsOnlyForPrivateCabinet bool
	CreatedAt               time.Time
	ExpiresAt               time.Time
}

type UserPromoCode struct {
	UserID      int32
	PromoCodeID int32
	ExpiresAt   time.Time
	CreatedAt   time.Time
}

type UserPromoCodeAndPromo struct {
	UserPromoCode *UserPromoCode
	PromoCode     *PromoCode
}

type UserAndPromoCodes struct {
	UserUUID   uuid.UUID
	UserID     int32
	PromoCodes []UserPromoCodeAndPromo
}
