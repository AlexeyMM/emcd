package handler_test

import (
	"context"
	"testing"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	coinValidatorMock "code.emcdtech.com/emcd/service/coin/repository/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"code.emcdtech.com/emcd/blockchain/address/internal/handler"
	externalMock "code.emcdtech.com/emcd/blockchain/address/mocks/external"
	"code.emcdtech.com/emcd/blockchain/address/mocks/internal_/service"
	addressPb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
	"code.emcdtech.com/emcd/blockchain/address/repository"
)

func TestAddressHandler_DeletePersonalAddress(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serviceAddressMock := service.NewMockAddressService(t)
	coinValidatorRepoMock := coinValidatorMock.NewMockCoinValidatorRepository(t)
	coinProtoCli := externalMock.NewMockCoinServiceClient(t)
	addressHandler := handler.NewAddressHandler(serviceAddressMock, coinValidatorRepoMock, coinProtoCli, IsNetworkOldWay)

	userUuid := uuid.New()
	network := nodeCommon.EthNetworkId

	t.Run("success personal address", func(t *testing.T) {
		req := &addressPb.DeletePersonalAddressRequest{
			UserUuid: userUuid.String(),
			Network:  network.ToString(),
		}

		t.Run("delete", func(t *testing.T) {
			serviceAddressMock.EXPECT().
				DeletePersonalAddress(ctx, userUuid, network).
				Return(nil).Once()

			_, err := addressHandler.DeletePersonalAddress(ctx, req)
			require.NoError(t, err)

		})
	})

	t.Run("error parse user uuid", func(t *testing.T) {
		req := &addressPb.DeletePersonalAddressRequest{
			UserUuid: userUuid.String() + "fake",
			Network:  network.ToString(),
		}

		_, err := addressHandler.DeletePersonalAddress(ctx, req)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAddr1091)

	})

	t.Run("error parse network", func(t *testing.T) {
		req := &addressPb.DeletePersonalAddressRequest{
			UserUuid: userUuid.String(),
			Network:  network.ToString() + "fake",
		}

		_, err := addressHandler.DeletePersonalAddress(ctx, req)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAddr1091)

	})

	t.Run("error mock service create processing address", func(t *testing.T) {
		req := &addressPb.DeletePersonalAddressRequest{
			UserUuid: userUuid.String(),
			Network:  network.ToString(),
		}

		errMock := newMockError()

		serviceAddressMock.EXPECT().
			DeletePersonalAddress(ctx, userUuid, network).
			Return(errMock).Once()

		_, err := addressHandler.DeletePersonalAddress(ctx, req)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAddr1092)
	})
}
