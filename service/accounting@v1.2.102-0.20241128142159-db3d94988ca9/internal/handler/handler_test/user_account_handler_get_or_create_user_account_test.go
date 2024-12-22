package handler_test

import (
	"code.emcdtech.com/emcd/service/accounting/internal/handler"
	"code.emcdtech.com/emcd/service/accounting/internal/handler/mapping"
	userAccountServiceMock "code.emcdtech.com/emcd/service/accounting/mocks/internal_/service"
	"code.emcdtech.com/emcd/service/accounting/model/enum"
	userAccountPb "code.emcdtech.com/emcd/service/accounting/protocol/user_account"
	"code.emcdtech.com/emcd/service/accounting/repository"
	coinValidatorRepoMock "code.emcdtech.com/emcd/service/coin/repository/mocks"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUserAccountHandlerService_GetOrCreateUserAccount(t *testing.T) {
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

	t.Run("success", func(t *testing.T) {
		req := &userAccountPb.UserAccountOneRequest{
			UserId:        userId,
			UserIdNew:     userIdNew.String(),
			AccountTypeId: accountTypeId.ToInt32(),
			Minpay:        0,
			Address:       nil,
			Img1:          nil,
			Img2:          nil,
			Fee:           nil,
			CoinNew:       coinNew,
		}

		coinValidatorMock.EXPECT().
			GetIdByCode(coinNew).
			Return(coinIdLegacy, true).Twice()

		userAccount, errMapping := mapping.MapProtoOneRequestToModelUserAccount(coinValidatorMock, req)
		require.NoError(t, errMapping)
		require.NotNil(t, userAccount)

		userAccountCli.EXPECT().
			GetOrCreateUserAccount(ctx, userAccount).
			Return(userAccount, nil).Once()

		resp, err := userAccountHandler.GetOrCreateUserAccount(ctx, req)
		require.NotNil(t, resp)
		require.NoError(t, err)

	})

	t.Run("error mapping uuid", func(t *testing.T) {
		req := &userAccountPb.UserAccountOneRequest{
			UserId:        userId,
			UserIdNew:     userIdNew.String() + "fake",
			AccountTypeId: accountTypeId.ToInt32(),
			Minpay:        0,
			Address:       nil,
			Img1:          nil,
			Img2:          nil,
			Fee:           nil,
			CoinNew:       coinNew,
		}

		resp, err := userAccountHandler.GetOrCreateUserAccount(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAcc1021)

	})

	t.Run("error mapping account type", func(t *testing.T) {
		req := &userAccountPb.UserAccountOneRequest{
			UserId:        userId,
			UserIdNew:     userIdNew.String(),
			AccountTypeId: accountTypeId.ToInt32() + 99,
			Minpay:        0,
			Address:       nil,
			Img1:          nil,
			Img2:          nil,
			Fee:           nil,
			CoinNew:       coinNew,
		}

		resp, err := userAccountHandler.GetOrCreateUserAccount(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAcc1021)
	})

	t.Run("error mock coin validate", func(t *testing.T) {
		req := &userAccountPb.UserAccountOneRequest{
			UserId:        userId,
			UserIdNew:     userIdNew.String(),
			AccountTypeId: accountTypeId.ToInt32(),
			Minpay:        0,
			Address:       nil,
			Img1:          nil,
			Img2:          nil,
			Fee:           nil,
			CoinNew:       coinNew,
		}

		coinValidatorMock.EXPECT().
			GetIdByCode(coinNew).
			Return(0, false).Once()

		resp, err := userAccountHandler.GetOrCreateUserAccount(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAcc1021)

	})

	t.Run("error mock repo", func(t *testing.T) {
		req := &userAccountPb.UserAccountOneRequest{
			UserId:        userId,
			UserIdNew:     userIdNew.String(),
			AccountTypeId: accountTypeId.ToInt32(),
			Minpay:        0,
			Address:       nil,
			Img1:          nil,
			Img2:          nil,
			Fee:           nil,
			CoinNew:       coinNew,
		}

		coinValidatorMock.EXPECT().
			GetIdByCode(coinNew).
			Return(coinIdLegacy, true).Twice()

		userAccount, errMapping := mapping.MapProtoOneRequestToModelUserAccount(coinValidatorMock, req)
		require.NoError(t, errMapping)
		require.NotNil(t, userAccount)

		userAccountCli.EXPECT().
			GetOrCreateUserAccount(ctx, userAccount).
			Return(nil, errors.New("error mock repo"))

		resp, err := userAccountHandler.GetOrCreateUserAccount(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAcc1022)

	})
}
