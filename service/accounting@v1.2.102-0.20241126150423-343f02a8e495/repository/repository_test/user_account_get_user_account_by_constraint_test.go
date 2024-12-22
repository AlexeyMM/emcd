package repository_test

import (
	"code.emcdtech.com/emcd/service/accounting/internal/handler/mapping"
	"code.emcdtech.com/emcd/service/accounting/internal/utils"
	userAccountProtoMock "code.emcdtech.com/emcd/service/accounting/mocks/protocol/user_account"
	"code.emcdtech.com/emcd/service/accounting/model"
	"code.emcdtech.com/emcd/service/accounting/model/enum"
	userAccountPb "code.emcdtech.com/emcd/service/accounting/protocol/user_account"
	"code.emcdtech.com/emcd/service/accounting/repository"
	coinValidatorRepoMock "code.emcdtech.com/emcd/service/coin/repository/mocks"
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUserAccountRepository_GetUserAccountByConstraint(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		userAccountServiceMock := userAccountProtoMock.NewMockUserAccountServiceClient(t)
		coinValidatorMock := coinValidatorRepoMock.NewMockCoinValidatorRepository(t)
		userAccountRepository := repository.NewUserAccountRepository(userAccountServiceMock, coinValidatorMock)

		userIdNew := uuid.New()
		coinNew := "btc"
		accountTypeId := enum.WalletAccountTypeID

		req := &userAccountPb.UserAccountConstraintRequest{
			UserIdNew:     userIdNew.String(),
			CoinNew:       coinNew,
			AccountTypeId: accountTypeId.ToInt32(),
		}

		userAccount := &model.UserAccount{
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
		}

		userAccountResponse := mapping.MapModelUserAccountToProtoResponse(userAccount)

		// only for warning
		coinValidatorMock.EXPECT().
			GetIdByCode(coinNew).
			Return(0, false)

		userAccountServiceMock.EXPECT().
			GetUserAccountByConstraint(ctx, req).
			Return(userAccountResponse, nil)

		resp, err := userAccountRepository.GetUserAccountByConstraint(ctx, userIdNew, coinNew, accountTypeId)
		require.NotNil(t, resp)
		require.NoError(t, err)

	})

	t.Run("error mock handler", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		userAccountServiceMock := userAccountProtoMock.NewMockUserAccountServiceClient(t)
		coinValidatorMock := coinValidatorRepoMock.NewMockCoinValidatorRepository(t)
		userAccountRepository := repository.NewUserAccountRepository(userAccountServiceMock, coinValidatorMock)

		userIdNew := uuid.New()
		coinNew := "btc"
		accountTypeId := enum.WalletAccountTypeID

		req := &userAccountPb.UserAccountConstraintRequest{
			UserIdNew:     userIdNew.String(),
			CoinNew:       coinNew,
			AccountTypeId: accountTypeId.ToInt32(),
		}

		retError := newMockError()

		userAccountServiceMock.EXPECT().
			GetUserAccountByConstraint(ctx, req).
			Return(nil, retError)

		resp, err := userAccountRepository.GetUserAccountByConstraint(ctx, userIdNew, coinNew, accountTypeId)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, retError)

	})
}
