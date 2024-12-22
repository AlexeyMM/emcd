package service_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	nodeRepoMock "code.emcdtech.com/emcd/blockchain/node/repository_external/mocks"
	userAccountModel "code.emcdtech.com/emcd/service/accounting/model"
	userAccountModelEnum "code.emcdtech.com/emcd/service/accounting/model/enum"
	userAccountRepositoryMock "code.emcdtech.com/emcd/service/accounting/repository/repository_mock"
	profilePb "code.emcdtech.com/emcd/service/profile/protocol/profile"
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
	"code.emcdtech.com/emcd/blockchain/address/model/enum"
	addressPb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
)

func TestAddressService_CreateOldAddress(t *testing.T) {
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

	t.Run("success old rpc based none", func(t *testing.T) {
		network := nodeCommon.BtcNetworkId
		addressType := addressPb.AddressType_ADDRESS_TYPE_BASED_NONE
		coin := "btc"
		userId := int32(1)
		userAccountId := int32(1)
		userUuid := uuid.New()

		addressStr := "address"

		sendConsumerNetworkReloadAddr(t, ctx, rabbitServiceMock, network)

		profileReq := &profilePb.GetByUserIDV2Request{
			UserID: userUuid.String(),
		}

		profile := &profilePb.GetByUserIDV2Response{
			Profile: &profilePb.ProfileV2{
				User: &profilePb.UserV2{
					ID:    userUuid.String(),
					OldID: userId,
				},
				Commissions: nil,
			},
		}

		userAccount := &userAccountModel.UserAccount{
			ID:            0,
			UserID:        userId,
			CoinID:        0, // will ignore
			AccountTypeID: userAccountModelEnum.NewAccountTypeIdWrapper(userAccountModelEnum.WalletAccountTypeID),
			Minpay:        0.0,
			Address:       sql.NullString{},
			ChangedAt:     sql.NullTime{},
			Img1:          sql.NullFloat64{},
			Img2:          sql.NullFloat64{},
			IsActive:      utils.BoolToBoolNull(true),
			CreatedAt:     sql.NullTime{},
			UpdatedAt:     sql.NullTime{},
			Fee:           utils.Float64ToFloat64Null(0.0),
			UserIDNew:     utils.UuidToUuidNull(userUuid),
			CoinNew:       utils.StringToStringNull(coin),
		}

		userAccountResp := userModelCopy(userAccount)
		userAccountResp.ID = userAccountId

		userAccounts := userAccountModel.UserAccounts{userAccount}
		userAccountsResp := userAccountModel.UserAccounts{userAccountResp}

		address := &model.AddressOld{
			Id:          addressUuid,
			Address:     addressStr,
			UserUuid:    userUuid,
			AddressType: enum.NewAddressTypeWrapper(addressType),
			Network:     nodeCommon.NewNetworkEnumWrapper(network),
			CreatedAt:   time.Time{},
		}

		matchAddress := mock.MatchedBy(func(addressOther *model.AddressOld) bool {

			return address.Id == addressOther.Id &&
				address.Address == addressOther.Address &&
				address.UserUuid == addressOther.UserUuid &&
				address.AddressType == addressOther.AddressType &&
				address.Network == addressOther.Network
		})

		profileProtoMock.EXPECT().
			GetByUserIDV2(ctx, profileReq).
			Return(profile, nil).Once()

		userAccountRepoMock.EXPECT().
			CreateUserAccounts(ctx, userId, userUuid, userAccounts).
			Return(userAccountsResp, nil).Once()

		nodeAddressRepoMock.EXPECT().
			GenerateAddress(ctx, network, &userAccountId, &userUuid).
			Return(addressStr, nil).Once()

		addressRepoMock.EXPECT().
			AddOldAddress(ctx, matchAddress).
			Return(nil).Once()

		resp, err := s.CreateOldAddress(ctx, addressUuid, userUuid, addressType, network, coin)
		require.NotNil(t, resp)
		require.NoError(t, err)

	})

	t.Run("success old rpc based id", func(t *testing.T) {
		network := nodeCommon.EthNetworkId
		addressType := addressPb.AddressType_ADDRESS_TYPE_BASED_ID
		coin := "eth"
		userId := int32(1)
		userAccountId := int32(0)
		userUuid := uuid.New()

		addressStr := "address"

		sendConsumerNetworkReloadAddr(t, ctx, rabbitServiceMock, network)

		profileReq := &profilePb.GetByUserIDV2Request{
			UserID: userUuid.String(),
		}

		profile := &profilePb.GetByUserIDV2Response{
			Profile: &profilePb.ProfileV2{
				User: &profilePb.UserV2{
					ID:    userUuid.String(),
					OldID: userId,
				},
				Commissions: nil,
			},
		}

		userAccount := &userAccountModel.UserAccount{
			ID:            0,
			UserID:        userId,
			CoinID:        0, // will ignore
			AccountTypeID: userAccountModelEnum.NewAccountTypeIdWrapper(userAccountModelEnum.WalletAccountTypeID),
			Minpay:        0.0,
			Address:       sql.NullString{},
			ChangedAt:     sql.NullTime{},
			Img1:          sql.NullFloat64{},
			Img2:          sql.NullFloat64{},
			IsActive:      utils.BoolToBoolNull(true),
			CreatedAt:     sql.NullTime{},
			UpdatedAt:     sql.NullTime{},
			Fee:           utils.Float64ToFloat64Null(0.0),
			UserIDNew:     utils.UuidToUuidNull(userUuid),
			CoinNew:       utils.StringToStringNull(coin),
		}

		userAccountResp := userModelCopy(userAccount)
		userAccountResp.ID = userAccountId

		userAccounts := userAccountModel.UserAccounts{userAccount}
		userAccountsResp := userAccountModel.UserAccounts{userAccountResp}

		address := &model.AddressOld{
			Id:          addressUuid,
			Address:     addressStr,
			UserUuid:    userUuid,
			AddressType: enum.NewAddressTypeWrapper(addressType),
			Network:     nodeCommon.NewNetworkEnumWrapper(network),
			CreatedAt:   time.Time{},
		}

		matchAddress := mock.MatchedBy(func(addressOther *model.AddressOld) bool {

			return address.Id == addressOther.Id &&
				address.Address == addressOther.Address &&
				address.UserUuid == addressOther.UserUuid &&
				address.AddressType == addressOther.AddressType &&
				address.Network == addressOther.Network
		})

		profileProtoMock.EXPECT().
			GetByUserIDV2(ctx, profileReq).
			Return(profile, nil).Once()

		userAccountRepoMock.EXPECT().
			CreateUserAccounts(ctx, userId, userUuid, userAccounts).
			Return(userAccountsResp, nil).Once()

		nodeAddressRepoMock.EXPECT().
			GenerateAddress(ctx, network, &userAccountId, &userUuid).
			Return(addressStr, nil).Once()

		addressRepoMock.EXPECT().
			AddOldAddress(ctx, matchAddress).
			Return(nil).Once()

		resp, err := s.CreateOldAddress(ctx, addressUuid, userUuid, addressType, network, coin)
		require.NotNil(t, resp)
		require.NoError(t, err)

	})

	t.Run("success old direct", func(t *testing.T) {
		network := nodeCommon.TrxNetworkId
		addressType := addressPb.AddressType_ADDRESS_TYPE_DIRECT
		coin := "usdt"
		userId := int32(1)
		userAccountId := int32(0)
		userUuid := uuid.New()

		addressStr := "address"

		sendConsumerNetworkReloadAddr(t, ctx, rabbitServiceMock, network)

		profileReq := &profilePb.GetByUserIDV2Request{
			UserID: userUuid.String(),
		}

		profile := &profilePb.GetByUserIDV2Response{
			Profile: &profilePb.ProfileV2{
				User: &profilePb.UserV2{
					ID:    userUuid.String(),
					OldID: userId,
				},
				Commissions: nil,
			},
		}

		userAccount := &userAccountModel.UserAccount{
			ID:            0,
			UserID:        userId,
			CoinID:        0, // will ignore
			AccountTypeID: userAccountModelEnum.NewAccountTypeIdWrapper(userAccountModelEnum.WalletAccountTypeID),
			Minpay:        0.0,
			Address:       sql.NullString{},
			ChangedAt:     sql.NullTime{},
			Img1:          sql.NullFloat64{},
			Img2:          sql.NullFloat64{},
			IsActive:      utils.BoolToBoolNull(true),
			CreatedAt:     sql.NullTime{},
			UpdatedAt:     sql.NullTime{},
			Fee:           utils.Float64ToFloat64Null(0.0),
			UserIDNew:     utils.UuidToUuidNull(userUuid),
			CoinNew:       utils.StringToStringNull(coin),
		}

		userAccountResp := userModelCopy(userAccount)
		userAccountResp.ID = userAccountId

		userAccounts := userAccountModel.UserAccounts{userAccount}
		userAccountsResp := userAccountModel.UserAccounts{userAccountResp}

		address := &model.AddressOld{
			Id:          addressUuid,
			Address:     addressStr,
			UserUuid:    userUuid,
			AddressType: enum.NewAddressTypeWrapper(addressType),
			Network:     nodeCommon.NewNetworkEnumWrapper(network),
			CreatedAt:   time.Time{},
		}

		matchAddress := mock.MatchedBy(func(addressOther *model.AddressOld) bool {

			return address.Id == addressOther.Id &&
				address.Address == addressOther.Address &&
				address.UserUuid == addressOther.UserUuid &&
				address.AddressType == addressOther.AddressType &&
				address.Network == addressOther.Network
		})

		profileProtoMock.EXPECT().
			GetByUserIDV2(ctx, profileReq).
			Return(profile, nil).Once()

		userAccountRepoMock.EXPECT().
			CreateUserAccounts(ctx, userId, userUuid, userAccounts).
			Return(userAccountsResp, nil).Once()

		nodeAddressRepoMock.EXPECT().
			GenerateAddress(ctx, network, &userAccountId, &userUuid).
			Return(addressStr, nil).Once()

		addressRepoMock.EXPECT().
			AddOldAddress(ctx, matchAddress).
			Return(nil).Once()

		resp, err := s.CreateOldAddress(ctx, addressUuid, userUuid, addressType, network, coin)
		require.NotNil(t, resp)
		require.NoError(t, err)

	})

	t.Run("success old memo", func(t *testing.T) {
		network := nodeCommon.TonNetworkId
		addressType := addressPb.AddressType_ADDRESS_TYPE_MEMO
		coin := "usdt"
		userId := int32(1)
		userAccountId := int32(0)
		userUuid := uuid.New()

		addressStr := "address"

		sendConsumerNetworkReloadAddr(t, ctx, rabbitServiceMock, network)

		profileReq := &profilePb.GetByUserIDV2Request{
			UserID: userUuid.String(),
		}

		profile := &profilePb.GetByUserIDV2Response{
			Profile: &profilePb.ProfileV2{
				User: &profilePb.UserV2{
					ID:    userUuid.String(),
					OldID: userId,
				},
				Commissions: nil,
			},
		}

		userAccount := &userAccountModel.UserAccount{
			ID:            0,
			UserID:        userId,
			CoinID:        0, // will ignore
			AccountTypeID: userAccountModelEnum.NewAccountTypeIdWrapper(userAccountModelEnum.WalletAccountTypeID),
			Minpay:        0.0,
			Address:       sql.NullString{},
			ChangedAt:     sql.NullTime{},
			Img1:          sql.NullFloat64{},
			Img2:          sql.NullFloat64{},
			IsActive:      utils.BoolToBoolNull(true),
			CreatedAt:     sql.NullTime{},
			UpdatedAt:     sql.NullTime{},
			Fee:           utils.Float64ToFloat64Null(0.0),
			UserIDNew:     utils.UuidToUuidNull(userUuid),
			CoinNew:       utils.StringToStringNull(coin),
		}

		userAccountResp := userModelCopy(userAccount)
		userAccountResp.ID = userAccountId

		userAccounts := userAccountModel.UserAccounts{userAccount}
		userAccountsResp := userAccountModel.UserAccounts{userAccountResp}

		address := &model.AddressOld{
			Id:          addressUuid,
			Address:     addressStr,
			UserUuid:    userUuid,
			AddressType: enum.NewAddressTypeWrapper(addressType),
			Network:     nodeCommon.NewNetworkEnumWrapper(network),
			CreatedAt:   time.Time{},
		}

		matchAddress := mock.MatchedBy(func(addressOther *model.AddressOld) bool {

			return address.Id == addressOther.Id &&
				// address.Address == addressOther.Address && // for memo
				address.UserUuid == addressOther.UserUuid &&
				address.AddressType == addressOther.AddressType &&
				address.Network == addressOther.Network
		})

		profileProtoMock.EXPECT().
			GetByUserIDV2(ctx, profileReq).
			Return(profile, nil).Once()

		userAccountRepoMock.EXPECT().
			CreateUserAccounts(ctx, userId, userUuid, userAccounts).
			Return(userAccountsResp, nil).Once()

		addressRepoMock.EXPECT().
			AddOldAddress(ctx, matchAddress).
			Return(nil).Once()

		resp, err := s.CreateOldAddress(ctx, addressUuid, userUuid, addressType, network, coin)
		require.NotNil(t, resp)
		require.NoError(t, err)

	})
}
