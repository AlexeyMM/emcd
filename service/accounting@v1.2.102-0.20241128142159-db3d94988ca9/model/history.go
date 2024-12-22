package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type (
	Income struct {
		Diff          int             `json:"difficulty,omitempty"`     // Diff only for BTC, LTC
		ChangePercent float64         `json:"change_percent,omitempty"` // ChangePercent only for BTC, LTC
		Time          float64         `json:"time"`
		Income        decimal.Decimal `json:"income"`
		Code          int             `json:"code"`
		HashRate      *int            `json:"hashrate"`
	}

	Payout struct {
		Time   float64 `json:"time"`
		Amount float64 `json:"amount"`
		Tx     *string `json:"tx"`
		TxId   *string `json:"txId"`
	}

	Wallet struct {
		TxID                  *string   `json:"tx_id" db:"tx_id"`
		FiatStatus            *string   `json:"fiat_status" db:"fiat_status"`
		Address               *string   `json:"address" db:"address"`
		Comment               *string   `json:"comment" db:"comment"`
		CoinholdType          string    `json:"coinhold_type" db:"coinhold_type"`
		ExchangeToCoinID      *int      `json:"exchange_to_coin_id" db:"exchange_to_coin_id"`
		CoinholdID            *int      `json:"coinhold_id" db:"coinhold_id"`
		OrderID               *int      `json:"order_id" db:"order_id"`
		CreatedAt             *int      `json:"created_at" db:"created_at"`
		Amount                *float64  `json:"amount" db:"amount"`
		Fee                   *float64  `json:"fee" db:"fee"`
		FiatAmount            *float64  `json:"fiat_amount" db:"fiat_amount"`
		ExchangeAmountReceive *float64  `json:"exchange_amount_receive" db:"exchange_amount_receive"`
		ExchangeAmountSent    *float64  `json:"exchange_amount_sent" db:"exchange_amount_sent"`
		ExchangeRate          *float64  `json:"exchange_rate" db:"exchange_rate"`
		ExchangeIsSuccess     *bool     `json:"exchange_is_success" db:"exchange_is_success"`
		Date                  time.Time `json:"date" db:"date"`
		TokenID               int       `json:"token_id" db:"token_id"`
		CoinID                int       `json:"coin_id" db:"coin_id"`
		Status                int       `json:"status" db:"status"`
		Type                  int       `json:"type" db:"type"`
		ID                    int       `json:"id" db:"id"`
		P2PStatus             int       `json:"-" db:"p2p_status"`
		P2POrderID            int       `json:"-" db:"p2p_order_id"`
		ReferralEmail         *string   `json:"referral_email"`
		ReferralType          *int      `json:"referral_type"`
		NetworkID             string    `json:"network_id"`
		CoinStrID             string    `json:"coin_str_id"`
	}
)

type HistoryInput struct {
	Type                HistoryType `json:"type"`
	CoinCode            string      `json:"coin_code"`
	From                string      `json:"from"`
	To                  string      `json:"to"`
	Limit               int32       `json:"limit"`
	Offset              int32       `json:"offset"`
	CoinholdID          int64       `json:"coinhold_id"`
	UserID              int64       `json:"-"`
	TransactionTypesIDs []int64     `json:"transaction_types_ids"`
	AccountTypeIDs      []int64     `json:"account_type_ids"`
	CoinsIDs            []int64     `json:"coins_ids"`
}

type HistoryOutput struct {
	TotalCount    int              `json:"total_count"`
	IncomesSum    *decimal.Decimal `json:"incomes_sum,omitempty"`
	PayoutsSum    *decimal.Decimal `json:"payouts_sum,omitempty"`
	HasNewIncome  *bool            `json:"hasNewIncome,omitempty"`
	HasNewPayouts *bool            `json:"hasNewPayouts,omitempty"`
	Incomes       []*Income        `json:"incomes"`
	Payouts       []*Payout        `json:"payouts"`
	Wallets       []*Wallet        `json:"wallets"`
}

type HistoryType string

const (
	HistoryIncome   HistoryType = "income"
	HistoryPayout   HistoryType = "payout"
	HistoryWallet   HistoryType = "wallet"
	HistoryCoinhold HistoryType = "coinhold"
)
