package model

import "github.com/shopspring/decimal"

const (
	DefaultAccuracy         = 4
	DefaultWithdrawAccuracy = 8
)

type Coin struct {
	Title    string
	Accuracy int
	Info     *CoinInfo
	Networks []*Network
}

type CoinInfo struct {
	Rating  int
	IconURL string
}

type Network struct {
	Title                      string
	WithdrawFee                *WithdrawFee
	WithdrawMin                decimal.Decimal
	DepositMin                 decimal.Decimal
	AccuracyWithdrawAndDeposit int
	WithdrawSupported          bool
}

type WithdrawFee struct {
	Fee           decimal.Decimal
	PercentageFee decimal.Decimal
}
