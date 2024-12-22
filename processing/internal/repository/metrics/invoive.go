package metrics

import (
	"context"
	"strconv"
	"time"

	transactor "code.emcdtech.com/emcd/sdk/pg"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/prometheus/client_golang/prometheus"

	"code.emcdtech.com/b2b/processing/internal/repository"
	"code.emcdtech.com/b2b/processing/model"
)

type Invoice struct {
	repo      repository.Invoice
	histogram *prometheus.HistogramVec
}

func NewInvoice(
	repo repository.Invoice,
	histogram *prometheus.HistogramVec,
) *Invoice {
	return &Invoice{
		repo:      repo,
		histogram: histogram,
	}
}

func (i *Invoice) WithinTransaction(ctx context.Context, txFn func(ctx context.Context) error) error {
	start := time.Now()
	err := i.repo.WithinTransaction(ctx, txFn)
	duration := time.Since(start).Seconds()

	i.histogram.WithLabelValues("invoice.WithinTransaction", strconv.FormatBool(err == nil)).Observe(duration)

	return err
}

func (i *Invoice) WithinTransactionWithOptions(
	ctx context.Context,
	txFn func(ctx context.Context) error,
	opts pgx.TxOptions,
) error {
	start := time.Now()
	err := i.repo.WithinTransactionWithOptions(ctx, txFn, opts)
	duration := time.Since(start).Seconds()

	i.histogram.WithLabelValues("invoice.WithinTransactionWithOptions", strconv.FormatBool(err == nil)).Observe(duration)

	return err
}

func (i *Invoice) Runner(ctx context.Context) transactor.PgxQueryRunner {
	return i.repo.Runner(ctx)
}

func (i *Invoice) CreateInvoice(ctx context.Context, invoice *model.Invoice) error {
	start := time.Now()
	err := i.repo.CreateInvoice(ctx, invoice)
	duration := time.Since(start).Seconds()

	i.histogram.WithLabelValues("invoice.CreateInvoice", strconv.FormatBool(err == nil)).Observe(duration)

	return err
}

func (i *Invoice) CreateInvoiceForm(ctx context.Context, form *model.InvoiceForm) error {
	start := time.Now()
	err := i.repo.CreateInvoiceForm(ctx, form)
	duration := time.Since(start).Seconds()

	i.histogram.WithLabelValues("invoice.CreateInvoiceForm", strconv.FormatBool(err == nil)).Observe(duration)

	return err
}

func (i *Invoice) GetInvoiceForm(ctx context.Context, id uuid.UUID) (*model.InvoiceForm, error) {
	start := time.Now()
	form, err := i.repo.GetInvoiceForm(ctx, id)
	duration := time.Since(start).Seconds()

	i.histogram.WithLabelValues("invoice.GetInvoiceForm", strconv.FormatBool(err == nil)).Observe(duration)

	return form, err
}

func (i *Invoice) GetInvoice(ctx context.Context, id uuid.UUID) (*model.Invoice, error) {
	start := time.Now()
	invoice, err := i.repo.GetInvoice(ctx, id)
	duration := time.Since(start).Seconds()

	i.histogram.WithLabelValues("invoice.GetInvoice", strconv.FormatBool(err == nil)).Observe(duration)

	return invoice, err
}

func (i *Invoice) GetActiveInvoiceByDepositAddressForUpdate(
	ctx context.Context,
	address string,
) (*model.Invoice, error) {
	start := time.Now()
	invoice, err := i.repo.GetActiveInvoiceByDepositAddressForUpdate(ctx, address)
	duration := time.Since(start).Seconds()

	i.histogram.WithLabelValues(
		"invoice.GetActiveInvoiceByDepositAddressForUpdate",
		strconv.FormatBool(err == nil)).Observe(duration)

	return invoice, err
}

func (i *Invoice) UpdateStatus(ctx context.Context, id uuid.UUID, status model.InvoiceStatus) error {
	start := time.Now()
	err := i.repo.UpdateStatus(ctx, id, status)
	duration := time.Since(start).Seconds()

	i.histogram.WithLabelValues("invoice.UpdateStatus", strconv.FormatBool(err == nil)).Observe(duration)

	return err
}

func (i *Invoice) SetInvoicesExpired(ctx context.Context) error {
	start := time.Now()
	err := i.repo.SetInvoicesExpired(ctx)
	duration := time.Since(start).Seconds()

	i.histogram.WithLabelValues("invoice.SetInvoicesExpired", strconv.FormatBool(err == nil)).Observe(duration)

	return err
}

func (i *Invoice) CountInvoiceByStatus(ctx context.Context) (map[model.InvoiceStatus]int, error) {
	start := time.Now()
	statuses, err := i.repo.CountInvoiceByStatus(ctx)
	duration := time.Since(start).Seconds()

	i.histogram.WithLabelValues("invoice.CountInvoiceByStatus", strconv.FormatBool(err == nil)).Observe(duration)

	return statuses, err
}
