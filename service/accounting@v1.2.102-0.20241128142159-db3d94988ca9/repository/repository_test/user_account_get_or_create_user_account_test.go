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

func TestUserAccountRepository_GetOrCreateUserAccount(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	userAccountServiceMock := userAccountProtoMock.NewMockUserAccountServiceClient(t)
	coinValidatorMock := coinValidatorRepoMock.NewMockCoinValidatorRepository(t)
	userAccountRepository := repository.NewUserAccountRepository(userAccountServiceMock, coinValidatorMock)

	t.Run("success", func(t *testing.T) {
		userId := int32(1)
		userIdNew := uuid.New()
		coinIdLegacy := int32(1)
		coinNew := "btc"

		userAccount := &model.UserAccount{
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
		}

		userAccountRequest := mapping.MapModelUserAccountToProtoOneRequest(userAccount)

		userAccountResponse := mapping.MapModelUserAccountToProtoResponse(userAccount)

		// only for warning
		coinValidatorMock.EXPECT().
			GetIdByCode(coinNew).
			Return(0, false)

		userAccountServiceMock.EXPECT().
			GetOrCreateUserAccount(ctx, userAccountRequest).
			Return(userAccountResponse, nil)

		resp, err := userAccountRepository.GetOrCreateUserAccount(ctx, userAccount)
		require.NotNil(t, resp)
		require.NoError(t, err)
		require.NotNil(t, resp)
	})

	t.Run("error mock handler", func(t *testing.T) {
		userId := int32(1)
		userIdNew := uuid.New()
		coinIdLegacy := int32(1)
		coinNew := "btc"

		userAccount := &model.UserAccount{
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
		}

		retError := newMockError()

		userAccountRequest := mapping.MapModelUserAccountToProtoOneRequest(userAccount)
		require.NotNil(t, userAccountRequest)

		userAccountServiceMock.EXPECT().
			GetOrCreateUserAccount(ctx, userAccountRequest).
			Return(nil, retError)

		resp, err := userAccountRepository.GetOrCreateUserAccount(ctx, userAccount)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, retError)

	})
}
