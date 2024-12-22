package repository_test

import (
	"context"
	"testing"
	"time"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"code.emcdtech.com/emcd/blockchain/address/common/utils"
	"code.emcdtech.com/emcd/blockchain/address/internal/handler/mapping"
	addressPbMock "code.emcdtech.com/emcd/blockchain/address/mocks/protocol/address"
	"code.emcdtech.com/emcd/blockchain/address/model"
	"code.emcdtech.com/emcd/blockchain/address/model/enum"
	addressPb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
	"code.emcdtech.com/emcd/blockchain/address/repository"
)

func TestAddressRepository_GetOrCreateAddress(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handlerMock := addressPbMock.NewMockAddressServiceClient(t)
	addressRepo := repository.NewAddressRepository(handlerMock)

	userUuid := uuid.New()
	network := nodeCommon.EthNetworkId
	coin := utils.StringToPtr("eth")

	t.Run("success new", func(t *testing.T) {
		addressType := addressPb.AddressType_ADDRESS_TYPE_DERIVED

		req := &addressPb.CreateAddressRequest{
			UserUuid: userUuid.String(),
			Network:  network.ToString(),
			Coin:     coin,
		}

		addressNew := &model.Address{
			Id:           uuid.UUID{},
			Address:      "",
			UserUuid:     userUuid,
			AddressType:  enum.NewAddressTypeWrapper(addressType),
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

		resp := mapping.MapModelAddressNewToProto(addressNew)

		handlerMock.EXPECT().
			GetOrCreateAddress(
				ctx,
				req).
			Return(resp, nil).
			Once()

		addressNewResp, addressOldResp, err := addressRepo.GetOrCreateAddress(ctx, userUuid, network, coin)
		require.NotNil(t, addressNewResp)
		require.Nil(t, addressOldResp)
		require.NoError(t, err)

		require.Equal(t, addressNewResp, addressNew)

	})

	t.Run("success old", func(t *testing.T) {
		addressType := addressPb.AddressType_ADDRESS_TYPE_BASED_ID
		userAccountId := int32(0)

		req := &addressPb.CreateAddressRequest{
			UserUuid: userUuid.String(),
			Network:  network.ToString(),
			Coin:     coin,
		}

		addressOld := &model.AddressOld{
			Id:            uuid.UUID{},
			Address:       "",
			UserUuid:      userUuid,
			AddressType:   enum.NewAddressTypeWrapper(addressType),
			Network:       nodeCommon.NewNetworkEnumWrapper(network),
			UserAccountId: userAccountId,
			Coin:          *coin,
			CreatedAt:     time.Time{},
		}

		resp := mapping.MapModelAddressOldToProto(addressOld)

		handlerMock.EXPECT().
			GetOrCreateAddress(
				ctx,
				req).
			Return(resp, nil).
			Once()

		addressNewResp, addressOldResp, err := addressRepo.GetOrCreateAddress(ctx, userUuid, network, coin)
		require.Nil(t, addressNewResp)
		require.NotNil(t, addressOldResp)
		require.NoError(t, err)

		require.Equal(t, addressOldResp, addressOld)

	})

	t.Run("error service mock", func(t *testing.T) {
		req := &addressPb.CreateAddressRequest{
			UserUuid: userUuid.String(),
			Network:  network.ToString(),
			Coin:     coin,
		}

		mockErr := newMockError()

		handlerMock.EXPECT().
			GetOrCreateAddress(
				ctx,
				req).
			Return(nil, mockErr).
			Once()

		addressNewResp, addressOldResp, err := addressRepo.GetOrCreateAddress(ctx, userUuid, network, coin)
		require.Nil(t, addressNewResp)
		require.Nil(t, addressOldResp)
		require.Error(t, err)

		require.ErrorIs(t, err, mockErr)

	})
}
