package repository

import (
	walletPb "code.emcdtech.com/emcd/blockchain/wallet/protocol/emcd_wallet_pb"
	"context"
	"strings"
)

type Wallet interface {
	GetAddress(ctx context.Context, userID int32, coin string) (*walletPb.GetAddressResponse, error)
}

type wallet struct {
	client walletPb.WalletServiceClient
}

func NewWallet(clientProto walletPb.WalletServiceClient) Wallet {
	return &wallet{
		client: clientProto,
	}
}

func (a *wallet) GetAddress(ctx context.Context, userID int32, coin string) (*walletPb.GetAddressResponse, error) {
	result, err := a.client.GetAddress(ctx, &walletPb.GetAddressRequest{UserId: userID, CoinId: strings.ToLower(coin)})
	if err != nil {
		return nil, err
	}

	return result, nil
}
