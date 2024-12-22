package repository_test

import (
	"context"
	"database/sql"
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

func TestAddressRepository_GetPersonalAddressesByFilter(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handlerMock := addressPbMock.NewMockAddressServiceClient(t)
	addressRepo := repository.NewAddressRepository(handlerMock)
	network := nodeCommon.EthNetworkId

	t.Run("success", func(t *testing.T) {
		filter := &model.AddressPersonalFilter{
			Id:         nil,
			Address:    nil,
			UserUuid:   nil,
			Network:    nil,
			IsDeleted:  nil,
			Pagination: nil,
		}

		filterReq := mapping.MapModelAddressPersonalFilterToProto(filter)

		addressPersonal := &model.AddressPersonal{
			Id:        uuid.UUID{},
			Address:   "",
			UserUuid:  uuid.UUID{},
			Network:   nodeCommon.NewNetworkEnumWrapper(network),
			DeletedAt: sql.NullTime{},
			UpdatedAt: time.Time{},
			CreatedAt: time.Time{},
		}

		addressesPersonal := model.AddressesPersonal{addressPersonal}

		resp := mapping.MapModelAddressesPersonalToProto(nil, addressesPersonal)

		handlerMock.EXPECT().
			GetPersonalAddressesByFilter(
				ctx,
				filterReq).
			Return(resp, nil).
			Once()

		_, addressPersonalResp, err := addressRepo.GetPersonalAddressesByFilter(ctx, filter)
		require.NotNil(t, addressPersonalResp)
		require.NoError(t, err)

		require.Equal(t, addressPersonalResp, addressesPersonal)

	})

	t.Run("error service mock", func(t *testing.T) {
		filter := &model.AddressPersonalFilter{
			Id:         nil,
			Address:    nil,
			UserUuid:   nil,
			Network:    nil,
			IsDeleted:  nil,
			Pagination: nil,
		}

		filterReq := mapping.MapModelAddressPersonalFilterToProto(filter)

		mockErr := newMockError()

		handlerMock.EXPECT().
			GetPersonalAddressesByFilter(
				ctx,
				filterReq).
			Return(nil, mockErr).
			Once()

		_, addressOldResp, err := addressRepo.GetPersonalAddressesByFilter(ctx, filter)
		require.Nil(t, addressOldResp)
		require.Error(t, err)

		require.ErrorIs(t, err, mockErr)

	})
}
