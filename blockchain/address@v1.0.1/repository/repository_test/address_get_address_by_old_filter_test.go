package repository_test

import (
	"context"
	"testing"
	"time"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"code.emcdtech.com/emcd/blockchain/address/internal/handler/mapping"
	addressPbMock "code.emcdtech.com/emcd/blockchain/address/mocks/protocol/address"
	"code.emcdtech.com/emcd/blockchain/address/model"
	"code.emcdtech.com/emcd/blockchain/address/model/enum"
	"code.emcdtech.com/emcd/blockchain/address/repository"
)

func TestAddressRepository_GetAddressesOldByFilter(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handlerMock := addressPbMock.NewMockAddressServiceClient(t)
	addressRepo := repository.NewAddressRepository(handlerMock)
	network := nodeCommon.EthNetworkId

	// addressUuid := uuid.New()

	t.Run("success", func(t *testing.T) {
		filter := &model.AddressOldFilter{
			Id:            nil,
			Address:       nil,
			UserUuid:      nil,
			AddressType:   nil,
			Network:       nil,
			UserAccountId: nil,
			Coin:          nil,
			Pagination:    nil,
		}

		filterReq := mapping.MapModelAddressOldFilterToProto(filter)

		addressOld := &model.AddressOld{
			Id:            uuid.UUID{},
			Address:       "",
			UserUuid:      uuid.UUID{},
			AddressType:   enum.AddressTypeWrapper{},
			Network:       nodeCommon.NewNetworkEnumWrapper(network),
			UserAccountId: 0,
			Coin:          "",
			CreatedAt:     time.Time{},
		}

		addressesOld := model.AddressesOld{addressOld}

		resp := mapping.MapModelAddressesToProto(nil, nil, addressesOld)

		handlerMock.EXPECT().
			GetAddressesOldByFilter(
				ctx,
				filterReq).
			Return(resp, nil).
			Once()

		_, addressOldResp, err := addressRepo.GetAddressesOldByFilter(ctx, filter)
		require.NotNil(t, addressOldResp)
		require.NoError(t, err)

		require.Equal(t, addressOldResp, addressesOld)

	})

	t.Run("error service mock", func(t *testing.T) {
		filter := &model.AddressOldFilter{
			Id:            nil,
			Address:       nil,
			UserUuid:      nil,
			AddressType:   nil,
			Network:       nil,
			UserAccountId: nil,
			Coin:          nil,
			Pagination:    nil,
		}

		filterReq := mapping.MapModelAddressOldFilterToProto(filter)

		mockErr := newMockError()

		handlerMock.EXPECT().
			GetAddressesOldByFilter(
				ctx,
				filterReq).
			Return(nil, mockErr).
			Once()

		_, addressOldResp, err := addressRepo.GetAddressesOldByFilter(ctx, filter)
		require.Nil(t, addressOldResp)
		require.Error(t, err)

		require.ErrorIs(t, err, mockErr)

	})
}
