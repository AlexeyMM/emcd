package handler_test

import (
	"code.emcdtech.com/emcd/service/accounting/internal/handler"
	"code.emcdtech.com/emcd/service/accounting/internal/handler/mapping"
	"code.emcdtech.com/emcd/service/accounting/internal/utils"
	userAccountServiceMock "code.emcdtech.com/emcd/service/accounting/mocks/internal_/service"
	"code.emcdtech.com/emcd/service/accounting/model"
	"code.emcdtech.com/emcd/service/accounting/model/enum"
	userAccountPb "code.emcdtech.com/emcd/service/accounting/protocol/user_account"
	"code.emcdtech.com/emcd/service/accounting/repository"
	coinValidatorRepoMock "code.emcdtech.com/emcd/service/coin/repository/mocks"
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUserAccountHandlerService_GetUserAccountsByFilter(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		userAccountCli := userAccountServiceMock.NewMockUserAccountService(t)
		coinValidatorMock := coinValidatorRepoMock.NewMockCoinValidatorRepository(t)
		userAccountHandler := handler.NewUserAccountHandler(userAccountCli, coinValidatorMock)

		userIdNew := uuid.UUID{}
		coinNew := "btc"
		accountTypeId := enum.WalletAccountTypeID

		req := &userAccountPb.UserAccountFilter{
			Id:            nil,
			UserId:        nil,
			AccountTypeId: utils.Int32ToPtr(accountTypeId.ToInt32()),
			UserIdNew:     utils.StringToPtr(userIdNew.String()),
			CoinNew:       utils.StringToPtr(coinNew),
			Pagination:    nil,
		}

		retUserAccounts := model.UserAccounts{&model.UserAccount{
			ID:            0,
			UserID:        0,
			CoinID:        0,
			AccountTypeID: enum.NewAccountTypeIdWrapper(accountTypeId),
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

		coinValidatorMock.EXPECT().
			IsValidCode(coinNew).
			Return(true)

		filter, err := mapping.MapProtoToModelUserAccountFilter(coinValidatorMock, req)

		require.NoError(t, err)

		userAccountCli.EXPECT().
			GetUserAccountsByFilter(ctx, filter).
			Return(nil, retUserAccounts, nil)

		resp, err := userAccountHandler.GetUserAccountsByFilter(ctx, req)
		require.NotNil(t, resp)
		require.NoError(t, err)

	})

	t.Run("error parse account type", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		userAccountCli := userAccountServiceMock.NewMockUserAccountService(t)
		coinValidatorMock := coinValidatorRepoMock.NewMockCoinValidatorRepository(t)
		userAccountHandler := handler.NewUserAccountHandler(userAccountCli, coinValidatorMock)

		userIdNew := uuid.UUID{}
		coinNew := "btc"
		accountTypeId := enum.WalletAccountTypeID

		req := &userAccountPb.UserAccountFilter{
			Id:            nil,
			UserId:        nil,
			AccountTypeId: utils.Int32ToPtr(accountTypeId.ToInt32() + 99),
			UserIdNew:     utils.StringToPtr(userIdNew.String()),
			CoinNew:       utils.StringToPtr(coinNew),
			Pagination:    nil,
		}

		resp, err := userAccountHandler.GetUserAccountsByFilter(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAcc1031)

	})

	t.Run("error parse uuid", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		userAccountCli := userAccountServiceMock.NewMockUserAccountService(t)
		coinValidatorMock := coinValidatorRepoMock.NewMockCoinValidatorRepository(t)
		userAccountHandler := handler.NewUserAccountHandler(userAccountCli, coinValidatorMock)

		userIdNew := uuid.UUID{}
		coinNew := "btc"
		accountTypeId := enum.WalletAccountTypeID

		req := &userAccountPb.UserAccountFilter{
			Id:            nil,
			UserId:        nil,
			AccountTypeId: utils.Int32ToPtr(accountTypeId.ToInt32()),
			UserIdNew:     utils.StringToPtr(userIdNew.String() + "fake"),
			CoinNew:       utils.StringToPtr(coinNew),
			Pagination:    nil,
		}

		resp, err := userAccountHandler.GetUserAccountsByFilter(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAcc1031)

	})

	t.Run("error mock coin validate", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		userAccountCli := userAccountServiceMock.NewMockUserAccountService(t)
		coinValidatorMock := coinValidatorRepoMock.NewMockCoinValidatorRepository(t)
		userAccountHandler := handler.NewUserAccountHandler(userAccountCli, coinValidatorMock)

		userIdNew := uuid.UUID{}
		coinNew := "btc"
		accountTypeId := enum.WalletAccountTypeID

		req := &userAccountPb.UserAccountFilter{
			Id:            nil,
			UserId:        nil,
			AccountTypeId: utils.Int32ToPtr(accountTypeId.ToInt32()),
			UserIdNew:     utils.StringToPtr(userIdNew.String()),
			CoinNew:       utils.StringToPtr(coinNew),
			Pagination:    nil,
		}

		coinValidatorMock.EXPECT().
			IsValidCode(coinNew).
			Return(false)

		resp, err := userAccountHandler.GetUserAccountsByFilter(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAcc1031)

	})

	t.Run("error mock repo", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		userAccountCli := userAccountServiceMock.NewMockUserAccountService(t)
		coinValidatorMock := coinValidatorRepoMock.NewMockCoinValidatorRepository(t)
		userAccountHandler := handler.NewUserAccountHandler(userAccountCli, coinValidatorMock)

		userIdNew := uuid.UUID{}
		coinNew := "btc"
		accountTypeId := enum.WalletAccountTypeID

		req := &userAccountPb.UserAccountFilter{
			Id:            nil,
			UserId:        nil,
			AccountTypeId: utils.Int32ToPtr(accountTypeId.ToInt32()),
			UserIdNew:     utils.StringToPtr(userIdNew.String()),
			CoinNew:       utils.StringToPtr(coinNew),
			Pagination:    nil,
		}

		coinValidatorMock.EXPECT().
			IsValidCode(coinNew).
			Return(true)

		filter, err := mapping.MapProtoToModelUserAccountFilter(coinValidatorMock, req)

		require.NoError(t, err)

		userAccountCli.EXPECT().
			GetUserAccountsByFilter(ctx, filter).
			Return(nil, nil, errors.New("error mock repo"))

		resp, err := userAccountHandler.GetUserAccountsByFilter(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAcc1032)

	})
}
