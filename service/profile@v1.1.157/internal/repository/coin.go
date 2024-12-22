package repository

import (
	"context"
	"fmt"

	pb "code.emcdtech.com/emcd/service/coin/protocol/coin"
	"code.emcdtech.com/emcd/service/profile/internal/model"
)

type Coin interface {
	GetMiningCoins(ctx context.Context) ([]string, error)
	GetByCode(ctx context.Context, code string) (*model.Coin, error)
	GetByLegacyID(ctx context.Context, coinID int) (*model.Coin, error)
}

type СoinRepositoryGrpc struct {
	coinCli pb.CoinServiceClient
}

func NewCoin(coinCli pb.CoinServiceClient) *СoinRepositoryGrpc {
	return &СoinRepositoryGrpc{
		coinCli: coinCli,
	}
}

func (e *СoinRepositoryGrpc) GetMiningCoins(ctx context.Context) ([]string, error) {
	req := &pb.GetCoinsRequest{
		Limit:  9999999,
		Offset: 0,
		WlId:   nil,
	}
	if resp, err := e.coinCli.GetCoins(ctx, req); err != nil {
		return nil, fmt.Errorf("rpc coin.GetCoins: %w", err)
	} else {
		var coinCodesMining []string
		for _, coin := range resp.Coins {
			for _, network := range coin.Networks {
				if network.IsMining {
					coinCodesMining = append(coinCodesMining, coin.Id)

					break
				}
			}
		}

		return coinCodesMining, nil
	}
}

func (e *СoinRepositoryGrpc) GetByCode(ctx context.Context, code string) (*model.Coin, error) {
	resp, err := e.coinCli.GetCoin(ctx, &pb.GetCoinRequest{CoinId: code})
	if err != nil {
		return nil, fmt.Errorf("rpc coin.GetByCode(CoidId: %s): %w", code, err)
	}
	return e.pbCoinToCoin(resp.GetCoin()), nil

}

func (e *СoinRepositoryGrpc) getCoinIDFromLegacyCoinID(ctx context.Context, coinID int) (string, error) {
	resp, err := e.coinCli.GetCoinIDFromLegacyID(ctx, &pb.GetCoinIDFromLegacyIDRequest{LegacyCoinId: int32(coinID)})
	if err != nil {
		return "", fmt.Errorf("rpc coin.GetCoinIDFromLegacyID (CoinId: %d): %w", coinID, err)
	}
	return resp.GetCoinId(), nil
}

func (e *СoinRepositoryGrpc) GetCoin(ctx context.Context, coinID string) (*model.Coin, error) {
	resp, err := e.coinCli.GetCoin(ctx, &pb.GetCoinRequest{CoinId: coinID})
	if err != nil {
		return nil, fmt.Errorf("rpc coin.GetCoin (CoinId: %s): %w", coinID, err)
	}
	if resp == nil {

	}
	// TODO похоже в сервисах не возвращают ошибку "not found"
	if resp.GetCoin() == nil {
		return nil, fmt.Errorf("rpc coin.GetCoin (CoinId: %s): not found", coinID)
	}
	return e.pbCoinToCoin(resp.GetCoin()), nil
}

func (e *СoinRepositoryGrpc) GetByLegacyID(ctx context.Context, coinID int) (*model.Coin, error) {
	legacyCoinId, err := e.getCoinIDFromLegacyCoinID(ctx, coinID)
	if err != nil {
		return nil, fmt.Errorf("get coin id (legacy coin id %d): %w", coinID, err)
	}
	coin, err := e.GetCoin(ctx, legacyCoinId)
	if err != nil {
		return nil, fmt.Errorf("get coin (coin id %s): %w", legacyCoinId, err)
	}
	return coin, nil
}

func (e *СoinRepositoryGrpc) pbCoinToCoin(coin *pb.Coin) *model.Coin {
	return &model.Coin{
		ID:   int(coin.GetLegacyCoinId()),
		Name: coin.GetTitle(),
		Code: coin.GetId(),
	}
}
