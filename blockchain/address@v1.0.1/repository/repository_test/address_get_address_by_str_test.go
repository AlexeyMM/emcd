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

func TestAddressRepository_GetAddressByStr(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handlerMock := addressPbMock.NewMockAddressServiceClient(t)
	addressRepo := repository.NewAddressRepository(handlerMock)

	addressStr := "address"
	network := nodeCommon.EthNetworkId
	coin := utils.StringToPtr("eth")

	t.Run("success new", func(t *testing.T) {
		addressType := addressPb.AddressType_ADDRESS_TYPE_DERIVED

		req := &addressPb.AddressStrId{
			Address: addressStr,
		}

		addressNew := &model.Address{
			Id:           uuid.UUID{},
			Address:      addressStr,
			UserUuid:     uuid.UUID{},
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

		resp := mapping.MapModelOneOfAddressToProto(ctx, addressesNew, nil)

		handlerMock.EXPECT().
			GetAddressByStr(
				mock.Anything,
				req).
			Return(resp, nil).
			Once()

		addressNewResp, addressOldResp, err := addressRepo.GetAddressByStr(ctx, addressStr)
		require.NotNil(t, addressNewResp)
		require.Nil(t, addressOldResp)
		require.NoError(t, err)

		require.Equal(t, addressNewResp, addressNew)

	})

	t.Run("success old", func(t *testing.T) {
		addressType := addressPb.AddressType_ADDRESS_TYPE_BASED_ID
		userAccountId := int32(0)

		req := &addressPb.AddressStrId{
			Address: addressStr,
		}

		addressOld := &model.AddressOld{
			Id:            uuid.UUID{},
			Address:       addressStr,
			UserUuid:      uuid.UUID{},
			AddressType:   enum.NewAddressTypeWrapper(addressType),
			Network:       nodeCommon.NewNetworkEnumWrapper(network),
			UserAccountId: userAccountId,
			Coin:          *coin,
			CreatedAt:     time.Time{},
		}

		addressesOld := model.AddressesOld{addressOld}

		resp := mapping.MapModelOneOfAddressToProto(ctx, nil, addressesOld)

		handlerMock.EXPECT().
			GetAddressByStr(
				mock.Anything,
				req).
			Return(resp, nil).
			Once()

		addressNewResp, addressOldResp, err := addressRepo.GetAddressByStr(ctx, addressStr)
		require.Nil(t, addressNewResp)
		require.NotNil(t, addressOldResp)
		require.NoError(t, err)

		require.Equal(t, addressOldResp, addressOld)

	})

	t.Run("error service mock", func(t *testing.T) {
		req := &addressPb.AddressStrId{
			Address: addressStr,
		}

		mockErr := newMockError()

		handlerMock.EXPECT().
			GetAddressByStr(
				mock.Anything,
				req).
			Return(nil, mockErr).
			Once()

		addressNewResp, addressOldResp, err := addressRepo.GetAddressByStr(ctx, addressStr)
		require.Nil(t, addressNewResp)
		require.Nil(t, addressOldResp)
		require.Error(t, err)

		require.ErrorIs(t, err, mockErr)

	})
}
