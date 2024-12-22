package service

import (
	"context"
	"fmt"
	"strconv"

	referralPb "code.emcdtech.com/emcd/service/accounting/protocol/referral"
	"github.com/google/uuid"

	"code.emcdtech.com/emcd/service/accounting/internal/repository"
	"code.emcdtech.com/emcd/service/accounting/model"
)

type Referral interface {
	GetReferralsStatistic(ctx context.Context, req *referralPb.GetReferralsStatisticRequest) (*referralPb.GetReferralsStatisticResponse, error)
}

func NewReferral(statistic repository.ReferralStatistic, accounts repository.UserAccountRepo) Referral {
	return &referral{
		statistic:     statistic,
		usersAccounts: accounts,
	}
}

type referral struct {
	statistic     repository.ReferralStatistic
	usersAccounts repository.UserAccountRepo
}

func (s *referral) GetReferralsStatistic(ctx context.Context, req *referralPb.GetReferralsStatisticRequest) (*referralPb.GetReferralsStatisticResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, fmt.Errorf("get referrals statistic service: parsing user_id: %w", err)
	}
	accounts, err := s.usersAccounts.FindUserAccountByUserIdLegacy(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get referrals statistic service: getting account ids for user: %s: %w", userID, err)
	}
	accountIDs := make([]int64, 0, len(accounts))
	for _, account := range accounts {
		accountIDs = append(accountIDs, int64(account.ID))
	}

	in := model.ReferralsStatisticInput{
		UserID:              userID,
		AccountIDs:          accountIDs,
		TransactionTypesIDs: req.TransactionTypesIds,
	}
	stats, err := s.statistic.GetReferralsStatistic(ctx, &in)
	if err != nil {
		return nil, fmt.Errorf("get referrals service getReferralsStatistic for user: %s: %w", userID, err)
	}

	return &referralPb.GetReferralsStatisticResponse{
		YesterdayIncome: s.convertResponse(stats.Income.Yesterday),
		ThisMonthIncome: s.convertResponse(stats.Income.ThisMonth),
	}, nil
}

func (s *referral) convertResponse(in []*model.ReferralIncome) []*referralPb.ReferralIncome {

	res := make([]*referralPb.ReferralIncome, len(in))
	for i, inc := range in {
		amount := strconv.FormatFloat(inc.Amount, 'f', -1, 64)
		res[i] = &referralPb.ReferralIncome{CoinId: inc.CoinID, Amount: amount}
	}
	return res
}
