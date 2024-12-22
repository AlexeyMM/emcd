package handler_test

import (
	"context"
	"testing"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	coinValidatorMock "code.emcdtech.com/emcd/service/coin/repository/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"code.emcdtech.com/emcd/blockchain/address/common/utils"
	"code.emcdtech.com/emcd/blockchain/address/internal/handler"
	"code.emcdtech.com/emcd/blockchain/address/internal/handler/mapping"
	externalMock "code.emcdtech.com/emcd/blockchain/address/mocks/external"
	"code.emcdtech.com/emcd/blockchain/address/mocks/internal_/service"
	"code.emcdtech.com/emcd/blockchain/address/model"
	addressPb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
	"code.emcdtech.com/emcd/blockchain/address/repository"
)

var IsNetworkOldWay = map[nodeCommon.NetworkEnum]bool{
	nodeCommon.BtcNetworkId: false,
	nodeCommon.EthNetworkId: true,
}

func TestAddressHandler_GetOrCreateAddress(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serviceAddressMock := service.NewMockAddressService(t)
	coinValidatorRepoMock := coinValidatorMock.NewMockCoinValidatorRepository(t)
	coinProtoCli := externalMock.NewMockCoinServiceClient(t)
	addressHandler := handler.NewAddressHandler(serviceAddressMock, coinValidatorRepoMock, coinProtoCli, IsNetworkOldWay)

	userUuid := uuid.New()
	addressUuid := mock.MatchedBy(func(_ any) bool { return true })

	coin := "eth"

	t.Run("success new address", func(t *testing.T) {
		network := nodeCommon.EthNetworkId

		addressType := addressPb.AddressType_ADDRESS_TYPE_DERIVED
		req := &addressPb.CreateAddressRequest{
			UserUuid: userUuid.String(),
			Network:  network.ToString(),
			Coin:     nil,
		}

		t.Run("get", func(t *testing.T) {
			addressNewGet := &model.Address{NetworkGroup: nodeCommon.NewNetworkGroupEnumWrapper(network.Group())}
			addressesNewGet := model.Addresses{addressNewGet}
			serviceAddressMock.EXPECT().
				GetNewAddressByConstraint(ctx, userUuid, addressType, network.Group()).
				Return(addressesNewGet, nil).Once()

			addressNewGetProto := mapping.MapModelAddressNewToProto(addressNewGet)

			resp, err := addressHandler.GetOrCreateAddress(ctx, req)
			require.NotNil(t, resp)
			require.NoError(t, err)

			require.Equal(t, resp, addressNewGetProto)

		})

		t.Run("create", func(t *testing.T) {
			addressesNewGet := model.Addresses{}
			serviceAddressMock.EXPECT().
				GetNewAddressByConstraint(ctx, userUuid, addressType, network.Group()).
				Return(addressesNewGet, nil).Once()

			addressNewCreate := &model.Address{NetworkGroup: nodeCommon.NewNetworkGroupEnumWrapper(network.Group())}
			serviceAddressMock.EXPECT().
				CreateNewAddress(ctx, addressUuid, userUuid, addressType, network.Group()).
				Return(addressNewCreate, nil).Once()

			addressNewCreateProto := mapping.MapModelAddressNewToProto(addressNewCreate)

			resp, err := addressHandler.GetOrCreateAddress(ctx, req)
			require.NotNil(t, resp)
			require.NoError(t, err)

			require.Equal(t, resp, addressNewCreateProto)

		})
	})

	t.Run("success old address", func(t *testing.T) {
		network := nodeCommon.EthNetworkId
		addressType := addressPb.AddressType_ADDRESS_TYPE_BASED_ID

		req := &addressPb.CreateAddressRequest{
			UserUuid: userUuid.String(),
			Network:  network.ToString(),
			Coin:     &coin,
		}

		t.Run("get", func(t *testing.T) {
			coinValidatorRepoMock.EXPECT().
				IsValidCode(coin).
				Return(true).Once()

			addressOldGet := &model.AddressOld{Network: nodeCommon.NewNetworkEnumWrapper(network)}
			addressesOldGet := model.AddressesOld{addressOldGet}
			serviceAddressMock.EXPECT().
				GetOldAddressByConstraint(ctx, userUuid, network, coin).
				Return(addressesOldGet, nil).Once()

			addressOldGetProto := mapping.MapModelAddressOldToProto(addressOldGet)

			resp, err := addressHandler.GetOrCreateAddress(ctx, req)
			require.NotNil(t, resp)
			require.NoError(t, err)

			require.Equal(t, resp, addressOldGetProto)

		})

		t.Run("create", func(t *testing.T) {
			coinValidatorRepoMock.EXPECT().
				IsValidCode(coin).
				Return(true).Once()

			addressesOldGet := model.AddressesOld{}
			serviceAddressMock.EXPECT().
				GetOldAddressByConstraint(ctx, userUuid, network, coin).
				Return(addressesOldGet, nil).Once()

			addressOldCreate := &model.AddressOld{Network: nodeCommon.NewNetworkEnumWrapper(network)}
			serviceAddressMock.EXPECT().
				CreateOldAddress(ctx, addressUuid, userUuid, addressType, network, coin).
				Return(addressOldCreate, nil).Once()

			addressOldCreateProto := mapping.MapModelAddressOldToProto(addressOldCreate)

			resp, err := addressHandler.GetOrCreateAddress(ctx, req)
			require.NotNil(t, resp)
			require.NoError(t, err)

			require.Equal(t, resp, addressOldCreateProto)

		})
	})

	t.Run("error parse user uuid", func(t *testing.T) {
		networkGroup := nodeCommon.EthNetworkGroupId

		req := &addressPb.CreateAddressRequest{
			UserUuid: userUuid.String() + "fake",
			Network:  networkGroup.ToString(),
			Coin:     nil,
		}

		resp, err := addressHandler.GetOrCreateAddress(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAddr1011)

	})

	t.Run("error parse network", func(t *testing.T) {
		networkGroup := nodeCommon.EthNetworkGroupId

		req := &addressPb.CreateAddressRequest{
			UserUuid: userUuid.String(),
			Network:  networkGroup.ToString() + "fake",
			Coin:     nil,
		}

		resp, err := addressHandler.GetOrCreateAddress(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAddr1011)

	})

	t.Run("error mock coin validator", func(t *testing.T) {
		networkGroup := nodeCommon.EthNetworkGroupId

		req := &addressPb.CreateAddressRequest{
			UserUuid: userUuid.String(),
			Network:  networkGroup.ToString(),
			Coin:     utils.StringToPtr(coin),
		}

		coinValidatorRepoMock.EXPECT().
			IsValidCode(coin).
			Return(false).Once()

		resp, err := addressHandler.GetOrCreateAddress(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAddr1011)

	})

	// t.Run("error unsupported network for new way", func(t *testing.T) {
	// 	networkGroup := nodeCommon.EthNetworkId
	//
	// 	req := &addressPb.CreateAddressRequest{
	// 		UserUuid: userUuid.String(),
	// 		Network:  networkGroup.ToString() + "fake",
	// 		Coin:     nil,
	// 	}
	//
	// 	resp, err := addressHandler.GetOrCreateAddress(ctx, req)
	// 	require.Nil(t, resp)
	// 	require.Error(t, err)
	//
	// 	require.ErrorIs(t, err, repository.ErrAddr1012)
	//
	// })

	t.Run("error mock service new address", func(t *testing.T) {
		networkGroup := nodeCommon.EthNetworkGroupId

		addressType := addressPb.AddressType_ADDRESS_TYPE_DERIVED
		req := &addressPb.CreateAddressRequest{
			UserUuid: userUuid.String(),
			Network:  networkGroup.ToString(),
			Coin:     nil,
		}

		errMock := newMockError()

		t.Run("get constraint", func(t *testing.T) {
			serviceAddressMock.EXPECT().
				GetNewAddressByConstraint(ctx, userUuid, addressType, networkGroup).
				Return(nil, errMock).Once()

			resp, err := addressHandler.GetOrCreateAddress(ctx, req)
			require.Nil(t, resp)
			require.Error(t, err)

			require.ErrorIs(t, err, repository.ErrAddr1013)

		})

		t.Run("create", func(t *testing.T) {
			addressesNewGet := model.Addresses{}
			serviceAddressMock.EXPECT().
				GetNewAddressByConstraint(ctx, userUuid, addressType, networkGroup).
				Return(addressesNewGet, nil).Once()

			serviceAddressMock.EXPECT().
				CreateNewAddress(ctx, addressUuid, userUuid, addressType, networkGroup).
				Return(nil, errMock).Once()

			resp, err := addressHandler.GetOrCreateAddress(ctx, req)
			require.Nil(t, resp)
			require.Error(t, err)

			require.ErrorIs(t, err, repository.ErrAddr1014)
		})
	})

	t.Run("error mock old address", func(t *testing.T) {
		network := nodeCommon.EthNetworkId
		addressType := addressPb.AddressType_ADDRESS_TYPE_BASED_ID

		req := &addressPb.CreateAddressRequest{
			UserUuid: userUuid.String(),
			Network:  network.Group().ToString(),
			Coin:     &coin,
		}

		coinValidatorRepoMock.EXPECT().
			IsValidCode(coin).
			Return(true).Twice()

		errMock := newMockError()

		t.Run("get constraint", func(t *testing.T) {
			serviceAddressMock.EXPECT().
				GetOldAddressByConstraint(ctx, userUuid, network, coin).
				Return(nil, errMock).Once()

			resp, err := addressHandler.GetOrCreateAddress(ctx, req)
			require.Nil(t, resp)
			require.Error(t, err)

			require.ErrorIs(t, err, repository.ErrAddr1016)

		})

		t.Run("create", func(t *testing.T) {
			addressesOldGet := model.AddressesOld{}
			serviceAddressMock.EXPECT().
				GetOldAddressByConstraint(ctx, userUuid, network, coin).
				Return(addressesOldGet, nil).Once()

			serviceAddressMock.EXPECT().
				CreateOldAddress(ctx, addressUuid, userUuid, addressType, network, coin).
				Return(nil, mockErrorImp{}).Once()

			resp, err := addressHandler.GetOrCreateAddress(ctx, req)
			require.Nil(t, resp)
			require.Error(t, err)

			require.ErrorIs(t, err, repository.ErrAddr1017)

		})
	})
}
