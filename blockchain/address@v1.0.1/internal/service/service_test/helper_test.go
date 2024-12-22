package service_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	userAccountModel "code.emcdtech.com/emcd/service/accounting/model"
	"github.com/stretchr/testify/require"

	"code.emcdtech.com/emcd/blockchain/address/internal/service"
	serviceMock "code.emcdtech.com/emcd/blockchain/address/mocks/internal_/service"
)

const reloadAddressesCommand = "reload_addresses"

func sendConsumerNetworkReloadAddr(t *testing.T, ctx context.Context, rabbitService *serviceMock.MockRabbitService, network nodeCommon.NetworkEnum) {
	message := &service.AdminMessage{
		Command:      reloadAddressesCommand,
		TxID:         "",
		Hash:         "",
		ToBlockScore: 0,
	}

	if msgByte, err := json.Marshal(message); err != nil {
		require.NoError(t, err)

	} else {
		routingKey := fmt.Sprintf("%s.admin", network.ToString()) // TODO: id?

		rabbitService.EXPECT().
			Publish(ctx, routingKey, msgByte).Return(nil).Once()

	}
}

func sendConsumerNetworkGroupReloadAddr(t *testing.T, ctx context.Context, rabbitService *serviceMock.MockRabbitService, networkGroup nodeCommon.NetworkGroupEnum) {
	for _, network := range networkGroup.GetNetworks() {
		sendConsumerNetworkReloadAddr(t, ctx, rabbitService, network)

	}
}

func userModelCopy(u *userAccountModel.UserAccount) *userAccountModel.UserAccount {

	return &userAccountModel.UserAccount{
		ID:            u.ID,
		UserID:        u.UserID,
		CoinID:        u.CoinID,
		AccountTypeID: u.AccountTypeID,
		Minpay:        u.Minpay,
		Address:       u.Address,
		ChangedAt:     u.ChangedAt,
		Img1:          u.Img1,
		Img2:          u.Img2,
		IsActive:      u.IsActive,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
		Fee:           u.Fee,
		UserIDNew:     u.UserIDNew,
		CoinNew:       u.CoinNew,
	}
}
