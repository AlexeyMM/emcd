package service_test

import (
	repositoryMock "code.emcdtech.com/emcd/service/accounting/internal/repository/mock"
	"code.emcdtech.com/emcd/service/accounting/internal/service"
	"code.emcdtech.com/emcd/service/accounting/internal/utils"
	"code.emcdtech.com/emcd/service/accounting/model"
	accountingPb "code.emcdtech.com/emcd/service/accounting/protocol/accounting"
	"context"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPayoutsHistoryGetHistory_Success(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	payoutHistoryRepo := repositoryMock.NewMockPayoutsHistory(t)
	payoutHistoryService := service.NewPayoutsHistory(payoutHistoryRepo)

	req := &accountingPb.GetHistoryRequest{
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

	reqIn := &model.HistoryInput{
		CoinCode:            req.CoinCode,
		From:                req.From,
		To:                  req.To,
		Limit:               req.Limit,
		Offset:              req.Offset,
		UserID:              req.UserID,
		TransactionTypesIDs: req.TransactionTypesIDs,
	}

	retHistoryOut := &model.HistoryOutput{
		TotalCount:    0,
		IncomesSum:    nil,
		PayoutsSum:    utils.DecimalToPtr(decimal.NewFromInt(1)),
		HasNewIncome:  nil,
		HasNewPayouts: nil,
		Incomes:       nil,
		Payouts:       nil,
		Wallets:       nil,
	}

	retBool := true

	payoutHistoryRepo.EXPECT().
		GetNewPayouts(ctx, reqIn).
		Return(retHistoryOut, nil)
	payoutHistoryRepo.EXPECT().
		GetPayoutsIsViewed(ctx, int(req.UserID), req.CoinCode).
		Return(retBool, nil)

	result, err := payoutHistoryService.GetHistory(ctx, req)
	require.NotNil(t, result)
	require.NoError(t, err)

}

func TestPayoutsHistoryGetHistory_ErrorGetNewPayouts(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	payoutHistoryRepo := repositoryMock.NewMockPayoutsHistory(t)
	payoutHistoryService := service.NewPayoutsHistory(payoutHistoryRepo)

	req := &accountingPb.GetHistoryRequest{
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

	reqIn := &model.HistoryInput{
		CoinCode:            req.CoinCode,
		From:                req.From,
		To:                  req.To,
		Limit:               req.Limit,
		Offset:              req.Offset,
		UserID:              req.UserID,
		TransactionTypesIDs: req.TransactionTypesIDs,
	}

	retError := errors.New("error")

	payoutHistoryRepo.EXPECT().
		GetNewPayouts(ctx, reqIn).
		Return(nil, retError)

	result, err := payoutHistoryService.GetHistory(ctx, req)
	require.Nil(t, result)
	require.Error(t, err)

	require.EqualError(t, err, fmt.Errorf("GetNewPayouts: error").Error())

}

func TestPayoutsHistoryGetHistory_ErrorGetPayoutsIsViewed(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	payoutHistoryRepo := repositoryMock.NewMockPayoutsHistory(t)
	payoutHistoryService := service.NewPayoutsHistory(payoutHistoryRepo)

	req := &accountingPb.GetHistoryRequest{
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

	reqIn := &model.HistoryInput{
		CoinCode:            req.CoinCode,
		From:                req.From,
		To:                  req.To,
		Limit:               req.Limit,
		Offset:              req.Offset,
		UserID:              req.UserID,
		TransactionTypesIDs: req.TransactionTypesIDs,
	}

	retHistoryOut := &model.HistoryOutput{
		TotalCount:    0,
		IncomesSum:    nil,
		PayoutsSum:    utils.DecimalToPtr(decimal.NewFromInt(1)),
		HasNewIncome:  nil,
		HasNewPayouts: nil,
		Incomes:       nil,
		Payouts:       nil,
		Wallets:       nil,
	}

	retBool := true

	retError := errors.New("error")

	payoutHistoryRepo.EXPECT().
		GetNewPayouts(ctx, reqIn).
		Return(retHistoryOut, nil)
	payoutHistoryRepo.EXPECT().
		GetPayoutsIsViewed(ctx, int(req.UserID), req.CoinCode).
		Return(retBool, retError)

	result, err := payoutHistoryService.GetHistory(ctx, req)
	require.Nil(t, result)
	require.Error(t, err)

	require.EqualError(t, err, fmt.Errorf("error: Get Payouts is viewed error").Error())

}
