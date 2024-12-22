package model

import "time"

type Coin struct {
	ID                             int64
	Name                           string
	Description                    string
	Code                           string
	Rate                           float64
	CreatedAt                      time.Time
	UpdatedAt                      time.Time
	IsMining                       bool
	IsWallet                       bool
	IsDeposit                      bool
	IsHedge                        bool
	IsP2p                          bool
	IsWithdrawalsDisabled          bool
	WithdrawalsDisabledDescription string
	IsFreeWithdraw                 bool
	WithdrawMinLimit               float64
	WithdrawFee                    float64
	Title                          string
	PictureUrl                     string
	TokensNetworkTitle             string
}
