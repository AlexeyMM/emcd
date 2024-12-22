package repository

import (
	"context"
	"database/sql"

	coinValidatorRepo "code.emcdtech.com/emcd/service/coin/repository"
	"github.com/google/uuid"

	"code.emcdtech.com/emcd/service/accounting/internal/handler/mapping"
	"code.emcdtech.com/emcd/service/accounting/internal/utils"
	"code.emcdtech.com/emcd/service/accounting/model"
	"code.emcdtech.com/emcd/service/accounting/model/enum"
	userAccountPb "code.emcdtech.com/emcd/service/accounting/protocol/user_account"
)

// UserAccountRepository is helper interface for request via grpc from other microservices and using internal structs as is.
type UserAccountRepository interface {
	CreateUserAccounts(ctx context.Context, userId int32, userIdNew uuid.UUID, req model.UserAccounts) (model.UserAccounts, error)
	CreateDefaultUserAccounts(ctx context.Context, userId int32, userIdNew uuid.UUID, minPay, fee float64) (model.UserAccounts, error)
	GetOrCreateUserAccount(ctx context.Context, req *model.UserAccount) (*model.UserAccount, error)
	GetOrCreateUserAccountByArgs(ctx context.Context, userId int32, userIdNew uuid.UUID, coinIdNew string, accountTypeId enum.AccountTypeId, minPay, fee float64) (*model.UserAccount, error)
	GetUserAccountsByFilter(ctx context.Context, filter *model.UserAccountFilter) (*uint64, model.UserAccounts, error)
	GetUserAccountById(ctx context.Context, userId int32) (*model.UserAccount, error)
	GetUserAccountByConstraint(ctx context.Context, userIdNew uuid.UUID, coinIdNew string, accountTypeId enum.AccountTypeId) (*model.UserAccount, error)
	GetUserAccountsByUuid(ctx context.Context, userIdNew uuid.UUID) (model.UserAccounts, error)
}

type userAccountImp struct {
	handler       userAccountPb.UserAccountServiceClient
	coinValidator coinValidatorRepo.CoinValidatorRepository
}

// NewUserAccountRepository - only for external using
func NewUserAccountRepository(
	handler userAccountPb.UserAccountServiceClient,
	coinValidator coinValidatorRepo.CoinValidatorRepository,
) UserAccountRepository {

	return &userAccountImp{
		handler:       handler,
		coinValidator: coinValidator,
	}
}

func (r *userAccountImp) CreateUserAccounts(ctx context.Context, userId int32, userIdNew uuid.UUID, req model.UserAccounts) (model.UserAccounts, error) {
	if userAccounts, err := mapping.MapModelUserAccountsToProtoMultiRequest(userId, userIdNew, req); err != nil {

		return nil, err
	} else if resp, err := r.handler.CreateUserAccounts(ctx, userAccounts); err != nil {

		return nil, err
	} else {

		return mapping.MapProtoMultiResponseToModelUserAccounts(ctx, r.coinValidator, resp.UserAccounts), nil
	}
}

func (r *userAccountImp) CreateDefaultUserAccounts(ctx context.Context, userId int32, userIdNew uuid.UUID, minPay, fee float64) (model.UserAccounts, error) {
	userAccountRequest := GenerateDefaultUsersAccounts(userId, userIdNew, minPay, fee)

	return r.CreateUserAccounts(ctx, userId, userIdNew, userAccountRequest)
}

func (r *userAccountImp) GetOrCreateUserAccount(ctx context.Context, req *model.UserAccount) (*model.UserAccount, error) {
	userAccount := mapping.MapModelUserAccountToProtoOneRequest(req)

	if userAccountResp, err := r.handler.GetOrCreateUserAccount(ctx, userAccount); err != nil {

		return nil, err
	} else {

		return mapping.MapProtoResponseToModelUserAccount(ctx, r.coinValidator, userAccountResp), nil
	}
}

func (r *userAccountImp) GetOrCreateUserAccountByArgs(ctx context.Context, userId int32, userIdNew uuid.UUID, coinIdNew string, accountTypeId enum.AccountTypeId, minPay, fee float64) (*model.UserAccount, error) {
	req := &model.UserAccount{
		ID:            0,
		UserID:        userId,
		CoinID:        0,
		AccountTypeID: enum.NewAccountTypeIdWrapper(accountTypeId),
		Minpay:        minPay,
		Address:       sql.NullString{},
		ChangedAt:     sql.NullTime{},
		Img1:          sql.NullFloat64{},
		Img2:          sql.NullFloat64{},
		IsActive:      sql.NullBool{},
		CreatedAt:     sql.NullTime{},
		UpdatedAt:     sql.NullTime{},
		Fee:           sql.NullFloat64{Float64: fee, Valid: true},
		UserIDNew:     uuid.NullUUID{UUID: userIdNew, Valid: true},
		CoinNew:       sql.NullString{String: coinIdNew, Valid: true},
	}

	return r.GetOrCreateUserAccount(ctx, req)
}

func (r *userAccountImp) GetUserAccountsByFilter(ctx context.Context, filter *model.UserAccountFilter) (*uint64, model.UserAccounts, error) {
	filterProto := mapping.MapModelUserAccountFilterToProto(filter)

	if userAccountMultiResponse, err := r.handler.GetUserAccountsByFilter(ctx, filterProto); err != nil {

		return nil, nil, err
	} else {

		return userAccountMultiResponse.TotalCount, mapping.MapProtoMultiResponseToModelUserAccounts(ctx, r.coinValidator, userAccountMultiResponse.UserAccounts), nil
	}
}

func (r *userAccountImp) GetUserAccountById(ctx context.Context, userId int32) (*model.UserAccount, error) {
	if userAccountResponse, err := r.handler.GetUserAccountById(ctx, &userAccountPb.UserAccountId{Id: userId}); err != nil {

		return nil, err
	} else {

		return mapping.MapProtoResponseToModelUserAccount(ctx, r.coinValidator, userAccountResponse), nil
	}
}

func (r *userAccountImp) GetUserAccountByConstraint(ctx context.Context, userIdNew uuid.UUID, coinIdNew string, accountTypeId enum.AccountTypeId) (*model.UserAccount, error) {
	req := &userAccountPb.UserAccountConstraintRequest{
		UserIdNew:     userIdNew.String(),
		CoinNew:       coinIdNew,
		AccountTypeId: accountTypeId.ToInt32(),
	}

	if userAccountResponse, err := r.handler.GetUserAccountByConstraint(ctx, req); err != nil {

		return nil, err
	} else {

		return mapping.MapProtoResponseToModelUserAccount(ctx, r.coinValidator, userAccountResponse), nil
	}
}

func (r *userAccountImp) GetUserAccountsByUuid(ctx context.Context, userIdNew uuid.UUID) (model.UserAccounts, error) {
	if userAccountMultiResponse, err := r.handler.GetUserAccountsByUuid(ctx, &userAccountPb.UserAccountUuid{Uuid: userIdNew.String()}); err != nil {

		return nil, err
	} else {

		return mapping.MapProtoMultiResponseToModelUserAccounts(ctx, r.coinValidator, userAccountMultiResponse.UserAccounts), nil
	}
}

func newUserAccount(userIdLegacy int32, userIdNew uuid.UUID, coinCode string, accountTypeId enum.AccountTypeId, minPay, fee float64) *model.UserAccount {

	return &model.UserAccount{
		ID:            0,
		UserID:        userIdLegacy,
		CoinID:        0, // will ignore
		AccountTypeID: enum.NewAccountTypeIdWrapper(accountTypeId),
		Minpay:        minPay,
		Address:       sql.NullString{},
		ChangedAt:     sql.NullTime{},
		Img1:          sql.NullFloat64{},
		Img2:          sql.NullFloat64{},
		IsActive:      utils.BoolToBoolNull(true),
		CreatedAt:     sql.NullTime{},
		UpdatedAt:     sql.NullTime{},
		Fee:           utils.Float64ToFloat64Null(fee),
		UserIDNew:     utils.UuidToUuidNull(userIdNew),
		CoinNew:       utils.StringToStringNull(coinCode),
	}
}

func GenerateDefaultUsersAccounts(userIdLegacy int32, userIdNew uuid.UUID, minPay, fee float64) model.UserAccounts {

	return model.UserAccounts{
		// BTC users accounts
		newUserAccount(userIdLegacy, userIdNew, "btc", enum.WalletAccountTypeID, minPay, fee),
		newUserAccount(userIdLegacy, userIdNew, "btc", enum.MiningAccountTypeID, minPay, fee),
		newUserAccount(userIdLegacy, userIdNew, "btc", enum.ReferralAccountTypeID, minPay, fee),
		newUserAccount(userIdLegacy, userIdNew, "btc", enum.BlockUserAccountTypeID, minPay, fee),

		// BCH users accounts
		newUserAccount(userIdLegacy, userIdNew, "bch", enum.WalletAccountTypeID, minPay, fee),
		newUserAccount(userIdLegacy, userIdNew, "bch", enum.MiningAccountTypeID, minPay, fee),
		newUserAccount(userIdLegacy, userIdNew, "bch", enum.ReferralAccountTypeID, minPay, fee),
		newUserAccount(userIdLegacy, userIdNew, "bch", enum.BlockUserAccountTypeID, minPay, fee),

		// LTC users accounts
		newUserAccount(userIdLegacy, userIdNew, "ltc", enum.WalletAccountTypeID, minPay, fee),
		newUserAccount(userIdLegacy, userIdNew, "ltc", enum.MiningAccountTypeID, minPay, fee),
		newUserAccount(userIdLegacy, userIdNew, "ltc", enum.ReferralAccountTypeID, minPay, fee),
		newUserAccount(userIdLegacy, userIdNew, "ltc", enum.BlockUserAccountTypeID, minPay, fee),

		// DASH users accounts
		newUserAccount(userIdLegacy, userIdNew, "dash", enum.WalletAccountTypeID, minPay, fee),
		newUserAccount(userIdLegacy, userIdNew, "dash", enum.MiningAccountTypeID, minPay, fee),
		newUserAccount(userIdLegacy, userIdNew, "dash", enum.ReferralAccountTypeID, minPay, fee),
		newUserAccount(userIdLegacy, userIdNew, "dash", enum.BlockUserAccountTypeID, minPay, fee),

		// ETH users accounts
		newUserAccount(userIdLegacy, userIdNew, "eth", enum.WalletAccountTypeID, minPay, fee),
		newUserAccount(userIdLegacy, userIdNew, "eth", enum.MiningAccountTypeID, minPay, fee),
		newUserAccount(userIdLegacy, userIdNew, "eth", enum.ReferralAccountTypeID, minPay, fee),
		newUserAccount(userIdLegacy, userIdNew, "eth", enum.BlockUserAccountTypeID, minPay, fee),

		// ETC users accounts
		newUserAccount(userIdLegacy, userIdNew, "etc", enum.WalletAccountTypeID, minPay, fee),
		newUserAccount(userIdLegacy, userIdNew, "etc", enum.MiningAccountTypeID, minPay, fee),
		newUserAccount(userIdLegacy, userIdNew, "etc", enum.ReferralAccountTypeID, minPay, fee),
		newUserAccount(userIdLegacy, userIdNew, "etc", enum.BlockUserAccountTypeID, minPay, fee),

		// DOGE users accounts
		newUserAccount(userIdLegacy, userIdNew, "doge", enum.WalletAccountTypeID, minPay, fee),
		newUserAccount(userIdLegacy, userIdNew, "doge", enum.MiningAccountTypeID, minPay, fee),
		newUserAccount(userIdLegacy, userIdNew, "doge", enum.ReferralAccountTypeID, minPay, fee),
		newUserAccount(userIdLegacy, userIdNew, "doge", enum.BlockUserAccountTypeID, minPay, fee),

		// USDT users accounts
		newUserAccount(userIdLegacy, userIdNew, "usdt", enum.WalletAccountTypeID, minPay, fee),
		newUserAccount(userIdLegacy, userIdNew, "usdt", enum.BlockUserAccountTypeID, minPay, fee),

		// USDC users accounts
		newUserAccount(userIdLegacy, userIdNew, "usdc", enum.WalletAccountTypeID, minPay, fee),
		newUserAccount(userIdLegacy, userIdNew, "usdc", enum.BlockUserAccountTypeID, minPay, fee),
	}
}
