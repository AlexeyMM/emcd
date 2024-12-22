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

func TestAddressHandler_GetAddressesOldByFilter(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serviceAddressMock := service.NewMockAddressService(t)
	coinValidatorRepoMock := coinValidatorMock.NewMockCoinValidatorRepository(t)
	coinProtoCli := externalMock.NewMockCoinServiceClient(t)
	addressHandler := handler.NewAddressHandler(serviceAddressMock, coinValidatorRepoMock, coinProtoCli, IsNetworkOldWay)

	addressType := addressPb.AddressType_ADDRESS_TYPE_DIRECT
	network := nodeCommon.BtcNetworkId

	t.Run("success", func(t *testing.T) {
		req := &addressPb.AddressOldFilter{
			AddressUuid:   nil,
			Address:       nil,
			UserUuid:      nil,
			AddressType:   nil,
			Network:       nil,
			UserAccountId: nil,
			Coin:          nil,
			Pagination:    nil,
		}

		filter, err := mapping.MapProtoToModelAddressOldFilter(coinValidatorRepoMock, req)
		require.NoError(t, err)

		addressOldResp := &model.AddressOld{
			Id:            uuid.UUID{},
			Address:       "",
			UserUuid:      uuid.UUID{},
			AddressType:   enum.NewAddressTypeWrapper(addressType),
			Network:       nodeCommon.NewNetworkEnumWrapper(network),
			UserAccountId: 0,
			Coin:          "",
			CreatedAt:     time.Time{},
		}
		addressesOldResp := model.AddressesOld{addressOldResp}
		addressesOldRespProto := mapping.MapModelAddressesToProto(nil, nil, addressesOldResp)

		serviceAddressMock.EXPECT().
			GetOldAddressesByFilter(
				ctx,
				filter,
			).
			Return(nil, addressesOldResp, nil).
			Once()

		resp, err := addressHandler.GetAddressesOldByFilter(ctx, req)
		require.NotNil(t, resp)
		require.NoError(t, err)

		require.Len(t, resp.Addresses, 1)
		require.Equal(t, resp, addressesOldRespProto)

	})

	t.Run("success with pagination", func(t *testing.T) {
		req := &addressPb.AddressOldFilter{
			AddressUuid:   nil,
			Address:       nil,
			UserUuid:      nil,
			AddressType:   nil,
			Network:       nil,
			UserAccountId: nil,
			Coin:          nil,
			Pagination: &addressPb.AddressPagination{
				Limit:  0,
				Offset: 0,
			},
		}

		filter, err := mapping.MapProtoToModelAddressOldFilter(coinValidatorRepoMock, req)
		require.NoError(t, err)

		addressOldResp := &model.AddressOld{
			Id:            uuid.UUID{},
			Address:       "",
			UserUuid:      uuid.UUID{},
			AddressType:   enum.NewAddressTypeWrapper(addressType),
			Network:       nodeCommon.NewNetworkEnumWrapper(network),
			UserAccountId: 0,
			Coin:          "",
			CreatedAt:     time.Time{},
		}
		addressesOldResp := model.AddressesOld{addressOldResp}
		totalCount := uint64(99)
		addressesOldRespProto := mapping.MapModelAddressesToProto(&totalCount, nil, addressesOldResp)

		serviceAddressMock.EXPECT().
			GetOldAddressesByFilter(
				ctx,
				filter,
			).
			Return(&totalCount, addressesOldResp, nil).
			Once()

		resp, err := addressHandler.GetAddressesOldByFilter(ctx, req)
		require.NotNil(t, resp)
		require.NoError(t, err)

		require.Len(t, resp.Addresses, 1)
		require.Equal(t, resp, addressesOldRespProto)

	})

	t.Run("error parse uuid", func(t *testing.T) {
		req := &addressPb.AddressOldFilter{
			AddressUuid:   utils.StringToPtr("fake"),
			Address:       nil,
			UserUuid:      nil,
			AddressType:   nil,
			Network:       nil,
			UserAccountId: nil,
			Coin:          nil,
			Pagination:    nil,
		}

		resp, err := addressHandler.GetAddressesOldByFilter(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAddr1061)

	})

	t.Run("error parse uuid", func(t *testing.T) {
		req := &addressPb.AddressOldFilter{
			AddressUuid:   nil,
			Address:       nil,
			UserUuid:      utils.StringToPtr("fake"),
			AddressType:   nil,
			Network:       nil,
			UserAccountId: nil,
			Coin:          nil,
			Pagination:    nil,
		}

		resp, err := addressHandler.GetAddressesOldByFilter(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAddr1061)

	})

	t.Run("error parse network", func(t *testing.T) {
		req := &addressPb.AddressOldFilter{
			AddressUuid:   nil,
			Address:       nil,
			UserUuid:      nil,
			AddressType:   nil,
			Network:       utils.StringToPtr("fake"),
			UserAccountId: nil,
			Coin:          nil,
			Pagination:    nil,
		}

		resp, err := addressHandler.GetAddressesOldByFilter(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAddr1061)

	})

	t.Run("error mock coin validator", func(t *testing.T) {
		coin := "fake"
		req := &addressPb.AddressOldFilter{
			AddressUuid:   nil,
			Address:       nil,
			UserUuid:      nil,
			AddressType:   nil,
			Network:       nil,
			UserAccountId: nil,
			Coin:          &coin,
			Pagination:    nil,
		}

		coinValidatorRepoMock.EXPECT().
			IsValidCode(coin).
			Return(false)

		resp, err := addressHandler.GetAddressesOldByFilter(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAddr1061)

	})

	t.Run("error mock new", func(t *testing.T) {
		req := &addressPb.AddressOldFilter{
			AddressUuid:   nil,
			Address:       nil,
			UserUuid:      nil,
			AddressType:   nil,
			Network:       nil,
			UserAccountId: nil,
			Coin:          nil,
			Pagination:    nil,
		}

		filter, err := mapping.MapProtoToModelAddressOldFilter(coinValidatorRepoMock, req)
		require.NoError(t, err)

		errMock := newMockError()

		serviceAddressMock.EXPECT().
			GetOldAddressesByFilter(
				ctx,
				filter,
			).
			Return(nil, nil, errMock).
			Once()

		resp, err := addressHandler.GetAddressesOldByFilter(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAddr1062)

	})
}
