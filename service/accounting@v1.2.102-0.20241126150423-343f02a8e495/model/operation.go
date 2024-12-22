package model

import (
	"code.emcdtech.com/emcd/service/accounting/model/enum"
	"github.com/google/uuid"
	"time"

	"github.com/shopspring/decimal"
)

type Operation struct {
	ID            int
	Type          TransactionType
	TransactionID int64
	AccountID     int64
	CoinID        int64
	Amount        decimal.Decimal
	CreatedAt     time.Time
	Hash          string
	TokenID       int
}

type OperationSelection struct {
	Amount            decimal.Decimal
	AccountID         int64
	OperationCoinID   int64
	UserAccountCoinID int64
	AccountTypeID     enum.AccountTypeId
	IsActive          bool
	Type              TransactionType
	CreatedAt         time.Time
	TransactionID     int64
}

type OperationSelectionWithBlock struct {
	Amount               decimal.Decimal
	AccountID            int64
	UserAccountCoinID    int64
	AccountTypeID        enum.AccountTypeId
	IsActive             bool
	UnblockTransactionID int64
	UnblockToAccountID   int64
	OperationCoinID      int64
	Type                 TransactionType
	CreatedAt            time.Time
	TransactionID        int64
}

type OperationWithTransaction struct {
	Id                   int64
	AccountID            int64
	CoinID               int64
	TokenID              int64
	Amount               decimal.Decimal
	Type                 TransactionType
	TransactionID        int64
	ActionID             string
	Comment              string
	Fee                  decimal.Decimal
	FromReferralId       int64
	GasPrice             decimal.Decimal
	Hash                 string
	Hashrate             int64
	ReceiverAccountID    int64
	ReceiverAddress      string
	SenderAccountID      int64
	TransactionBlockID   int64
	BlockedTill          time.Time
	UnblockToAccountId   int64
	UnblockTransactionId int64
	CreatedAt            time.Time
}

type OperationWithTransactionQuery struct {
	UserID               int64
	CoinID               int64
	TokenID              int64
	ActionID             *uuid.UUID
	AccountType          int32
	OperationTypes       []int32
	DateFrom             string
	DateTo               string
	Amount               float64
	Hash                 string
	ReceiverAccountID    int64
	ReceiverAddress      string
	SenderAccountID      int64
	TransactionBlockID   int64
	UnblockToAccountId   int64
	UnblockTransactionId int64
	FromReferralId       int64
	Limit                int32
	Offset               int32
	SortField            string
	Asc                  bool
}
