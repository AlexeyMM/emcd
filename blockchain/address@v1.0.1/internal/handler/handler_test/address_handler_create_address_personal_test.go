package handler_test

import (
	"context"
	"testing"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	coinPb "code.emcdtech.com/emcd/service/coin/protocol/coin"
	coinValidatorMock "code.emcdtech.com/emcd/service/coin/repository/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"code.emcdtech.com/emcd/blockchain/address/common/utils"
	"code.emcdtech.com/emcd/blockchain/address/internal/handler"
	"code.emcdtech.com/emcd/blockchain/address/internal/handler/mapping"
	externalMock "code.emcdtech.com/emcd/blockchain/address/mocks/external"
	"code.emcdtech.com/emcd/blockchain/address/mocks/internal_/service"
	"code.emcdtech.com/emcd/blockchain/address/model"
	addressPb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
	"code.emcdtech.com/emcd/blockchain/address/repository"
)

func TestAddressHandler_AddOrUpdatePersonalAddress(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serviceAddressMock := service.NewMockAddressService(t)
	coinValidatorRepoMock := coinValidatorMock.NewMockCoinValidatorRepository(t)
	coinProtoCliMock := externalMock.NewMockCoinServiceClient(t)
	addressHandler := handler.NewAddressHandler(serviceAddressMock, coinValidatorRepoMock, coinProtoCliMock, IsNetworkOldWay)

	userUuid := uuid.New()
	addressStr := uuid.NewString()
	network := nodeCommon.EthNetworkId
	defaultMinPayout := 1.0
	minPayout := 1.0

	var minPayoutNil *float64

	reqCoin := &coinPb.GetCoinRequest{
		CoinId: network.ToString(),
	}

	coinResp := &coinPb.GetCoinResponse{
		Coin: &coinPb.Coin{
			Id:                    "",
			IsActive:              false,
			Title:                 "",
			Description:           "",
			SortPriorityMining:    0,
			SortPriorityWallet:    0,
			MediaUrl:              "",
			IsWithdrawalsDisabled: false,
			Networks: []*coinPb.CoinNetwork{{
				CoinId:                  "",
				NetworkId:               "",
				IsActive:                false,
				Title:                   "",
				Description:             "",
				ContractAddress:         "",
				Decimals:                0,
				MinpayMining:            defaultMinPayout,
				WithdrawalFee:           0,
				WithdrawalMinLimit:      0,
				WithdrawalFeeTtlSeconds: 0,
				IsMining:                true, // true
				IsWallet:                false,
				IsFreeWithdraw:          false,
				IsWithdrawalsDisabled:   false,
				HashDivisorPowerOfTen:   0,
				ExplorerUrl:             "",
				Priority:                0,
			}},
			LegacyCoinId:     0,
			MiningRewardType: "",
		},
	}

	t.Run("success personal address", func(t *testing.T) {
		t.Run("create", func(t *testing.T) {
			req := &addressPb.CreatePersonalAddressRequest{
				Address:   addressStr,
				UserUuid:  userUuid.String(),
				Network:   network.ToString(),
				MinPayout: &minPayout,
			}

			serviceAddressMock.EXPECT().
				GetPersonalAddressByConstraint(ctx, userUuid, network).
				Return(nil, nil).Once()

			coinProtoCliMock.EXPECT().GetCoin(ctx, reqCoin).
				Return(coinResp, nil).Once()

			addressPersonalCreate := &model.AddressPersonal{}
			serviceAddressMock.EXPECT().
				CreatePersonalAddress(ctx, addressStr, userUuid, network, minPayout).
				Return(addressPersonalCreate, nil).Once()

			addressPersonalCreateProto := mapping.MapModelAddressPersonalToProto(addressPersonalCreate)

			resp, err := addressHandler.AddOrUpdatePersonalAddress(ctx, req)
			require.NotNil(t, resp)
			require.NoError(t, err)

			require.Equal(t, resp, addressPersonalCreateProto)

		})

		t.Run("update with min pay", func(t *testing.T) {
			req := &addressPb.CreatePersonalAddressRequest{
				Address:   addressStr,
				UserUuid:  userUuid.String(),
				Network:   network.ToString(),
				MinPayout: &minPayout,
			}

			coinProtoCliMock.EXPECT().GetCoin(ctx, reqCoin).
				Return(coinResp, nil).Once()

			addressPersonalCreate := &model.AddressPersonal{}
			addressesPersonalCreate := model.AddressesPersonal{addressPersonalCreate}
			serviceAddressMock.EXPECT().
				GetPersonalAddressByConstraint(ctx, userUuid, network).
				Return(addressesPersonalCreate, nil).Once()

			serviceAddressMock.EXPECT().
				UpdatePersonalAddress(ctx, addressPersonalCreate, addressStr, &minPayout).
				Return(addressPersonalCreate, nil).Once()

			addressPersonalCreateProto := mapping.MapModelAddressPersonalToProto(addressPersonalCreate)

			resp, err := addressHandler.AddOrUpdatePersonalAddress(ctx, req)
			require.NotNil(t, resp)
			require.NoError(t, err)

			require.Equal(t, resp, addressPersonalCreateProto)

		})

		t.Run("update without min pay", func(t *testing.T) {
			req := &addressPb.CreatePersonalAddressRequest{
				Address:   addressStr,
				UserUuid:  userUuid.String(),
				Network:   network.ToString(),
				MinPayout: minPayoutNil,
			}

			addressPersonalCreate := &model.AddressPersonal{}
			addressesPersonalCreate := model.AddressesPersonal{addressPersonalCreate}
			serviceAddressMock.EXPECT().
				GetPersonalAddressByConstraint(ctx, userUuid, network).
				Return(addressesPersonalCreate, nil).Once()

			serviceAddressMock.EXPECT().
				UpdatePersonalAddress(ctx, addressPersonalCreate, addressStr, minPayoutNil).
				Return(addressPersonalCreate, nil).Once()

			addressPersonalCreateProto := mapping.MapModelAddressPersonalToProto(addressPersonalCreate)

			resp, err := addressHandler.AddOrUpdatePersonalAddress(ctx, req)
			require.NotNil(t, resp)
			require.NoError(t, err)

			require.Equal(t, resp, addressPersonalCreateProto)

		})
	})

	t.Run("error parse empty address", func(t *testing.T) {
		req := &addressPb.CreatePersonalAddressRequest{
			Address:   "",
			UserUuid:  userUuid.String(),
			Network:   network.ToString(),
			MinPayout: nil,
		}

		resp, err := addressHandler.AddOrUpdatePersonalAddress(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAddr1081)

	})

	t.Run("error parse user uuid", func(t *testing.T) {
		req := &addressPb.CreatePersonalAddressRequest{
			Address:   addressStr,
			UserUuid:  userUuid.String() + "fake",
			Network:   network.ToString(),
			MinPayout: nil,
		}

		resp, err := addressHandler.AddOrUpdatePersonalAddress(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAddr1081)

	})

	t.Run("error parse network", func(t *testing.T) {
		req := &addressPb.CreatePersonalAddressRequest{
			Address:   addressStr,
			UserUuid:  userUuid.String(),
			Network:   network.ToString() + "fake",
			MinPayout: nil,
		}

		resp, err := addressHandler.AddOrUpdatePersonalAddress(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAddr1081)

	})

	t.Run("error mock get personal by constraint", func(t *testing.T) {
		errMock := newMockError()
		serviceAddressMock.EXPECT().
			GetPersonalAddressByConstraint(ctx, userUuid, network).
			Return(nil, errMock).Once()

		req := &addressPb.CreatePersonalAddressRequest{
			Address:   addressStr,
			UserUuid:  userUuid.String(),
			Network:   network.ToString(),
			MinPayout: nil,
		}

		resp, err := addressHandler.AddOrUpdatePersonalAddress(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAddr1082)

	})

	t.Run("empty min payout for create personal by constraint", func(t *testing.T) {
		serviceAddressMock.EXPECT().
			GetPersonalAddressByConstraint(ctx, userUuid, network).
			Return(nil, nil).Once()

		req := &addressPb.CreatePersonalAddressRequest{
			Address:   addressStr,
			UserUuid:  userUuid.String(),
			Network:   network.ToString(),
			MinPayout: nil,
		}

		resp, err := addressHandler.AddOrUpdatePersonalAddress(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAddr1083)

	})

	t.Run("error mock get default min payout for create personal by constraint", func(t *testing.T) {
		serviceAddressMock.EXPECT().
			GetPersonalAddressByConstraint(ctx, userUuid, network).
			Return(nil, nil).Once()

		errMock := newMockError()
		coinProtoCliMock.EXPECT().GetCoin(ctx, reqCoin).
			Return(nil, errMock).Once()

		req := &addressPb.CreatePersonalAddressRequest{
			Address:   addressStr,
			UserUuid:  userUuid.String(),
			Network:   network.ToString(),
			MinPayout: utils.Float64ToPtr(minPayout),
		}

		resp, err := addressHandler.AddOrUpdatePersonalAddress(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAddr1084)

	})

	t.Run("min payout less then default for create personal by constraint", func(t *testing.T) {
		serviceAddressMock.EXPECT().
			GetPersonalAddressByConstraint(ctx, userUuid, network).
			Return(nil, nil).Once()

		coinProtoCliMock.EXPECT().GetCoin(ctx, reqCoin).
			Return(coinResp, nil).Once()

		req := &addressPb.CreatePersonalAddressRequest{
			Address:   addressStr,
			UserUuid:  userUuid.String(),
			Network:   network.ToString(),
			MinPayout: utils.Float64ToPtr(minPayout - 0.1),
		}

		resp, err := addressHandler.AddOrUpdatePersonalAddress(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAddr1085)

	})

	t.Run("mock error for create personal by constraint", func(t *testing.T) {
		serviceAddressMock.EXPECT().
			GetPersonalAddressByConstraint(ctx, userUuid, network).
			Return(nil, nil).Once()

		coinProtoCliMock.EXPECT().GetCoin(ctx, reqCoin).
			Return(coinResp, nil).Once()

		errMock := newMockError()
		serviceAddressMock.EXPECT().
			CreatePersonalAddress(ctx, addressStr, userUuid, network, minPayout).
			Return(nil, errMock).Once()

		req := &addressPb.CreatePersonalAddressRequest{
			Address:   addressStr,
			UserUuid:  userUuid.String(),
			Network:   network.ToString(),
			MinPayout: utils.Float64ToPtr(minPayout),
		}

		resp, err := addressHandler.AddOrUpdatePersonalAddress(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAddr1086)

	})

	t.Run("error mock coin update for update personal by constraint", func(t *testing.T) {
		addressPersonalCreate := &model.AddressPersonal{}
		addressesPersonalCreate := model.AddressesPersonal{addressPersonalCreate}
		serviceAddressMock.EXPECT().
			GetPersonalAddressByConstraint(ctx, userUuid, network).
			Return(addressesPersonalCreate, nil).Once()

		errMock := newMockError()
		coinProtoCliMock.EXPECT().GetCoin(ctx, reqCoin).
			Return(nil, errMock).Once()

		req := &addressPb.CreatePersonalAddressRequest{
			Address:   addressStr,
			UserUuid:  userUuid.String(),
			Network:   network.ToString(),
			MinPayout: &minPayout,
		}

		resp, err := addressHandler.AddOrUpdatePersonalAddress(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAddr1084)

	})

	t.Run("min payout less then default min payout for update personal by constraint", func(t *testing.T) {
		addressPersonalCreate := &model.AddressPersonal{}
		addressesPersonalCreate := model.AddressesPersonal{addressPersonalCreate}
		serviceAddressMock.EXPECT().
			GetPersonalAddressByConstraint(ctx, userUuid, network).
			Return(addressesPersonalCreate, nil).Once()

		coinProtoCliMock.EXPECT().GetCoin(ctx, reqCoin).
			Return(coinResp, nil).Once()

		req := &addressPb.CreatePersonalAddressRequest{
			Address:   addressStr,
			UserUuid:  userUuid.String(),
			Network:   network.ToString(),
			MinPayout: utils.Float64ToPtr(minPayout - 0.1),
		}

		resp, err := addressHandler.AddOrUpdatePersonalAddress(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAddr1085)

	})

	t.Run("error mock update for update personal by constraint", func(t *testing.T) {
		addressPersonalCreate := &model.AddressPersonal{}
		addressesPersonalCreate := model.AddressesPersonal{addressPersonalCreate}
		serviceAddressMock.EXPECT().
			GetPersonalAddressByConstraint(ctx, userUuid, network).
			Return(addressesPersonalCreate, nil).Once()

		errMock := newMockError()
		serviceAddressMock.EXPECT().UpdatePersonalAddress(ctx, addressPersonalCreate, addressStr, minPayoutNil).
			Return(nil, errMock).Once()

		req := &addressPb.CreatePersonalAddressRequest{
			Address:   addressStr,
			UserUuid:  userUuid.String(),
			Network:   network.ToString(),
			MinPayout: nil,
		}

		resp, err := addressHandler.AddOrUpdatePersonalAddress(ctx, req)
		require.Nil(t, resp)
		require.Error(t, err)

		require.ErrorIs(t, err, repository.ErrAddr1087)

	})
}
