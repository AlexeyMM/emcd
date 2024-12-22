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

func TestUserAccountServiceGet_UserAccountByConstraint(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		userAccountRepo := repositoryMock.NewMockUserAccountRepo(t)
		userAccountService := service.NewUserAccountService(userAccountRepo)

		userIdNew := uuid.UUID{}
		coinId := int32(1)
		coinNew := "btc"
		userAccountID := enum.WalletAccountTypeID

		filter := &model.UserAccountFilter{
			ID:              nil,
			UserID:          nil,
			AccountTypeID:   userAccountID.ToPtr(),
			UserIDNew:       &userIdNew,
			CoinNew:         utils.StringToPtr(coinNew),
			Pagination:      nil,
			UserIDNewIsNull: nil,
			CoinNewIsNull:   nil,
		}

		userAccountsModel := model.UserAccounts{{
			ID:            0,
			UserID:        int32(1),
			CoinID:        coinId,
			AccountTypeID: enum.NewAccountTypeIdWrapper(userAccountID),
			Minpay:        0,
			Address:       sql.NullString{},
			ChangedAt:     sql.NullTime{},
			Img1:          sql.NullFloat64{},
			Img2:          sql.NullFloat64{},
			IsActive:      sql.NullBool{},
			CreatedAt:     sql.NullTime{},
			UpdatedAt:     sql.NullTime{},
			Fee:           sql.NullFloat64{},
			UserIDNew:     utils.UuidToUuidNull(userIdNew),
			CoinNew:       utils.StringToStringNull(coinNew),
		},
		}

		userAccountRepo.EXPECT().
			FindUserAccountByFilter(ctx, filter).
			Return(nil, userAccountsModel, nil)

		result, err := userAccountService.GetUserAccountByConstraint(ctx, userIdNew, coinNew, userAccountID)
		require.NotNil(t, result)
		require.NoError(t, err)

	})

	t.Run("error many return", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		userAccountRepo := repositoryMock.NewMockUserAccountRepo(t)
		userAccountService := service.NewUserAccountService(userAccountRepo)

		userIdNew := uuid.UUID{}
		coinId := int32(1)
		coinNew := "btc"
		userAccountID := enum.WalletAccountTypeID

		filter := &model.UserAccountFilter{
			ID:              nil,
			UserID:          nil,
			AccountTypeID:   userAccountID.ToPtr(),
			UserIDNew:       &userIdNew,
			CoinNew:         utils.StringToPtr(coinNew),
			Pagination:      nil,
			UserIDNewIsNull: nil,
			CoinNewIsNull:   nil,
		}

		userAccountsModel := model.UserAccounts{
			{
				ID:            1,
				UserID:        1,
				CoinID:        coinId,
				AccountTypeID: enum.NewAccountTypeIdWrapper(userAccountID),
				Minpay:        0,
				Address:       sql.NullString{},
				ChangedAt:     sql.NullTime{},
				Img1:          sql.NullFloat64{},
				Img2:          sql.NullFloat64{},
				IsActive:      sql.NullBool{},
				CreatedAt:     sql.NullTime{},
				UpdatedAt:     sql.NullTime{},
				Fee:           sql.NullFloat64{},
				UserIDNew:     utils.UuidToUuidNull(userIdNew),
				CoinNew:       utils.StringToStringNull(coinNew),
			},
			{
				ID:            1,
				UserID:        1,
				CoinID:        coinId,
				AccountTypeID: enum.NewAccountTypeIdWrapper(userAccountID),
				Minpay:        0,
				Address:       sql.NullString{},
				ChangedAt:     sql.NullTime{},
				Img1:          sql.NullFloat64{},
				Img2:          sql.NullFloat64{},
				IsActive:      sql.NullBool{},
				CreatedAt:     sql.NullTime{},
				UpdatedAt:     sql.NullTime{},
				Fee:           sql.NullFloat64{},
				UserIDNew:     utils.UuidToUuidNull(userIdNew),
				CoinNew:       utils.StringToStringNull(coinNew),
			},
		}

		userAccountRepo.EXPECT().
			FindUserAccountByFilter(ctx, filter).
			Return(nil, userAccountsModel, nil)

		result, err := userAccountService.GetUserAccountByConstraint(ctx, userIdNew, coinNew, userAccountID)
		require.Nil(t, result)
		require.Error(t, err)

		require.EqualError(t, err, "user account is not uniq")

	})

	t.Run("error empty return", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		userAccountRepo := repositoryMock.NewMockUserAccountRepo(t)
		userAccountService := service.NewUserAccountService(userAccountRepo)

		userIdNew := uuid.UUID{}
		coinNew := "btc"
		userAccountID := enum.WalletAccountTypeID

		filter := &model.UserAccountFilter{
			ID:              nil,
			UserID:          nil,
			AccountTypeID:   userAccountID.ToPtr(),
			UserIDNew:       &userIdNew,
			CoinNew:         utils.StringToPtr(coinNew),
			Pagination:      nil,
			UserIDNewIsNull: nil,
			CoinNewIsNull:   nil,
		}

		userAccountsModel := model.UserAccounts{}

		userAccountRepo.EXPECT().
			FindUserAccountByFilter(ctx, filter).
			Return(nil, userAccountsModel, nil)

		result, err := userAccountService.GetUserAccountByConstraint(ctx, userIdNew, coinNew, userAccountID)
		require.Nil(t, result)
		require.Error(t, err)

		require.EqualError(t, err, "user account is not found")

	})

	t.Run("error mock return", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		userAccountRepo := repositoryMock.NewMockUserAccountRepo(t)
		userAccountService := service.NewUserAccountService(userAccountRepo)

		userIdNew := uuid.UUID{}
		coinNew := "btc"
		userAccountID := enum.WalletAccountTypeID

		filter := &model.UserAccountFilter{
			ID:              nil,
			UserID:          nil,
			AccountTypeID:   userAccountID.ToPtr(),
			UserIDNew:       &userIdNew,
			CoinNew:         utils.StringToPtr(coinNew),
			Pagination:      nil,
			UserIDNewIsNull: nil,
			CoinNewIsNull:   nil,
		}

		retError := newMockError()

		userAccountRepo.EXPECT().
			FindUserAccountByFilter(ctx, filter).
			Return(nil, nil, retError)

		result, err := userAccountService.GetUserAccountByConstraint(ctx, userIdNew, coinNew, userAccountID)
		require.Nil(t, result)
		require.Error(t, err)

		require.ErrorIs(t, err, retError)

	})
}
