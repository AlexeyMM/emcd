package merchant

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"code.emcdtech.com/b2b/processing/internal/controller/common"
	"code.emcdtech.com/b2b/processing/model"
	"code.emcdtech.com/b2b/processing/protocol/commonpb"
	"code.emcdtech.com/b2b/processing/protocol/merchantpb"
)

func (c *Controller) CreateInvoice(
	ctx context.Context,
	request *merchantpb.CreateInvoiceRequest,
) (*commonpb.Invoice, error) {
	m, err := convertProtoCreateInvoiceRequestToModel(request)
	if err != nil {
		return nil, err
	}

	invoice, err := c.merchantInvoiceService.CreateInvoice(ctx, m)
	if err != nil {
		return nil, err
	}

	return common.ConvertModelInvoiceToProto(invoice), nil
}

func convertProtoCreateInvoiceRequestToModel(
	request *merchantpb.CreateInvoiceRequest,
) (*model.CreateInvoiceRequest, error) {
	amount, err := decimal.NewFromString(request.GetAmount())
	if err != nil {
		return nil, &model.Error{
			Code:    model.ErrorCodeInvalidArgument,
			Message: fmt.Sprintf("invalid amount: %q", request.GetAmount()),
			Inner:   err,
		}
	}

	merchantID, err := uuid.Parse(request.GetMerchantId())
	if err != nil {
		return nil, &model.Error{
			Code:    model.ErrorCodeInvalidArgument,
			Message: fmt.Sprintf("invalid merchant id: %q", request.GetMerchantId()),
			Inner:   err,
		}
	}

	return &model.CreateInvoiceRequest{
		ExternalID:  request.GetExternalId(),
		Title:       request.GetTitle(),
		Description: request.GetDescription(),
		CoinID:      request.GetCoinId(),
		NetworkID:   request.GetNetworkId(),
		Amount:      amount,
		BuyerEmail:  request.GetBuyerEmail(),
		CheckoutURL: request.GetCheckoutUrl(),
		MerchantID:  merchantID,
		ExpiresAt:   request.GetExpiresAt().AsTime(),
	}, nil
}
