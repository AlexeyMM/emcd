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

func TestUserAccountHandlerService_GetUserAccountsByUuid(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		userAccountCli := userAccountServiceMock.NewMockUserAccountService(t)
		coinValidatorMock := coinValidatorRepoMock.NewMockCoinValidatorRepository(t)
		userAccountHandler := handler.NewUserAccountHandler(userAccountCli, coinValidatorMock)

		userIdNew := uuid.New()

		req := &userAccountPb.UserAccountUuid{
			Uuid: userIdNew.String(),
		}

		retUserAccounts := model.UserAccounts{&model.UserAccount{
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
			UserIDNew:     utils.UuidToUuidNull(userIdNew),
			CoinNew:       utils.StringToStringNull(""),
		},
		}

		userAccountCli.EXPECT().
			GetUserAccountsByUuid(ctx, userIdNew).
			Return(retUserAccounts, nil)

		resp, err := userAccountHandler.GetUserAccountsByUuid(ctx, req)
		require.NotNil(t, resp)
		require.NoError(t, err)

	})

	t.Run("error parse uuid", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		userAccountCli := userAccountServiceMock.NewMockUserAccountService(t)
		coinValidatorMock := coinValidatorRepoMock.NewMockCoinValidatorRepository(t)
		userAccountHandler := handler.NewUserAccountHandler(userAccountCli, coinValidatorMock)

		userIdNew := uuid.New()

		req := &userAccountPb.UserAccountUuid{
			Uuid: userIdNew.String() + "fake",
		}

		resp, err := userAccountHandler.GetUserAccountsByUuid(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAcc1061)

	})

	t.Run("error mock repo", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		userAccountCli := userAccountServiceMock.NewMockUserAccountService(t)
		coinValidatorMock := coinValidatorRepoMock.NewMockCoinValidatorRepository(t)
		userAccountHandler := handler.NewUserAccountHandler(userAccountCli, coinValidatorMock)

		userIdNew := uuid.New()

		req := &userAccountPb.UserAccountUuid{
			Uuid: userIdNew.String(),
		}

		userAccountCli.EXPECT().
			GetUserAccountsByUuid(ctx, userIdNew).
			Return(nil, errors.New("error mock repo")).Once()

		resp, err := userAccountHandler.GetUserAccountsByUuid(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAcc1062)

	})
}
