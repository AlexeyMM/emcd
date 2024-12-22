package merchant

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/types/known/timestamppb"

	"code.emcdtech.com/b2b/processing/model"
	"code.emcdtech.com/b2b/processing/pkg/gokit"
	"code.emcdtech.com/b2b/processing/protocol/merchantpb"
)

const externalIDMaxLen = 36

func (c *Controller) CreateInvoiceForm(
	ctx context.Context,
	request *merchantpb.CreateInvoiceFormRequest,
) (*merchantpb.InvoiceForm, error) {
	form, err := convertProtoInvoiceFormToModel(request)
	if err != nil {
		return nil, err
	}

	createdForm, err := c.merchantInvoiceService.CreateInvoiceForm(ctx, form)
	if err != nil {
		return nil, err
	}

	return convertModelInvoiceFormToProto(createdForm), nil
}

func convertProtoInvoiceFormToModel(req *merchantpb.CreateInvoiceFormRequest) (*model.InvoiceForm, error) {
	merchantID, err := uuid.Parse(req.GetMerchantId())
	if err != nil {
		return nil, &model.Error{
			Code:    model.ErrorCodeInvalidArgument,
			Message: fmt.Sprintf("invalid merchant id: %q", req.GetMerchantId()),
			Inner:   err,
		}
	}

	// TODO: check with protovalidate
	if req.ExpiresAt != nil && req.ExternalId == nil {
		return nil, &model.Error{
			Code:    model.ErrorCodeInvalidArgument,
			Message: "expires_at can only be specified together with external_id",
		}
	}

	form := &model.InvoiceForm{
		ID:          uuid.New(),
		MerchantID:  merchantID,
		CheckoutURL: req.GetCheckoutUrl(),
	}

	if req.Title != nil {
		form.Title = req.Title
	}

	if req.Description != nil {
		form.Description = req.Description
	}

	if req.CoinId != nil {
		form.CoinID = req.CoinId
	}

	if req.NetworkId != nil {
		form.NetworkID = req.NetworkId
	}

	if req.Amount != nil {
		amount, err := decimal.NewFromString(*req.Amount)
		if err != nil {
			return nil, &model.Error{
				Code:    model.ErrorCodeInvalidArgument,
				Message: "invalid amount format",
				Inner:   err,
			}
		}

		form.Amount = &amount
	}

	if req.BuyerEmail != nil {
		form.BuyerEmail = req.BuyerEmail
	}

	if req.ExternalId != nil {
		// TODO: use protovalidate for this
		if len(*req.ExternalId) > externalIDMaxLen {
			return nil, &model.Error{
				Code:    model.ErrorCodeInvalidArgument,
				Message: "external_id too long, max length is 36",
			}
		}

		form.ExternalID = req.ExternalId
	}

	if req.ExpiresAt != nil {
		form.ExpiresAt = gokit.Ptr(req.ExpiresAt.AsTime())
	}

	return form, nil
}

func convertModelInvoiceFormToProto(form *model.InvoiceForm) *merchantpb.InvoiceForm {
	result := &merchantpb.InvoiceForm{
		Id:          form.ID.String(),
		MerchantId:  form.MerchantID.String(),
		CheckoutUrl: form.CheckoutURL,
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
		result.Amount = gokit.Ptr(form.Amount.String())
	}

	if form.BuyerEmail != nil {
		result.BuyerEmail = form.BuyerEmail
	}

	if form.ExpiresAt != nil {
		result.ExpiresAt = timestamppb.New(*form.ExpiresAt)
	}

	return result
}
