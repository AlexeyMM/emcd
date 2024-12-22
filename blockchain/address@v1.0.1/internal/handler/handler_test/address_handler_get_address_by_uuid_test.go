package handler_test

import (
	"context"
	"testing"
	"time"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	coinValidatorMock "code.emcdtech.com/emcd/service/coin/repository/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"code.emcdtech.com/emcd/blockchain/address/internal/handler"
	"code.emcdtech.com/emcd/blockchain/address/internal/handler/mapping"
	externalMock "code.emcdtech.com/emcd/blockchain/address/mocks/external"
	"code.emcdtech.com/emcd/blockchain/address/mocks/internal_/service"
	"code.emcdtech.com/emcd/blockchain/address/model"
	"code.emcdtech.com/emcd/blockchain/address/model/enum"
	addressPb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
	"code.emcdtech.com/emcd/blockchain/address/repository"
)

func TestAddressHandler_GetAddressByUuid(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serviceAddressMock := service.NewMockAddressService(t)
	coinValidatorRepoMock := coinValidatorMock.NewMockCoinValidatorRepository(t)
	coinProtoCli := externalMock.NewMockCoinServiceClient(t)
	addressHandler := handler.NewAddressHandler(serviceAddressMock, coinValidatorRepoMock, coinProtoCli, IsNetworkOldWay)

	addressUuid := uuid.New()
	// userUuid := uuid.New()
	addressType := addressPb.AddressType_ADDRESS_TYPE_DIRECT
	network := nodeCommon.BtcNetworkId
	// coin := "btc"

	t.Run("success new", func(t *testing.T) {
		req := &addressPb.AddressUuid{
			AddressUuid: addressUuid.String(),
		}

		addressNewResp := &model.Address{
			Id:           addressUuid,
			Address:      "",
			UserUuid:     uuid.UUID{},
			AddressType:  enum.NewAddressTypeWrapper(addressType),
			NetworkGroup: nodeCommon.NewNetworkGroupEnumWrapper(network.Group()),
			CreatedAt:    time.Time{},
		}
		addressesNewResp := model.Addresses{addressNewResp}
		addressesNewRespProto := mapping.MapModelOneOfAddressToProto(ctx, addressesNewResp, nil)

		serviceAddressMock.EXPECT().
			GetNewAddressByUuid(
				ctx,
				addressUuid,
			).
			Return(addressesNewResp, nil).
			Once()

		serviceAddressMock.EXPECT().
			GetOldAddressByUuid(
				ctx,
				addressUuid,
			).
			Return(nil, nil).
			Once()

		resp, err := addressHandler.GetAddressByUuid(ctx, req)
		require.NotNil(t, resp)
		require.NoError(t, err)

		require.Equal(t, resp, addressesNewRespProto)

	})

	t.Run("success old", func(t *testing.T) {
		req := &addressPb.AddressUuid{
			AddressUuid: addressUuid.String(),
		}

		serviceAddressMock.EXPECT().
			GetNewAddressByUuid(
				ctx,
				addressUuid,
			).
			Return(nil, nil).
			Once()

		addressOldResp := &model.AddressOld{
			Id:            addressUuid,
			Address:       "",
			UserUuid:      uuid.UUID{},
			AddressType:   enum.NewAddressTypeWrapper(addressType),
			Network:       nodeCommon.NewNetworkEnumWrapper(network),
			UserAccountId: 0,
			Coin:          "",
			CreatedAt:     time.Time{},
		}
		addressesOldResp := model.AddressesOld{addressOldResp}
		addressesOldRespProto := mapping.MapModelOneOfAddressToProto(ctx, nil, addressesOldResp)

		serviceAddressMock.EXPECT().
			GetOldAddressByUuid(
				ctx,
				addressUuid,
			).
			Return(addressesOldResp, nil).
			Once()

		resp, err := addressHandler.GetAddressByUuid(ctx, req)
		require.NotNil(t, resp)
		require.NoError(t, err)

		require.Equal(t, resp, addressesOldRespProto)

	})

	t.Run("error parse uuid", func(t *testing.T) {
		req := &addressPb.AddressUuid{
			AddressUuid: addressUuid.String() + "fake",
		}

		resp, err := addressHandler.GetAddressByUuid(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAddr1031)

	})

	t.Run("error mock new", func(t *testing.T) {
		req := &addressPb.AddressUuid{
			AddressUuid: addressUuid.String(),
		}

		errMock := newMockError()

		serviceAddressMock.EXPECT().
			GetNewAddressByUuid(
				ctx,
				addressUuid,
			).
			Return(nil, errMock).
			Once()

		resp, err := addressHandler.GetAddressByUuid(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAddr1032)

	})

	t.Run("error mock old", func(t *testing.T) {
		req := &addressPb.AddressUuid{
			AddressUuid: addressUuid.String(),
		}

		serviceAddressMock.EXPECT().
			GetNewAddressByUuid(
				ctx,
				addressUuid,
			).
			Return(nil, nil).
			Once()

		errMock := newMockError()

		serviceAddressMock.EXPECT().
			GetOldAddressByUuid(
				ctx,
				addressUuid,
			).
			Return(nil, errMock).
			Once()

		resp, err := addressHandler.GetAddressByUuid(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAddr1033)

	})

}
