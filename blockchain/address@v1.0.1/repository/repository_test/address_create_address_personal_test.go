package repository_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"code.emcdtech.com/emcd/blockchain/address/internal/handler/mapping"
	addressPbMock "code.emcdtech.com/emcd/blockchain/address/mocks/protocol/address"
	"code.emcdtech.com/emcd/blockchain/address/model"
	addressPb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
	"code.emcdtech.com/emcd/blockchain/address/repository"
)

func TestAddressRepository_AddOrUpdatePersonalAddress(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handlerMock := addressPbMock.NewMockAddressServiceClient(t)
	addressRepo := repository.NewAddressRepository(handlerMock)

	userUuid := uuid.New()
	addressStr := uuid.NewString()
	network := nodeCommon.EthNetworkId

	t.Run("success personal", func(t *testing.T) {
		req := &addressPb.CreatePersonalAddressRequest{
			Address:  addressStr,
			UserUuid: userUuid.String(),
			Network:  network.ToString(),
		}

		addressNew := &model.AddressPersonal{
			Id:        uuid.UUID{},
			Address:   "",
			UserUuid:  userUuid,
			Network:   nodeCommon.NewNetworkEnumWrapper(network),
			DeletedAt: sql.NullTime{},
			UpdatedAt: time.Time{},
			CreatedAt: time.Time{},
		}

		resp := mapping.MapModelAddressPersonalToProto(addressNew)

		handlerMock.EXPECT().
			AddOrUpdatePersonalAddress(
				ctx,
				req).
			Return(resp, nil).
			Once()

		addressNewResp, err := addressRepo.AddOrUpdatePersonalAddress(ctx, addressStr, userUuid, network, nil)
		require.NotNil(t, addressNewResp)
		require.NoError(t, err)

		require.Equal(t, addressNewResp, addressNew)

	})

	t.Run("error service mock", func(t *testing.T) {
		req := &addressPb.CreatePersonalAddressRequest{
			Address:  addressStr,
			UserUuid: userUuid.String(),
			Network:  network.ToString(),
		}

		mockErr := newMockError()

		handlerMock.EXPECT().
			AddOrUpdatePersonalAddress(
				mock.Anything,
				req).
			Return(nil, mockErr).
			Once()

		addressNewResp, err := addressRepo.AddOrUpdatePersonalAddress(ctx, addressStr, userUuid, network, nil)
		require.Nil(t, addressNewResp)
		require.Error(t, err)

		require.ErrorIs(t, err, mockErr)

	})
}
