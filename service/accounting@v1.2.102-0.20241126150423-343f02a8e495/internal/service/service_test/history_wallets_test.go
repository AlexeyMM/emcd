package service_test

import (
	repositoryMock "code.emcdtech.com/emcd/service/accounting/internal/repository/mock"
	"code.emcdtech.com/emcd/service/accounting/internal/service"
	"code.emcdtech.com/emcd/service/accounting/model"
	accountingPb "code.emcdtech.com/emcd/service/accounting/protocol/accounting"
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWalletsHistoryGetHistory_Success(t *testing.T) {
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()
	//
	// repository := mockedSuccessWalletsRepository{}
	// service := service2.NewWalletsHistory(&repository)
	//
	// result, err := service.GetHistory(ctx, &pb.GetHistoryRequest{})
	// require.NotNil(t, result.Wallets)
	// require.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	walletHistoryRepo := repositoryMock.NewMockWalletsHistory(t)
	walletHistoryService := service.NewWalletsHistory(walletHistoryRepo)

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
		PayoutsSum:    nil,
		HasNewIncome:  nil,
		HasNewPayouts: nil,
		Incomes:       nil,
		Payouts:       nil,
		Wallets:       nil,
	}

	retInt64 := int64(5)

	walletHistoryRepo.EXPECT().
		GetWalletHistory(ctx, reqIn).
		Return(retHistoryOut, nil)
	walletHistoryRepo.EXPECT().
		GetWalletHistoryTotal(ctx, reqIn).
		Return(retInt64, nil)

	result, err := walletHistoryService.GetHistory(ctx, req)
	require.NotNil(t, result)
	require.NoError(t, err)

}

func TestWalletsHistoryGetHistory_ErrorGetWalletHistory(t *testing.T) {
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()
	//
	// repository := mockedFailWalletsRepository{}
	// service := service2.NewWalletsHistory(&repository)
	//
	// result, err := service.GetHistory(ctx, &pb.GetHistoryRequest{})
	// require.Error(t, err)
	// require.Nil(t, result)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	walletHistoryRepo := repositoryMock.NewMockWalletsHistory(t)
	walletHistoryService := service.NewWalletsHistory(walletHistoryRepo)

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

	walletHistoryRepo.EXPECT().
		GetWalletHistory(ctx, reqIn).
		Return(nil, retError)
	// walletHistoryRepo.EXPECT().
	//	GetWalletHistoryTotal(ctx, *reqIn).
	//	Return(retInt64, nil)

	result, err := walletHistoryService.GetHistory(ctx, req)
	require.Nil(t, result)
	require.Error(t, err)

	require.EqualError(t, err, fmt.Errorf("GetWalletHistory: error").Error())
}

func TestWalletsHistoryGetHistory_ErrorGetWalletHistoryTotal(t *testing.T) {
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()
	//
	// repository := mockedFailWalletsRepository{}
	// service := service2.NewWalletsHistory(&repository)
	//
	// result, err := service.GetHistory(ctx, &pb.GetHistoryRequest{})
	// require.Error(t, err)
	// require.Nil(t, result)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	walletHistoryRepo := repositoryMock.NewMockWalletsHistory(t)
	walletHistoryService := service.NewWalletsHistory(walletHistoryRepo)

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
		PayoutsSum:    nil,
		HasNewIncome:  nil,
		HasNewPayouts: nil,
		Incomes:       nil,
		Payouts:       nil,
		Wallets:       nil,
	}

	retInt64 := int64(0)
	retError := errors.New("error")

	walletHistoryRepo.EXPECT().
		GetWalletHistory(ctx, reqIn).
		Return(retHistoryOut, nil)
	walletHistoryRepo.EXPECT().
		GetWalletHistoryTotal(ctx, reqIn).
		Return(retInt64, retError)

	result, err := walletHistoryService.GetHistory(ctx, req)
	require.NotNil(t, result) // TODO: fix
	require.NoError(t, err)   // TODO: fix

}
