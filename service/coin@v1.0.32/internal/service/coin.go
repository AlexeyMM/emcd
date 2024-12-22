package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"code.emcdtech.com/emcd/blockchain/node/proto/node"
	wlPb "code.emcdtech.com/emcd/service/whitelabel/protocol/whitelabel"
	ratePb "code.emcdtech.com/emcd/services/rate/proto"

	"code.emcdtech.com/emcd/service/coin/internal/config"
	"code.emcdtech.com/emcd/service/coin/internal/repository"
	"code.emcdtech.com/emcd/service/coin/model"
	pb "code.emcdtech.com/emcd/service/coin/protocol/coin"
)

const ethNetworkID = "eth"

type Coin struct {
	repo               repository.Coin
	nodeClient         node.NodeServiceClient
	wlCli              repository.WhiteLabel
	rateCli            ratePb.RateClient
	nodeAuthKey        string
	feeMultiplierErc20 float64
}

func NewCoin(
	repo repository.Coin,
	nodeClient node.NodeServiceClient,
	wlClient repository.WhiteLabel,
	rateCli ratePb.RateClient,
	cfg config.Config,
) *Coin {
	return &Coin{
		repo:               repo,
		nodeClient:         nodeClient,
		wlCli:              wlClient,
		rateCli:            rateCli,
		nodeAuthKey:        cfg.NodeAuthKey,
		feeMultiplierErc20: cfg.FeeMultiplierErc20,
	}
}

func (c *Coin) GetCoinIDFromLegacyID(ctx context.Context, legacyCoinID int32) (string, error) {
	coin, err := c.repo.GetCoinFromLegacyID(ctx, legacyCoinID)
	if err != nil {
		return "", fmt.Errorf("get coin from legacy id: %w", err)
	}
	return coin.ID, nil
}

func (c *Coin) GetCoin(ctx context.Context, coinID string) (*model.Coin, error) {
	coin, err := c.repo.GetCoin(ctx, coinID)
	if err != nil {
		return nil, fmt.Errorf("get coin: %w", err)
	}
	return coin, nil
}

func (c *Coin) GetCoins(ctx context.Context, wlID *string, limit, offset int32) ([]*model.Coin, int32, error) {
	networks, err := c.repo.GetCoinsNetworks(ctx)
	if err != nil {
		return nil, 0, err
	}

	coins, totalCount, err := c.repo.GetCoins(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	for i, coinItem := range coins {
		for _, network := range networks {
			if network.CoinID == coinItem.ID {
				coins[i].Networks = append(coins[i].Networks, network)
			}
		}
	}

	if wlID != nil {
		_, err = uuid.Parse(*wlID)
		if err != nil {
			return nil, 0, fmt.Errorf("wrong whitelabel ID %s", *wlID)
		}

		coinsResp, err := c.wlCli.GetCoins(ctx, &wlPb.GetCoinsRequest{
			WlId: *wlID,
		})
		if err != nil {
			return nil, 0, err
		}

		if len(coinsResp.Coins) == 0 {
			return nil, 0, fmt.Errorf("no coins for whitelabel %s", *wlID)
		}

		wlCoins := make([]*model.Coin, 0)
		for _, c := range coins {
			for _, wlC := range coinsResp.Coins {
				if strings.EqualFold(wlC.CoinId, c.ID) {
					wlCoins = append(wlCoins, c)
				}
			}
		}
		return wlCoins, int32(len(wlCoins)), nil
	}

	return coins, totalCount, nil
}

func (c *Coin) GetWithdrawalFee(ctx context.Context, req *pb.RequestGetWithdrawalFee) (float64, time.Time, error) {
	if req.GetNetworkId() != ethNetworkID {
		return 0, time.Time{}, errors.WithStack(fmt.Errorf("GetWithdrawalFee unimplemented for: %+v", req))
	}

	feeResp, err := c.nodeClient.GetEvmTxPaymentFee(ctx, &node.RequestGetEvmTxFee{
		AuthKey:   c.nodeAuthKey,
		CoinID:    req.GetCoinId(),
		NetworkID: req.GetNetworkId(),
		ToAddress: req.GetToAddress(),
		Amount:    req.GetAmount(),
	})
	if err != nil {
		return 0, time.Time{}, errors.WithStack(err)
	}

	coinNetwork, err := c.repo.GetCoinNetwork(ctx, req.GetCoinId(), req.GetNetworkId())
	if err != nil {
		return 0, time.Time{}, errors.WithStack(err)
	}

	fee := feeResp.GetFee()
	if req.GetCoinId() != ethNetworkID {
		rateResp, err := c.rateCli.GetCrossRate(ctx, &ratePb.RequestGetCrossRate{CoinFrom: req.GetCoinId(), CoinTo: ethNetworkID})
		if err != nil {
			return 0, time.Time{}, errors.WithStack(err)
		}

		fee /= rateResp.Rate
	}

	fee *= c.feeMultiplierErc20
	validTill := time.Now().Add(time.Second * time.Duration(coinNetwork.WithdrawalFeeTtlSeconds))

	return fee, validTill, nil
}
