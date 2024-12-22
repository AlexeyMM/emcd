package buyer

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"code.emcdtech.com/b2b/processing/model"
	"code.emcdtech.com/b2b/processing/protocol/buyerpb"
)

func (c *Controller) CalculateInvoicePayment(ctx context.Context,
	request *buyerpb.CalculateInvoicePaymentRequest,
) (*buyerpb.CalculateInvoicePaymentResponse, error) {
	merchantID, err := uuid.Parse(request.GetMerchantId())
	if err != nil {
		return nil, &model.Error{
			Code:    model.ErrorCodeInvalidArgument,
			Message: fmt.Sprintf("invalid merchantId: %q", request.GetMerchantId()),
		}
	}

	paymentAmount, err := decimal.NewFromString(request.GetRawAmount())
	if err != nil {
		return nil, &model.Error{
			Code:    model.ErrorCodeInvalidArgument,
			Message: fmt.Sprintf("invalid paymentAmount: %q", request.GetRawAmount()),
		}
	}

	amount, fee, err := c.buyerInvoiceService.CalculateInvoicePayment(ctx, merchantID, paymentAmount)
	if err != nil {
		return nil, err
	}

	return &buyerpb.CalculateInvoicePaymentResponse{
		Amount:   amount.String(),
		BuyerFee: fee.String(),
	}, nil
}
