package repository_test

import (
	"context"
	"testing"

	coinValidatorRepoMock "code.emcdtech.com/emcd/service/coin/repository/mocks"
	"github.com/stretchr/testify/mock"

	"code.emcdtech.com/emcd/service/accounting/internal/handler/mapping"
	userAccountProtoMock "code.emcdtech.com/emcd/service/accounting/mocks/protocol/user_account"
	"code.emcdtech.com/emcd/service/accounting/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestUserAccountRepository_CreateDefaultUserAccounts(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		userAccountServiceMock := userAccountProtoMock.NewMockUserAccountServiceClient(t)
		coinValidatorMock := coinValidatorRepoMock.NewMockCoinValidatorRepository(t)
		userAccountRepository := repository.NewUserAccountRepository(userAccountServiceMock, coinValidatorMock)

		userId := int32(1)
		userIdNew := uuid.New()

		userAccountRequest := repository.GenerateDefaultUsersAccounts(userId, userIdNew, 0.0, 0.0)
		request, err := mapping.MapModelUserAccountsToProtoMultiRequest(userId, userIdNew, userAccountRequest)
		require.NoError(t, err)

		response := mapping.MapModelUserAccountsToProtoMultiResponse(nil, userAccountRequest)

		// only for warning
		coinValidatorMock.EXPECT().
			GetIdByCode(mock.Anything).
			Return(0, true)

		userAccountServiceMock.EXPECT().
			CreateUserAccounts(ctx, request).
			Return(response, nil)

		resp, err := userAccountRepository.CreateDefaultUserAccounts(ctx, userId, userIdNew, 0.0, 0.0)
		require.NotNil(t, resp)
		require.NoError(t, err)

		require.Len(t, resp, len(userAccountRequest))
	})

	t.Run("error mock coin validate", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		userAccountServiceMock := userAccountProtoMock.NewMockUserAccountServiceClient(t)
		coinValidatorMock := coinValidatorRepoMock.NewMockCoinValidatorRepository(t)
		userAccountRepository := repository.NewUserAccountRepository(userAccountServiceMock, coinValidatorMock)

		userId := int32(1)
		userIdNew := uuid.New()

		userAccountRequest := repository.GenerateDefaultUsersAccounts(userId, userIdNew, 0.0, 0.0)
		request, err := mapping.MapModelUserAccountsToProtoMultiRequest(userId, userIdNew, userAccountRequest)
		require.NoError(t, err)

		retError := newMockError()

		userAccountServiceMock.EXPECT().
			CreateUserAccounts(ctx, request).
			Return(nil, retError)

		resp, err := userAccountRepository.CreateDefaultUserAccounts(ctx, userId, userIdNew, 0.0, 0.0)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, retError)

	})
}
