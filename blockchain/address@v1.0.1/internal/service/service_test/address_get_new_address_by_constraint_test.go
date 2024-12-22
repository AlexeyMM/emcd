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

	"code.emcdtech.com/emcd/blockchain/address/common/utils"
	"code.emcdtech.com/emcd/blockchain/address/internal/service"
	externalMock "code.emcdtech.com/emcd/blockchain/address/mocks/external"
	repositoryMock "code.emcdtech.com/emcd/blockchain/address/mocks/internal_/repository"
	serviceMock "code.emcdtech.com/emcd/blockchain/address/mocks/internal_/service"
	"code.emcdtech.com/emcd/blockchain/address/model"
	"code.emcdtech.com/emcd/blockchain/address/model/enum"
	addressPb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
)

func TestAddressService_GetNewAddressByConstraint(t *testing.T) {
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
	addressStr := ""
	userUuid := uuid.New()
	networkGroup := nodeCommon.EthNetworkGroupId
	addressType := addressPb.AddressType_ADDRESS_TYPE_DERIVED

	t.Run("success new get", func(t *testing.T) {
		addressFilter := &model.AddressFilter{
			Id:           nil,
			Address:      nil,
			UserUuid:     &userUuid,
			IsProcessing: utils.BoolToPtr(false),
			AddressType:  addressType.Enum(),
			NetworkGroup: &networkGroup,
			CreatedAtGt:  nil,
			Pagination:   nil,
		}

		address := &model.Address{
			Id:             addressUuid,
			Address:        addressStr,
			UserUuid:       userUuid,
			ProcessingUuid: userUuid,
			AddressType:    enum.NewAddressTypeWrapper(addressType),
			NetworkGroup:   nodeCommon.NewNetworkGroupEnumWrapper(networkGroup),
			CreatedAt:      time.Time{},
		}

		addresses := model.Addresses{address}

		addressRepoMock.EXPECT().
			GetNewAddresses(ctx, addressFilter).
			Return(nil, addresses, nil).Once()

		resp, err := s.GetNewAddressByConstraint(ctx, userUuid, addressType, networkGroup)
		require.NotEmpty(t, resp)
		require.NoError(t, err)

		require.Len(t, resp, 1)

	})

	t.Run("error mock repo", func(t *testing.T) {
		addressFilter := &model.AddressFilter{
			Id:           nil,
			Address:      nil,
			UserUuid:     &userUuid,
			IsProcessing: utils.BoolToPtr(false),
			AddressType:  addressType.Enum(),
			NetworkGroup: &networkGroup,
			CreatedAtGt:  nil,
			Pagination:   nil,
		}

		errMock := newMockError()

		addressRepoMock.EXPECT().
			GetNewAddresses(ctx, addressFilter).
			Return(nil, nil, errMock).Once()

		resp, err := s.GetNewAddressByConstraint(ctx, userUuid, addressType, networkGroup)
		require.Empty(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, errMock)

	})

	t.Run("error more one len", func(t *testing.T) {
		addressFilter := &model.AddressFilter{
			Id:           nil,
			Address:      nil,
			UserUuid:     &userUuid,
			IsProcessing: utils.BoolToPtr(false),
			AddressType:  addressType.Enum(),
			NetworkGroup: &networkGroup,
			CreatedAtGt:  nil,
			Pagination:   nil,
		}

		addresses := model.Addresses{&model.Address{}, &model.Address{}}

		addressRepoMock.EXPECT().
			GetNewAddresses(ctx, addressFilter).
			Return(nil, addresses, nil).Once()

		resp, err := s.GetNewAddressByConstraint(ctx, userUuid, addressType, networkGroup)
		require.Empty(t, resp)
		require.Error(t, err)

		require.ErrorContains(t, err, "can not determinate new address by constraint")

	})
}
