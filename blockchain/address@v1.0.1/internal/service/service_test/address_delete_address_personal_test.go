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

	"code.emcdtech.com/emcd/blockchain/address/common/utils"
	"code.emcdtech.com/emcd/blockchain/address/internal/service"
	externalMock "code.emcdtech.com/emcd/blockchain/address/mocks/external"
	repositoryMock "code.emcdtech.com/emcd/blockchain/address/mocks/internal_/repository"
	serviceMock "code.emcdtech.com/emcd/blockchain/address/mocks/internal_/service"
	"code.emcdtech.com/emcd/blockchain/address/model"
)

func TestAddressService_DeletePersonalAddress(t *testing.T) {
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

	t.Run("success delete personal address", func(t *testing.T) {
		network := nodeCommon.EthNetworkId

		address := &model.AddressPersonal{
			Id:        addressUuid,
			Address:   addressStr,
			UserUuid:  userUuid,
			Network:   nodeCommon.NewNetworkEnumWrapper(network),
			CreatedAt: time.Time{},
		}

		filter := &model.AddressPersonalFilter{
			Id:         nil,
			Address:    nil,
			UserUuid:   &userUuid,
			Network:    network.ToPtr(),
			IsDeleted:  nil,
			Pagination: nil,
		}

		t.Run("success not found way", func(t *testing.T) {
			addressRepoMock.EXPECT().
				GetPersonalAddresses(ctx, filter).
				Return(nil, nil, nil).Once()

			err := s.DeletePersonalAddress(ctx, userUuid, network)
			require.NoError(t, err)

		})

		t.Run("success exists way", func(t *testing.T) {
			addreses := model.AddressesPersonal{address}

			addressRepoMock.EXPECT().
				GetPersonalAddresses(ctx, filter).
				Return(nil, addreses, nil).Once()

			partial := &model.AddressPersonalPartial{
				Address:   utils.StringToPtr(""),
				DeletedAt: &sql.NullTime{Time: time.Time{}, Valid: true},
				UpdatedAt: nil,
			}

			matchPartial := mock.MatchedBy(func(partialOther *model.AddressPersonalPartial) bool {

				return *partial.Address == *partialOther.Address &&
					partial.DeletedAt.Valid == partialOther.DeletedAt.Valid
			})

			addressRepoMock.EXPECT().
				UpdatePersonalAddress(ctx, address, matchPartial).
				Return(nil)

			err := s.DeletePersonalAddress(ctx, userUuid, network)
			require.NoError(t, err)

		})
	})
}
