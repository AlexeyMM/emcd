package service

import (
	"context"

	"code.emcdtech.com/b2b/processing/model"
)

type Fee interface {
	ChargeFeeForInvoice(ctx context.Context, invoice *model.Invoice) error
}
