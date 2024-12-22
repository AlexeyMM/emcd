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

func TestAddressHandler_GetAddressByStr(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serviceAddressMock := service.NewMockAddressService(t)
	coinValidatorRepoMock := coinValidatorMock.NewMockCoinValidatorRepository(t)
	coinProtoCli := externalMock.NewMockCoinServiceClient(t)
	addressHandler := handler.NewAddressHandler(serviceAddressMock, coinValidatorRepoMock, coinProtoCli, IsNetworkOldWay)

	addressStr := "address"
	addressType := addressPb.AddressType_ADDRESS_TYPE_DIRECT
	network := nodeCommon.BtcNetworkId
	// coin := "btc"

	t.Run("success new", func(t *testing.T) {
		req := &addressPb.AddressStrId{
			Address: addressStr,
		}

		addressNewResp := &model.Address{
			Id:           uuid.UUID{},
			Address:      addressStr,
			UserUuid:     uuid.UUID{},
			AddressType:  enum.NewAddressTypeWrapper(addressType),
			NetworkGroup: nodeCommon.NewNetworkGroupEnumWrapper(network.Group()),
			CreatedAt:    time.Time{},
		}
		addressesNewResp := model.Addresses{addressNewResp}
		addressesNewRespProto := mapping.MapModelOneOfAddressToProto(ctx, addressesNewResp, nil)

		serviceAddressMock.EXPECT().
			GetNewAddressByStr(
				ctx,
				addressStr,
			).
			Return(addressesNewResp, nil).
			Once()

		serviceAddressMock.EXPECT().
			GetOldAddressByStr(
				ctx,
				addressStr,
			).
			Return(nil, nil).
			Once()

		resp, err := addressHandler.GetAddressByStr(ctx, req)
		require.NotNil(t, resp)
		require.NoError(t, err)

		require.Equal(t, resp, addressesNewRespProto)

	})

	t.Run("success old", func(t *testing.T) {
		req := &addressPb.AddressStrId{
			Address: addressStr,
		}

		serviceAddressMock.EXPECT().
			GetNewAddressByStr(
				ctx,
				addressStr,
			).
			Return(nil, nil).
			Once()

		addressOldResp := &model.AddressOld{
			Id:            uuid.UUID{},
			Address:       addressStr,
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
			GetOldAddressByStr(
				ctx,
				addressStr,
			).
			Return(addressesOldResp, nil).
			Once()

		resp, err := addressHandler.GetAddressByStr(ctx, req)
		require.NotNil(t, resp)
		require.NoError(t, err)

		require.Equal(t, resp, addressesOldRespProto)

	})

	t.Run("error mock new", func(t *testing.T) {
		req := &addressPb.AddressStrId{
			Address: addressStr,
		}

		errMock := newMockError()

		serviceAddressMock.EXPECT().
			GetNewAddressByStr(
				ctx,
				addressStr,
			).
			Return(nil, errMock).
			Once()

		resp, err := addressHandler.GetAddressByStr(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAddr1041)

	})

	t.Run("error mock old", func(t *testing.T) {
		req := &addressPb.AddressStrId{
			Address: addressStr,
		}

		serviceAddressMock.EXPECT().
			GetNewAddressByStr(
				ctx,
				addressStr,
			).
			Return(nil, nil).
			Once()

		errMock := newMockError()

		serviceAddressMock.EXPECT().
			GetOldAddressByStr(
				ctx,
				addressStr,
			).
			Return(nil, errMock).
			Once()

		resp, err := addressHandler.GetAddressByStr(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAddr1042)

	})
}
