package transaction

import (
	"context"
	"errors"
	"fmt"
	"time"

	"code.emcdtech.com/emcd/sdk/log"
	"github.com/prometheus/client_golang/prometheus"

	"code.emcdtech.com/b2b/processing/internal/repository"
	"code.emcdtech.com/b2b/processing/internal/service"
	"code.emcdtech.com/b2b/processing/model"
)

type Service struct {
	transactionRepo           repository.Transaction
	invoiceRepo               repository.Invoice
	feeService                service.Fee
	invoiceExecutionHistogram *prometheus.HistogramVec
}

func NewService(
	transactionRepo repository.Transaction,
	invoiceRepo repository.Invoice,
	feeService service.Fee,
	invoiceExecutionHistogram *prometheus.HistogramVec,
) *Service {
	return &Service{
		transactionRepo:           transactionRepo,
		invoiceRepo:               invoiceRepo,
		feeService:                feeService,
		invoiceExecutionHistogram: invoiceExecutionHistogram,
	}
}

func (s *Service) ProcessTransaction(ctx context.Context, transaction *model.Transaction) error {
	var invoice *model.Invoice

	err := s.invoiceRepo.WithinTransaction(ctx, func(ctx context.Context) error {
		var err error

		invoice, err = s.invoiceRepo.GetActiveInvoiceByDepositAddressForUpdate(ctx, transaction.Address)
		if err != nil {
			return fmt.Errorf("GetActiveInvoiceByDepositAddressForUpdate: %w", err)
		}

		if transaction.CoinID != invoice.CoinID {
			log.SWarn(ctx, "transaction coin ID does not match invoice coin ID", map[string]any{
				"transaction": transaction,
				"invoice":     invoice,
			})

			return nil
		}

		transaction.InvoiceID = invoice.ID
		if err := s.transactionRepo.SaveTransaction(ctx, transaction); err != nil {
			return fmt.Errorf("updateTransaction: %w", err)
		}

		if transaction.IsConfirmed {
			invoice.PaidAmount = invoice.PaidAmount.Add(transaction.Amount)
		}

		newStatus := invoice.Status

		switch {
		case !transaction.IsConfirmed && invoice.Status == model.InvoiceStatusWaitingForDeposit:
			newStatus = model.InvoiceStatusPaymentConfirmation
		case transaction.IsConfirmed && invoice.PaidAmount.GreaterThanOrEqual(invoice.RequiredPayment):
			newStatus = model.InvoiceStatusPaymentAccepted
		case transaction.IsConfirmed && invoice.PaidAmount.LessThan(invoice.RequiredPayment):
			newStatus = model.InvoiceStatusPartiallyPaid
		}

		if newStatus != invoice.Status {
			invoice.Status = newStatus

			err = s.updateStatus(ctx, invoice, newStatus)
			if err != nil {
				return fmt.Errorf("updateStatus: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		if errors.Is(err, &model.Error{Code: model.ErrorCodeNoSuchInvoice}) {
			log.SWarn(ctx, "funds received on an address that isn't currently used in any invoice", map[string]any{
				"transaction": transaction,
			})

			return nil // this isn't an error for us, just a curious event
		}

		return err
	}

	if invoice.Status == model.InvoiceStatusPaymentAccepted {
		if err := s.feeService.ChargeFeeForInvoice(ctx, invoice); err != nil {
			return fmt.Errorf("chargeFeeForInvoice: %w", err)
		}

		// Note: hmm, do we need a worker to look for stale "PaymentAccepted" invoices? Or will message processing retry cover it for good?

		if err = s.updateStatus(ctx, invoice, model.InvoiceStatusFinished); err != nil {
			return fmt.Errorf("updateStatus to finsihed : %w", err)
		}
	}

	return nil
}

func (s *Service) updateStatus(ctx context.Context, invoice *model.Invoice, newStatus model.InvoiceStatus) error {
	if err := s.invoiceRepo.UpdateStatus(ctx, invoice.ID, newStatus); err != nil {
		return fmt.Errorf("updateStatus: %w", err)
	}

	s.invoiceExecutionHistogram.WithLabelValues(string(newStatus)).
		Observe(time.Now().UTC().Sub(invoice.CreatedAt).Seconds())

	return nil
}
