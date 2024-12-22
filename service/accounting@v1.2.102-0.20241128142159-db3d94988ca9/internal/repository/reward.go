package repository

import (
	"code.emcdtech.com/emcd/service/accounting/model"
	"code.emcdtech.com/emcd/service/referral/protocol/reward"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Reward interface {
	GetReward(ctx context.Context, userID uuid.UUID, product, coin string, amount decimal.Decimal) ([]*model.UserIncome, error)
	Calculate(ctx context.Context, userID uuid.UUID, product, coin string, amount decimal.Decimal) ([]model.ReferralCalculation, error)
}

type rewardRepo struct {
	cli reward.RewardServiceClient
}

func NewReward(cli reward.RewardServiceClient) Reward {
	return &rewardRepo{
		cli: cli,
	}
}

func (r *rewardRepo) Calculate(ctx context.Context, userID uuid.UUID, product, coin string, amount decimal.Decimal) ([]model.ReferralCalculation, error) {
	resp, err := r.cli.Calculate(ctx, &reward.CalculateRequest{
		UserId:  userID.String(),
		Amount:  amount.String(),
		Coin:    coin,
		Product: product,
	})
	if err != nil {
		return nil, err
	}

	out := make([]model.ReferralCalculation, 0)
	for _, tx := range resp.Txs {

		am, err := decimal.NewFromString(tx.Amount)
		if err != nil {
			return nil, errors.New("invalid amount: " + tx.Amount)
		}

		uId, err := uuid.Parse(tx.UserId)
		if err != nil {
			return nil, errors.New("invalid user id: " + tx.UserId)
		}

		out = append(out, model.ReferralCalculation{
			UserId: uId,
			Type:   tx.Type,
			Amount: am,
		})
	}

	return out, nil
}

func (r *rewardRepo) GetReward(ctx context.Context, userID uuid.UUID, product, coin string, amount decimal.Decimal) ([]*model.UserIncome, error) {
	resp, err := r.cli.Calculate(ctx, &reward.CalculateRequest{
		UserId:  userID.String(),
		Amount:  amount.String(),
		Coin:    coin,
		Product: product,
	})
	if err != nil {
		return nil, err
	}

	incomes := make([]*model.UserIncome, 0, len(resp.Txs)-1)
	isWl := containsType(resp.Txs, model.WlFee)

	for i := range resp.Txs {
		if resp.Txs[i].Type == model.Normal {
			continue
		}
		income, err := r.toUserIncome(resp.Txs[i], isWl)
		if err != nil {
			return nil, fmt.Errorf("toUserIncome: %w", err)
		}
		if income != nil {
			incomes = append(incomes, income)
		}
	}
	return incomes, nil
}

func (r *rewardRepo) toUserIncome(tr *reward.Transaction, isWl bool) (*model.UserIncome, error) {
	id, err := uuid.Parse(tr.UserId)
	if err != nil {
		return nil, fmt.Errorf("parse user id: %s. %w", tr.UserId, err)
	}
	incomeAmount, err := decimal.NewFromString(tr.Amount)
	if err != nil {
		return nil, fmt.Errorf("parse amount: %s. %w", tr.Amount, err)
	}
	var trType model.TransactionType
	switch tr.Type {
	case model.Fee:
		trType = model.UserPaysPoolComsTrTypeID
	case model.WlFee:
		trType = model.UserPaysWlComsTrTypeID
	case model.RefFee:
		if isWl {
			trType = model.WlPaysUserReferralsTrTypeID
		} else {
			trType = model.PoolPaysUsersReferralsTrTypeID
		}
	}
	return &model.UserIncome{
		UserID: id,
		Amount: incomeAmount,
		Type:   trType,
	}, nil
}

func containsType(trs []*reward.Transaction, targetType string) bool {
	for i := range trs {
		if trs[i].Type == targetType {
			return true
		}
	}
	return false
}
