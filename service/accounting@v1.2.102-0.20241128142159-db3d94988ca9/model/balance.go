package model

import "github.com/shopspring/decimal"

type CoinSummary struct {
	TotalAmount decimal.Decimal
	CoinID      string
}

type Balance struct {
	WalletBalance               decimal.Decimal
	MiningBalance               decimal.Decimal
	CoinholdsBalance            decimal.Decimal
	P2pBalance                  decimal.Decimal
	BlockedBalanceP2p           decimal.Decimal
	BlockedBalanceCoinhold      decimal.Decimal
	BlockedBalanceFreeWithdraw  decimal.Decimal
	BlockedBalanceMiningPayouts decimal.Decimal
	CoinID                      string
}
