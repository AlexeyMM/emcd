package handler

import (
	sdkLog "code.emcdtech.com/emcd/sdk/log"
	"code.emcdtech.com/emcd/service/accounting/internal/repository"
	"code.emcdtech.com/emcd/service/accounting/model"
	accountingPb "code.emcdtech.com/emcd/service/accounting/protocol/accounting"
	"context"
	"errors"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *AccountingHandler) GetUserAccount(ctx context.Context, req *accountingPb.GetUserAccountRequest) (*accountingPb.GetUserAccountResponse, error) {
	userAccount, err := h.userAccountService.GetUserAccount(ctx, int32(req.Id))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			sdkLog.Info(ctx, "GetUserAccount (account id: %d): %v", req.Id, err)
			return nil, status.Errorf(codes.NotFound, "account not found: %v", req.Id)
		}
		sdkLog.Error(ctx, "GetUserAccount (account id: %d): %v", req.Id, err.Error())
		return nil, status.Errorf(codes.Internal, "internal error: %v", req.Id)
	}
	return &accountingPb.GetUserAccountResponse{
		UserAccount: h.userAccountToPBUserAccount(userAccount),
	}, nil
}

func (h *AccountingHandler) GetUserAccounts(ctx context.Context, req *accountingPb.GetUserAccountsRequest) (*accountingPb.GetUserAccountsResponse, error) {
	userId, err := uuid.Parse(req.GetUserId())
	if err != nil {
		sdkLog.Error(ctx, "parse string (user id: %s): %v", req.GetUserId(), err)
		return nil, status.Errorf(codes.InvalidArgument, "invalid user id: %v", req.GetUserId())
	}

	userAccounts, err := h.userAccountService.GetUserAccounts(ctx, userId)
	if err != nil {
		sdkLog.Error(ctx, "GetUserAccounts (user id: %s): %v", userId, err)
		return nil, status.Errorf(codes.Internal, "internal error: %v", userId)
	}

	pbUserAccounts := make([]*accountingPb.UserAccount, 0, len(userAccounts))
	for _, userAccount := range userAccounts {
		pbUserAccount := h.userAccountToPBUserAccount(userAccount)
		pbUserAccounts = append(pbUserAccounts, pbUserAccount)
	}

	return &accountingPb.GetUserAccountsResponse{
		UserAccounts: pbUserAccounts,
	}, nil
}

func (h *AccountingHandler) userAccountToPBUserAccount(userAccount *model.UserAccount) *accountingPb.UserAccount {
	return &accountingPb.UserAccount{
		Id:     int64(userAccount.ID),
		UserId: userAccount.UserIDNew.UUID.String(),
		CoinId: userAccount.CoinNew.String,
		TypeId: int32(userAccount.AccountTypeID.ToInt()),
	}
}
