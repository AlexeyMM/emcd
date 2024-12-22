package repository

import (
	"context"

	transactor "code.emcdtech.com/emcd/sdk/pg"
	"github.com/google/uuid"

	"code.emcdtech.com/b2b/processing/model"
)

type Invoice interface {
	transactor.PgxTransactor
	CreateInvoice(ctx context.Context, invoice *model.Invoice) error
	CreateInvoiceForm(ctx context.Context, form *model.InvoiceForm) error
	GetInvoiceForm(ctx context.Context, id uuid.UUID) (*model.InvoiceForm, error)
	GetInvoice(ctx context.Context, id uuid.UUID) (*model.Invoice, error)
	GetActiveInvoiceByDepositAddressForUpdate(ctx context.Context, address string) (*model.Invoice, error)
	// UpdateStatus do it only with metrics
	// s.invoiceExecutionHistogram.WithLabelValues(string(newStatus)).
	// Observe(time.Now().UTC().Sub(invoice.CreatedAt).Seconds())
	UpdateStatus(ctx context.Context, id uuid.UUID, status model.InvoiceStatus) error
	SetInvoicesExpired(ctx context.Context) error
	CountInvoiceByStatus(ctx context.Context) (map[model.InvoiceStatus]int, error)
}
