package invoicestatusmeter

import (
	"context"
	"time"

	"code.emcdtech.com/emcd/sdk/log"
	"github.com/prometheus/client_golang/prometheus"

	"code.emcdtech.com/b2b/processing/internal/repository"
)

const tickInterval = 1 * time.Minute

type Worker struct {
	invoiceRepo          repository.Invoice
	invoiceStatusesGauge *prometheus.GaugeVec
}

func NewWorker(invoiceRepo repository.Invoice, invoiceStatusesGauge *prometheus.GaugeVec) *Worker {
	return &Worker{
		invoiceRepo:          invoiceRepo,
		invoiceStatusesGauge: invoiceStatusesGauge,
	}
}

func (w *Worker) Run(ctx context.Context) error {
	ticker := time.NewTicker(tickInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			statuses, err := w.invoiceRepo.CountInvoiceByStatus(ctx)
			if err != nil {
				log.SError(ctx, "failed to set invoices expired", map[string]any{"error": err})

				continue
			}

			for status, count := range statuses {
				w.invoiceStatusesGauge.WithLabelValues(string(status)).Set(float64(count))
			}
		}
	}
}
