package repository_test

import (
	"context"
	"testing"
	"time"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"code.emcdtech.com/emcd/blockchain/address/common/utils"
	"code.emcdtech.com/emcd/blockchain/address/internal/handler/mapping"
	addressPbMock "code.emcdtech.com/emcd/blockchain/address/mocks/protocol/address"
	"code.emcdtech.com/emcd/blockchain/address/model"
	"code.emcdtech.com/emcd/blockchain/address/model/enum"
	addressPb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
	"code.emcdtech.com/emcd/blockchain/address/repository"
)

func TestAddressRepository_GetAddressByUserUuid(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handlerMock := addressPbMock.NewMockAddressServiceClient(t)
	addressRepo := repository.NewAddressRepository(handlerMock)

	userUuid := uuid.New()
	network := nodeCommon.EthNetworkId
	coin := utils.StringToPtr("eth")

	t.Run("success new", func(t *testing.T) {
		addressType := addressPb.AddressType_ADDRESS_TYPE_DERIVED

		req := &addressPb.UserUuid{
			UserUuid: userUuid.String(),
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
		addressesNew := model.Addresses{addressNew}

		resp := mapping.MapModelAddressesToProto(nil, addressesNew, nil)

		handlerMock.EXPECT().
			GetAddressesByUserUuid(
				mock.Anything,
				req).
			Return(resp, nil).
			Once()

		addressNewResp, addressOldResp, err := addressRepo.GetAddressesByUserUuid(ctx, userUuid)
		require.NotNil(t, addressNewResp)
		require.Nil(t, addressOldResp)
		require.NoError(t, err)

		require.Equal(t, addressNewResp, addressesNew)

	})

	t.Run("success old", func(t *testing.T) {
		addressType := addressPb.AddressType_ADDRESS_TYPE_BASED_ID
		userAccountId := int32(0)

		req := &addressPb.UserUuid{
			UserUuid: userUuid.String(),
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

		addressesOld := model.AddressesOld{addressOld}

		resp := mapping.MapModelAddressesToProto(nil, nil, addressesOld)

		handlerMock.EXPECT().
			GetAddressesByUserUuid(
				mock.Anything,
				req).
			Return(resp, nil).
			Once()

		addressNewResp, addressOldResp, err := addressRepo.GetAddressesByUserUuid(ctx, userUuid)
		require.Nil(t, addressNewResp)
		require.NotNil(t, addressOldResp)
		require.NoError(t, err)

		require.Equal(t, addressOldResp, addressesOld)

	})

	t.Run("error service mock", func(t *testing.T) {
		req := &addressPb.UserUuid{
			UserUuid: userUuid.String(),
		}

		mockErr := newMockError()

		handlerMock.EXPECT().
			GetAddressesByUserUuid(
				mock.Anything,
				req).
			Return(nil, mockErr).
			Once()

		addressNewResp, addressOldResp, err := addressRepo.GetAddressesByUserUuid(ctx, userUuid)
		require.Nil(t, addressNewResp)
		require.Nil(t, addressOldResp)
		require.Error(t, err)

		require.ErrorIs(t, err, mockErr)

	})
}
