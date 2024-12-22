package model

import "time"

type Block struct {
	ID                   int64
	BlockTransactionID   int64
	UnblockTransactionID int64
	UnblockToAccountID   int64
	BlockedTill          time.Time
}
