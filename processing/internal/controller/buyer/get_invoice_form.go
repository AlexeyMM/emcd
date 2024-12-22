package buyer

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"code.emcdtech.com/b2b/processing/model"
	"code.emcdtech.com/b2b/processing/protocol/buyerpb"
)

func (c *Controller) GetInvoiceForm(
	ctx context.Context,
	request *buyerpb.GetInvoiceFormRequest,
) (*buyerpb.InvoiceForm, error) {
	formID := uuid.MustParse(request.GetId())

	form, err := c.buyerInvoiceService.GetInvoiceForm(ctx, formID)
	if err != nil {
		return nil, err
	}

	return convertModelInvoiceFormToProto(form), nil
}

func convertModelInvoiceFormToProto(form *model.InvoiceForm) *buyerpb.InvoiceForm {
	result := &buyerpb.InvoiceForm{
		Id:         form.ID.String(),
		MerchantId: form.MerchantID.String(),
	}

	if form.Title != nil {
		result.Title = form.Title
	}

	if form.Description != nil {
		result.Description = form.Description
	}

	if form.CoinID != nil {
		result.CoinId = form.CoinID
	}

	if form.NetworkID != nil {
		result.NetworkId = form.NetworkID
	}

	if form.Amount != nil {
		amount := form.Amount.String()
		result.Amount = &amount
	}

	if form.BuyerEmail != nil {
		result.BuyerEmail = form.BuyerEmail
	}

	if form.ExpiresAt != nil {
		result.ExpiresAt = timestamppb.New(*form.ExpiresAt)
	}

	return result
}
