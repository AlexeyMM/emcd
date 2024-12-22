package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type OutboxTransactions interface {
	List(ctx context.Context, tx pgx.Tx, limit uint) ([]int64, error)
	Delete(ctx context.Context, tx pgx.Tx, ids ...int64) error
}

type outboxTransactions struct {
}

func NewOutboxTransactions() OutboxTransactions {
	return &outboxTransactions{}
}

func (r *outboxTransactions) List(ctx context.Context, tx pgx.Tx, limit uint) ([]int64, error) {
	const selectListOutboxTransactionsSQL = `
select transaction_id from outbox_transactions order by transaction_id limit $1 for update;
`
	rows, err := tx.Query(ctx, selectListOutboxTransactionsSQL, limit)
	if err != nil {
		return nil, fmt.Errorf("execute selectListSQL: %w", err)
	}
	defer rows.Close()

	ids := make([]int64, 0, limit)
	for rows.Next() {
		var id int64
		err = rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (r *outboxTransactions) Delete(ctx context.Context, tx pgx.Tx, ids ...int64) error {
	const deleteOutboxTransactionsSQL = `
delete from outbox_transactions where transaction_id = ANY($1);
`
	_, err := tx.Exec(ctx, deleteOutboxTransactionsSQL, ids)
	if err != nil {
		return fmt.Errorf("execute deleteOutboxTransactionsSQL: %w", err)
	}
	return nil
}
