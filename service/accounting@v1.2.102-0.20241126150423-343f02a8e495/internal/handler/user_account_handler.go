package handler

import (
	sdkLog "code.emcdtech.com/emcd/sdk/log"
	"code.emcdtech.com/emcd/service/accounting/internal/handler/mapping"
	"code.emcdtech.com/emcd/service/accounting/internal/service"
	userAccountPb "code.emcdtech.com/emcd/service/accounting/protocol/user_account"
	"code.emcdtech.com/emcd/service/accounting/repository"
	coinValidatorRepo "code.emcdtech.com/emcd/service/coin/repository"
	"context"
)

type UserAccountHandler struct {
	userAccountService service.UserAccountService
	coinValidator      coinValidatorRepo.CoinValidatorRepository
	userAccountPb.UnimplementedUserAccountServiceServer
}

func NewUserAccountHandler(
	userAccountService service.UserAccountService,
	coinValidator coinValidatorRepo.CoinValidatorRepository,
) *UserAccountHandler {

	return &UserAccountHandler{
		userAccountService:                    userAccountService,
		coinValidator:                         coinValidator,
		UnimplementedUserAccountServiceServer: userAccountPb.UnimplementedUserAccountServiceServer{},
	}
}

func (h *UserAccountHandler) CreateUserAccounts(ctx context.Context, req *userAccountPb.UserAccountMultiRequest) (*userAccountPb.UserAccountMultiResponse, error) {
	if userId, userIdNew, userAccounts, err := mapping.MapProtoMultiRequestToModelUserAccounts(h.coinValidator, req); err != nil {
		sdkLog.Error(ctx, err.Error())

		return nil, repository.ErrAcc1011
	} else if userAccountsResponse, err := h.userAccountService.CreateUserAccounts(ctx, userId, userIdNew, userAccounts); err != nil {
		sdkLog.Error(ctx, err.Error())

		return nil, repository.ErrAcc1012
	} else {

		return mapping.MapModelUserAccountsToProtoMultiResponse(nil, userAccountsResponse), nil
	}
}

func (h *UserAccountHandler) GetOrCreateUserAccount(ctx context.Context, req *userAccountPb.UserAccountOneRequest) (*userAccountPb.UserAccountResponse, error) {
	if userAccount, err := mapping.MapProtoOneRequestToModelUserAccount(h.coinValidator, req); err != nil {
		sdkLog.Error(ctx, err.Error())

		return nil, repository.ErrAcc1021
	} else if userAccountResponse, err := h.userAccountService.GetOrCreateUserAccount(ctx, userAccount); err != nil {
		sdkLog.Error(ctx, err.Error())

		return nil, repository.ErrAcc1022
	} else {

		return mapping.MapModelUserAccountToProtoResponse(userAccountResponse), nil
	}
}

func (h *UserAccountHandler) GetUserAccountsByFilter(ctx context.Context, req *userAccountPb.UserAccountFilter) (*userAccountPb.UserAccountMultiResponse, error) {
	if filter, err := mapping.MapProtoToModelUserAccountFilter(h.coinValidator, req); err != nil {
		sdkLog.Error(ctx, err.Error())

		return nil, repository.ErrAcc1031
	} else if totalCount, userAccountsResponse, err := h.userAccountService.GetUserAccountsByFilter(ctx, filter); err != nil {
		sdkLog.Error(ctx, err.Error())

		return nil, repository.ErrAcc1032
	} else {

		return mapping.MapModelUserAccountsToProtoMultiResponse(totalCount, userAccountsResponse), nil
	}
}

func (h *UserAccountHandler) GetUserAccountById(ctx context.Context, req *userAccountPb.UserAccountId) (*userAccountPb.UserAccountResponse, error) {
	if userAccountResponse, err := h.userAccountService.GetUserAccountById(ctx, req.Id); err != nil {
		sdkLog.Error(ctx, err.Error())

		return nil, repository.ErrAcc1041
	} else {

		return mapping.MapModelUserAccountToProtoResponse(userAccountResponse), nil
	}
}

func (h *UserAccountHandler) GetUserAccountByConstraint(ctx context.Context, req *userAccountPb.UserAccountConstraintRequest) (*userAccountPb.UserAccountResponse, error) {
	if userUuid, coinNew, userAccountTypeId, err := mapping.MapProtoUserAccountConstraintRequestToArgs(h.coinValidator, req); err != nil {
		sdkLog.Error(ctx, err.Error())

		return nil, repository.ErrAcc1051
	} else if userAccountResponse, err := h.userAccountService.GetUserAccountByConstraint(ctx, userUuid, coinNew, userAccountTypeId); err != nil {
		sdkLog.Error(ctx, err.Error())

		return nil, repository.ErrAcc1052
	} else {

		return mapping.MapModelUserAccountToProtoResponse(userAccountResponse), nil
	}
}

func (h *UserAccountHandler) GetUserAccountsByUuid(ctx context.Context, req *userAccountPb.UserAccountUuid) (*userAccountPb.UserAccountMultiResponse, error) {
	if userUuid, err := mapping.MapProtoUuidToUuid(req.Uuid); err != nil {
		sdkLog.Error(ctx, err.Error())

		return nil, repository.ErrAcc1061
	} else if userAccountsResponse, err := h.userAccountService.GetUserAccountsByUuid(ctx, userUuid); err != nil {
		sdkLog.Error(ctx, err.Error())

		return nil, repository.ErrAcc1062
	} else {

		return mapping.MapModelUserAccountsToProtoMultiResponse(nil, userAccountsResponse), nil
	}
}
