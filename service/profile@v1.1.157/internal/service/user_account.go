package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	accountingModel "code.emcdtech.com/emcd/service/accounting/model"
	userAccountModel "code.emcdtech.com/emcd/service/accounting/model"
	userAccountModelEnum "code.emcdtech.com/emcd/service/accounting/model/enum"
	userAccountRepository "code.emcdtech.com/emcd/service/accounting/repository"
	coinValidatorRepo "code.emcdtech.com/emcd/service/coin/repository"
	"github.com/google/uuid"

	"code.emcdtech.com/emcd/service/profile/internal/utils"
)

var (
	ErrUserAccountListIsEmpty     = errors.New("empty list response of user accounts")
	ErrUserAccountListMoreThenOne = errors.New("list contain more than one user account")
	ErrUserAccountListUnexpected  = errors.New("unexpected error")
)

type UserAccountService interface {
	CreateUserAccount(
		ctx context.Context,
		userId int32,
		userIdNew uuid.UUID,
		coinId, userAccountTypeId int32,
		minpay float64,
	) (int, error)
	GetUserAccountIdByLegacyParams(ctx context.Context, userId, coinId, userAccountTypeId int32) (int, error)
	CreateUserAccounts(ctx context.Context, userId int32, userIdNew uuid.UUID, userAccounts accountingModel.UserAccounts) (accountingModel.UserAccounts, error)
}

type userAccountServiceImp struct {
	userAccountRepo userAccountRepository.UserAccountRepository
	coinValidator   coinValidatorRepo.CoinValidatorRepository
}

func NewUserAccountService(
	userAccountRepo userAccountRepository.UserAccountRepository,
	coinValidator coinValidatorRepo.CoinValidatorRepository,
) UserAccountService {
	return &userAccountServiceImp{
		userAccountRepo: userAccountRepo,
		coinValidator:   coinValidator,
	}
}

func (s *userAccountServiceImp) CreateUserAccount(
	ctx context.Context,
	userId int32,
	userIdNew uuid.UUID,
	coinId, userAccountTypeId int32,
	minpay float64,
) (int, error) {
	userAccountTypeIdEnum := userAccountModelEnum.NewAccountTypeId(userAccountTypeId)
	if coinCode, ok := s.coinValidator.GetCodeById(coinId); !ok {
		return 0, fmt.Errorf("unknown coin id: %d", coinId)
	} else if err := userAccountTypeIdEnum.Validate(); err != nil {
		return 0, fmt.Errorf("unknown account type: %w", err)
	} else {
		userAccount := &userAccountModel.UserAccount{
			ID:            0,
			UserID:        userId,
			CoinID:        0, // will ignore
			AccountTypeID: userAccountModelEnum.NewAccountTypeIdWrapper(userAccountTypeIdEnum),
			Minpay:        minpay,
			Address:       sql.NullString{},
			ChangedAt:     sql.NullTime{},
			Img1:          sql.NullFloat64{},
			Img2:          sql.NullFloat64{},
			IsActive:      sql.NullBool{},
			CreatedAt:     sql.NullTime{},
			UpdatedAt:     sql.NullTime{},
			Fee:           sql.NullFloat64{},
			UserIDNew:     utils.UuidToUuidNull(userIdNew),
			CoinNew:       utils.StringToStringNull(coinCode),
		}

		if userAccountsResponse, err := s.userAccountRepo.CreateUserAccounts(ctx, userId, userIdNew, userAccountModel.UserAccounts{userAccount}); err != nil {
			return 0, err
		} else if len(userAccountsResponse) != 1 {
			return 0, ErrUserAccountListUnexpected
		} else {
			return int(userAccountsResponse[0].ID), nil
		}
	}
}

func (s *userAccountServiceImp) GetUserAccountIdByLegacyParams(
	ctx context.Context,
	userId, coinId, userAccountTypeId int32,
) (int, error) {
	userAccountTypeIdEnum := userAccountModelEnum.NewAccountTypeId(userAccountTypeId)
	if coinCode, ok := s.coinValidator.GetCodeById(coinId); !ok {
		return 0, fmt.Errorf("unknown coin id: %d", coinId)
	} else if err := userAccountTypeIdEnum.Validate(); err != nil {
		return 0, fmt.Errorf("unknown account type: %w", err)
	} else {

		filter := &userAccountModel.UserAccountFilter{
			ID:              nil,
			UserID:          &userId,
			AccountTypeID:   userAccountTypeIdEnum.ToPtr(),
			UserIDNew:       nil,
			CoinNew:         &coinCode,
			Pagination:      nil,
			UserIDNewIsNull: nil,
			CoinNewIsNull:   nil,
		}

		if _, userAccounts, err := s.userAccountRepo.GetUserAccountsByFilter(ctx, filter); err != nil {
			return 0, err
		} else if len(userAccounts) == 0 {
			return 0, ErrUserAccountListIsEmpty
		} else if len(userAccounts) > 1 {
			return 0, ErrUserAccountListMoreThenOne
		} else {
			return int(userAccounts[0].ID), nil
		}
	}
}

func (s *userAccountServiceImp) CreateUserAccounts(ctx context.Context, userId int32, userIdNew uuid.UUID, userAccounts accountingModel.UserAccounts) (accountingModel.UserAccounts, error) {
	userAccountsResponse, err := s.userAccountRepo.CreateUserAccounts(ctx, userId, userIdNew, userAccounts)
	if err != nil {
		return nil, err
	}
	return userAccountsResponse, err
}
