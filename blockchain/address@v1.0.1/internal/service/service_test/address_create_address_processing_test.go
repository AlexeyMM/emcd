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
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"code.emcdtech.com/emcd/blockchain/address/internal/repository"
	"code.emcdtech.com/emcd/blockchain/address/internal/service"
	externalMock "code.emcdtech.com/emcd/blockchain/address/mocks/external"
	repositoryMock "code.emcdtech.com/emcd/blockchain/address/mocks/internal_/repository"
	serviceMock "code.emcdtech.com/emcd/blockchain/address/mocks/internal_/service"
	"code.emcdtech.com/emcd/blockchain/address/model"
	"code.emcdtech.com/emcd/blockchain/address/model/enum"
	addressPb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
)

func TestAddressService_CreateProcessingAddress(t *testing.T) {
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
	processingUuid := uuid.New()
	masterKeyId := uint32(0)

	t.Run("success new derived", func(t *testing.T) {
		networkGroup := nodeCommon.EthNetworkGroupId
		sendConsumerNetworkGroupReloadAddr(t, ctx, rabbitServiceMock, networkGroup)

		addressType := addressPb.AddressType_ADDRESS_TYPE_DERIVED

		address := &model.Address{
			Id:             addressUuid,
			Address:        addressStr,
			UserUuid:       userUuid,
			ProcessingUuid: processingUuid,
			AddressType:    enum.NewAddressTypeWrapper(addressType),
			NetworkGroup:   nodeCommon.NewNetworkGroupEnumWrapper(networkGroup),
			CreatedAt:      time.Time{},
		}

		matchDerivedFunc := mock.MatchedBy(func(_ repository.DerivedFunc) bool { return true })

		matchAddress := mock.MatchedBy(func(addressOther *model.Address) bool {

			return address.Id == addressOther.Id &&
				address.Address == addressOther.Address &&
				address.UserUuid == addressOther.UserUuid &&
				address.AddressType == addressOther.AddressType &&
				address.NetworkGroup == addressOther.NetworkGroup
		})

		addressRepoMock.EXPECT().
			AddNewDerivedAddress(ctx, matchAddress, masterKeyId, matchDerivedFunc).
			Return(nil)

		resp, err := s.CreateProcessingAddress(ctx, addressUuid, userUuid, processingUuid, addressType, networkGroup)
		require.NotNil(t, resp)
		require.NoError(t, err)

	})

	t.Run("success new direct", func(t *testing.T) {
		networkGroup := nodeCommon.TrxNetworkGroupId
		addressType := addressPb.AddressType_ADDRESS_TYPE_DIRECT
		var userAccountId *int32

		addressStr := "address"

		sendConsumerNetworkGroupReloadAddr(t, ctx, rabbitServiceMock, networkGroup)

		address := &model.Address{
			Id:             addressUuid,
			Address:        addressStr,
			UserUuid:       userUuid,
			ProcessingUuid: processingUuid,
			AddressType:    enum.NewAddressTypeWrapper(addressType),
			NetworkGroup:   nodeCommon.NewNetworkGroupEnumWrapper(networkGroup),
			CreatedAt:      time.Time{},
		}

		matchAddress := mock.MatchedBy(func(addressOther *model.Address) bool {

			return address.Id == addressOther.Id &&
				address.Address == addressOther.Address &&
				address.UserUuid == addressOther.UserUuid &&
				address.AddressType == addressOther.AddressType &&
				address.NetworkGroup == addressOther.NetworkGroup
		})

		nodeAddressRepoMock.EXPECT().
			GenerateAddress(ctx, networkGroup.GetNetworks()[0], userAccountId, &userUuid).
			Return(addressStr, nil).Once()

		addressRepoMock.EXPECT().
			AddNewCommonAddress(ctx, matchAddress).
			Return(nil).Once()

		resp, err := s.CreateProcessingAddress(ctx, addressUuid, userUuid, processingUuid, addressType, networkGroup)
		require.NotNil(t, resp)
		require.NoError(t, err)

	})

	t.Run("success new memo", func(t *testing.T) {
		networkGroup := nodeCommon.TonNetworkGroupId
		addressType := addressPb.AddressType_ADDRESS_TYPE_MEMO

		addressStr := "address"

		sendConsumerNetworkGroupReloadAddr(t, ctx, rabbitServiceMock, networkGroup)

		address := &model.Address{
			Id:             addressUuid,
			Address:        addressStr,
			UserUuid:       userUuid,
			ProcessingUuid: processingUuid,
			AddressType:    enum.NewAddressTypeWrapper(addressType),
			NetworkGroup:   nodeCommon.NewNetworkGroupEnumWrapper(networkGroup),
			CreatedAt:      time.Time{},
		}

		matchAddress := mock.MatchedBy(func(addressOther *model.Address) bool {

			return address.Id == addressOther.Id &&
				// address.Address == addressOther.Address && // for memo
				address.UserUuid == addressOther.UserUuid &&
				address.AddressType == addressOther.AddressType &&
				address.NetworkGroup == addressOther.NetworkGroup
		})

		addressRepoMock.EXPECT().
			AddNewCommonAddress(ctx, matchAddress).
			Return(nil).Once()

		resp, err := s.CreateProcessingAddress(ctx, addressUuid, userUuid, processingUuid, addressType, networkGroup)
		require.NotNil(t, resp)
		require.NoError(t, err)

	})
}
