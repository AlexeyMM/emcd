package service_test

import (
	"code.emcdtech.com/emcd/service/accounting/internal/service"
	"code.emcdtech.com/emcd/service/accounting/internal/utils"
	repositoryMock "code.emcdtech.com/emcd/service/accounting/mocks/internal_/repository"
	"code.emcdtech.com/emcd/service/accounting/model"
	"code.emcdtech.com/emcd/service/accounting/model/enum"
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUserAccountService_CreateUserAccounts(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	userAccountRepo := repositoryMock.NewMockUserAccountRepo(t)
	userAccountService := service.NewUserAccountService(userAccountRepo)

	t.Run("success", func(t *testing.T) {
		userId := int32(1)
		userIdNew := uuid.UUID{}
		userAccountId := enum.WalletAccountTypeID
		coinId := int32(1)
		coinNew := "btc"

		userAccountsModel := model.UserAccounts{{
			ID:            0,
			UserID:        userId,
			CoinID:        coinId,
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
		},
		}

		txFn := func(ctx context.Context) error {
			if err := userAccountRepo.AddUserAccounts(ctx, userId, userIdNew, userAccountsModel); err != nil {

				return fmt.Errorf("failed add user accounts list: %w", err)
			} else {

				return nil
			}
		}

		userAccountRepo.EXPECT().
			AddUserAccounts(ctx, userId, userIdNew, userAccountsModel).
			Return(nil).Once()
		userAccountRepo.EXPECT().
			WithinTransaction(ctx, mock.AnythingOfType("func(context.Context) error")).
			RunAndReturn(func(ctx context.Context, _ func(context.Context) error) error { return txFn(ctx) }).Once()

		result, err := userAccountService.CreateUserAccounts(ctx, userId, userIdNew, userAccountsModel)
		require.NotNil(t, result)
		require.NoError(t, err)

	})

	t.Run("error mock return", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		userAccountRepo := repositoryMock.NewMockUserAccountRepo(t)
		userAccountService := service.NewUserAccountService(userAccountRepo)

		userId := int32(1)
		userIdNew := uuid.UUID{}
		userAccountId := enum.WalletAccountTypeID
		coinId := int32(1)
		coinNew := "btc"

		userAccountsModel := model.UserAccounts{{
			ID:            0,
			UserID:        userId,
			CoinID:        coinId,
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
		},
		}

		txFn := func(ctx context.Context) error {
			if err := userAccountRepo.AddUserAccounts(ctx, userId, userIdNew, userAccountsModel); err != nil {

				return fmt.Errorf("failed add user accounts list: %w", err)
			} else {

				return nil
			}
		}

		retError := newMockError()

		userAccountRepo.EXPECT().
			AddUserAccounts(ctx, userId, userIdNew, userAccountsModel).
			Return(retError).Once()
		userAccountRepo.EXPECT().
			WithinTransaction(ctx, mock.AnythingOfType("func(context.Context) error")).
			RunAndReturn(func(ctx context.Context, _ func(context.Context) error) error { return txFn(ctx) }).Once()

		result, err := userAccountService.CreateUserAccounts(ctx, userId, userIdNew, userAccountsModel)
		require.Nil(t, result)
		require.Error(t, err)

		require.ErrorIs(t, err, retError)

	})
}
