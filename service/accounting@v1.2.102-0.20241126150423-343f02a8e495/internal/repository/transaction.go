package repository

import (
	"code.emcdtech.com/emcd/service/accounting/model"
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"time"
)

const listTransactionsSQL = `
SELECT t.id,
	   t.type,
	   t.created_at,
	   t.sender_account_id,
	   t.receiver_account_id,
	   t.coin_id,
	   t.token_id,
	   t.amount,
	   t.comment,
	   t.hash,
	   t.receiver_address,
	   t.hashrate,
	   t.from_referral_id,
	   t.action_id
FROM emcd.transactions t
WHERE (  (t.receiver_account_id = ANY($1) and $2=true) or $2=false  )
  	  and (  (t.type = ANY ($3) and $4=true) or $4=false  ) 
  	  and t.created_at between $5 and $6
	  and t.id > $7
ORDER BY t.id
LIMIT $8;
`

type Transaction interface {
	ListTransactions(
		ctx context.Context,
		receiverAccountIds []int,
		types []int,
		from, to time.Time,
		limit int,
		fromTransactionID int64,
	) ([]*model.Transaction, int64, error)
}

type TransactionStore struct {
	pool *pgxpool.Pool
}

func NewTransactionStore(
	pool *pgxpool.Pool,
) *TransactionStore {
	return &TransactionStore{
		pool: pool,
	}
}

func (s *TransactionStore) ListTransactions(
	ctx context.Context,
	receiverAccountIDs []int,
	types []int,
	from, to time.Time,
	limit int,
	fromTransactionID int64,
) ([]*model.Transaction, int64, error) {
	rows, err := s.pool.Query(ctx, listTransactionsSQL,
		receiverAccountIDs,
		len(receiverAccountIDs) > 0,
		types,
		len(types) > 0,
		from,
		to,
		fromTransactionID,
		limit,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("r.pool.Query: %w", err)
	}
	defer rows.Close()

	var (
		result   []*model.Transaction
		maxTrxId int64
	)

	for rows.Next() {
		t, err := s.rowToTransaction(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("rowToTransaction: %w", err)
		}
		if maxTrxId < t.ID {
			maxTrxId = t.ID
		}
		result = append(result, t)
	}

	return result, maxTrxId, nil
}

func (s *TransactionStore) rowToTransaction(rows pgx.Rows) (*model.Transaction, error) {
	var (
		comment           sql.NullString
		hash              sql.NullString
		receiverAddress   sql.NullString
		receiverAccountID sql.NullInt64
		hashrate          sql.NullInt64
		tokenID           sql.NullInt64
		fromReferralId    sql.NullInt64
		amount            decimal.NullDecimal
		actionID          sql.NullString
	)

	t := new(model.Transaction)

	err := rows.Scan(&t.ID, &t.Type, &t.CreatedAt, &t.SenderAccountID, &receiverAccountID,
		&t.CoinID, &tokenID, &amount, &comment, &hash, &receiverAddress, &hashrate, &fromReferralId, &actionID)
	if err != nil {
		return nil, fmt.Errorf("rows.Scan: %w", err)
	}
	t.Comment = comment.String
	t.Hash = hash.String
	t.Hashrate = hashrate.Int64
	t.ReceiverAddress = receiverAddress.String
	t.ReceiverAccountID = receiverAccountID.Int64
	t.TokenID = tokenID.Int64
	t.Amount = amount.Decimal
	t.FromReferralId = fromReferralId.Int64
	t.ActionID = actionID.String

	return t, nil
}
