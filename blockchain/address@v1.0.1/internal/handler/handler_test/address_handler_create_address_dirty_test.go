package handler_test

import (
	"context"
	"testing"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	coinValidatorMock "code.emcdtech.com/emcd/service/coin/repository/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"code.emcdtech.com/emcd/blockchain/address/internal/handler"
	"code.emcdtech.com/emcd/blockchain/address/internal/handler/mapping"
	externalMock "code.emcdtech.com/emcd/blockchain/address/mocks/external"
	"code.emcdtech.com/emcd/blockchain/address/mocks/internal_/service"
	"code.emcdtech.com/emcd/blockchain/address/model"
	addressPb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
	"code.emcdtech.com/emcd/blockchain/address/repository"
)

func TestAddressHandler_CreateOrUpdateDirtyAddress(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serviceAddressMock := service.NewMockAddressService(t)
	coinValidatorRepoMock := coinValidatorMock.NewMockCoinValidatorRepository(t)
	coinProtoCliMock := externalMock.NewMockCoinServiceClient(t)
	addressHandler := handler.NewAddressHandler(serviceAddressMock, coinValidatorRepoMock, coinProtoCliMock, IsNetworkOldWay)

	addressStr := uuid.NewString()
	network := nodeCommon.EthNetworkId
	isDirty := false

	addressDirtyMatch := mock.MatchedBy(func(address *model.AddressDirty) bool {

		return address.Address == addressStr &&
			address.Network.NetworkEnum == network &&
			address.IsDirty == isDirty
	})

	t.Run("dirty address", func(t *testing.T) {
		t.Run("success create", func(t *testing.T) {
			req := &addressPb.DirtyAddressForm{
				Address: addressStr,
				Network: network.ToString(),
				IsDirty: isDirty,
			}

			addressDirty, err := mapping.MapProtoAddressFormToModelDirty(req)
			require.NoError(t, err)

			serviceAddressMock.EXPECT().
				CreateOrUpdateDirtyAddress(ctx, addressDirtyMatch).
				Return(addressDirty, nil).Once()

			resp, err := addressHandler.CreateOrUpdateDirtyAddress(ctx, req)
			require.NotNil(t, resp)
			require.NoError(t, err)

			require.Equal(t, resp, req)

		})

		t.Run("error parse network", func(t *testing.T) {
			req := &addressPb.DirtyAddressForm{
				Address: addressStr,
				Network: network.ToString() + "fake",
				IsDirty: isDirty,
			}

			resp, err := addressHandler.CreateOrUpdateDirtyAddress(ctx, req)
			require.Nil(t, resp)
			require.Error(t, err)

			require.ErrorIs(t, err, repository.ErrAddr1301)
		})

		t.Run("error mock", func(t *testing.T) {
			req := &addressPb.DirtyAddressForm{
				Address: addressStr,
				Network: network.ToString(),
				IsDirty: isDirty,
			}

			errMock := newMockError()

			serviceAddressMock.EXPECT().
				CreateOrUpdateDirtyAddress(ctx, addressDirtyMatch).
				Return(nil, errMock).Once()

			resp, err := addressHandler.CreateOrUpdateDirtyAddress(ctx, req)
			require.Nil(t, resp)
			require.Error(t, err)

			require.ErrorIs(t, err, repository.ErrAddr1302)
		})

	})
}
