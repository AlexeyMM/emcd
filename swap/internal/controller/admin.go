package controller

import (
	"context"
	"fmt"

	"code.emcdtech.com/emcd/sdk/log"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/types/known/timestamppb"

	"code.emcdtech.com/b2b/swap/internal/controller/mapping"
	"code.emcdtech.com/b2b/swap/internal/service"
	"code.emcdtech.com/b2b/swap/model"
	"code.emcdtech.com/b2b/swap/package/gokit"
	"code.emcdtech.com/b2b/swap/protocol/swapAdmin"
)

type Admin struct {
	adminSrv service.Admin
	swapSrv  service.Swap
	swapAdmin.UnimplementedAdminServiceServer
}

func NewAdmin(adminSrv service.Admin, swapSrv service.Swap) *Admin {
	return &Admin{
		adminSrv: adminSrv,
		swapSrv:  swapSrv,
	}
}

func (a *Admin) GetBalanceByCoin(ctx context.Context, req *swapAdmin.GetBalanceByCoinRequest) (*swapAdmin.GetBalanceByCoinResponse, error) {
	amount, err := a.adminSrv.GetBalanceByCoin(ctx, req.AccountType.String(), req.Coin)
	if err != nil {
		log.Error(ctx, "admin: getBalanceByCoin: %s", err.Error())
		return nil, err
	}
	return &swapAdmin.GetBalanceByCoinResponse{
		Amount: amount.String(),
	}, nil
}

func (a *Admin) TransferBetweenAccountTypes(ctx context.Context, req *swapAdmin.TransferBetweenAccountTypesRequest) (*swapAdmin.TransferBetweenAccountTypesResponse, error) {
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		log.Error(ctx, "admin: transferBetweenAccountTypes: convert amount: %s", err.Error())
	}

	err = a.adminSrv.TransferBetweenAccountTypes(ctx, req.FromAccountType.String(), req.ToAccountType.String(), req.Coin, amount)
	if err != nil {
		log.Error(ctx, "admin: transferBetweenAccountTypes: %s", err.Error())
		return nil, err
	}

	return &swapAdmin.TransferBetweenAccountTypesResponse{}, nil
}

func (a *Admin) PlaceOrderForUSDT(ctx context.Context, req *swapAdmin.PlaceOrderForUSDTRequest) (*swapAdmin.PlaceOrderForUSDTResponse, error) {
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		log.Error(ctx, "admin: placeBuyOrderForUSDT: convert amount: %s", err.Error())
		return nil, err
	}

	orderID, err := a.adminSrv.PlaceOrderForUSDT(ctx, req.Coin, model.Direction(req.Direction), amount)
	if err != nil {
		log.Error(ctx, "admin: placeBuyOrderForUSDT: %s", err.Error())
		return nil, err
	}

	return &swapAdmin.PlaceOrderForUSDTResponse{
		OrderId: orderID.String(),
	}, nil
}

func (a *Admin) CheckOrder(ctx context.Context, req *swapAdmin.CheckOrderRequest) (*swapAdmin.CheckOrderResponse, error) {
	orderID, err := uuid.Parse(req.Id)
	if err != nil {
		log.Error(ctx, "admin: checkOrder: parse uuid: %s", err.Error())
		return nil, err
	}

	status, err := a.adminSrv.CheckOrder(ctx, orderID)
	if err != nil {
		log.Error(ctx, "admin: checkOrder: %s", err.Error())
		return nil, err
	}

	return &swapAdmin.CheckOrderResponse{
		Status: swapAdmin.OrderStatus(status),
	}, nil
}

func (a *Admin) Withdraw(ctx context.Context, req *swapAdmin.WithdrawRequest) (*swapAdmin.WithdrawResponse, error) {
	swapID, err := uuid.Parse(req.SwapId)
	if err != nil {
		log.Error(ctx, "admin: withdraw: parse uuid: %s", err.Error())
		return nil, err
	}

	withdrawalID, err := a.adminSrv.Withdraw(ctx, swapID)
	if err != nil {
		log.Error(ctx, "admin: withdraw: %s", err.Error())
		return nil, err
	}

	return &swapAdmin.WithdrawResponse{
		WithdrawalId: int32(withdrawalID),
	}, nil
}

func (a *Admin) GetWithdrawalLink(ctx context.Context, req *swapAdmin.GetWithdrawalLinkRequest) (*swapAdmin.GetWithdrawalLinkResponse, error) {
	link, err := a.adminSrv.GetWithdrawalLink(ctx, int(req.WithdrawalId))
	if err != nil {
		log.Error(ctx, "admin: getWithdrawalLink: %s", err.Error())
		return nil, err
	}

	return &swapAdmin.GetWithdrawalLinkResponse{
		Link: link,
	}, nil
}

func (a *Admin) RequestAQuote(ctx context.Context, req *swapAdmin.RequestAQuoteRequest) (*swapAdmin.RequestAQuoteResponse, error) {
	accountType, err := getConvertAccountType(req.AccountType.String())
	if err != nil {
		log.Error(ctx, "admin: requestAQuote: get accountType: %s", err.Error())
		return nil, err
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		log.Error(ctx, "admin: requestAQuote: convert amount: %s", err.Error())
		return nil, err
	}

	quote, err := a.adminSrv.RequestAQuote(ctx, req.From, req.To, accountType, amount)
	if err != nil {
		log.Error(ctx, "admin: requestAQuote: %s", err.Error())
		return nil, err
	}

	return &swapAdmin.RequestAQuoteResponse{
		Id:          quote.ID,
		Rate:        quote.Rate.String(),
		FromAmount:  quote.FromAmount.String(),
		ToAmount:    quote.ToAmount.String(),
		ExpiredTime: timestamppb.New(quote.ExpiredTime),
	}, nil
}

func (a *Admin) ConfirmAQuote(ctx context.Context, req *swapAdmin.ConfirmAQuoteRequest) (*swapAdmin.ConfirmAQuoteResponse, error) {
	status, err := a.adminSrv.ConfirmAQuote(ctx, req.Id)
	if err != nil {
		log.Error(ctx, "admin: confirmAQuote: %s", err.Error())
		return nil, err
	}
	return &swapAdmin.ConfirmAQuoteResponse{
		Status: status,
	}, nil
}

func (a *Admin) GetConvertStatus(ctx context.Context, req *swapAdmin.GetConvertStatusRequest) (*swapAdmin.GetConvertStatusResponse, error) {
	accountType, err := getConvertAccountType(req.AccountType.String())
	if err != nil {
		log.Error(ctx, "admin: getConvertStatus: getConvertAccountType: %s", err.Error())
		return nil, err
	}

	status, err := a.adminSrv.GetConvertStatus(ctx, req.Id, accountType)
	if err != nil {
		log.Error(ctx, "admin: getConvertStatus: %s", err.Error())
		return nil, err
	}
	return &swapAdmin.GetConvertStatusResponse{
		Status: status,
	}, nil
}

func (a *Admin) ChangeManualSwapStatus(ctx context.Context, req *swapAdmin.ChangeManualSwapStatusRequest) (*swapAdmin.ChangeManualSwapStatusResponse, error) {
	swapID, err := uuid.Parse(req.SwapId)
	if err != nil {
		log.Error(ctx, "admin: changeManualSwapStatus: parse uuid: %s", err.Error())
		return nil, err
	}

	err = a.adminSrv.ChangeManualSwapStatus(ctx, swapID, model.Status(req.Status))
	if err != nil {
		return nil, fmt.Errorf("admin: changeManualSwapStatus: %s", err.Error())
	}

	return &swapAdmin.ChangeManualSwapStatusResponse{}, nil
}

func (a *Admin) GetSwaps(ctx context.Context, req *swapAdmin.GetSwapsRequest) (*swapAdmin.GetSwapsResponse, error) {
	var filter = &model.SwapFilter{
		Offset: gokit.Ptr(int(req.Offset)),
		Limit:  gokit.Ptr(int(req.Limit)),
	}

	if req.Id != nil {
		id, err := uuid.Parse(req.GetId())
		if err != nil {
			log.Error(ctx, "admin: get swaps: parse uuid: %s", err.Error())
			return nil, err
		}
		filter.ID = &id
	}

	if req.AddressFrom != nil {
		filter.AddressFrom = req.AddressFrom
	}

	if req.Email != nil {
		filter.Email = req.Email
	}

	if req.UserId != nil {
		id, err := uuid.Parse(req.GetId())
		if err != nil {
			log.Error(ctx, "admin: get swaps: parse user id uuid: %s", err.Error())
			return nil, err
		}
		filter.UserID = &id
	}

	if req.From != nil {
		filter.StartTimeFrom = gokit.Ptr(req.From.AsTime())
	}

	if req.To != nil {
		filter.StartTimeTo = gokit.Ptr(req.To.AsTime())
	}

	swaps, total, err := a.swapSrv.GetSwaps(ctx, filter)
	if err != nil {
		log.Error(ctx, "admin: getSwaps: %s", err.Error())
		return nil, err
	}

	var adminSwaps = make([]*swapAdmin.Swap, len(swaps))
	for i, swap := range swaps {
		adminSwaps[i] = mapping.MapModelSwapToProtoAdminSwap(swap)
	}

	return &swapAdmin.GetSwapsResponse{
		Swaps: adminSwaps,
		Total: int64(total),
	}, nil
}

func (a *Admin) GetSwapStatusHistory(ctx context.Context, req *swapAdmin.GetSwapStatusHistoryRequest) (*swapAdmin.GetSwapStatusHistoryResponse, error) {
	id, err := uuid.Parse(req.GetSwapId())
	if err != nil {
		log.Error(ctx, "admin: get swap id: parse uuid: %s", err.Error())
		return nil, err
	}
	items, err := a.adminSrv.GetSwapStatusHistory(ctx, id)
	if err != nil {
		log.Error(ctx, "admin: getSwapStatusHistory: %s", err.Error())
		return nil, err
	}

	var history = make([]*swapAdmin.GetSwapStatusHistoryResponse_HistoryItem, len(items))
	for i, item := range history {
		history[i] = &swapAdmin.GetSwapStatusHistoryResponse_HistoryItem{
			Status: item.Status,
			SetAt:  item.SetAt,
		}
	}

	return &swapAdmin.GetSwapStatusHistoryResponse{StatusHistory: history}, nil
}

func (a *Admin) SetDestinationAddress(ctx context.Context, req *swapAdmin.SetDestinationAddressRequest) (*swapAdmin.SetDestinationAddressResponse, error) {
	swapID, err := uuid.Parse(req.SwapId)
	if err != nil {
		log.Error(ctx, "admin: setAddress: parse uuid: %s", err.Error())
		return nil, err
	}

	err = a.swapSrv.Update(
		ctx,
		&model.SwapFilter{
			ID: &swapID,
		},
		&model.SwapPartial{
			AddressTo: req.Address,
			TagTo:     req.Tag,
		})
	if err != nil {
		log.Error(ctx, "admin: setAddress: update: %s", err.Error())
		return nil, err
	}

	return &swapAdmin.SetDestinationAddressResponse{}, nil
}

func getConvertAccountType(accountType string) (string, error) {
	// Не выделяли новые типы под конвертацию. Интуитивно это одно и то же, названия другие
	switch accountType {
	case model.Fund:
		return "eb_convert_funding", nil
	case model.UNIFIED:
		return "eb_convert_uta", nil
	default:
		return "", fmt.Errorf("unknown account type: %s", accountType)
	}
}
