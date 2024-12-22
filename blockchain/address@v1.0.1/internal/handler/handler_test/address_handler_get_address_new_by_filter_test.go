package handler_test

import (
	"context"
	"testing"
	"time"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	coinValidatorMock "code.emcdtech.com/emcd/service/coin/repository/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"code.emcdtech.com/emcd/blockchain/address/common/utils"
	"code.emcdtech.com/emcd/blockchain/address/internal/handler"
	"code.emcdtech.com/emcd/blockchain/address/internal/handler/mapping"
	externalMock "code.emcdtech.com/emcd/blockchain/address/mocks/external"
	"code.emcdtech.com/emcd/blockchain/address/mocks/internal_/service"
	"code.emcdtech.com/emcd/blockchain/address/model"
	"code.emcdtech.com/emcd/blockchain/address/model/enum"
	addressPb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
	"code.emcdtech.com/emcd/blockchain/address/repository"
)

func TestAddressHandler_GetAddressesNewByFilter(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serviceAddressMock := service.NewMockAddressService(t)
	coinValidatorRepoMock := coinValidatorMock.NewMockCoinValidatorRepository(t)
	coinProtoCli := externalMock.NewMockCoinServiceClient(t)
	addressHandler := handler.NewAddressHandler(serviceAddressMock, coinValidatorRepoMock, coinProtoCli, IsNetworkOldWay)

	addressType := addressPb.AddressType_ADDRESS_TYPE_DIRECT
	network := nodeCommon.BtcNetworkId

	t.Run("success", func(t *testing.T) {
		req := &addressPb.AddressNewFilter{
			AddressUuid:  nil,
			Address:      nil,
			UserUuid:     nil,
			AddressType:  nil,
			NetworkGroup: nil,
			Pagination:   nil,
		}

		filter, err := mapping.MapProtoToModelAddressNewFilter(req)
		require.NoError(t, err)

		addressNewResp := &model.Address{
			Id:           uuid.UUID{},
			Address:      "",
			UserUuid:     uuid.UUID{},
			AddressType:  enum.NewAddressTypeWrapper(addressType),
			NetworkGroup: nodeCommon.NewNetworkGroupEnumWrapper(network.Group()),
			CreatedAt:    time.Time{},
		}
		addressesNewResp := model.Addresses{addressNewResp}
		addressesNewRespProto := mapping.MapModelAddressesToProto(nil, addressesNewResp, nil)

		serviceAddressMock.EXPECT().
			GetNewAddressesByFilter(
				ctx,
				filter,
			).
			Return(nil, addressesNewResp, nil).
			Once()

		resp, err := addressHandler.GetAddressesNewByFilter(ctx, req)
		require.NotNil(t, resp)
		require.NoError(t, err)

		require.Len(t, resp.Addresses, 1)
		require.Equal(t, resp, addressesNewRespProto)

	})

	t.Run("success with pagination", func(t *testing.T) {
		req := &addressPb.AddressNewFilter{
			AddressUuid:  nil,
			Address:      nil,
			UserUuid:     nil,
			AddressType:  nil,
			NetworkGroup: nil,
			Pagination: &addressPb.AddressPagination{
				Limit:  0,
				Offset: 0,
			},
		}

		filter, err := mapping.MapProtoToModelAddressNewFilter(req)
		require.NoError(t, err)

		addressNewResp := &model.Address{
			Id:           uuid.UUID{},
			Address:      "",
			UserUuid:     uuid.UUID{},
			AddressType:  enum.NewAddressTypeWrapper(addressType),
			NetworkGroup: nodeCommon.NewNetworkGroupEnumWrapper(network.Group()),
			CreatedAt:    time.Time{},
		}
		addressesNewResp := model.Addresses{addressNewResp}
		totalCount := uint64(99)
		addressesNewRespProto := mapping.MapModelAddressesToProto(&totalCount, addressesNewResp, nil)

		serviceAddressMock.EXPECT().
			GetNewAddressesByFilter(
				ctx,
				filter,
			).
			Return(&totalCount, addressesNewResp, nil).
			Once()

		resp, err := addressHandler.GetAddressesNewByFilter(ctx, req)
		require.NotNil(t, resp)
		require.NoError(t, err)

		require.Len(t, resp.Addresses, 1)
		require.Equal(t, resp, addressesNewRespProto)

	})

	t.Run("error parse uuid", func(t *testing.T) {
		req := &addressPb.AddressNewFilter{
			AddressUuid:  utils.StringToPtr(uuid.New().String() + "fake"),
			Address:      nil,
			UserUuid:     nil,
			AddressType:  nil,
			NetworkGroup: nil,
			Pagination:   nil,
		}

		resp, err := addressHandler.GetAddressesNewByFilter(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAddr1071)

	})

	t.Run("error parse uuid", func(t *testing.T) {
		req := &addressPb.AddressNewFilter{
			AddressUuid:  nil,
			Address:      nil,
			UserUuid:     utils.StringToPtr(uuid.New().String() + "fake"),
			AddressType:  nil,
			NetworkGroup: nil,
			Pagination:   nil,
		}

		resp, err := addressHandler.GetAddressesNewByFilter(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAddr1071)

	})

	t.Run("error parse network group", func(t *testing.T) {
		req := &addressPb.AddressNewFilter{
			AddressUuid:  nil,
			Address:      nil,
			UserUuid:     nil,
			AddressType:  nil,
			NetworkGroup: utils.StringToPtr("fake"),
			Pagination:   nil,
		}

		resp, err := addressHandler.GetAddressesNewByFilter(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAddr1071)

	})

	t.Run("error mock new", func(t *testing.T) {
		req := &addressPb.AddressNewFilter{
			AddressUuid:  nil,
			Address:      nil,
			UserUuid:     nil,
			AddressType:  nil,
			NetworkGroup: nil,
			Pagination:   nil,
		}

		filter, err := mapping.MapProtoToModelAddressNewFilter(req)
		require.NoError(t, err)

		errMock := newMockError()

		serviceAddressMock.EXPECT().
			GetNewAddressesByFilter(
				ctx,
				filter,
			).
			Return(nil, nil, errMock).
			Once()

		resp, err := addressHandler.GetAddressesNewByFilter(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAddr1072)

	})
}
