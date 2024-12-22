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

func TestAddressRepository_GetAddressesNewByFilter(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handlerMock := addressPbMock.NewMockAddressServiceClient(t)
	addressRepo := repository.NewAddressRepository(handlerMock)
	network := nodeCommon.EthNetworkId

	// addressUuid := uuid.New()

	t.Run("success", func(t *testing.T) {
		filter := &model.AddressFilter{
			Id:           nil,
			Address:      nil,
			UserUuid:     nil,
			AddressType:  nil,
			NetworkGroup: nil,
			Pagination:   nil,
		}

		filterReq := mapping.MapModelAddressNewFilterToProto(filter)

		addressNew := &model.Address{
			Id:           uuid.UUID{},
			Address:      "",
			UserUuid:     uuid.UUID{},
			AddressType:  enum.AddressTypeWrapper{},
			NetworkGroup: nodeCommon.NewNetworkGroupEnumWrapper(network.Group()),
			CreatedAt:    time.Time{},
		}
		derivedAddress := &model.AddressDerived{
			AddressUuid:   addressNew.Id,
			NetworkGroup:  addressNew.NetworkGroup,
			MasterKeyId:   0,
			DerivedOffset: 0,
		}
		addressNew.SetAddressDerived(derivedAddress)

		addressesNew := model.Addresses{addressNew}

		resp := mapping.MapModelAddressesToProto(nil, addressesNew, nil)

		handlerMock.EXPECT().
			GetAddressesNewByFilter(
				ctx,
				filterReq).
			Return(resp, nil).
			Once()

		_, addressOldResp, err := addressRepo.GetAddressesNewByFilter(ctx, filter)
		require.NotNil(t, addressOldResp)
		require.NoError(t, err)

		require.Equal(t, addressOldResp, addressesNew)

	})

	t.Run("error service mock", func(t *testing.T) {
		filter := &model.AddressFilter{
			Id:           nil,
			Address:      nil,
			UserUuid:     nil,
			AddressType:  nil,
			NetworkGroup: nil,
			Pagination:   nil,
		}

		filterReq := mapping.MapModelAddressNewFilterToProto(filter)

		mockErr := newMockError()

		handlerMock.EXPECT().
			GetAddressesNewByFilter(
				ctx,
				filterReq).
			Return(nil, mockErr).
			Once()

		_, addressOldResp, err := addressRepo.GetAddressesNewByFilter(ctx, filter)
		require.Nil(t, addressOldResp)
		require.Error(t, err)

		require.ErrorIs(t, err, mockErr)

	})
}
