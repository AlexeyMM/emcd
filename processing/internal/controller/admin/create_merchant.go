package admin

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"code.emcdtech.com/b2b/processing/model"
	"code.emcdtech.com/b2b/processing/protocol/adminpb"
)

func (c *Controller) CreateMerchant(
	ctx context.Context,
	request *adminpb.CreateMerchantRequest,
) (*adminpb.CreateMerchantResponse, error) {
	m, err := convertProtoCreateMerchantRequestToModel(request)
	if err != nil {
		return nil, err
	}

	if err := c.merchantAdminService.CreateMerchant(ctx, m); err != nil {
		return nil, err
	}

	return &adminpb.CreateMerchantResponse{}, nil
}

func convertProtoCreateMerchantRequestToModel(request *adminpb.CreateMerchantRequest) (*model.Merchant, error) {
	id := uuid.MustParse(request.GetUserId())
	upperFee := decimal.RequireFromString(request.Tariff.GetUpperFee())
	lowerFee := decimal.RequireFromString(request.Tariff.GetLowerFee())

	minPay := decimal.RequireFromString(request.Tariff.GetMinPay())
	maxPay := decimal.RequireFromString(request.Tariff.GetMaxPay())

	if maxPay.LessThan(minPay) {
		return nil, &model.Error{
			Code:    model.ErrorCodeInvalidArgument,
			Message: "max pay must be greater than min pay",
		}
	}

	return &model.Merchant{
		ID: id,
		Tariff: &model.Tariff{
			UpperFee: upperFee,
			LowerFee: lowerFee,
			MinPay:   minPay,
			MaxPay:   maxPay,
		},
	}, nil
}
