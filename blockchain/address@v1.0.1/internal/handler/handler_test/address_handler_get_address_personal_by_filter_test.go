package handler_test

import (
	"context"
	"database/sql"
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
	addressPb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
	"code.emcdtech.com/emcd/blockchain/address/repository"
)

func TestAddressHandler_GetPersonalAddressesByFilter(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serviceAddressMock := service.NewMockAddressService(t)
	coinValidatorRepoMock := coinValidatorMock.NewMockCoinValidatorRepository(t)
	coinProtoCli := externalMock.NewMockCoinServiceClient(t)
	addressHandler := handler.NewAddressHandler(serviceAddressMock, coinValidatorRepoMock, coinProtoCli, IsNetworkOldWay)

	network := nodeCommon.EthNetworkId

	t.Run("success", func(t *testing.T) {
		req := &addressPb.AddressPersonalFilter{
			AddressUuid: nil,
			Address:     nil,
			UserUuid:    nil,
			Network:     nil,
			IsDeleted:   nil,
			Pagination:  nil,
		}

		filter, err := mapping.MapProtoToModelAddressPersonalFilter(req)
		require.NoError(t, err)

		addressPersonalResp := &model.AddressPersonal{
			Id:        uuid.UUID{},
			Address:   "",
			UserUuid:  uuid.UUID{},
			Network:   nodeCommon.NewNetworkEnumWrapper(network),
			DeletedAt: sql.NullTime{},
			UpdatedAt: time.Time{},
			CreatedAt: time.Time{},
		}
		addressesPersonalResp := model.AddressesPersonal{addressPersonalResp}
		addressesPersonalRespProto := mapping.MapModelAddressesPersonalToProto(nil, addressesPersonalResp)

		serviceAddressMock.EXPECT().
			GetPersonalAddressesByFilter(
				ctx,
				filter,
			).
			Return(nil, addressesPersonalResp, nil).
			Once()

		resp, err := addressHandler.GetPersonalAddressesByFilter(ctx, req)
		require.NotNil(t, resp)
		require.NoError(t, err)

		require.Len(t, resp.Addresses, 1)
		require.Equal(t, resp, addressesPersonalRespProto)

	})

	t.Run("success with pagination", func(t *testing.T) {
		req := &addressPb.AddressPersonalFilter{
			AddressUuid: nil,
			Address:     nil,
			UserUuid:    nil,
			Network:     nil,
			IsDeleted:   nil,
			Pagination: &addressPb.AddressPagination{
				Limit:  0,
				Offset: 0,
			},
		}

		filter, err := mapping.MapProtoToModelAddressPersonalFilter(req)
		require.NoError(t, err)

		addressPersonalResp := &model.AddressPersonal{
			Id:        uuid.UUID{},
			Address:   "",
			UserUuid:  uuid.UUID{},
			Network:   nodeCommon.NewNetworkEnumWrapper(network),
			DeletedAt: sql.NullTime{},
			UpdatedAt: time.Time{},
			CreatedAt: time.Time{},
		}
		addressesPersonalResp := model.AddressesPersonal{addressPersonalResp}
		totalCount := uint64(99)
		addressesPersonalRespProto := mapping.MapModelAddressesPersonalToProto(&totalCount, addressesPersonalResp)

		serviceAddressMock.EXPECT().
			GetPersonalAddressesByFilter(
				ctx,
				filter,
			).
			Return(&totalCount, addressesPersonalResp, nil).
			Once()

		resp, err := addressHandler.GetPersonalAddressesByFilter(ctx, req)
		require.NotNil(t, resp)
		require.NoError(t, err)

		require.Len(t, resp.Addresses, 1)
		require.Equal(t, resp, addressesPersonalRespProto)

	})

	t.Run("error parse uuid", func(t *testing.T) {
		req := &addressPb.AddressPersonalFilter{
			AddressUuid: utils.StringToPtr("fake"),
			Address:     nil,
			UserUuid:    nil,
			Network:     nil,
			IsDeleted:   nil,
			Pagination:  nil,
		}

		resp, err := addressHandler.GetPersonalAddressesByFilter(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAddr1101)

	})

	t.Run("error parse uuid", func(t *testing.T) {
		req := &addressPb.AddressPersonalFilter{
			AddressUuid: nil,
			Address:     nil,
			UserUuid:    utils.StringToPtr("fake"),
			Network:     nil,
			IsDeleted:   nil,
			Pagination:  nil,
		}

		resp, err := addressHandler.GetPersonalAddressesByFilter(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAddr1101)

	})

	t.Run("error parse network", func(t *testing.T) {
		req := &addressPb.AddressPersonalFilter{
			AddressUuid: nil,
			Address:     nil,
			UserUuid:    nil,
			Network:     utils.StringToPtr("fake"),
			IsDeleted:   nil,
			Pagination:  nil,
		}

		resp, err := addressHandler.GetPersonalAddressesByFilter(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAddr1101)

	})

	t.Run("error mock new", func(t *testing.T) {
		req := &addressPb.AddressPersonalFilter{
			AddressUuid: nil,
			Address:     nil,
			UserUuid:    nil,
			Network:     nil,
			IsDeleted:   nil,
			Pagination:  nil,
		}

		filter, err := mapping.MapProtoToModelAddressPersonalFilter(req)
		require.NoError(t, err)

		errMock := newMockError()

		serviceAddressMock.EXPECT().
			GetPersonalAddressesByFilter(
				ctx,
				filter,
			).
			Return(nil, nil, errMock).
			Once()

		resp, err := addressHandler.GetPersonalAddressesByFilter(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAddr1102)

	})
}
