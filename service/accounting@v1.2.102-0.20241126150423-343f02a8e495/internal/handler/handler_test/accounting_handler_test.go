package handler_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"code.emcdtech.com/emcd/service/accounting/internal/handler"
	"code.emcdtech.com/emcd/service/accounting/model"
	"code.emcdtech.com/emcd/service/accounting/model/enum"
	accountingPb "code.emcdtech.com/emcd/service/accounting/protocol/accounting"
)

func TestAccountingService_ViewBalanceSuccess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := handler.NewAccountingHandler(&mockedSuccessBalanceService{}, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	req := &accountingPb.ViewBalanceRequest{
		UserID:        1,
		CoinID:        "1",
		AccountTypeID: int64(enum.WalletAccountTypeID),
	}

	resp, err := service.ViewBalance(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestAccountingService_ViewBalanceFail(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := handler.NewAccountingHandler(&mockedFailBalanceService{}, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	req := &accountingPb.ViewBalanceRequest{
		UserID:        1,
		CoinID:        "1",
		AccountTypeID: int64(enum.WalletAccountTypeID),
	}

	resp, err := service.ViewBalance(ctx, req)
	require.Error(t, err)
	require.Nil(t, resp)
}

func TestAccountingService_ChangeBalanceSuccess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := handler.NewAccountingHandler(&mockedSuccessBalanceService{}, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	req := &accountingPb.ChangeBalanceRequest{
		Transactions: []*accountingPb.Transaction{},
	}

	resp, err := service.ChangeBalance(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestAccountingService_ChangeBalanceFail(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := handler.NewAccountingHandler(&mockedFailBalanceService{}, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	req := &accountingPb.ChangeBalanceRequest{
		Transactions: []*accountingPb.Transaction{},
	}

	resp, err := service.ChangeBalance(ctx, req)
	require.Error(t, err)
	require.Nil(t, resp)
}

func TestAccountingService_FindOperationsSuccess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := handler.NewAccountingHandler(&mockedSuccessBalanceService{}, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	req := &accountingPb.FindOperationsRequest{
		UserID: 1,
	}

	resp, err := service.FindOperations(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestAccountingService_FindOperationsFail(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := handler.NewAccountingHandler(&mockedFailBalanceService{}, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	req := &accountingPb.FindOperationsRequest{
		UserID: 1,
	}

	resp, err := service.FindOperations(ctx, req)
	require.Error(t, err)
	require.Nil(t, resp)
}

func TestAccountingService_FindTransactionsSuccess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := handler.NewAccountingHandler(&mockedSuccessBalanceService{}, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	req := &accountingPb.FindTransactionsRequest{}

	resp, err := service.FindTransactions(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestAccountingService_FindTransactionsFail(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := handler.NewAccountingHandler(&mockedFailBalanceService{}, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	req := &accountingPb.FindTransactionsRequest{}

	resp, err := service.FindTransactions(ctx, req)
	require.Error(t, err)
	require.Nil(t, resp)
}

func TestAccountingService_FindTransactionsWithBlocksSuccess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := handler.NewAccountingHandler(&mockedSuccessBalanceService{}, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	req := &accountingPb.FindTransactionsWithBlocksRequest{}

	resp, err := service.FindTransactionsWithBlocks(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestAccountingService_FindTransactionsWithBlocksFail(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := handler.NewAccountingHandler(&mockedFailBalanceService{}, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	req := &accountingPb.FindTransactionsWithBlocksRequest{}

	resp, err := service.FindTransactionsWithBlocks(ctx, req)
	require.Error(t, err)
	require.Nil(t, resp)
}

func TestAccountingService_GetTransactionByIDSuccess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := handler.NewAccountingHandler(&mockedSuccessBalanceService{}, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	req := &accountingPb.GetTransactionByIDRequest{}

	resp, err := service.GetTransactionByID(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestAccountingService_GetTransactionByIDFail(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := handler.NewAccountingHandler(&mockedFailBalanceService{}, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	req := &accountingPb.GetTransactionByIDRequest{}

	resp, err := service.GetTransactionByID(ctx, req)
	require.Error(t, err)
	require.Nil(t, resp)
}

func TestAccountingService_FindLastBlockTimeBalancesSuccess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := handler.NewAccountingHandler(&mockedSuccessBalanceService{}, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	req := &accountingPb.FindLastBlockTimeBalancesRequest{}

	resp, err := service.FindLastBlockTimeBalances(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestAccountingService_FindLastBlockTimeBalancesFail(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := handler.NewAccountingHandler(&mockedFailBalanceService{}, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	req := &accountingPb.FindLastBlockTimeBalancesRequest{}

	resp, err := service.FindLastBlockTimeBalances(ctx, req)
	require.Error(t, err)
	require.Nil(t, resp)
}

func TestAccountingService_FindBalancesDiffMiningSuccess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := handler.NewAccountingHandler(&mockedSuccessBalanceService{}, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	req := &accountingPb.FindBalancesDiffMiningRequest{}

	resp, err := service.FindBalancesDiffMining(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestAccountingService_FindBalancesDiffMiningFail(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := handler.NewAccountingHandler(&mockedFailBalanceService{}, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	req := &accountingPb.FindBalancesDiffMiningRequest{}

	resp, err := service.FindBalancesDiffMining(ctx, req)
	require.Error(t, err)
	require.Nil(t, resp)
}

func TestAccountingService_FindBalancesDiffWalletSuccess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := handler.NewAccountingHandler(&mockedSuccessBalanceService{}, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	req := &accountingPb.FindBalancesDiffWalletRequest{}

	resp, err := service.FindBalancesDiffWallet(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestAccountingService_FindBalancesDiffWalletFail(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := handler.NewAccountingHandler(&mockedFailBalanceService{}, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	req := &accountingPb.FindBalancesDiffWalletRequest{}

	resp, err := service.FindBalancesDiffWallet(ctx, req)
	require.Error(t, err)
	require.Nil(t, resp)
}

func TestAccountingService_FindBatchOperationsSuccess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := handler.NewAccountingHandler(&mockedSuccessBalanceService{}, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	req := &accountingPb.FindBatchOperationsRequest{}

	resp, err := service.FindBatchOperations(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestAccountingService_FindBatchOperationsFail(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := handler.NewAccountingHandler(&mockedFailBalanceService{}, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	req := &accountingPb.FindBatchOperationsRequest{}

	resp, err := service.FindBatchOperations(ctx, req)
	require.Error(t, err)
	require.Nil(t, resp)
}

func TestAccountingService_GetIncomesHistorySuccess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := handler.NewAccountingHandler(nil, nil, &mockedSuccessIncomesService{}, nil, nil, nil, nil, nil, nil, nil)
	req := &accountingPb.GetHistoryRequest{Type: string(model.HistoryIncome)}

	resp, err := service.GetHistory(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestAccountingService_GetIncomesHistoryFail(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := handler.NewAccountingHandler(nil, nil, &mockedFailIncomesService{}, nil, nil, nil, nil, nil, nil, nil)
	req := &accountingPb.GetHistoryRequest{Type: string(model.HistoryIncome)}

	resp, err := service.GetHistory(ctx, req)
	require.Error(t, err)
	require.Nil(t, resp)
}

func TestAccountingService_GetPayoutsHistorySuccess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := handler.NewAccountingHandler(nil, nil, nil, &mockedSuccessPayoutsService{}, nil, nil, nil, nil, nil, nil)
	req := &accountingPb.GetHistoryRequest{Type: string(model.HistoryPayout)}

	resp, err := service.GetHistory(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestAccountingService_GetPayoutsHistoryFail(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := handler.NewAccountingHandler(nil, nil, nil, &mockedFailPayoutsService{}, nil, nil, nil, nil, nil, nil)
	req := &accountingPb.GetHistoryRequest{Type: string(model.HistoryPayout)}

	resp, err := service.GetHistory(ctx, req)
	require.Error(t, err)
	require.Nil(t, resp)
}

func TestAccountingService_GetWalletsHistorySuccess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := handler.NewAccountingHandler(nil, &mockedSuccessWalletsService{}, nil, nil, nil, nil, nil, nil, nil, nil)
	req := &accountingPb.GetHistoryRequest{Type: string(model.HistoryWallet)}

	resp, err := service.GetHistory(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestAccountingService_GetWalletsHistoryFail(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := handler.NewAccountingHandler(nil, &mockedFailWalletsService{}, nil, nil, nil, nil, nil, nil, nil, nil)
	req := &accountingPb.GetHistoryRequest{Type: string(model.HistoryWallet)}

	resp, err := service.GetHistory(ctx, req)
	require.Error(t, err)
	require.Nil(t, resp)
}
