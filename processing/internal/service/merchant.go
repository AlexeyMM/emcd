package service

import (
	"context"

	"code.emcdtech.com/b2b/processing/model"
)

type MerchantInvoiceService interface {
	CreateInvoice(ctx context.Context, req *model.CreateInvoiceRequest) (*model.Invoice, error)
	CreateInvoiceForm(ctx context.Context, form *model.InvoiceForm) (*model.InvoiceForm, error)
}
