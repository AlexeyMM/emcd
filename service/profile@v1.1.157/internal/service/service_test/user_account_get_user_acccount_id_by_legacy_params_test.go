package service_test

/*
import (
	userAccountModel "code.emcdtech.com/emcd/service/accounting/model"
	userAccountModelEnum "code.emcdtech.com/emcd/service/accounting/model/enum"
	userAccountRepositoryMock "code.emcdtech.com/emcd/service/accounting/repository/repository_mock"
	coinValidatorRepoMock "code.emcdtech.com/emcd/service/coin/repository/mocks"
	"code.emcdtech.com/emcd/service/profile/internal/service"
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)


func TestAccountingUserAccountService_GetUserAccountByLegacyParams(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		coinValidatorMock := coinValidatorRepoMock.NewMockCoinValidatorRepository(t)
		userAccountRepoMock := userAccountRepositoryMock.NewMockUserAccountRepository(t)
		userAccountService := service.NewUserAccountService(userAccountRepoMock, coinValidatorMock)

		userId := int32(1)
		coinId := int32(1)
		coinCode := "btc"
		accountTypeId := userAccountModelEnum.WalletAccountTypeID

		filter := &userAccountModel.UserAccountFilter{
			ID:              nil,
			UserID:          &userId,
			AccountTypeID:   accountTypeId.ToPtr(),
			UserIDNew:       nil,
			CoinNew:         &coinCode,
			Pagination:      nil,
			UserIDNewIsNull: nil,
			CoinNewIsNull:   nil,
		}

		userAccountsResponse := userAccountModel.UserAccounts{{
			ID:            99,
			UserID:        userId,
			CoinID:        0,
			AccountTypeID: userAccountModelEnum.NewAccountTypeIdWrapper(userAccountModelEnum.WalletAccountTypeID),
			Minpay:        0,
			Address:       sql.NullString{},
			ChangedAt:     sql.NullTime{},
			Img1:          sql.NullFloat64{},
			Img2:          sql.NullFloat64{},
			IsActive:      sql.NullBool{},
			CreatedAt:     sql.NullTime{},
			UpdatedAt:     sql.NullTime{},
			Fee:           sql.NullFloat64{},
			UserIDNew:     uuid.NullUUID{},
			CoinNew:       sql.NullString{},
		}}

		coinValidatorMock.EXPECT().
			GetCodeById(coinId).
			Return(coinCode, true)

		userAccountRepoMock.EXPECT().GetUserAccountsByFilter(
			mock.Anything,
			filter,
		).Return(nil, userAccountsResponse, nil)

		resp, err := userAccountService.GetUserAccountIdByLegacyParams(ctx, userId, coinId, accountTypeId.ToInt32())
		require.NotEmpty(t, resp)
		require.NoError(t, err)

		require.Equal(t, resp, 99)

	})

	t.Run("error validate coin", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		coinValidatorMock := coinValidatorRepoMock.NewMockCoinValidatorRepository(t)
		userAccountRepoMock := userAccountRepositoryMock.NewMockUserAccountRepository(t)
		userAccountService := service.NewUserAccountService(userAccountRepoMock, coinValidatorMock)

		userId := int32(1)
		coinId := int32(1)
		accountTypeId := userAccountModelEnum.WalletAccountTypeID

		coinValidatorMock.EXPECT().GetCodeById(coinId).
			Return("", false)

		resp, err := userAccountService.GetUserAccountIdByLegacyParams(ctx, userId, coinId, accountTypeId.ToInt32())
		require.Empty(t, resp)
		require.Error(t, err)

		require.ErrorContains(t, err, "unknown coin id")

	})

	t.Run("error validate account type", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		coinValidatorMock := coinValidatorRepoMock.NewMockCoinValidatorRepository(t)
		userAccountRepoMock := userAccountRepositoryMock.NewMockUserAccountRepository(t)
		userAccountService := service.NewUserAccountService(userAccountRepoMock, coinValidatorMock)

		userId := int32(1)
		coinId := int32(1)
		coinCode := "btc"
		accountTypeId := userAccountModelEnum.WalletAccountTypeID + 99

		coinValidatorMock.EXPECT().
			GetCodeById(coinId).
			Return(coinCode, true)

		resp, err := userAccountService.GetUserAccountIdByLegacyParams(ctx, userId, coinId, accountTypeId.ToInt32())
		require.Empty(t, resp)
		require.Error(t, err)

		require.ErrorContains(t, err, "unknown account type")

	})

	t.Run("error return empty", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		coinValidatorMock := coinValidatorRepoMock.NewMockCoinValidatorRepository(t)
		userAccountRepoMock := userAccountRepositoryMock.NewMockUserAccountRepository(t)
		userAccountService := service.NewUserAccountService(userAccountRepoMock, coinValidatorMock)

		userId := int32(1)
		coinId := int32(1)
		coinCode := "btc"
		accountTypeId := userAccountModelEnum.WalletAccountTypeID

		filter := &userAccountModel.UserAccountFilter{
			ID:              nil,
			UserID:          &userId,
			AccountTypeID:   accountTypeId.ToPtr(),
			UserIDNew:       nil,
			CoinNew:         &coinCode,
			Pagination:      nil,
			UserIDNewIsNull: nil,
			CoinNewIsNull:   nil,
		}

		coinValidatorMock.EXPECT().
			GetCodeById(coinId).
			Return(coinCode, true)

		userAccountRepoMock.EXPECT().GetUserAccountsByFilter(
			mock.Anything,
			filter,
		).Return(nil, userAccountModel.UserAccounts{}, nil)

		resp, err := userAccountService.GetUserAccountIdByLegacyParams(ctx, userId, coinId, accountTypeId.ToInt32())
		require.Empty(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, service.ErrUserAccountListIsEmpty)

	})

	t.Run("error return many", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		coinValidatorMock := coinValidatorRepoMock.NewMockCoinValidatorRepository(t)
		userAccountRepoMock := userAccountRepositoryMock.NewMockUserAccountRepository(t)
		userAccountService := service.NewUserAccountService(userAccountRepoMock, coinValidatorMock)

		userId := int32(1)
		coinId := int32(1)
		coinCode := "btc"
		accountTypeId := userAccountModelEnum.WalletAccountTypeID

		filter := &userAccountModel.UserAccountFilter{
			ID:              nil,
			UserID:          &userId,
			AccountTypeID:   accountTypeId.ToPtr(),
			UserIDNew:       nil,
			CoinNew:         &coinCode,
			Pagination:      nil,
			UserIDNewIsNull: nil,
			CoinNewIsNull:   nil,
		}

		coinValidatorMock.EXPECT().
			GetCodeById(coinId).
			Return(coinCode, true)

		userAccountRepoMock.EXPECT().GetUserAccountsByFilter(
			mock.Anything,
			filter,
		).Return(nil, userAccountModel.UserAccounts{{}, {}}, nil)

		resp, err := userAccountService.GetUserAccountIdByLegacyParams(ctx, userId, coinId, accountTypeId.ToInt32())
		require.Empty(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, service.ErrUserAccountListMoreThenOne)

	})
}
*/
