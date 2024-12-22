package service_test

import (
	"context"
	"testing"
	"time"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	nodeRepoMock "code.emcdtech.com/emcd/blockchain/node/repository_external/mocks"
	userAccountRepositoryMock "code.emcdtech.com/emcd/service/accounting/repository/repository_mock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"code.emcdtech.com/emcd/blockchain/address/internal/service"
	externalMock "code.emcdtech.com/emcd/blockchain/address/mocks/external"
	repositoryMock "code.emcdtech.com/emcd/blockchain/address/mocks/internal_/repository"
	serviceMock "code.emcdtech.com/emcd/blockchain/address/mocks/internal_/service"
	"code.emcdtech.com/emcd/blockchain/address/model"
	"code.emcdtech.com/emcd/blockchain/address/model/enum"
	addressPb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
)

func TestAddressService_GetOldAddressesByFilter(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	addressRepoMock := repositoryMock.NewMockAddressRepository(t)
	userAccountRepoMock := userAccountRepositoryMock.NewMockUserAccountRepository(t)
	nodeAddressRepoMock := nodeRepoMock.NewMockAddressNodeRepository(t)
	profileProtoMock := externalMock.NewProfileServiceClient(t)
	rabbitServiceMock := serviceMock.NewMockRabbitService(t)

	var masterKeysIdMap = map[nodeCommon.NetworkGroupEnum][]string{
		nodeCommon.EthNetworkGroupId: {""},
	}

	s := service.NewAddressService(addressRepoMock,
		userAccountRepoMock,
		nodeAddressRepoMock,
		profileProtoMock,
		masterKeysIdMap,
		true,
		rabbitServiceMock,
	)

	addressUuid := uuid.New()
	// addressStr := ""
	// userUuid := uuid.New()
	network := nodeCommon.TrxNetworkId
	addressType := addressPb.AddressType_ADDRESS_TYPE_DIRECT

	t.Run("success", func(t *testing.T) {
		addressFilter := &model.AddressOldFilter{
			Id:            &addressUuid,
			Address:       nil,
			UserUuid:      nil,
			AddressType:   nil,
			Network:       nil,
			UserAccountId: nil,
			Coin:          nil,
			Pagination:    nil,
		}

		address := &model.AddressOld{
			Id:          addressUuid,
			Address:     "",
			UserUuid:    uuid.UUID{},
			AddressType: enum.NewAddressTypeWrapper(addressType),
			Network:     nodeCommon.NewNetworkEnumWrapper(network),
			CreatedAt:   time.Time{},
		}

		addresses := model.AddressesOld{address}

		addressRepoMock.EXPECT().
			GetOldAddresses(ctx, addressFilter).
			Return(nil, addresses, nil).Once()

		_, addressesResp, err := s.GetOldAddressesByFilter(ctx, addressFilter)
		require.NotEmpty(t, addressesResp)
		require.NoError(t, err)

		require.Len(t, addressesResp, 1)

	})

	t.Run("error mock repo", func(t *testing.T) {
		addressFilter := &model.AddressOldFilter{
			Id:            &addressUuid,
			Address:       nil,
			UserUuid:      nil,
			AddressType:   nil,
			Network:       nil,
			UserAccountId: nil,
			Coin:          nil,
			Pagination:    nil,
		}

		errMock := newMockError()

		addressRepoMock.EXPECT().
			GetOldAddresses(ctx, addressFilter).
			Return(nil, nil, errMock).Once()

		_, resp, err := s.GetOldAddressesByFilter(ctx, addressFilter)
		require.Empty(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, errMock)

	})
}
