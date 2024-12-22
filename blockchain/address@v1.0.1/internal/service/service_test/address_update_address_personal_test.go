package service_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	nodeRepoMock "code.emcdtech.com/emcd/blockchain/node/repository_external/mocks"
	userAccountRepositoryMock "code.emcdtech.com/emcd/service/accounting/repository/repository_mock"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"code.emcdtech.com/emcd/blockchain/address/internal/service"
	externalMock "code.emcdtech.com/emcd/blockchain/address/mocks/external"
	repositoryMock "code.emcdtech.com/emcd/blockchain/address/mocks/internal_/repository"
	serviceMock "code.emcdtech.com/emcd/blockchain/address/mocks/internal_/service"
	"code.emcdtech.com/emcd/blockchain/address/model"
)

func TestAddressService_UpdatePersonalAddress(t *testing.T) {
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

	addressUuid := uuid.New()
	addressStr := ""
	userUuid := uuid.New()
	minPayout := 0.0
	network := nodeCommon.EthNetworkId

	address := &model.AddressPersonal{
		Id:        addressUuid,
		Address:   addressStr,
		UserUuid:  userUuid,
		Network:   nodeCommon.NewNetworkEnumWrapper(network),
		MinPayout: minPayout,
		CreatedAt: time.Time{},
	}

	partial := &model.AddressPersonalPartial{
		Address:   &addressStr,
		MinPayout: &minPayout,
		DeletedAt: &sql.NullTime{Time: time.Time{}, Valid: false},
		UpdatedAt: nil,
	}

	matchPartial := mock.MatchedBy(func(partialOther *model.AddressPersonalPartial) bool {

		return *partial.Address == *partialOther.Address &&
			*partial.MinPayout == *partialOther.MinPayout &&
			partial.DeletedAt.Valid == partialOther.DeletedAt.Valid
	})

	t.Run("success update personal address", func(t *testing.T) {
		addressRepoMock.EXPECT().
			UpdatePersonalAddress(ctx, address, matchPartial).
			Return(nil).Once()

		resp, err := s.UpdatePersonalAddress(ctx, address, addressStr, &minPayout)
		require.NotNil(t, resp)
		require.NoError(t, err)

	})

	t.Run("error mock update personal address", func(t *testing.T) {
		errMock := newMockError()

		addressRepoMock.EXPECT().
			UpdatePersonalAddress(ctx, address, matchPartial).
			Return(errMock).Once()

		resp, err := s.UpdatePersonalAddress(ctx, address, addressStr, &minPayout)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, errMock)

	})
}
