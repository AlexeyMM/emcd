package model

import (
	"github.com/shopspring/decimal"
	"time"
)

type PayoutForBlock struct {
	AccountID int64
	UserID    int64
	Balance   decimal.Decimal
	Address   string
}

type PayoutBlockTransaction struct {
	ID      int64
	Balance decimal.Decimal
}

type FreePayout struct {
	AccountId         int64
	UserId            int64
	Username          string
	ID                int64
	ActionID          string
	Amount            decimal.Decimal
	CoinID            int64
	Comment           string
	CreatedAt         time.Time
	Fee               decimal.Decimal
	FromReferralID    int64
	GasPrice          decimal.Decimal
	Hash              string
	Hashrate          int64
	IsViewer          bool
	ReceiverAccountID int64
	ReceiverAddress   string
	SenderAccountID   int64
	TokenID           int64
	Type              int64
}

type PayoutWithCalculation struct {
	ID          int64
	AccountID2  int64
	UserID      int64
	RefID       int64
	CoinID      int64
	Username    string
	Minpay      decimal.Decimal
	MasterID    int64
	Address     string
	Balance     decimal.Decimal
	BlockID     int64
	BlockCreate time.Time
	Calc        PayoutCalculationData
}

type PayoutCalculationData struct {
	Coinhold    decimal.Decimal
	Incomes     decimal.Decimal
	Hashrate    decimal.Decimal
	FeeAndMore  decimal.Decimal
	Ref         decimal.Decimal
	Other       decimal.Decimal
	Types       string
	AccountID   int64
	LastPay     time.Time
	IncomeFirst time.Time
	IncomeLast  time.Time
}

type PayoutBlockStatus struct {
	ToAccountId          int64
	Type                 int64
	ReceiverAddress      string
	UnblockTransactionId int64
	Amount               decimal.Decimal
}

type IncomeWithFee struct {
	TransactionId int64
	Amount        decimal.Decimal
	Hashrate      int64
	CreatedAt     time.Time
	Fee           decimal.Decimal
}

type CheckIncomeOperationsQuery struct {
	CreatedAt time.Time
	Coin      string
	UserID    int64
	AccountID int64
	LastPayAt time.Time
}

type AveragePaidQuery struct {
	CoinID            int64
	Days              int64
	TransactionTypeID int64
	AccountTypeID     int64
	Username          string
}

type OtherOperationsWithTransaction struct {
	TransactionID     int64
	SenderID          int64
	ReceiverID        int64
	Hash              string
	OperationID       int64
	SenderUserID      int64
	ReceiverUserID    int64
	Amount            decimal.Decimal
	TransactionTypeID int64
	CreatedAt         time.Time
	Comment           string
}

type CheckOtherQuery struct {
	AccountID      int64
	Types          []int64
	LastPayAt      time.Time
	BlockCreatedAt *time.Time
}

type ServiceUserBlock struct {
	Address     string
	SuAccountID int64
	UserID      int64
	Username    string
	Amount      decimal.Decimal
	BlockID     int64
}
