package handler_test

import (
	sdkError "code.emcdtech.com/emcd/sdk/error"
	"code.emcdtech.com/emcd/service/accounting/internal/handler"
	"code.emcdtech.com/emcd/service/accounting/internal/handler/mapping"
	userAccountServiceMock "code.emcdtech.com/emcd/service/accounting/mocks/internal_/service"
	"code.emcdtech.com/emcd/service/accounting/model/enum"
	userAccountPb "code.emcdtech.com/emcd/service/accounting/protocol/user_account"
	coinValidatorRepoMock "code.emcdtech.com/emcd/service/coin/repository/mocks"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUserAccountHandlerService_CreateUserAccounts(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		userAccountCli := userAccountServiceMock.NewMockUserAccountService(t)
		coinValidatorMock := coinValidatorRepoMock.NewMockCoinValidatorRepository(t)
		userAccountHandler := handler.NewUserAccountHandler(userAccountCli, coinValidatorMock)

		userId := int32(1)
		userIdNew := uuid.UUID{}
		coinIdLegacy := int32(1)
		coinNew := "btc"
		accountTypeId := enum.WalletAccountTypeID

		req := &userAccountPb.UserAccountMultiRequest{
			UserId:    userId,
			UserIdNew: userIdNew.String(),
			UserAccounts: []*userAccountPb.UserAccountRequest{
				{
					AccountTypeId: accountTypeId.ToInt32(),
					Minpay:        0,
					Address:       nil,
					Img1:          nil,
					Img2:          nil,
					Fee:           nil,
					CoinNew:       coinNew,
				},
			},
		}

		coinValidatorMock.EXPECT().
			GetIdByCode(coinNew).
			Return(coinIdLegacy, true)

		userIdMapped, userIdNewMapped, userAccounts, errMapping := mapping.MapProtoMultiRequestToModelUserAccounts(coinValidatorMock, req)
		require.NoError(t, errMapping)
		require.Equal(t, userId, userIdMapped)
		require.Equal(t, userIdNew.String(), userIdNewMapped.String())

		userAccountCli.EXPECT().
			CreateUserAccounts(ctx, userIdMapped, userIdNewMapped, userAccounts).
			Return(userAccounts, nil)

		resp, err := userAccountHandler.CreateUserAccounts(ctx, req)
		require.NotNil(t, resp)
		require.NoError(t, err)

		require.Len(t, resp.UserAccounts, len(userAccounts))

	})

	t.Run("error mapping uuid", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		userAccountCli := userAccountServiceMock.NewMockUserAccountService(t)
		coinValidatorMock := coinValidatorRepoMock.NewMockCoinValidatorRepository(t)
		userAccountHandler := handler.NewUserAccountHandler(userAccountCli, coinValidatorMock)

		userId := int32(1)
		userIdNew := uuid.UUID{}
		coinNew := "btc"
		accountTypeId := enum.WalletAccountTypeID

		req := &userAccountPb.UserAccountMultiRequest{
			UserId:    userId,
			UserIdNew: userIdNew.String() + "fake",
			UserAccounts: []*userAccountPb.UserAccountRequest{
				{
					AccountTypeId: accountTypeId.ToInt32(),
					Minpay:        0,
					Address:       nil,
					Img1:          nil,
					Img2:          nil,
					Fee:           nil,
					CoinNew:       coinNew,
				},
			},
		}

		resp, err := userAccountHandler.CreateUserAccounts(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		// require.Error(t, errMapping)
		require.EqualError(t, err, sdkError.NewError("acc-1011", "failed mapping proto request").Error())

	})

	t.Run("error mapping account type", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		userAccountCli := userAccountServiceMock.NewMockUserAccountService(t)
		coinValidatorMock := coinValidatorRepoMock.NewMockCoinValidatorRepository(t)
		userAccountHandler := handler.NewUserAccountHandler(userAccountCli, coinValidatorMock)

		userId := int32(1)
		userIdNew := uuid.UUID{}
		coinNew := "btc"
		accountTypeId := enum.WalletAccountTypeID

		req := &userAccountPb.UserAccountMultiRequest{
			UserId:    userId,
			UserIdNew: userIdNew.String(),
			UserAccounts: []*userAccountPb.UserAccountRequest{
				{
					AccountTypeId: accountTypeId.ToInt32() + 99,
					Minpay:        0,
					Address:       nil,
					Img1:          nil,
					Img2:          nil,
					Fee:           nil,
					CoinNew:       coinNew,
				},
			},
		}

		resp, err := userAccountHandler.CreateUserAccounts(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.EqualError(t, err, sdkError.NewError("acc-1011", "failed mapping proto request").Error())
	})

	t.Run("error mock coin validate", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		userAccountCli := userAccountServiceMock.NewMockUserAccountService(t)
		coinValidatorMock := coinValidatorRepoMock.NewMockCoinValidatorRepository(t)
		userAccountHandler := handler.NewUserAccountHandler(userAccountCli, coinValidatorMock)

		userId := int32(1)
		userIdNew := uuid.UUID{}
		coinNew := "btc"
		accountTypeId := enum.WalletAccountTypeID

		req := &userAccountPb.UserAccountMultiRequest{
			UserId:    userId,
			UserIdNew: userIdNew.String(),
			UserAccounts: []*userAccountPb.UserAccountRequest{
				{
					AccountTypeId: accountTypeId.ToInt32(),
					Minpay:        0,
					Address:       nil,
					Img1:          nil,
					Img2:          nil,
					Fee:           nil,
					CoinNew:       coinNew,
				},
			},
		}

		coinValidatorMock.EXPECT().
			GetIdByCode(coinNew).
			Return(0, false)

		resp, err := userAccountHandler.CreateUserAccounts(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.EqualError(t, err, sdkError.NewError("acc-1011", "failed mapping proto request").Error())

	})

	t.Run("error mock repo", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		userAccountCli := userAccountServiceMock.NewMockUserAccountService(t)
		coinValidatorMock := coinValidatorRepoMock.NewMockCoinValidatorRepository(t)
		userAccountHandler := handler.NewUserAccountHandler(userAccountCli, coinValidatorMock)

		userId := int32(1)
		userIdNew := uuid.UUID{}
		coinNew := "btc"
		coinIdLegacy := int32(1)
		accountTypeId := enum.WalletAccountTypeID

		req := &userAccountPb.UserAccountMultiRequest{
			UserId:    userId,
			UserIdNew: userIdNew.String(),
			UserAccounts: []*userAccountPb.UserAccountRequest{
				{
					AccountTypeId: accountTypeId.ToInt32(),
					Minpay:        0,
					Address:       nil,
					Img1:          nil,
					Img2:          nil,
					Fee:           nil,
					CoinNew:       coinNew,
				},
			},
		}

		coinValidatorMock.EXPECT().
			GetIdByCode(coinNew).
			Return(coinIdLegacy, true)

		userIdMapped, userIdNewMapped, userAccounts, errMapping := mapping.MapProtoMultiRequestToModelUserAccounts(coinValidatorMock, req)
		require.NoError(t, errMapping)
		require.Equal(t, userId, userIdMapped)
		require.Equal(t, userIdNew.String(), userIdNewMapped.String())

		userAccountCli.EXPECT().
			CreateUserAccounts(ctx, userIdMapped, userIdNewMapped, userAccounts).
			Return(nil, errors.New("error mock repo"))

		resp, err := userAccountHandler.CreateUserAccounts(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.EqualError(t, err, sdkError.NewError("acc-1012", "failed create user accounts").Error())

	})
}
