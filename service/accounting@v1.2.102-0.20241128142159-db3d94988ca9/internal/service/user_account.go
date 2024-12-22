package service

import (
	"code.emcdtech.com/emcd/sdk/log"
	"code.emcdtech.com/emcd/service/accounting/internal/repository"
	"code.emcdtech.com/emcd/service/accounting/internal/utils"
	"code.emcdtech.com/emcd/service/accounting/model"
	"code.emcdtech.com/emcd/service/accounting/model/enum"
	"context"
	"fmt"
	"github.com/google/uuid"
)

type UserAccountService interface {
	CreateUserAccounts(ctx context.Context, userId int32, userIdNew uuid.UUID, userAccounts model.UserAccounts) (model.UserAccounts, error)
	GetOrCreateUserAccount(ctx context.Context, userAccount *model.UserAccount) (*model.UserAccount, error)
	GetUserAccountsByFilter(ctx context.Context, filter *model.UserAccountFilter) (*uint64, model.UserAccounts, error)
	GetUserAccountById(ctx context.Context, id int32) (*model.UserAccount, error)
	GetUserAccountByConstraint(ctx context.Context, userIdNew uuid.UUID, coinIdNew string, userAccountId enum.AccountTypeId) (*model.UserAccount, error)
	GetUserAccountsByUuid(ctx context.Context, userUuid uuid.UUID) (model.UserAccounts, error)
}

type userAccountImpl struct {
	repo repository.UserAccountRepo
}

func NewUserAccountService(repo repository.UserAccountRepo) UserAccountService {

	return &userAccountImpl{
		repo: repo,
	}
}

func (s *userAccountImpl) CreateUserAccounts(ctx context.Context, userId int32, userIdNew uuid.UUID, userAccounts model.UserAccounts) (model.UserAccounts, error) {
	if err := s.repo.WithinTransaction(ctx, func(ctx context.Context) error {
		if err := s.repo.AddUserAccounts(ctx, userId, userIdNew, userAccounts); err != nil {

			return fmt.Errorf("failed add user accounts list: %w", err)
		} else {

			return nil
		}
	}); err != nil {

		return nil, err
	} else {

		return userAccounts, nil
	}
}

func (s *userAccountImpl) GetOrCreateUserAccount(ctx context.Context, userAccount *model.UserAccount) (*model.UserAccount, error) {
	if err := s.repo.WithinTransaction(ctx, func(ctx context.Context) error {
		if err := s.repo.AddUserAccount(ctx, userAccount); err != nil {

			return fmt.Errorf("failed add user account: %w", err)
		} else {

			return nil
		}
	}); err != nil {

		return nil, err
	} else {

		return userAccount, nil
	}
}

func (s *userAccountImpl) GetUserAccountsByFilter(ctx context.Context, filter *model.UserAccountFilter) (*uint64, model.UserAccounts, error) {
	if totalCount, userAccounts, err := s.repo.FindUserAccountByFilter(ctx, filter); err != nil {

		return nil, nil, fmt.Errorf("failed find user accounts list by filter: %w", err)
	} else {

		return totalCount, userAccounts, nil
	}
}

func (s *userAccountImpl) GetUserAccountById(ctx context.Context, id int32) (*model.UserAccount, error) {
	filter := &model.UserAccountFilter{
		ID:              &id,
		UserID:          nil,
		AccountTypeID:   nil,
		UserIDNew:       nil,
		CoinNew:         nil,
		IsActive:        nil,
		Pagination:      nil,
		UserIDNewIsNull: nil,
		CoinNewIsNull:   nil,
	}

	if _, userAccounts, err := s.GetUserAccountsByFilter(ctx, filter); err != nil {

		return nil, err
	} else if len(userAccounts) == 0 {
		const errMsg = "user account is not found"
		log.Error(ctx, "%s: %d", errMsg, id)

		return nil, fmt.Errorf("%s", errMsg)
	} else if len(userAccounts) > 1 {
		const errMsg = "user account is not uniq id"
		log.Error(ctx, "%s: %d", errMsg, id)

		return nil, fmt.Errorf("%s", errMsg)
	} else {

		return userAccounts[0], nil
	}
}

func (s *userAccountImpl) GetUserAccountByConstraint(ctx context.Context, userIdNew uuid.UUID, coinIdNew string, userAccountId enum.AccountTypeId) (*model.UserAccount, error) {

	filter := &model.UserAccountFilter{
		ID:              nil,
		UserID:          nil,
		AccountTypeID:   userAccountId.ToPtr(),
		UserIDNew:       &userIdNew,
		CoinNew:         utils.StringToPtr(coinIdNew),
		IsActive:        nil,
		Pagination:      nil,
		UserIDNewIsNull: nil,
		CoinNewIsNull:   nil,
	}

	if _, userAccounts, err := s.GetUserAccountsByFilter(ctx, filter); err != nil {

		return nil, err
	} else if len(userAccounts) == 0 {
		const errMsg = "user account is not found"
		log.Error(ctx, "%s: %s, %s, %s", errMsg, userIdNew.String(), coinIdNew, userAccountId.ToString())

		return nil, fmt.Errorf("%s", errMsg)
	} else if len(userAccounts) > 1 {
		const errMsg = "user account is not uniq"
		log.Error(ctx, "%s: %s, %s, %s", errMsg, userIdNew.String(), coinIdNew, userAccountId.ToString())

		return nil, fmt.Errorf("%s", errMsg)
	} else {

		return userAccounts[0], nil
	}
}

func (s *userAccountImpl) GetUserAccountsByUuid(ctx context.Context, userUuid uuid.UUID) (model.UserAccounts, error) {
	filter := &model.UserAccountFilter{
		ID:              nil,
		UserID:          nil,
		AccountTypeID:   nil,
		UserIDNew:       &userUuid,
		CoinNew:         nil,
		IsActive:        nil,
		Pagination:      nil,
		UserIDNewIsNull: nil,
		CoinNewIsNull:   nil,
	}

	if _, userAccounts, err := s.GetUserAccountsByFilter(ctx, filter); err != nil {

		return nil, err
	} else {

		return userAccounts, err
	}
}
