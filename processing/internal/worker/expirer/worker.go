package expirer

import (
	"context"
	"time"

	"code.emcdtech.com/b2b/processing/internal/repository"
	"code.emcdtech.com/emcd/sdk/log"
)

const tickInterval = 1 * time.Minute

type Worker struct {
	invoiceRepo repository.Invoice
}

func NewWorker(invoiceRepo repository.Invoice) *Worker {
	return &Worker{
		invoiceRepo: invoiceRepo,
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
			if err := w.invoiceRepo.SetInvoicesExpired(ctx); err != nil {
				log.SError(ctx, "failed to set invoices expired", map[string]any{"error": err})

				continue
			}
		}
	}
}
