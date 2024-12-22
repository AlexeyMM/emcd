package repository

import (
	businessErr "code.emcdtech.com/emcd/sdk/error"
	"code.emcdtech.com/emcd/service/accounting/model"
	referralPb "code.emcdtech.com/emcd/service/accounting/protocol/referral"
	"context"
	"fmt"
	"strconv"
)

type ReferralStatisticRepository interface {
	GetReferralsStatistic(ctx context.Context, client referralPb.AccountingReferralServiceClient, data model.ReferralsStatisticInput) (model.ReferralsStatisticOutput, error)
}

type referralStatisticRepository struct{}

func NewReferralStatisticRepository() ReferralStatisticRepository {
	return &referralStatisticRepository{}
}

func (r *referralStatisticRepository) GetReferralsStatistic(ctx context.Context, client referralPb.AccountingReferralServiceClient, data model.ReferralsStatisticInput) (model.ReferralsStatisticOutput, error) {
	if len(data.TransactionTypesIDs) == 0 {
		return model.ReferralsStatisticOutput{}, fmt.Errorf("referral accounting: empty transaction types filter")
	}
	var req = referralPb.GetReferralsStatisticRequest{
		UserId:              data.UserID.String(),
		TransactionTypesIds: data.TransactionTypesIDs,
	}
	resp, err := client.GetReferralsStatistic(ctx, &req)
	if err != nil {
		return model.ReferralsStatisticOutput{}, businessErr.NewError("", fmt.Sprintf("referral_accounting: %s", err.Error()))
	}

	var income model.AggregatedReferralIncome
	income.Yesterday, err = r.convertResponse(resp.YesterdayIncome)
	if err != nil {
		return model.ReferralsStatisticOutput{}, fmt.Errorf("referral accounting: unable to parse amount string (yesterday): %w", err)
	}
	income.ThisMonth, err = r.convertResponse(resp.ThisMonthIncome)
	if err != nil {
		return model.ReferralsStatisticOutput{}, fmt.Errorf("referral accounting: unable to parse amount string (this month): %w", err)
	}
	return model.ReferralsStatisticOutput{
		Income: income,
	}, nil
}

func (*referralStatisticRepository) convertResponse(in []*referralPb.ReferralIncome) ([]*model.ReferralIncome, error) {
	res := make([]*model.ReferralIncome, len(in))
	for i, inc := range in {
		amount, err := strconv.ParseFloat(inc.Amount, 64)
		if err != nil {
			return nil, err
		}
		res[i] = &model.ReferralIncome{CoinID: inc.CoinId, Amount: amount}
	}
	return res, nil
}
