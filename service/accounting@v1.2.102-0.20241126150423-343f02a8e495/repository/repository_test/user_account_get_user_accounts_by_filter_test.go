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

func TestUserAccountRepositoryGetUserAccountsByFilter_Success(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		userAccountServiceMock := userAccountProtoMock.NewMockUserAccountServiceClient(t)
		coinValidatorMock := coinValidatorRepoMock.NewMockCoinValidatorRepository(t)
		userAccountRepository := repository.NewUserAccountRepository(userAccountServiceMock, coinValidatorMock)

		userId := int32(1)
		userIdNew := uuid.New()
		userAccountId := enum.WalletAccountTypeID
		coinIdLegacy := int32(1)
		coinNew := "btc"

		filter := &model.UserAccountFilter{
			ID:              nil,
			UserID:          nil,
			AccountTypeID:   userAccountId.ToPtr(),
			UserIDNew:       utils.UuidToPtr(userIdNew),
			CoinNew:         &coinNew,
			IsActive:        nil,
			Pagination:      nil,
			UserIDNewIsNull: nil,
			CoinNewIsNull:   nil,
		}

		userAccounts := model.UserAccounts{
			{
				ID:            0,
				UserID:        userId,
				CoinID:        coinIdLegacy,
				AccountTypeID: enum.NewAccountTypeIdWrapper(userAccountId),
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

		filterProto := mapping.MapModelUserAccountFilterToProto(filter)

		// only for warning
		coinValidatorMock.EXPECT().
			GetIdByCode(coinNew).
			Return(0, false)

		userAccountServiceMock.EXPECT().
			GetUserAccountsByFilter(ctx, filterProto).
			Return(userAccountResponse, nil)

		totalCount, resp, err := userAccountRepository.GetUserAccountsByFilter(ctx, filter)
		require.Nil(t, totalCount)
		require.NotNil(t, resp)
		require.NoError(t, err)

		require.Len(t, resp, len(userAccounts))

	})

	t.Run("success pagination", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		userAccountServiceMock := userAccountProtoMock.NewMockUserAccountServiceClient(t)
		coinValidatorMock := coinValidatorRepoMock.NewMockCoinValidatorRepository(t)
		userAccountRepository := repository.NewUserAccountRepository(userAccountServiceMock, coinValidatorMock)

		userIdNew := uuid.New()
		userAccountId := enum.WalletAccountTypeID
		coinNew := "btc"
		totalCountRet := uint64(5)

		filter := &model.UserAccountFilter{
			ID:            nil,
			UserID:        nil,
			AccountTypeID: userAccountId.ToPtr(),
			UserIDNew:     utils.UuidToPtr(userIdNew),
			CoinNew:       &coinNew,
			Pagination: &model.Pagination{
				Limit:  10,
				Offset: 0,
			},
			UserIDNewIsNull: nil,
			CoinNewIsNull:   nil,
		}

		userAccounts := model.UserAccounts{
			{
				ID:            0,
				UserID:        0,
				CoinID:        0,
				AccountTypeID: enum.NewAccountTypeIdWrapper(userAccountId),
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

		userAccountResponse := mapping.MapModelUserAccountsToProtoMultiResponse(&totalCountRet, userAccounts)

		filterProto := mapping.MapModelUserAccountFilterToProto(filter)

		// only for warning
		coinValidatorMock.EXPECT().
			GetIdByCode(coinNew).
			Return(0, false)

		userAccountServiceMock.EXPECT().
			GetUserAccountsByFilter(ctx, filterProto).
			Return(userAccountResponse, nil)

		totalCount, resp, err := userAccountRepository.GetUserAccountsByFilter(ctx, filter)
		require.NotNil(t, totalCount)
		require.NotNil(t, resp)
		require.NoError(t, err)

		require.Equal(t, *totalCount, totalCountRet)
		require.Len(t, resp, len(userAccounts))

	})

	t.Run("error mock handler", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		userAccountServiceMock := userAccountProtoMock.NewMockUserAccountServiceClient(t)
		coinValidatorMock := coinValidatorRepoMock.NewMockCoinValidatorRepository(t)
		userAccountRepository := repository.NewUserAccountRepository(userAccountServiceMock, coinValidatorMock)

		userIdNew := uuid.New()
		userAccountId := enum.WalletAccountTypeID
		coinNew := "btc"

		filter := &model.UserAccountFilter{
			ID:              nil,
			UserID:          nil,
			AccountTypeID:   userAccountId.ToPtr(),
			UserIDNew:       utils.UuidToPtr(userIdNew),
			CoinNew:         &coinNew,
			Pagination:      nil,
			UserIDNewIsNull: nil,
			CoinNewIsNull:   nil,
		}

		retError := newMockError()

		filterProto := mapping.MapModelUserAccountFilterToProto(filter)

		userAccountServiceMock.EXPECT().
			GetUserAccountsByFilter(ctx, filterProto).
			Return(nil, retError)

		totalCount, resp, err := userAccountRepository.GetUserAccountsByFilter(ctx, filter)
		require.Nil(t, totalCount)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorAs(t, err, &retError)

	})
}
