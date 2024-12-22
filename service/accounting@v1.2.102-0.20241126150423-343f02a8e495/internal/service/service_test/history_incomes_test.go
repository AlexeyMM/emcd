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

func TestIncomesHistoryGetHistory_Success(t *testing.T) {
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()
	//
	// repository := mockedSuccessIncomeRepository{}
	// service := service2.NewIncomesHistory(&repository)
	//
	// result, err := service.GetHistory(ctx, &pb.GetHistoryRequest{})
	// require.NotNil(t, result.Incomes)
	// require.NotNil(t, result.IncomesSum)
	// require.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	incomeHistoryRepo := repositoryMock.NewMockIncomesHistory(t)
	incomeHistoryService := service.NewIncomesHistory(incomeHistoryRepo)

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
		IncomesSum:    utils.DecimalToPtr(decimal.NewFromInt(1)),
		PayoutsSum:    nil,
		HasNewIncome:  nil,
		HasNewPayouts: nil,
		Incomes:       nil,
		Payouts:       nil,
		Wallets:       nil,
	}

	retBool := true

	incomeHistoryRepo.EXPECT().
		GetNewIncomes(ctx, reqIn).
		Return(retHistoryOut, nil)
	incomeHistoryRepo.EXPECT().
		GetIncomesIsViewed(ctx, int(req.UserID), req.CoinCode).
		Return(retBool, nil)

	result, err := incomeHistoryService.GetHistory(ctx, req)
	require.NotNil(t, result)
	require.NoError(t, err)

}

func TestIncomesHistoryGetHistory_ErrorGetNewIncomes(t *testing.T) {
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()
	//
	// repository := mockedFailIncomeRepository{}
	// service := service2.NewIncomesHistory(&repository)
	//
	// result, err := service.GetHistory(ctx, &pb.GetHistoryRequest{})
	// require.Error(t, err)
	// require.Nil(t, result)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	incomeHistoryRepo := repositoryMock.NewMockIncomesHistory(t)
	incomeHistoryService := service.NewIncomesHistory(incomeHistoryRepo)

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

	incomeHistoryRepo.EXPECT().
		GetNewIncomes(ctx, reqIn).
		Return(nil, retError)
	// incomeHistoryRepo.EXPECT().
	//	GetIncomesIsViewed(ctx, int(req.UserID), req.CoinCode).
	//	Return(retBool, nil)

	result, err := incomeHistoryService.GetHistory(ctx, req)
	require.Nil(t, result)
	require.Error(t, err)

	require.EqualError(t, err, fmt.Errorf("GetNewIncomes: error").Error())
}

func TestIncomesHistoryGetHistory_ErrorGetIncomesIsViewed(t *testing.T) {
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()
	//
	// repository := mockedFailIncomeRepository{}
	// service := service2.NewIncomesHistory(&repository)
	//
	// result, err := service.GetHistory(ctx, &pb.GetHistoryRequest{})
	// require.Error(t, err)
	// require.Nil(t, result)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	incomeHistoryRepo := repositoryMock.NewMockIncomesHistory(t)
	incomeHistoryService := service.NewIncomesHistory(incomeHistoryRepo)

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
		IncomesSum:    utils.DecimalToPtr(decimal.NewFromInt(1)),
		PayoutsSum:    nil,
		HasNewIncome:  nil,
		HasNewPayouts: nil,
		Incomes:       nil,
		Payouts:       nil,
		Wallets:       nil,
	}

	retBool := true

	retError := errors.New("error")

	incomeHistoryRepo.EXPECT().
		GetNewIncomes(ctx, reqIn).
		Return(retHistoryOut, nil)
	incomeHistoryRepo.EXPECT().
		GetIncomesIsViewed(ctx, int(req.UserID), req.CoinCode).
		Return(retBool, retError)

	result, err := incomeHistoryService.GetHistory(ctx, req)
	require.Nil(t, result)
	require.Error(t, err)

	require.EqualError(t, err, fmt.Errorf("GetIncomesIsViewed: error").Error())
}
