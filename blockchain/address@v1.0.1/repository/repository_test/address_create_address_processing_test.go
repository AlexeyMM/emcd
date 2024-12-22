package repository_test

import (
	"context"
	"testing"
	"time"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"code.emcdtech.com/emcd/blockchain/address/internal/handler/mapping"
	addressPbMock "code.emcdtech.com/emcd/blockchain/address/mocks/protocol/address"
	"code.emcdtech.com/emcd/blockchain/address/model"
	"code.emcdtech.com/emcd/blockchain/address/model/enum"
	addressPb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
	"code.emcdtech.com/emcd/blockchain/address/repository"
)

func TestAddressRepository_CreateProcessingAddress(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handlerMock := addressPbMock.NewMockAddressServiceClient(t)
	addressRepo := repository.NewAddressRepository(handlerMock)

	userUuid := uuid.New()
	processingUuid := uuid.New()
	network := nodeCommon.EthNetworkId

	t.Run("success new", func(t *testing.T) {
		addressType := addressPb.AddressType_ADDRESS_TYPE_DERIVED

		req := &addressPb.CreateProcessingAddressRequest{
			UserUuid:       userUuid.String(),
			Network:        network.ToString(),
			ProcessingUuid: processingUuid.String(),
		}

		addressNew := &model.Address{
			Id:             uuid.UUID{},
			Address:        "",
			UserUuid:       userUuid,
			ProcessingUuid: processingUuid,
			AddressType:    enum.NewAddressTypeWrapper(addressType),
			NetworkGroup:   nodeCommon.NewNetworkGroupEnumWrapper(network.Group()),
			CreatedAt:      time.Time{},
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
			CreateProcessingAddress(
				ctx,
				req).
			Return(resp, nil).
			Once()

		addressNewResp, addressOldResp, err := addressRepo.CreateProcessingAddress(ctx, userUuid, processingUuid, network)
		require.NotNil(t, addressNewResp)
		require.Nil(t, addressOldResp)
		require.NoError(t, err)

		require.Equal(t, addressNewResp, addressNew)

	})

	t.Run("error service mock", func(t *testing.T) {
		req := &addressPb.CreateProcessingAddressRequest{
			UserUuid:       userUuid.String(),
			Network:        network.ToString(),
			ProcessingUuid: processingUuid.String(),
		}

		mockErr := newMockError()

		handlerMock.EXPECT().
			CreateProcessingAddress(
				mock.Anything,
				req).
			Return(nil, mockErr).
			Once()

		addressNewResp, addressOldResp, err := addressRepo.CreateProcessingAddress(ctx, userUuid, processingUuid, network)
		require.Nil(t, addressNewResp)
		require.Nil(t, addressOldResp)
		require.Error(t, err)

		require.ErrorIs(t, err, mockErr)

	})
}
