package buyer

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"code.emcdtech.com/b2b/processing/internal/controller/common"
	"code.emcdtech.com/b2b/processing/model"
	"code.emcdtech.com/b2b/processing/pkg/gokit"
	"code.emcdtech.com/b2b/processing/protocol/buyerpb"
	"code.emcdtech.com/b2b/processing/protocol/commonpb"
)

func (c *Controller) SubmitInvoiceForm(
	ctx context.Context,
	request *buyerpb.InvoiceForm,
) (*commonpb.Invoice, error) {
	formID := uuid.MustParse(request.GetId())

	var amount *decimal.Decimal

	if request.Amount != nil {
		amount = gokit.Ptr(decimal.RequireFromString(*request.Amount))
	}

	form := &model.InvoiceForm{
		ID:          formID,
		Title:       request.Title,
		Description: request.Description,
		CoinID:      request.CoinId,
		NetworkID:   request.NetworkId,
		Amount:      amount,
		BuyerEmail:  request.BuyerEmail,
	}

	invoice, err := c.buyerInvoiceService.SubmitInvoiceForm(ctx, form)
	if err != nil {
		return nil, err
	}

	return common.ConvertModelInvoiceToProto(invoice), nil
}
