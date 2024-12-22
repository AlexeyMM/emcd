package model

import (
	"github.com/shopspring/decimal"
	"time"
)

type HashrateByDate struct {
	CoinId   int64
	Hashrate string
}

type SumCheckData struct {
	CoinId int64
	Sum    decimal.Decimal
}

type TransactionOperationsIntegrityData struct {
	Count       int64
	TrId        int64
	Op2Id       int64
	Op1Id       int64
	OpPairCheck bool
	TrNegChk    bool
	OpSumChk    bool
	DiffChk     bool
	TrDateChk   bool
	CoinChk     bool
	AccChk      bool
}

type CheckTransactionCoinsData struct {
	TrIds []int64
	OpIds []int64
}

type CheckFreezePayoutsBlocksData struct {
	TrId      int64
	Type      int64
	UserId    int64
	CreatedAt time.Time
}
