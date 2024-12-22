package handler_test

import (
	"context"
	"testing"
	"time"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	coinValidatorMock "code.emcdtech.com/emcd/service/coin/repository/mocks"
	"github.com/stretchr/testify/require"

	"code.emcdtech.com/emcd/blockchain/address/common/utils"
	"code.emcdtech.com/emcd/blockchain/address/internal/handler"
	"code.emcdtech.com/emcd/blockchain/address/internal/handler/mapping"
	externalMock "code.emcdtech.com/emcd/blockchain/address/mocks/external"
	"code.emcdtech.com/emcd/blockchain/address/mocks/internal_/service"
	"code.emcdtech.com/emcd/blockchain/address/model"
	addressPb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
	"code.emcdtech.com/emcd/blockchain/address/repository"
)

func TestAddressHandler_GetDirtyAddressesByFilter(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serviceAddressMock := service.NewMockAddressService(t)
	coinValidatorRepoMock := coinValidatorMock.NewMockCoinValidatorRepository(t)
	coinProtoCli := externalMock.NewMockCoinServiceClient(t)
	addressHandler := handler.NewAddressHandler(serviceAddressMock, coinValidatorRepoMock, coinProtoCli, IsNetworkOldWay)

	network := nodeCommon.EthNetworkId

	t.Run("success", func(t *testing.T) {
		req := &addressPb.DirtyAddressFilter{
			Address: nil,
			Network: nil,
		}

		filter, err := mapping.MapProtoToModelAddressDirtyFilter(req)
		require.NoError(t, err)

		addressDirtyResp := &model.AddressDirty{
			Address:   "",
			Network:   nodeCommon.NewNetworkEnumWrapper(network),
			IsDirty:   false,
			UpdatedAt: time.Time{},
			CreatedAt: time.Time{},
		}
		addressesDirtyResp := model.AddressesDirty{addressDirtyResp}
		addressesPersonalRespProto := mapping.MapModelDirtyAddressesToProto(addressesDirtyResp)

		serviceAddressMock.EXPECT().
			GetDirtyAddressesByFilter(
				ctx,
				filter,
			).
			Return(addressesDirtyResp, nil).
			Once()

		resp, err := addressHandler.GetDirtyAddressesByFilter(ctx, req)
		require.NotNil(t, resp)
		require.NoError(t, err)

		require.Len(t, resp.Addresses, 1)
		require.Equal(t, resp, addressesPersonalRespProto)

	})

	t.Run("error parse network", func(t *testing.T) {
		req := &addressPb.DirtyAddressFilter{
			Address: nil,
			Network: utils.StringToPtr("fake"),
		}

		resp, err := addressHandler.GetDirtyAddressesByFilter(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAddr1401)

	})

	t.Run("error mock new", func(t *testing.T) {
		req := &addressPb.DirtyAddressFilter{
			Address: nil,
			Network: nil,
		}

		filter, err := mapping.MapProtoToModelAddressDirtyFilter(req)
		require.NoError(t, err)

		errMock := newMockError()

		serviceAddressMock.EXPECT().
			GetDirtyAddressesByFilter(
				ctx,
				filter,
			).
			Return(nil, errMock).
			Once()

		resp, err := addressHandler.GetDirtyAddressesByFilter(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAddr1402)

	})
}
