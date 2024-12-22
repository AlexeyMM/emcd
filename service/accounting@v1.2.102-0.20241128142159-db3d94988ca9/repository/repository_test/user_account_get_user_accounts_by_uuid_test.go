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

func TestUserAccountRepositoryGetUserAccountsByUuid_Success(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		userAccountServiceMock := userAccountProtoMock.NewMockUserAccountServiceClient(t)
		coinValidatorMock := coinValidatorRepoMock.NewMockCoinValidatorRepository(t)
		userAccountRepository := repository.NewUserAccountRepository(userAccountServiceMock, coinValidatorMock)

		userIdNew := uuid.New()
		coinNew := "btc"

		userAccounts := model.UserAccounts{
			{
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
				CoinNew:       utils.StringToStringNull(coinNew),
			}}

		userAccountResponse := mapping.MapModelUserAccountsToProtoMultiResponse(nil, userAccounts)

		uuidRequest := &userAccountPb.UserAccountUuid{Uuid: userIdNew.String()}

		// only for warning
		coinValidatorMock.EXPECT().
			GetIdByCode(coinNew).
			Return(0, false)

		userAccountServiceMock.EXPECT().
			GetUserAccountsByUuid(ctx, uuidRequest).
			Return(userAccountResponse, nil)

		resp, err := userAccountRepository.GetUserAccountsByUuid(ctx, userIdNew)
		require.NotNil(t, resp)
		require.NoError(t, err)

		require.Len(t, resp, len(userAccounts))

	})

	t.Run("error mocker handler", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		userAccountServiceMock := userAccountProtoMock.NewMockUserAccountServiceClient(t)
		coinValidatorMock := coinValidatorRepoMock.NewMockCoinValidatorRepository(t)
		userAccountRepository := repository.NewUserAccountRepository(userAccountServiceMock, coinValidatorMock)

		userIdNew := uuid.New()

		retError := newMockError()

		uuidRequest := &userAccountPb.UserAccountUuid{Uuid: userIdNew.String()}

		userAccountServiceMock.EXPECT().
			GetUserAccountsByUuid(ctx, uuidRequest).
			Return(nil, retError)

		resp, err := userAccountRepository.GetUserAccountsByUuid(ctx, userIdNew)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorAs(t, err, &retError)

	})
}
