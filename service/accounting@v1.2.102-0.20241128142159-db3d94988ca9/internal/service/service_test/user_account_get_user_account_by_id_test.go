package service_test

import (
	"code.emcdtech.com/emcd/service/accounting/internal/service"
	"code.emcdtech.com/emcd/service/accounting/internal/utils"
	repositoryMock "code.emcdtech.com/emcd/service/accounting/mocks/internal_/repository"
	"code.emcdtech.com/emcd/service/accounting/model"
	"code.emcdtech.com/emcd/service/accounting/model/enum"
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUserAccountService_GetUserAccountById(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		userAccountRepo := repositoryMock.NewMockUserAccountRepo(t)
		userAccountService := service.NewUserAccountService(userAccountRepo)

		userAccountId := int32(1)
		userAccountTypeId := enum.WalletAccountTypeID

		filter := &model.UserAccountFilter{
			ID:              &userAccountId,
			UserID:          nil,
			AccountTypeID:   nil,
			UserIDNew:       nil,
			CoinNew:         nil,
			Pagination:      nil,
			UserIDNewIsNull: nil,
			CoinNewIsNull:   nil,
		}

		userAccountsModel := model.UserAccounts{{
			ID:            userAccountId,
			UserID:        1,
			CoinID:        1,
			AccountTypeID: enum.NewAccountTypeIdWrapper(userAccountTypeId),
			Minpay:        0,
			Address:       sql.NullString{},
			ChangedAt:     sql.NullTime{},
			Img1:          sql.NullFloat64{},
			Img2:          sql.NullFloat64{},
			IsActive:      sql.NullBool{},
			CreatedAt:     sql.NullTime{},
			UpdatedAt:     sql.NullTime{},
			Fee:           sql.NullFloat64{},
			UserIDNew:     utils.UuidToUuidNull(uuid.UUID{}),
			CoinNew:       utils.StringToStringNull("btc"),
		},
		}

		userAccountRepo.EXPECT().
			FindUserAccountByFilter(ctx, filter).
			Return(nil, userAccountsModel, nil)

		result, err := userAccountService.GetUserAccountById(ctx, userAccountId)
		require.NotNil(t, result)
		require.NoError(t, err)

		require.Equal(t, userAccountId, result.ID)

	})
	t.Run("error many return", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		userAccountRepo := repositoryMock.NewMockUserAccountRepo(t)
		userAccountService := service.NewUserAccountService(userAccountRepo)

		userAccountId := int32(1)

		filter := &model.UserAccountFilter{
			ID:              &userAccountId,
			UserID:          nil,
			AccountTypeID:   nil,
			UserIDNew:       nil,
			CoinNew:         nil,
			Pagination:      nil,
			UserIDNewIsNull: nil,
			CoinNewIsNull:   nil,
		}

		userAccountsModel := model.UserAccounts{
			{
				ID:            userAccountId,
				UserID:        1,
				CoinID:        1,
				AccountTypeID: enum.NewAccountTypeIdWrapper(enum.WalletAccountTypeID),
				Minpay:        0,
				Address:       sql.NullString{},
				ChangedAt:     sql.NullTime{},
				Img1:          sql.NullFloat64{},
				Img2:          sql.NullFloat64{},
				IsActive:      sql.NullBool{},
				CreatedAt:     sql.NullTime{},
				UpdatedAt:     sql.NullTime{},
				Fee:           sql.NullFloat64{},
				UserIDNew:     utils.UuidToUuidNull(uuid.UUID{}),
				CoinNew:       utils.StringToStringNull("btc"),
			},
			{
				ID:            userAccountId,
				UserID:        1,
				CoinID:        1,
				AccountTypeID: enum.NewAccountTypeIdWrapper(enum.WalletAccountTypeID),
				Minpay:        0,
				Address:       sql.NullString{},
				ChangedAt:     sql.NullTime{},
				Img1:          sql.NullFloat64{},
				Img2:          sql.NullFloat64{},
				IsActive:      sql.NullBool{},
				CreatedAt:     sql.NullTime{},
				UpdatedAt:     sql.NullTime{},
				Fee:           sql.NullFloat64{},
				UserIDNew:     utils.UuidToUuidNull(uuid.UUID{}),
				CoinNew:       utils.StringToStringNull("btc"),
			},
		}

		userAccountRepo.EXPECT().
			FindUserAccountByFilter(ctx, filter).
			Return(nil, userAccountsModel, nil)

		result, err := userAccountService.GetUserAccountById(ctx, userAccountId)
		require.Nil(t, result)
		require.Error(t, err)

		require.EqualError(t, err, "user account is not uniq id")

	})

	t.Run("error empty return", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		userAccountRepo := repositoryMock.NewMockUserAccountRepo(t)
		userAccountService := service.NewUserAccountService(userAccountRepo)

		userAccountId := int32(1)

		filter := &model.UserAccountFilter{
			ID:              &userAccountId,
			UserID:          nil,
			AccountTypeID:   nil,
			UserIDNew:       nil,
			CoinNew:         nil,
			Pagination:      nil,
			UserIDNewIsNull: nil,
			CoinNewIsNull:   nil,
		}

		userAccountsModel := model.UserAccounts{}

		userAccountRepo.EXPECT().
			FindUserAccountByFilter(ctx, filter).
			Return(nil, userAccountsModel, nil)

		result, err := userAccountService.GetUserAccountById(ctx, userAccountId)
		require.Nil(t, result)
		require.Error(t, err)

		require.EqualError(t, err, "user account is not found")

	})

	t.Run("error mock return", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		userAccountRepo := repositoryMock.NewMockUserAccountRepo(t)
		userAccountService := service.NewUserAccountService(userAccountRepo)

		userAccountId := int32(1)

		filter := &model.UserAccountFilter{
			ID:              &userAccountId,
			UserID:          nil,
			AccountTypeID:   nil,
			UserIDNew:       nil,
			CoinNew:         nil,
			Pagination:      nil,
			UserIDNewIsNull: nil,
			CoinNewIsNull:   nil,
		}

		retError := newMockError()

		userAccountRepo.EXPECT().
			FindUserAccountByFilter(ctx, filter).
			Return(nil, nil, retError)

		result, err := userAccountService.GetUserAccountById(ctx, userAccountId)
		require.Nil(t, result)
		require.Error(t, err)

		require.ErrorIs(t, err, retError)

	})
}
