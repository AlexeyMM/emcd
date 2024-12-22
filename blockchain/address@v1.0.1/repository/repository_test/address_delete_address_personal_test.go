package repository_test

import (
	"context"
	"testing"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"

	addressPbMock "code.emcdtech.com/emcd/blockchain/address/mocks/protocol/address"
	addressPb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
	"code.emcdtech.com/emcd/blockchain/address/repository"
)

func TestAddressRepository_DeletePersonalAddress(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handlerMock := addressPbMock.NewMockAddressServiceClient(t)
	addressRepo := repository.NewAddressRepository(handlerMock)

	userUuid := uuid.New()
	network := nodeCommon.EthNetworkId

	t.Run("success delete", func(t *testing.T) {
		req := &addressPb.DeletePersonalAddressRequest{
			UserUuid: userUuid.String(),
			Network:  network.ToString(),
		}

		handlerMock.EXPECT().
			DeletePersonalAddress(
				ctx,
				req).
			Return(&emptypb.Empty{}, nil).
			Once()

		err := addressRepo.DeletePersonalAddress(ctx, userUuid, network)
		require.NoError(t, err)

	})

	t.Run("error service mock", func(t *testing.T) {
		req := &addressPb.DeletePersonalAddressRequest{
			UserUuid: userUuid.String(),
			Network:  network.ToString(),
		}

		mockErr := newMockError()

		handlerMock.EXPECT().
			DeletePersonalAddress(
				mock.Anything,
				req).
			Return(nil, mockErr).
			Once()

		err := addressRepo.DeletePersonalAddress(ctx, userUuid, network)
		require.Error(t, err)

		require.ErrorIs(t, err, mockErr)

	})
}
