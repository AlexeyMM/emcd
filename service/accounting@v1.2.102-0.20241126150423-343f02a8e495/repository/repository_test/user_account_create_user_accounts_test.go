package repository_test

import (
	"code.emcdtech.com/emcd/service/accounting/internal/handler/mapping"
	"code.emcdtech.com/emcd/service/accounting/internal/utils"
	userAccountProtoMock "code.emcdtech.com/emcd/service/accounting/mocks/protocol/user_account"
	"code.emcdtech.com/emcd/service/accounting/model"
	"code.emcdtech.com/emcd/service/accounting/model/enum"
	"code.emcdtech.com/emcd/service/accounting/repository"
	coinValidatorRepoMock "code.emcdtech.com/emcd/service/coin/repository/mocks"

	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUserAccountRepository_CreateUserAccounts(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		userAccountServiceMock := userAccountProtoMock.NewMockUserAccountServiceClient(t)
		coinValidatorMock := coinValidatorRepoMock.NewMockCoinValidatorRepository(t)
		userAccountRepository := repository.NewUserAccountRepository(userAccountServiceMock, coinValidatorMock)

		userId := int32(1)
		userIdNew := uuid.New()
		coinIdLegacy := int32(1)
		coinNew := "btc"

		userAccounts := model.UserAccounts{
			{
				ID:            0,
				UserID:        userId,
				CoinID:        coinIdLegacy,
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

		userAccountRequest, errMapping := mapping.MapModelUserAccountsToProtoMultiRequest(userId, userIdNew, userAccounts)
		require.NotNil(t, userAccountRequest)
		require.NoError(t, errMapping)

		userAccountResponse := mapping.MapModelUserAccountsToProtoMultiResponse(nil, userAccounts)

		// only for warning
		coinValidatorMock.EXPECT().
			GetIdByCode(coinNew).
			Return(0, false)

		userAccountServiceMock.EXPECT().
			CreateUserAccounts(ctx, userAccountRequest).
			Return(userAccountResponse, nil)

		resp, err := userAccountRepository.CreateUserAccounts(ctx, userId, userIdNew, userAccounts)
		require.NotNil(t, resp)
		require.NoError(t, err)
		require.NotNil(t, resp)
	})

	t.Run("mismatch request params", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		userAccountServiceMock := userAccountProtoMock.NewMockUserAccountServiceClient(t)
		coinValidatorMock := coinValidatorRepoMock.NewMockCoinValidatorRepository(t)
		userAccountRepository := repository.NewUserAccountRepository(userAccountServiceMock, coinValidatorMock)

		userId := int32(1)
		userIdNew := uuid.New()
		coinIdLegacy := int32(1)
		coinNew := "btc"

		userAccounts := model.UserAccounts{
			{
				ID:            0,
				UserID:        userId,
				CoinID:        coinIdLegacy,
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

		userAccountRequest, errMapping := mapping.MapModelUserAccountsToProtoMultiRequest(userId, userIdNew, userAccounts)
		require.NotNil(t, userAccountRequest)
		require.NoError(t, errMapping)

		resp, err := userAccountRepository.CreateUserAccounts(ctx, userId, uuid.New(), userAccounts)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorContains(t, err, "mismatch user_id_new:")
	})

	t.Run("error mock handler", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		userAccountServiceMock := userAccountProtoMock.NewMockUserAccountServiceClient(t)
		coinValidatorMock := coinValidatorRepoMock.NewMockCoinValidatorRepository(t)
		userAccountRepository := repository.NewUserAccountRepository(userAccountServiceMock, coinValidatorMock)

		userId := int32(1)
		userIdNew := uuid.New()
		coinIdLegacy := int32(1)
		coinNew := "btc"

		userAccounts := model.UserAccounts{
			{
				ID:            0,
				UserID:        userId,
				CoinID:        coinIdLegacy,
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

		retError := newMockError()

		userAccountRequest, errMapping := mapping.MapModelUserAccountsToProtoMultiRequest(userId, userIdNew, userAccounts)
		require.NotNil(t, userAccountRequest)
		require.NoError(t, errMapping)

		userAccountServiceMock.EXPECT().
			CreateUserAccounts(ctx, userAccountRequest).
			Return(nil, retError)

		resp, err := userAccountRepository.CreateUserAccounts(ctx, userId, userIdNew, userAccounts)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, retError)

	})
}
