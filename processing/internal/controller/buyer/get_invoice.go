package buyer

import (
	"context"

	"github.com/google/uuid"

	"code.emcdtech.com/b2b/processing/internal/controller/common"
	"code.emcdtech.com/b2b/processing/protocol/buyerpb"
	"code.emcdtech.com/b2b/processing/protocol/commonpb"
)

func (c *Controller) GetInvoice(ctx context.Context, request *buyerpb.GetInvoiceRequest) (*commonpb.Invoice, error) {
	invoiceID := uuid.MustParse(request.GetId())

	invoice, err := c.buyerInvoiceService.GetInvoice(ctx, invoiceID)
	if err != nil {
		return nil, err
	}

	return common.ConvertModelInvoiceToProto(invoice), nil
}
