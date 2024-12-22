package metrics

import (
	"context"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"

	"code.emcdtech.com/b2b/processing/internal/repository"
	"code.emcdtech.com/b2b/processing/model"
)

type Transaction struct {
	repo      repository.Transaction
	histogram *prometheus.HistogramVec
}

func NewTransaction(repo repository.Transaction, histogram *prometheus.HistogramVec) *Transaction {
	return &Transaction{
		repo:      repo,
		histogram: histogram,
	}
}

func (t *Transaction) SaveTransaction(ctx context.Context, tx *model.Transaction) error {
	start := time.Now()
	err := t.repo.SaveTransaction(ctx, tx)
	duration := time.Since(start).Seconds()

	t.histogram.WithLabelValues("transaction.SaveTransaction", strconv.FormatBool(err == nil)).Observe(duration)

	return err
}

func (t *Transaction) GetInvoiceTransactions(ctx context.Context, invoiceID uuid.UUID) ([]*model.Transaction, error) {
	start := time.Now()
	tr, err := t.repo.GetInvoiceTransactions(ctx, invoiceID)
	duration := time.Since(start).Seconds()

	t.histogram.WithLabelValues("transaction.GetInvoiceTransactions", strconv.FormatBool(err == nil)).Observe(duration)

	return tr, err
}
