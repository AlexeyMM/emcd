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

func TestAddressHandler_GetAddressByUserUuid(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serviceAddressMock := service.NewMockAddressService(t)
	coinValidatorRepoMock := coinValidatorMock.NewMockCoinValidatorRepository(t)
	coinProtoCli := externalMock.NewMockCoinServiceClient(t)
	addressHandler := handler.NewAddressHandler(serviceAddressMock, coinValidatorRepoMock, coinProtoCli, IsNetworkOldWay)

	userUuid := uuid.New()
	addressType := addressPb.AddressType_ADDRESS_TYPE_DIRECT
	network := nodeCommon.BtcNetworkId
	// coin := "btc"

	t.Run("success new", func(t *testing.T) {
		req := &addressPb.UserUuid{
			UserUuid: userUuid.String(),
		}

		addressNewResp := &model.Address{
			Id:           uuid.UUID{},
			Address:      "",
			UserUuid:     userUuid,
			AddressType:  enum.NewAddressTypeWrapper(addressType),
			NetworkGroup: nodeCommon.NewNetworkGroupEnumWrapper(network.Group()),
			CreatedAt:    time.Time{},
		}
		addressesNewResp := model.Addresses{addressNewResp}
		addressesNewRespProto := mapping.MapModelAddressesToProto(nil, addressesNewResp, nil)

		serviceAddressMock.EXPECT().
			GetNewAddressesByUserUuid(
				ctx,
				userUuid,
			).
			Return(addressesNewResp, nil).
			Once()

		serviceAddressMock.EXPECT().
			GetOldAddressesByUserUuid(
				ctx,
				userUuid,
			).
			Return(nil, nil).
			Once()

		resp, err := addressHandler.GetAddressesByUserUuid(ctx, req)
		require.NotNil(t, resp)
		require.NoError(t, err)

		require.Len(t, resp.Addresses, 1)
		require.Equal(t, resp, addressesNewRespProto)

	})

	t.Run("success old", func(t *testing.T) {
		req := &addressPb.UserUuid{
			UserUuid: userUuid.String(),
		}

		serviceAddressMock.EXPECT().
			GetNewAddressesByUserUuid(
				ctx,
				userUuid,
			).
			Return(nil, nil).
			Once()

		addressOldResp := &model.AddressOld{
			Id:            uuid.UUID{},
			Address:       "",
			UserUuid:      userUuid,
			AddressType:   enum.NewAddressTypeWrapper(addressType),
			Network:       nodeCommon.NewNetworkEnumWrapper(network),
			UserAccountId: 0,
			Coin:          "",
			CreatedAt:     time.Time{},
		}
		addressesOldResp := model.AddressesOld{addressOldResp}
		addressesOldRespProto := mapping.MapModelAddressesToProto(nil, nil, addressesOldResp)

		serviceAddressMock.EXPECT().
			GetOldAddressesByUserUuid(
				ctx,
				userUuid,
			).
			Return(addressesOldResp, nil).
			Once()

		resp, err := addressHandler.GetAddressesByUserUuid(ctx, req)
		require.NotNil(t, resp)
		require.NoError(t, err)

		require.Len(t, resp.Addresses, 1)
		require.Equal(t, resp, addressesOldRespProto)

	})

	t.Run("error parse uuid", func(t *testing.T) {
		req := &addressPb.UserUuid{
			UserUuid: userUuid.String() + "fake",
		}

		resp, err := addressHandler.GetAddressesByUserUuid(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAddr1051)

	})

	t.Run("error mock new", func(t *testing.T) {
		req := &addressPb.UserUuid{
			UserUuid: userUuid.String(),
		}

		errMock := newMockError()

		serviceAddressMock.EXPECT().
			GetNewAddressesByUserUuid(
				ctx,
				userUuid,
			).
			Return(nil, errMock).
			Once()

		resp, err := addressHandler.GetAddressesByUserUuid(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAddr1052)

	})

	t.Run("error mock old", func(t *testing.T) {
		req := &addressPb.UserUuid{
			UserUuid: userUuid.String(),
		}

		serviceAddressMock.EXPECT().
			GetNewAddressesByUserUuid(
				ctx,
				userUuid,
			).
			Return(nil, nil).
			Once()

		errMock := newMockError()

		serviceAddressMock.EXPECT().
			GetOldAddressesByUserUuid(
				ctx,
				userUuid,
			).
			Return(nil, errMock).
			Once()

		resp, err := addressHandler.GetAddressesByUserUuid(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAddr1053)

	})

}
