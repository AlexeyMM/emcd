package service_test

/*
import (
	userAccountModel "code.emcdtech.com/emcd/service/accounting/model"
	userAccountModelEnum "code.emcdtech.com/emcd/service/accounting/model/enum"
	userAccountRepositoryMock "code.emcdtech.com/emcd/service/accounting/repository/repository_mock"
	coinValidatorRepoMock "code.emcdtech.com/emcd/service/coin/repository/mocks"
	"code.emcdtech.com/emcd/service/profile/internal/service"
	"code.emcdtech.com/emcd/service/profile/internal/utils"
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAccountingUserAccountService_CreateUserAccount(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		coinValidatorMock := coinValidatorRepoMock.NewMockCoinValidatorRepository(t)
		userAccountRepoMock := userAccountRepositoryMock.NewMockUserAccountRepository(t)
		userAccountService := service.NewUserAccountService(userAccountRepoMock, coinValidatorMock)

		userId := int32(1)
		userIdNew := uuid.New()
		coinId := int32(1)
		coinCode := "btc"
		accountTypeId := userAccountModelEnum.WalletAccountTypeID

		userAccount := &userAccountModel.UserAccount{
			ID:            0,
			UserID:        userId,
			CoinID:        0, // will ignore
			AccountTypeID: userAccountModelEnum.NewAccountTypeIdWrapper(accountTypeId),
			Minpay:        0.0,
			Address:       sql.NullString{},
			ChangedAt:     sql.NullTime{},
			Img1:          sql.NullFloat64{},
			Img2:          sql.NullFloat64{},
			IsActive:      sql.NullBool{},
			CreatedAt:     sql.NullTime{},
			UpdatedAt:     sql.NullTime{},
			Fee:           sql.NullFloat64{},
			UserIDNew:     utils.UuidToUuidNull(userIdNew),
			CoinNew:       utils.StringToStringNull(coinCode),
		}

		userAccountsResponse := userAccountModel.UserAccounts{userAccount}

		userAccountsMatch := mock.MatchedBy(func(userAccounts userAccountModel.UserAccounts) bool {
			isMatch := true

			for i := range userAccounts {
				this := userAccounts[i]
				other := userAccountsResponse[i]

				if this.AccountTypeID != other.AccountTypeID ||
					this.CoinNew != other.CoinNew ||
					this.UserIDNew.UUID.String() != other.UserIDNew.UUID.String() {
					isMatch = false

					break
				}
			}

			return isMatch
		})

		coinValidatorMock.EXPECT().
			GetCodeById(coinId).
			Return(coinCode, true)

		userAccountRepoMock.EXPECT().
			CreateUserAccounts(
				mock.Anything,
				userId,
				userIdNew,
				userAccountsMatch,
			).Return(userAccountsResponse, nil)

		_, err := userAccountService.CreateUserAccount(ctx, userId, userIdNew, coinId, accountTypeId.ToInt32(), 0)
		require.NoError(t, err)

	})

	t.Run("error mock coin validate", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		coinValidatorMock := coinValidatorRepoMock.NewMockCoinValidatorRepository(t)
		userAccountRepoMock := userAccountRepositoryMock.NewMockUserAccountRepository(t)
		userAccountService := service.NewUserAccountService(userAccountRepoMock, coinValidatorMock)

		userId := int32(1)
		userIdNew := uuid.New()
		coinId := int32(1)
		accountTypeId := userAccountModelEnum.WalletAccountTypeID

		coinValidatorMock.EXPECT().
			GetCodeById(coinId).
			Return("", false)

		_, err := userAccountService.CreateUserAccount(ctx, userId, userIdNew, coinId, accountTypeId.ToInt32(), 0)
		require.Error(t, err)

		require.ErrorContains(t, err, "unknown coin")

	})

	t.Run("error validate type account", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		coinValidatorMock := coinValidatorRepoMock.NewMockCoinValidatorRepository(t)
		userAccountRepoMock := userAccountRepositoryMock.NewMockUserAccountRepository(t)
		userAccountService := service.NewUserAccountService(userAccountRepoMock, coinValidatorMock)

		userId := int32(1)
		userIdNew := uuid.New()
		coinId := int32(1)
		coinCode := "btc"
		accountTypeId := userAccountModelEnum.WalletAccountTypeID + 99

		coinValidatorMock.EXPECT().
			GetCodeById(coinId).
			Return(coinCode, true)

		_, err := userAccountService.CreateUserAccount(ctx, userId, userIdNew, coinId, accountTypeId.ToInt32(), 0)
		require.Error(t, err)

		require.ErrorContains(t, err, "unknown account type")

	})

	t.Run("error mock repository", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		coinValidatorMock := coinValidatorRepoMock.NewMockCoinValidatorRepository(t)
		userAccountRepoMock := userAccountRepositoryMock.NewMockUserAccountRepository(t)
		userAccountService := service.NewUserAccountService(userAccountRepoMock, coinValidatorMock)

		userId := int32(1)
		userIdNew := uuid.New()
		coinId := int32(1)
		coinCode := "btc"
		accountTypeId := userAccountModelEnum.WalletAccountTypeID

		userAccount := &userAccountModel.UserAccount{
			ID:            0,
			UserID:        userId,
			CoinID:        0, // will ignore
			AccountTypeID: userAccountModelEnum.NewAccountTypeIdWrapper(accountTypeId),
			Minpay:        0.0,
			Address:       sql.NullString{},
			ChangedAt:     sql.NullTime{},
			Img1:          sql.NullFloat64{},
			Img2:          sql.NullFloat64{},
			IsActive:      sql.NullBool{},
			CreatedAt:     sql.NullTime{},
			UpdatedAt:     sql.NullTime{},
			Fee:           sql.NullFloat64{},
			UserIDNew:     utils.UuidToUuidNull(userIdNew),
			CoinNew:       utils.StringToStringNull(coinCode),
		}

		userAccountsResponse := userAccountModel.UserAccounts{userAccount}

		userAccountsMatch := mock.MatchedBy(func(userAccounts userAccountModel.UserAccounts) bool {
			isMatch := true

			for i := range userAccounts {
				this := userAccounts[i]
				other := userAccountsResponse[i]

				if this.AccountTypeID != other.AccountTypeID ||
					this.CoinNew != other.CoinNew ||
					this.UserIDNew.UUID.String() != other.UserIDNew.UUID.String() {
					isMatch = false

					break
				}
			}

			return isMatch
		})

		retError := newMockError()

		coinValidatorMock.EXPECT().
			GetCodeById(coinId).
			Return(coinCode, true)

		userAccountRepoMock.EXPECT().
			CreateUserAccounts(
				mock.Anything,
				userId,
				userIdNew,
				userAccountsMatch,
			).Return(nil, retError)

		_, err := userAccountService.CreateUserAccount(ctx, userId, userIdNew, coinId, accountTypeId.ToInt32(), 0)
		require.Error(t, err)

		require.ErrorIs(t, retError, retError)

	})

	t.Run("error return empty", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		coinValidatorMock := coinValidatorRepoMock.NewMockCoinValidatorRepository(t)
		userAccountRepoMock := userAccountRepositoryMock.NewMockUserAccountRepository(t)
		userAccountService := service.NewUserAccountService(userAccountRepoMock, coinValidatorMock)

		userId := int32(1)
		userIdNew := uuid.New()
		coinId := int32(1)
		coinCode := "btc"
		accountTypeId := userAccountModelEnum.WalletAccountTypeID

		userAccount := &userAccountModel.UserAccount{
			ID:            0,
			UserID:        userId,
			CoinID:        0, // will ignore
			AccountTypeID: userAccountModelEnum.NewAccountTypeIdWrapper(accountTypeId),
			Minpay:        0.0,
			Address:       sql.NullString{},
			ChangedAt:     sql.NullTime{},
			Img1:          sql.NullFloat64{},
			Img2:          sql.NullFloat64{},
			IsActive:      sql.NullBool{},
			CreatedAt:     sql.NullTime{},
			UpdatedAt:     sql.NullTime{},
			Fee:           sql.NullFloat64{},
			UserIDNew:     utils.UuidToUuidNull(userIdNew),
			CoinNew:       utils.StringToStringNull(coinCode),
		}

		userAccountsResponse := userAccountModel.UserAccounts{userAccount}

		userAccountsMatch := mock.MatchedBy(func(userAccounts userAccountModel.UserAccounts) bool {
			isMatch := true

			for i := range userAccounts {
				this := userAccounts[i]
				other := userAccountsResponse[i]

				if this.AccountTypeID != other.AccountTypeID ||
					this.CoinNew != other.CoinNew ||
					this.UserIDNew.UUID.String() != other.UserIDNew.UUID.String() {
					isMatch = false

					break
				}
			}

			return isMatch
		})

		coinValidatorMock.EXPECT().
			GetCodeById(coinId).
			Return(coinCode, true)

		userAccountRepoMock.EXPECT().
			CreateUserAccounts(
				mock.Anything,
				userId,
				userIdNew,
				userAccountsMatch,
			).Return(userAccountModel.UserAccounts{}, nil)

		_, err := userAccountService.CreateUserAccount(ctx, userId, userIdNew, coinId, accountTypeId.ToInt32(), 0)
		require.Error(t, err)

		require.ErrorIs(t, err, service.ErrUserAccountListUnexpected)
	})

}
*/
