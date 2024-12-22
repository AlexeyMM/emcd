package repository

import (
	"context"

	"code.emcdtech.com/b2b/processing/model"
	"github.com/google/uuid"
)

type Transaction interface {
	SaveTransaction(ctx context.Context, tx *model.Transaction) error
	GetInvoiceTransactions(ctx context.Context, invoiceID uuid.UUID) ([]*model.Transaction, error)
}
