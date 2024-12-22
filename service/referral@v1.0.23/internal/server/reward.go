// Package server implements a simple web server.
package server

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"code.emcdtech.com/emcd/sdk/log"

	"code.emcdtech.com/emcd/service/referral/internal/service"
	pb "code.emcdtech.com/emcd/service/referral/protocol/reward"
)

// RewardServer is a struct that represents a server for the RewardService.
// It implements the pb.UnimplementedRewardServiceServer interface and uses the service.Reward interface
type RewardServer struct {
	pb.UnimplementedRewardServiceServer
	service  service.Reward
	referral ReferralService
}

// NewReward creates a new instance of RewardServer with the provided service.Reward implementation.
// It returns a pointer to the newly created RewardServer.
func NewReward(s service.Reward, referral ReferralService) *RewardServer {
	return &RewardServer{
		service:  s,
		referral: referral,
	}
}

// Calculate calculates rewards based on the provided request parameters.
// It uses the Calculate method of the service implementation to perform the calculation.
// If an common occurs during the calculation, it logs the common message and returns the common.
// It creates a new CalculateResponse and populates it with the calculated transactions.
// It then returns the CalculateResponse and nil common.
func (r *RewardServer) Calculate(ctx context.Context, in *pb.CalculateRequest) (*pb.CalculateResponse, error) {
	userID, err := uuid.Parse(in.UserId)
	if err != nil {
		log.Error(ctx, "RewardServer.Calculate.ParseUserID: %v", err)
		return nil, err
	}

	amount, err := decimal.NewFromString(in.Amount)
	if err != nil {
		log.Error(ctx, "RewardServer.Calculate.ParseAmount: %v", err)
		return nil, err
	}

	txs, err := r.service.Calculate(ctx, userID, in.Product, in.Coin, amount)
	if err != nil {
		log.Error(ctx, "RewardServer.Calculate: %v", err)
		return nil, err
	}

	out := &pb.CalculateResponse{}

	for i := range txs {
		o := &pb.Transaction{
			UserId: txs[i].UserID.String(),
			Type:   txs[i].Type,
			Amount: txs[i].Amount.String(),
		}

		out.Txs = append(out.Txs, o)
	}

	return out, nil
}

// TODO: remove from reward when all services switched to referral.UpdateWithMultiplier
func (r *RewardServer) UpdateWithMultiplier(
	ctx context.Context,
	req *pb.UpdateWithMultiplierRequest,
) (*pb.UpdateWithMultiplierResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		log.Error(ctx, "reward.UpdateWithMultiplier: parse user id: %s. %v", req.UserId, err)
		return nil, err
	}
	err = r.referral.UpdateFeeWithMultiplier(ctx, userID, req.Product, req.Coins, decimal.NewFromFloat32(req.Multiplier))
	if err != nil {
		log.Error(ctx, "reward.UpdateWithMultiplier: %v", err)
		return nil, err
	}
	return &pb.UpdateWithMultiplierResponse{}, nil
}
