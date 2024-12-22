package handler_test

import (
	"context"
	"testing"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	coinValidatorMock "code.emcdtech.com/emcd/service/coin/repository/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"code.emcdtech.com/emcd/blockchain/address/internal/handler"
	"code.emcdtech.com/emcd/blockchain/address/internal/handler/mapping"
	externalMock "code.emcdtech.com/emcd/blockchain/address/mocks/external"
	"code.emcdtech.com/emcd/blockchain/address/mocks/internal_/service"
	"code.emcdtech.com/emcd/blockchain/address/model"
	addressPb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
	"code.emcdtech.com/emcd/blockchain/address/repository"
)

func TestAddressHandler_CreateProcessingAddress(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serviceAddressMock := service.NewMockAddressService(t)
	coinValidatorRepoMock := coinValidatorMock.NewMockCoinValidatorRepository(t)
	coinProtoCli := externalMock.NewMockCoinServiceClient(t)
	addressHandler := handler.NewAddressHandler(serviceAddressMock, coinValidatorRepoMock, coinProtoCli, IsNetworkOldWay)

	userUuid := uuid.New()
	processingUuid := uuid.New()
	addressUuid := mock.MatchedBy(func(_ any) bool { return true })
	network := nodeCommon.EthNetworkId

	t.Run("success processing address", func(t *testing.T) {
		addressType := addressPb.AddressType_ADDRESS_TYPE_DERIVED
		req := &addressPb.CreateProcessingAddressRequest{
			UserUuid:       userUuid.String(),
			Network:        network.ToString(),
			ProcessingUuid: processingUuid.String(),
		}

		t.Run("create", func(t *testing.T) {
			addressNewCreate := &model.Address{NetworkGroup: nodeCommon.NewNetworkGroupEnumWrapper(network.Group())}
			serviceAddressMock.EXPECT().
				CreateProcessingAddress(ctx, addressUuid, userUuid, processingUuid, addressType, network.Group()).
				Return(addressNewCreate, nil).Once()

			addressNewCreateProto := mapping.MapModelAddressNewToProto(addressNewCreate)

			resp, err := addressHandler.CreateProcessingAddress(ctx, req)
			require.NotNil(t, resp)
			require.NoError(t, err)

			require.Equal(t, resp, addressNewCreateProto)

		})
	})

	t.Run("error parse user uuid", func(t *testing.T) {
		req := &addressPb.CreateProcessingAddressRequest{
			UserUuid:       userUuid.String() + "fake",
			Network:        network.ToString(),
			ProcessingUuid: processingUuid.String(),
		}

		resp, err := addressHandler.CreateProcessingAddress(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAddr1021)

	})

	t.Run("error parse network", func(t *testing.T) {
		req := &addressPb.CreateProcessingAddressRequest{
			UserUuid:       userUuid.String(),
			Network:        network.ToString() + "fake",
			ProcessingUuid: processingUuid.String(),
		}

		resp, err := addressHandler.CreateProcessingAddress(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAddr1021)

	})

	t.Run("error parse processing uuid", func(t *testing.T) {
		req := &addressPb.CreateProcessingAddressRequest{
			UserUuid:       userUuid.String(),
			Network:        network.ToString(),
			ProcessingUuid: processingUuid.String() + "fake",
		}

		resp, err := addressHandler.CreateProcessingAddress(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAddr1021)

	})

	// t.Run("error unsupported network for new way", func(t *testing.T) {
	// 	network := nodeCommon.EthNetworkGroupId
	//
	// 	req := &addressPb.CreateProcessingAddressRequest{
	// 		UserUuid:       userUuid.String(),
	// 		Network:        network.ToString() + "fake",
	// 		ProcessingUuid: processingUuid.String(),
	// 	}
	//
	// 	resp, err := addressHandler.CreateProcessingAddress(ctx, req)
	// 	require.Nil(t, resp)
	// 	require.Error(t, err)
	//
	// 	require.ErrorIs(t, err, repository.ErrAddr1022)
	//
	// })

	t.Run("error mock service create processing address", func(t *testing.T) {
		network := nodeCommon.EthNetworkId

		addressType := addressPb.AddressType_ADDRESS_TYPE_DERIVED
		req := &addressPb.CreateProcessingAddressRequest{
			UserUuid:       userUuid.String(),
			Network:        network.ToString(),
			ProcessingUuid: processingUuid.String(),
		}

		errMock := newMockError()

		serviceAddressMock.EXPECT().
			CreateProcessingAddress(ctx, addressUuid, userUuid, processingUuid, addressType, network.Group()).
			Return(nil, errMock).Once()

		resp, err := addressHandler.CreateProcessingAddress(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAddr1023)
	})
}
