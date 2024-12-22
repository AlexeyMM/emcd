package repository_test

import (
	protoAccountingMock "code.emcdtech.com/emcd/service/accounting/mocks/protocol/accounting"
	"code.emcdtech.com/emcd/service/accounting/model"
	accountingPb "code.emcdtech.com/emcd/service/accounting/protocol/accounting"
	externalRepository "code.emcdtech.com/emcd/service/accounting/repository"
	"context"
	"errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestClientRepository_GetHistorySuccess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	historyRepository := externalRepository.NewHistoryRepository()
	protoAccountingClientMock := protoAccountingMock.NewMockAccountingServiceClient(t)

	req := &model.HistoryInput{
		Type:                "",
		CoinCode:            "",
		From:                "",
		To:                  "",
		Limit:               0,
		Offset:              0,
		CoinholdID:          0,
		UserID:              0,
		TransactionTypesIDs: nil,
		AccountTypeIDs:      nil,
		CoinsIDs:            nil,
	}

	protoAccountingClientMock.On("GetHistory", mock.Anything, &accountingPb.GetHistoryRequest{
		Type:                string(req.Type),
		CoinCode:            req.CoinCode,
		From:                req.From,
		To:                  req.To,
		Limit:               req.Limit,
		Offset:              req.Offset,
		CoinholdID:          req.CoinholdID,
		UserID:              req.UserID,
		TransactionTypesIDs: req.TransactionTypesIDs,
		AccountTypeIDs:      req.AccountTypeIDs,
		CoinsIDs:            req.CoinsIDs,
	}).Return(&accountingPb.GetHistoryResponse{
		TotalCount:    0,
		IncomesSum:    "",
		PayoutsSum:    "",
		HasNewIncome:  nil,
		HasNewPayouts: nil,
		Incomes:       nil,
		Payouts:       nil,
		Wallets:       nil,
	}, nil)

	history, err := historyRepository.GetHistory(ctx, protoAccountingClientMock, req)
	require.NoError(t, err)
	require.NotNil(t, history)
}

func TestClientRepository_GetHistoryFail(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	historyRepository := externalRepository.NewHistoryRepository()
	protoAccountingClientMock := protoAccountingMock.NewMockAccountingServiceClient(t)

	req := &model.HistoryInput{
		Type:                "",
		CoinCode:            "",
		From:                "",
		To:                  "",
		Limit:               0,
		Offset:              0,
		CoinholdID:          0,
		UserID:              0,
		TransactionTypesIDs: nil,
		AccountTypeIDs:      nil,
		CoinsIDs:            nil,
	}

	protoAccountingClientMock.On("GetHistory", mock.Anything, &accountingPb.GetHistoryRequest{
		Type:                string(req.Type),
		CoinCode:            req.CoinCode,
		From:                req.From,
		To:                  req.To,
		Limit:               req.Limit,
		Offset:              req.Offset,
		CoinholdID:          req.CoinholdID,
		UserID:              req.UserID,
		TransactionTypesIDs: req.TransactionTypesIDs,
		AccountTypeIDs:      req.AccountTypeIDs,
		CoinsIDs:            req.CoinsIDs,
	}).Return(nil, errors.New("error"))

	history, err := historyRepository.GetHistory(ctx, protoAccountingClientMock, req)
	require.Error(t, err)
	require.Empty(t, history)
}
