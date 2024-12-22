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
	"code.emcdtech.com/emcd/blockchain/address/repository"
)

func TestAddressRepository_CreateOrUpdateDirtyAddress(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handlerMock := addressPbMock.NewMockAddressServiceClient(t)
	addressRepo := repository.NewAddressRepository(handlerMock)

	addressStr := uuid.NewString()
	network := nodeCommon.EthNetworkId

	addressDirty := &model.AddressDirty{
		Address:   addressStr,
		Network:   nodeCommon.NewNetworkEnumWrapper(network),
		IsDirty:   false,
		UpdatedAt: time.Time{},
		CreatedAt: time.Time{},
	}

	t.Run("success dirty", func(t *testing.T) {
		resp := mapping.MapModelAddressDirtyToProto(addressDirty)
		req := resp

		handlerMock.EXPECT().
			CreateOrUpdateDirtyAddress(
				ctx,
				req).
			Return(resp, nil).
			Once()

		addressNewResp, err := addressRepo.CreateOrUpdateDirtyAddress(ctx, addressDirty)
		require.NotNil(t, addressNewResp)
		require.NoError(t, err)

		require.Equal(t, addressNewResp, addressDirty)

	})

	t.Run("error service mock", func(t *testing.T) {
		resp := mapping.MapModelAddressDirtyToProto(addressDirty)
		req := resp

		mockErr := newMockError()

		handlerMock.EXPECT().
			CreateOrUpdateDirtyAddress(
				ctx,
				req).
			Return(nil, mockErr).
			Once()

		addressNewResp, err := addressRepo.CreateOrUpdateDirtyAddress(ctx, addressDirty)
		require.Nil(t, addressNewResp)
		require.Error(t, err)

		require.ErrorIs(t, err, mockErr)

	})
}
