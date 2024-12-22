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
	addressPb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
	"code.emcdtech.com/emcd/blockchain/address/repository"
)

func TestAddressRepository_GetPersonalAddressByUserUuid(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handlerMock := addressPbMock.NewMockAddressServiceClient(t)
	addressRepo := repository.NewAddressRepository(handlerMock)

	network := nodeCommon.EthNetworkId
	userUuid := uuid.New()

	t.Run("success", func(t *testing.T) {
		filterReq := &addressPb.UserUuid{
			UserUuid: userUuid.String(),
		}

		addressPersonal := &model.AddressPersonal{
			Id:        uuid.UUID{},
			Address:   "",
			UserUuid:  userUuid,
			Network:   nodeCommon.NewNetworkEnumWrapper(network),
			DeletedAt: sql.NullTime{},
			UpdatedAt: time.Time{},
			CreatedAt: time.Time{},
		}
		addressesPersonal := model.AddressesPersonal{addressPersonal}

		resp := mapping.MapModelAddressesPersonalToProto(nil, addressesPersonal)

		handlerMock.EXPECT().
			GetPersonalAddressesByUserUuid(
				ctx,
				filterReq).
			Return(resp, nil).
			Once()

		addressPersonalResp, err := addressRepo.GetPersonalAddressByUserUuid(ctx, userUuid)
		require.NotNil(t, addressPersonalResp)
		require.NoError(t, err)

		require.Equal(t, addressPersonalResp, addressesPersonal)

	})

	t.Run("error service mock", func(t *testing.T) {
		filterReq := &addressPb.UserUuid{
			UserUuid: userUuid.String(),
		}

		mockErr := newMockError()

		handlerMock.EXPECT().
			GetPersonalAddressesByUserUuid(
				ctx,
				filterReq).
			Return(nil, mockErr).
			Once()

		addressResp, err := addressRepo.GetPersonalAddressByUserUuid(ctx, userUuid)
		require.Nil(t, addressResp)
		require.Error(t, err)

		require.ErrorIs(t, err, mockErr)

	})
}
