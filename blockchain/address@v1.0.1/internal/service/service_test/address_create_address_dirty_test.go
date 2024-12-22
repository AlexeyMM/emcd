package service_test

import (
	"context"
	"testing"
	"time"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	nodeRepoMock "code.emcdtech.com/emcd/blockchain/node/repository_external/mocks"
	userAccountRepositoryMock "code.emcdtech.com/emcd/service/accounting/repository/repository_mock"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/require"

	"code.emcdtech.com/emcd/blockchain/address/internal/service"
	externalMock "code.emcdtech.com/emcd/blockchain/address/mocks/external"
	repositoryMock "code.emcdtech.com/emcd/blockchain/address/mocks/internal_/repository"
	serviceMock "code.emcdtech.com/emcd/blockchain/address/mocks/internal_/service"
	"code.emcdtech.com/emcd/blockchain/address/model"
)

func TestAddressService_CreateOrUpdateDirtyAddress(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	addressRepoMock := repositoryMock.NewMockAddressRepository(t)
	userAccountRepoMock := userAccountRepositoryMock.NewMockUserAccountRepository(t)
	nodeAddressRepoMock := nodeRepoMock.NewMockAddressNodeRepository(t)
	profileProtoMock := externalMock.NewProfileServiceClient(t)
	rabbitServiceMock := serviceMock.NewMockRabbitService(t)

	seed, err := hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)
	require.NoError(t, err)

	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	require.NoError(t, err)

	var masterKeysIdMap = map[nodeCommon.NetworkGroupEnum][]string{
		nodeCommon.EthNetworkGroupId: {masterKey.String()},
	}

	s := service.NewAddressService(addressRepoMock,
		userAccountRepoMock,
		nodeAddressRepoMock,
		profileProtoMock,
		masterKeysIdMap,
		true,
		rabbitServiceMock,
	)

	addressStr := ""
	network := nodeCommon.EthNetworkId

	address := &model.AddressDirty{
		Address:   addressStr,
		Network:   nodeCommon.NewNetworkEnumWrapper(network),
		IsDirty:   true,
		UpdatedAt: time.Time{},
		CreatedAt: time.Time{},
	}

	t.Run("success create dirty address", func(t *testing.T) {
		addressRepoMock.EXPECT().
			AddOrUpdateDirtyAddress(ctx, address).
			Return(nil).Once()

		resp, err := s.CreateOrUpdateDirtyAddress(ctx, address)
		require.NotNil(t, resp)
		require.NoError(t, err)

	})

	t.Run("error mock repository", func(t *testing.T) {
		errMock := newMockError()

		addressRepoMock.EXPECT().
			AddOrUpdateDirtyAddress(ctx, address).
			Return(errMock).Once()

		resp, err := s.CreateOrUpdateDirtyAddress(ctx, address)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, errMock)
	})
}
