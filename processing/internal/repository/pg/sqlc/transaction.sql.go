// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: transaction.sql

package sqlc

import (
	"context"

	uuid "github.com/google/uuid"
	decimal "github.com/shopspring/decimal"
)

const CreateTransaction = `-- name: CreateTransaction :exec
INSERT INTO invoice_transaction (
    invoice_id,
    hash,
    amount,
    received_address,
    confirmation_status
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
`

type CreateTransactionParams struct {
	InvoiceID          uuid.UUID                     `db:"invoice_id"`
	Hash               string                        `db:"hash"`
	Amount             decimal.Decimal               `db:"amount"`
	ReceivedAddress    string                        `db:"received_address"`
	ConfirmationStatus TransactionConfirmationStatus `db:"confirmation_status"`
}

func (q *Queries) CreateTransaction(ctx context.Context, arg *CreateTransactionParams) error {
	_, err := q.db.Exec(ctx, CreateTransaction,
		arg.InvoiceID,
		arg.Hash,
		arg.Amount,
		arg.ReceivedAddress,
		arg.ConfirmationStatus,
	)
	return err
}

const GetInvoiceTransactions = `-- name: GetInvoiceTransactions :many
SELECT hash, invoice_id, amount, received_address, confirmation_status, created_at
FROM invoice_transaction
WHERE invoice_id = $1
ORDER BY created_at
`

func (q *Queries) GetInvoiceTransactions(ctx context.Context, invoiceID uuid.UUID) ([]InvoiceTransaction, error) {
	rows, err := q.db.Query(ctx, GetInvoiceTransactions, invoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []InvoiceTransaction
	for rows.Next() {
		var i InvoiceTransaction
		if err := rows.Scan(
			&i.Hash,
			&i.InvoiceID,
			&i.Amount,
			&i.ReceivedAddress,
			&i.ConfirmationStatus,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const SaveTransaction = `-- name: SaveTransaction :exec
INSERT INTO invoice_transaction (
    invoice_id,
    hash,
    amount,
    received_address,
    confirmation_status
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
) ON CONFLICT (hash) DO UPDATE SET
    invoice_id = EXCLUDED.invoice_id,
    amount = EXCLUDED.amount,
    received_address = EXCLUDED.received_address,
    confirmation_status = EXCLUDED.confirmation_status
`

type SaveTransactionParams struct {
	InvoiceID          uuid.UUID                     `db:"invoice_id"`
	Hash               string                        `db:"hash"`
	Amount             decimal.Decimal               `db:"amount"`
	ReceivedAddress    string                        `db:"received_address"`
	ConfirmationStatus TransactionConfirmationStatus `db:"confirmation_status"`
}

func (q *Queries) SaveTransaction(ctx context.Context, arg *SaveTransactionParams) error {
	_, err := q.db.Exec(ctx, SaveTransaction,
		arg.InvoiceID,
		arg.Hash,
		arg.Amount,
		arg.ReceivedAddress,
		arg.ConfirmationStatus,
	)
	return err
}
