package service

import (
	"context"

	"code.emcdtech.com/b2b/processing/model"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type BuyerInvoiceService interface {
	GetInvoiceForm(ctx context.Context, invoiceID uuid.UUID) (*model.InvoiceForm, error)
	SubmitInvoiceForm(ctx context.Context, invoiceForm *model.InvoiceForm) (*model.Invoice, error)
	GetInvoice(ctx context.Context, invoiceID uuid.UUID) (*model.Invoice, error)
	// CalculateInvoicePayment returns raw payment plus buyer fee, buyer fee
	CalculateInvoicePayment(ctx context.Context, merchantID uuid.UUID, rawPayment decimal.Decimal) (decimal.Decimal,
		decimal.Decimal, error)
}
