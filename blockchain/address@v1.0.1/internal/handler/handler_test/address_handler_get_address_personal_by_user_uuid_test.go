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

func TestAddressHandler_GetPersonalAddressesByUserUuid(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serviceAddressMock := service.NewMockAddressService(t)
	coinValidatorRepoMock := coinValidatorMock.NewMockCoinValidatorRepository(t)
	coinProtoCli := externalMock.NewMockCoinServiceClient(t)
	addressHandler := handler.NewAddressHandler(serviceAddressMock, coinValidatorRepoMock, coinProtoCli, IsNetworkOldWay)

	network := nodeCommon.EthNetworkId
	userUuid := uuid.New()

	t.Run("success", func(t *testing.T) {
		req := &addressPb.UserUuid{
			UserUuid: userUuid.String(),
		}

		reqFilter := &addressPb.AddressPersonalFilter{
			AddressUuid: nil,
			Address:     nil,
			UserUuid:    &req.UserUuid,
			Network:     nil,
			IsDeleted:   utils.BoolToPtr(false),
			Pagination:  nil,
		}

		filter, err := mapping.MapProtoToModelAddressPersonalFilter(reqFilter)
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

		resp, err := addressHandler.GetPersonalAddressesByUserUuid(ctx, req)
		require.NotNil(t, resp)
		require.NoError(t, err)

		require.NotNil(t, resp)
		require.Equal(t, resp, addressesPersonalRespProto)

	})

	t.Run("success empty", func(t *testing.T) {
		req := &addressPb.UserUuid{
			UserUuid: userUuid.String(),
		}

		reqFilter := &addressPb.AddressPersonalFilter{
			AddressUuid: nil,
			Address:     nil,
			UserUuid:    &req.UserUuid,
			Network:     nil,
			IsDeleted:   utils.BoolToPtr(false),
			Pagination:  nil,
		}

		filter, err := mapping.MapProtoToModelAddressPersonalFilter(reqFilter)
		require.NoError(t, err)

		addressesPersonalResp := model.AddressesPersonal{}

		serviceAddressMock.EXPECT().
			GetPersonalAddressesByFilter(
				ctx,
				filter,
			).
			Return(nil, addressesPersonalResp, nil).
			Once()

		resp, err := addressHandler.GetPersonalAddressesByUserUuid(ctx, req)
		require.Empty(t, resp)
		require.NoError(t, err)

	})

	t.Run("error parse user uuid", func(t *testing.T) {
		req := &addressPb.UserUuid{
			UserUuid: userUuid.String() + "fake",
		}

		resp, err := addressHandler.GetPersonalAddressesByUserUuid(ctx, req)
		require.Empty(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAddr1201)

	})

	t.Run("error mock new", func(t *testing.T) {
		req := &addressPb.UserUuid{
			UserUuid: userUuid.String(),
		}

		reqFilter := &addressPb.AddressPersonalFilter{
			AddressUuid: nil,
			Address:     nil,
			UserUuid:    &req.UserUuid,
			Network:     nil,
			IsDeleted:   utils.BoolToPtr(false),
			Pagination:  nil,
		}

		filter, err := mapping.MapProtoToModelAddressPersonalFilter(reqFilter)
		require.NoError(t, err)

		errMock := newMockError()

		serviceAddressMock.EXPECT().
			GetPersonalAddressesByFilter(
				ctx,
				filter,
			).
			Return(nil, nil, errMock).
			Once()

		resp, err := addressHandler.GetPersonalAddressesByUserUuid(ctx, req)
		require.Empty(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAddr1202)

	})

	t.Run("multiple return", func(t *testing.T) {
		req := &addressPb.UserUuid{
			UserUuid: userUuid.String(),
		}

		reqFilter := &addressPb.AddressPersonalFilter{
			AddressUuid: nil,
			Address:     nil,
			UserUuid:    &req.UserUuid,
			Network:     nil,
			IsDeleted:   utils.BoolToPtr(false),
			Pagination:  nil,
		}

		filter, err := mapping.MapProtoToModelAddressPersonalFilter(reqFilter)
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
		addressesPersonalResp := model.AddressesPersonal{addressPersonalResp, addressPersonalResp}

		serviceAddressMock.EXPECT().
			GetPersonalAddressesByFilter(
				ctx,
				filter,
			).
			Return(nil, addressesPersonalResp, nil).
			Once()

		resp, err := addressHandler.GetPersonalAddressesByUserUuid(ctx, req)
		require.Empty(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAddr1203)

	})

}
