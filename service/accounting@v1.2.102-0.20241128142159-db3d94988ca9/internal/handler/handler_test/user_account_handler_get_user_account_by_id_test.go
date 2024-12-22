package handler_test

import (
	"code.emcdtech.com/emcd/service/accounting/internal/handler"
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

func TestUserAccountHandler_ServiceGetUserAccountById(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		userAccountCli := userAccountServiceMock.NewMockUserAccountService(t)
		coinValidatorMock := coinValidatorRepoMock.NewMockCoinValidatorRepository(t)
		userAccountHandler := handler.NewUserAccountHandler(userAccountCli, coinValidatorMock)

		userAccountId := int32(1)

		req := &userAccountPb.UserAccountId{
			Id: userAccountId,
		}

		retUserAccount := &model.UserAccount{
			ID:            0,
			UserID:        0,
			CoinID:        0,
			AccountTypeID: enum.AccountTypeIdWrapper{},
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
			CoinNew:       utils.StringToStringNull(""),
		}

		userAccountCli.EXPECT().
			GetUserAccountById(ctx, userAccountId).
			Return(retUserAccount, nil)

		resp, err := userAccountHandler.GetUserAccountById(ctx, req)
		require.NotNil(t, resp)
		require.NoError(t, err)

	})

	t.Run("error mock repo", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		userAccountCli := userAccountServiceMock.NewMockUserAccountService(t)
		coinValidatorMock := coinValidatorRepoMock.NewMockCoinValidatorRepository(t)
		userAccountHandler := handler.NewUserAccountHandler(userAccountCli, coinValidatorMock)

		userAccountId := int32(1)

		req := &userAccountPb.UserAccountId{
			Id: userAccountId,
		}

		userAccountCli.EXPECT().
			GetUserAccountById(ctx, userAccountId).
			Return(nil, errors.New("mock error repo"))

		resp, err := userAccountHandler.GetUserAccountById(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAcc1041)

	})
}
