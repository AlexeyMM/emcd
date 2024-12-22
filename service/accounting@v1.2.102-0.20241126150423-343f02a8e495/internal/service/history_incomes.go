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

type IncomesHistory interface {
	GetHistory(ctx context.Context, params *accountingPb.GetHistoryRequest) (*accountingPb.GetHistoryResponse, error)
}

type incomesHistory struct {
	repo repository.IncomesHistory
}

func NewIncomesHistory(repo repository.IncomesHistory) IncomesHistory {
	return &incomesHistory{
		repo: repo,
	}
}

func (m *incomesHistory) GetHistory(ctx context.Context, params *accountingPb.GetHistoryRequest) (*accountingPb.GetHistoryResponse, error) {
	in := model.HistoryInput{
		Type:                model.HistoryType(""), // TODO: fix it
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
	result, err := m.repo.GetNewIncomes(ctx, &in)
	if err != nil {
		return nil, fmt.Errorf("GetNewIncomes: %w", err)
	}

	history := accountingPb.GetHistoryResponse{
		TotalCount:    int64(result.TotalCount),
		IncomesSum:    result.IncomesSum.String(),
		PayoutsSum:    "",
		HasNewIncome:  nil,
		HasNewPayouts: nil,
		Incomes:       nil,
		Payouts:       nil,
		Wallets:       nil,
	}

	history.Incomes = getPBIncomes(result.Incomes)
	hasNewIncome, err := m.repo.GetIncomesIsViewed(ctx, int(params.UserID), params.CoinCode)
	if err != nil {
		return nil, fmt.Errorf("GetIncomesIsViewed: %w", err)
	}
	history.HasNewIncome = &wrapperspb.BoolValue{Value: hasNewIncome}

	return &history, nil
}

func getPBIncomes(data []*model.Income) []*accountingPb.Income {
	var incomes = make([]*accountingPb.Income, 0, len(data))
	for _, v := range data {
		var income = accountingPb.Income{
			Diff:          int64(v.Diff),
			ChangePercent: strconv.FormatFloat(v.ChangePercent, 'f', -1, 64),
			Time:          strconv.FormatFloat(v.Time, 'f', -1, 64),
			Income:        v.Income.String(),
			Code:          int64(v.Code),
			HashRate:      NullInt64(v.HashRate),
		}
		incomes = append(incomes, &income)
	}

	return incomes
}
