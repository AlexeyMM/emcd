package server

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	businessErr "code.emcdtech.com/emcd/sdk/error"
	"code.emcdtech.com/emcd/sdk/log"

	"code.emcdtech.com/emcd/service/coin/internal/repository"
	"code.emcdtech.com/emcd/service/coin/internal/service"
	"code.emcdtech.com/emcd/service/coin/model"
	pb "code.emcdtech.com/emcd/service/coin/protocol/coin"
)

const (
	internal = "C000001"
)

type CoinService pb.CoinServiceServer

type coin struct {
	pb.UnsafeCoinServiceServer
	svc *service.Coin
}

func NewCoin(svc *service.Coin) CoinService {
	return &coin{
		svc: svc,
	}
}

func (c *coin) GetCoinIDFromLegacyID(
	ctx context.Context,
	req *pb.GetCoinIDFromLegacyIDRequest,
) (*pb.GetCoinIDFromLegacyIDResponse, error) {
	id, err := c.svc.GetCoinIDFromLegacyID(ctx, req.LegacyCoinId)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return &pb.GetCoinIDFromLegacyIDResponse{
		CoinId: id,
	}, nil
}

func (c *coin) GetCoin(
	ctx context.Context,
	req *pb.GetCoinRequest,
) (*pb.GetCoinResponse, error) {
	coin, err := c.svc.GetCoin(ctx, req.CoinId)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return &pb.GetCoinResponse{
		Coin: toProtoCoin(coin),
	}, nil
}

func (c *coin) GetCoins(ctx context.Context, req *pb.GetCoinsRequest) (*pb.GetCoinsResponse, error) {
	coins, totalCount, err := c.svc.GetCoins(ctx, req.WlId, req.Limit, req.Offset)
	if err != nil {
		log.Error(ctx, "GetCoint: %s", err.Error())
		return nil, businessErr.NewError(internal, err.Error())
	}

	pbCoins := make([]*pb.Coin, 0, len(coins))
	for _, coin := range coins {
		pbCoins = append(pbCoins, toProtoCoin(coin))
	}
	return &pb.GetCoinsResponse{Coins: pbCoins, TotalCount: totalCount}, nil
}

func (c *coin) GetWithdrawalFee(
	ctx context.Context,
	req *pb.RequestGetWithdrawalFee,
) (*pb.ResponseGetWithdrawalFee, error) {
	fee, validTill, err := c.svc.GetWithdrawalFee(ctx, req)
	if err != nil {
		log.Error(ctx, "GetWithdrawalFee: %s", err.Error())
		return nil, err
	}

	return &pb.ResponseGetWithdrawalFee{Fee: fee, ValidTill: validTill.String()}, nil
}

func toProtoCoin(coin *model.Coin) *pb.Coin {
	pbCoin := &pb.Coin{
		Id:                    coin.ID,
		Title:                 coin.Title,
		Description:           coin.Description,
		IsActive:              coin.IsActive,
		SortPriorityMining:    coin.SortPriorityMining,
		SortPriorityWallet:    coin.SortPriorityWallet,
		MediaUrl:              coin.MediaURL,
		IsWithdrawalsDisabled: coin.IsWithdrawalsDisabled,
		LegacyCoinId:          coin.LegacyCoinID,
		MiningRewardType:      coin.MiningRewardType,
		Networks:              toProtoNetworks(coin.Networks),
		SortPrioritySwap:      coin.SortPrioritySwap,
	}

	return pbCoin
}

func toProtoNetworks(networks []*model.CoinNetwork) []*pb.CoinNetwork {
	protoNetworks := make([]*pb.CoinNetwork, len(networks))
	for i, coinNetwork := range networks {
		protoNetworks[i] = &pb.CoinNetwork{
			CoinId:                  coinNetwork.CoinID,
			NetworkId:               coinNetwork.NetworkID,
			Title:                   coinNetwork.Title,
			IsActive:                coinNetwork.IsActive,
			Description:             coinNetwork.Description,
			ContractAddress:         coinNetwork.ContractAddress,
			Decimals:                coinNetwork.Decimals,
			MinpayMining:            coinNetwork.MinPayMining,
			WithdrawalFee:           coinNetwork.WithdrawalFee,
			WithdrawalMinLimit:      coinNetwork.WithdrawalMinLimit,
			WithdrawalFeeTtlSeconds: coinNetwork.WithdrawalFeeTtlSeconds,
			IsMining:                coinNetwork.IsMining,
			IsWallet:                coinNetwork.IsWallet,
			IsFreeWithdraw:          coinNetwork.IsFreeWithdraw,
			IsWithdrawalsDisabled:   coinNetwork.IsWithdrawalsDisabled,
			HashDivisorPowerOfTen:   coinNetwork.HashDivisorPowerOfTen,
			ExplorerUrl:             coinNetwork.ExplorerUrl,
			Priority:                coinNetwork.Priority,
		}
	}
	return protoNetworks
}
