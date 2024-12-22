package service_test

import (
	"context"
	"testing"
	"time"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	nodeRepoMock "code.emcdtech.com/emcd/blockchain/node/repository_external/mocks"
	userAccountRepositoryMock "code.emcdtech.com/emcd/service/accounting/repository/repository_mock"
	"github.com/stretchr/testify/require"

	"code.emcdtech.com/emcd/blockchain/address/internal/service"
	externalMock "code.emcdtech.com/emcd/blockchain/address/mocks/external"
	repositoryMock "code.emcdtech.com/emcd/blockchain/address/mocks/internal_/repository"
	serviceMock "code.emcdtech.com/emcd/blockchain/address/mocks/internal_/service"
	"code.emcdtech.com/emcd/blockchain/address/model"
)

func TestAddressService_GetDirtyAddressesByFilter(t *testing.T) {
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

	addressFilter := &model.AddressDirtyFilter{
		Address: nil,
		Network: nil,
	}

	t.Run("success", func(t *testing.T) {
		address := &model.AddressDirty{
			Address:   "",
			Network:   nodeCommon.NetworkEnumWrapper{},
			IsDirty:   false,
			UpdatedAt: time.Time{},
			CreatedAt: time.Time{},
		}

		addresses := model.AddressesDirty{address}

		addressRepoMock.EXPECT().
			GetDirtyAddresses(ctx, addressFilter).
			Return(addresses, nil).Once()

		addressesResp, err := s.GetDirtyAddressesByFilter(ctx, addressFilter)
		require.NotEmpty(t, addressesResp)
		require.NoError(t, err)

		require.Len(t, addressesResp, 1)

	})

	t.Run("error mock repo", func(t *testing.T) {
		errMock := newMockError()

		addressRepoMock.EXPECT().
			GetDirtyAddresses(ctx, addressFilter).
			Return(nil, errMock).Once()

		resp, err := s.GetDirtyAddressesByFilter(ctx, addressFilter)
		require.Empty(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, errMock)

	})
}
