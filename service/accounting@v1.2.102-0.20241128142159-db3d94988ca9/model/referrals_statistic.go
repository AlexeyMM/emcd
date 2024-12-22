package model

import "github.com/google/uuid"

type ReferralsStatisticInput struct {
	UserID              uuid.UUID
	AccountIDs          []int64
	TransactionTypesIDs []int64
}

type ReferralsStatisticOutput struct {
	Income AggregatedReferralIncome
}

type AggregatedReferralIncome struct {
	Yesterday []*ReferralIncome
	ThisMonth []*ReferralIncome
}

type ReferralIncome struct {
	CoinID int64
	Amount float64
}
