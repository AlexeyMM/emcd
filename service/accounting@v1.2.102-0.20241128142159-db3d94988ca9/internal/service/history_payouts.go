package service

import (
	"context"
	"fmt"
	"strconv"

	"google.golang.org/protobuf/types/known/wrapperspb"

	accountingPb "code.emcdtech.com/emcd/service/accounting/protocol/accounting"

	"code.emcdtech.com/emcd/service/accounting/internal/repository"
	"code.emcdtech.com/emcd/service/accounting/model"
)

type PayoutsHistory interface {
	GetHistory(ctx context.Context, params *accountingPb.GetHistoryRequest) (*accountingPb.GetHistoryResponse, error)
}

type payoutsHistory struct {
	repo repository.PayoutsHistory
}

func NewPayoutsHistory(repo repository.PayoutsHistory) PayoutsHistory {
	return &payoutsHistory{
		repo: repo,
	}
}

func (m *payoutsHistory) GetHistory(ctx context.Context, params *accountingPb.GetHistoryRequest) (*accountingPb.GetHistoryResponse, error) {
	in := model.HistoryInput{
		Type:                "", // TODO: why empty?
		CoinCode:            params.CoinCode,
		From:                params.From,
		To:                  params.To,
		Limit:               params.Limit,
		Offset:              params.Offset,
		CoinholdID:          0, // TODO: why empty?
		UserID:              params.UserID,
		TransactionTypesIDs: params.TransactionTypesIDs,
		AccountTypeIDs:      nil,
		CoinsIDs:            nil,
	}
	result, err := m.repo.GetNewPayouts(ctx, &in)
	if err != nil {
		return nil, fmt.Errorf("GetNewPayouts: %w", err)
	}

	var history = accountingPb.GetHistoryResponse{
		TotalCount:    int64(result.TotalCount),
		IncomesSum:    "",
		PayoutsSum:    result.PayoutsSum.String(),
		HasNewIncome:  nil,
		HasNewPayouts: nil,
		Incomes:       nil,
		Payouts:       nil,
		Wallets:       nil,
	}
	history.Payouts = getPBPayouts(result.Payouts)
	hasNewPayouts, err := m.repo.GetPayoutsIsViewed(ctx, int(params.UserID), params.CoinCode)
	if err != nil {
		return nil, fmt.Errorf("%w: Get Payouts is viewed error", err)
	}
	history.HasNewPayouts = &wrapperspb.BoolValue{Value: hasNewPayouts}

	return &history, nil
}

func getPBPayouts(data []*model.Payout) []*accountingPb.Payout {
	var payouts = make([]*accountingPb.Payout, 0, len(data))
	for _, v := range data {
		var payout = accountingPb.Payout{
			Time:   strconv.FormatFloat(v.Time, 'f', -1, 64),
			Amount: strconv.FormatFloat(v.Amount, 'f', -1, 64),
			Tx:     NullString(v.Tx),
			TxID:   NullString(v.TxId),
		}
		payouts = append(payouts, &payout)
	}

	return payouts
}
