package pg

import (
	"context"
	"fmt"

	"code.emcdtech.com/b2b/processing/internal/repository/pg/sqlc"
	"code.emcdtech.com/b2b/processing/model"
	transactor "code.emcdtech.com/emcd/sdk/pg"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Transaction struct {
	transactor.PgxTransactor
}

func NewTransaction(pool *pgxpool.Pool) *Transaction {
	return &Transaction{PgxTransactor: transactor.NewPgxTransactor(pool)}
}

func (r *Transaction) SaveTransaction(ctx context.Context, tx *model.Transaction) error {
	q := sqlc.New(r.Runner(ctx))

	confirmationStatus := sqlc.TransactionConfirmationStatusConfirming
	if tx.IsConfirmed {
		confirmationStatus = sqlc.TransactionConfirmationStatusConfirmed
	}

	err := q.SaveTransaction(ctx, &sqlc.SaveTransactionParams{
		Hash:               tx.Hash,
		InvoiceID:          tx.InvoiceID,
		Amount:             tx.Amount,
		ReceivedAddress:    tx.Address,
		ConfirmationStatus: confirmationStatus,
	})
	if err != nil {
		return fmt.Errorf("updateTransaction: %w", err)
	}

	return nil
}

func (r *Transaction) GetInvoiceTransactions(ctx context.Context, invoiceID uuid.UUID) ([]*model.Transaction, error) {
	q := sqlc.New(r.Runner(ctx))

	txs, err := q.GetInvoiceTransactions(ctx, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("getInvoiceTransactions: %w", err)
	}

	result := make([]*model.Transaction, len(txs))
	for i, tx := range txs {
		result[i] = &model.Transaction{
			InvoiceID:   tx.InvoiceID,
			Hash:        tx.Hash,
			Amount:      tx.Amount,
			Address:     tx.ReceivedAddress,
			IsConfirmed: tx.ConfirmationStatus == sqlc.TransactionConfirmationStatusConfirmed,
			CreatedAt:   tx.CreatedAt,
		}
	}

	return result, nil
}
