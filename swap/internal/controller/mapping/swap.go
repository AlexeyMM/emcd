package mapping

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/types/known/timestamppb"

	"code.emcdtech.com/b2b/swap/internal/service"

	"code.emcdtech.com/b2b/swap/model"
	swapSwap "code.emcdtech.com/b2b/swap/protocol/swap"
	"code.emcdtech.com/b2b/swap/protocol/swapAdmin"
)

func MapModelEstimateRequestToProto(req *model.EstimateRequest) *swapSwap.EstimateRequest {
	return &swapSwap.EstimateRequest{
		CoinFrom:    req.CoinFrom,
		CoinTo:      req.CoinTo,
		NetworkFrom: req.NetworkFrom,
		NetworkTo:   req.NetworkTo,
		AmountFrom:  req.AmountFrom.String(),
		AmountTo:    req.AmountTo.String(),
	}
}

func MapProtoEstimateRequestToModel(req *swapSwap.EstimateRequest) (*model.EstimateRequest, error) {
	amountFrom, err := decimal.NewFromString(req.AmountFrom)
	if err != nil {
		return nil, fmt.Errorf("convert amountFrom %s to decimal: %w", req.AmountFrom, err)
	}
	amountTo, err := decimal.NewFromString(req.AmountTo)
	if err != nil {
		return nil, fmt.Errorf("convert amountTo %s to decimal: %w", req.AmountTo, err)
	}

	return &model.EstimateRequest{
		CoinFrom:    req.CoinFrom,
		CoinTo:      req.CoinTo,
		NetworkFrom: req.NetworkFrom,
		NetworkTo:   req.NetworkTo,
		AmountFrom:  amountFrom,
		AmountTo:    amountTo,
	}, nil
}

func MapProtoEstimateResponseToModel(resp *swapSwap.EstimateResponse) (amountFrom, amountTo, rate decimal.Decimal, limits *model.Limits, err error) {
	amountFrom, err = decimal.NewFromString(resp.AmountFrom)
	if err != nil {
		return decimal.Zero, decimal.Zero, decimal.Zero, nil, fmt.Errorf("convert amountFrom string to decimal; error: %w", err)
	}
	amountTo, err = decimal.NewFromString(resp.AmountTo)
	if err != nil {
		return decimal.Zero, decimal.Zero, decimal.Zero, nil, fmt.Errorf("convert amountTo string to decimal; error: %w", err)
	}
	rate, err = decimal.NewFromString(resp.Rate)
	if err != nil {
		return decimal.Zero, decimal.Zero, decimal.Zero, nil, fmt.Errorf("convert rate string to decimal; error: %w", err)
	}
	minFrom, err := decimal.NewFromString(resp.MinFrom)
	if err != nil {
		return decimal.Zero, decimal.Zero, decimal.Zero, nil, fmt.Errorf("convert minFrom string to decimal; error: %w", err)
	}
	maxFrom, err := decimal.NewFromString(resp.MaxFrom)
	if err != nil {
		return decimal.Zero, decimal.Zero, decimal.Zero, nil, fmt.Errorf("convert maxFrom string to decimal; error: %w", err)
	}
	limits = &model.Limits{
		Min: minFrom,
		Max: maxFrom,
	}
	return amountFrom, amountTo, rate, limits, nil
}

func MapModelEstimateResponseToProto(amountFrom, amountTo, rate decimal.Decimal, limits *model.Limits) *swapSwap.EstimateResponse {
	return &swapSwap.EstimateResponse{
		AmountFrom: amountFrom.String(),
		AmountTo:   amountTo.String(),
		Rate:       rate.String(),
		MinFrom:    limits.Min.String(),
		MaxFrom:    limits.Max.String(),
	}
}

func MapModelSwapRequestToProto(req *model.SwapRequest) *swapSwap.PrepareSwapRequest {
	return &swapSwap.PrepareSwapRequest{
		CoinFrom:    req.CoinFrom,
		CoinTo:      req.CoinTo,
		NetworkFrom: req.NetworkFrom,
		NetworkTo:   req.NetworkTo,
		AmountFrom:  req.AmountFrom.String(),
		AmountTo:    req.AmountTo.String(),
		AddressTo: &swapSwap.AddressData{
			Address: req.AddressTo.Address,
			Tag:     req.AddressTo.Tag,
		},
		ParentId: req.PartnerID,
	}
}

func MapProtoSwapRequestToModel(req *swapSwap.PrepareSwapRequest) (*model.SwapRequest, error) {
	amountFrom, err := decimal.NewFromString(req.AmountFrom)
	if err != nil {
		return nil, fmt.Errorf("convert amountFrom string to decimal; error: %w", err)
	}

	amountTo, err := decimal.NewFromString(req.AmountTo)
	if err != nil {
		return nil, fmt.Errorf("convert amountTo string to decimal; error: %w", err)
	}

	return &model.SwapRequest{
		CoinFrom:    req.CoinFrom,
		CoinTo:      req.CoinTo,
		NetworkFrom: req.NetworkFrom,
		NetworkTo:   req.NetworkTo,
		AmountFrom:  amountFrom,
		AmountTo:    amountTo,
		AddressTo: &model.AddressData{
			Address: req.AddressTo.Address,
			Tag:     req.AddressTo.Tag,
		},
		PartnerID: req.ParentId,
	}, nil
}

func MapProtoSwapResponseToModel(resp *swapSwap.PrepareSwapResponse) (id uuid.UUID, depositAddress *model.AddressData, err error) {
	id, err = uuid.Parse(resp.Id)
	if err != nil {
		return uuid.Nil, nil, fmt.Errorf("convert id string to uuid; error: %w", err)
	}
	depositAddress = &model.AddressData{
		Address: resp.DepositAddress.Address,
		Tag:     resp.DepositAddress.Tag,
	}
	return id, depositAddress, nil
}

func MapSwapResponseToProto(swapID uuid.UUID, depositAddress *model.AddressData) *swapSwap.PrepareSwapResponse {
	return &swapSwap.PrepareSwapResponse{
		Id: swapID.String(),
		DepositAddress: &swapSwap.AddressData{
			Address: depositAddress.Address,
			Tag:     depositAddress.Tag,
		},
	}
}

func MapGetSwapByIDToProto(swap *service.SwapByID) *swapSwap.GetSwapByIDResponse {
	return &swapSwap.GetSwapByIDResponse{
		CoinFrom:    swap.CoinFrom,
		CoinTo:      swap.CoinTo,
		NetworkFrom: swap.NetworkFrom,
		NetworkTo:   swap.NetworkTo,
		AmountFrom:  swap.AmountFrom.String(),
		AmountTo:    swap.AmountTo.String(),
		AddressFrom: &swapSwap.AddressData{
			Address: swap.AddressFrom,
			Tag:     swap.TagFrom,
		},
		AddressTo: &swapSwap.AddressData{
			Address: swap.AddressTo,
			Tag:     swap.TagTo,
		},
		Status:       int32(swap.Status),
		StartTime:    int32(swap.StartTime.Unix()),
		SwapDuration: int32(swap.SwapDuration.Seconds()),
		Rate:         swap.Rate.String(),
	}
}

func MapProtoGetSwapByIDResponseToModel(resp *swapSwap.GetSwapByIDResponse) (*model.SwapByIDResponse, error) {
	amountFrom, err := decimal.NewFromString(resp.AmountFrom)
	if err != nil {
		return nil, fmt.Errorf("convert amountFrom string to decimal; error: %w", err)
	}

	amountTo, err := decimal.NewFromString(resp.AmountTo)
	if err != nil {
		return nil, fmt.Errorf("convert amountTo string to decimal; error: %w", err)
	}

	rate, err := decimal.NewFromString(resp.Rate)

	startTime := time.Unix(int64(resp.StartTime), 0)

	timeDuration := time.Duration(resp.SwapDuration) * time.Second

	if err != nil {
		return nil, fmt.Errorf("convert rate string to decimal; error: %w", err)
	}
	return &model.SwapByIDResponse{
		CoinFrom:    resp.CoinFrom,
		CoinTo:      resp.CoinTo,
		NetworkFrom: resp.NetworkFrom,
		NetworkTo:   resp.NetworkTo,
		Rate:        rate,
		AmountFrom:  amountFrom,
		AmountTo:    amountTo,
		AddressTo: &model.AddressData{
			Address: resp.AddressTo.Address,
			Tag:     resp.AddressTo.Tag,
		},
		AddressFrom: &model.AddressData{
			Address: resp.AddressFrom.Address,
			Tag:     resp.AddressFrom.Tag,
		},
		StartTime:    startTime,
		SwapDuration: timeDuration,
		Status:       model.Status(resp.Status),
	}, nil
}

func MapModelSwapToProtoAdminSwap(swap *model.Swap) *swapAdmin.Swap {
	return &swapAdmin.Swap{
		Id:          swap.ID.String(),
		UserId:      swap.UserID.String(),
		CoinFrom:    swap.CoinFrom,
		CoinTo:      swap.CoinTo,
		NetworkFrom: swap.NetworkFrom,
		NetworkTo:   swap.NetworkTo,
		AddressFrom: &swapSwap.AddressData{
			Address: swap.AddressFrom,
			Tag:     swap.TagFrom,
		},
		AddressTo: &swapSwap.AddressData{
			Address: swap.AddressTo,
			Tag:     swap.TagTo,
		},
		AmountFrom: swap.AmountFrom.String(),
		AmountTo:   swap.AmountTo.String(),
		Status:     swapSwap.SwapStatus(swap.Status),
		StartTime:  timestamppb.New(swap.StartTime),
		EndTime:    timestamppb.New(swap.EndTime),
		PartnerId:  swap.PartnerID,
	}
}
