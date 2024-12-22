package handler_test

import (
	"context"
	"errors"
	"time"

	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/types/known/timestamppb"

	"code.emcdtech.com/emcd/service/accounting/model"
	"code.emcdtech.com/emcd/service/accounting/model/enum"
	accountingPb "code.emcdtech.com/emcd/service/accounting/protocol/accounting"
)

// Balance mocks
type mockedSuccessBalanceService struct{}

func (m *mockedSuccessBalanceService) GetBalanceBeforeTransaction(_ context.Context, _ int64, _ int64) (decimal.Decimal, error) {
	return decimal.Decimal{}, nil
}

func (m *mockedSuccessBalanceService) FindOperationsAndTransactions(_ context.Context, _ *model.OperationWithTransactionQuery) ([]*model.OperationWithTransaction, int64, error) {
	return []*model.OperationWithTransaction{}, 0, nil
}

func (m *mockedSuccessBalanceService) GetBalanceByCoin(_ context.Context, _ int32, _ string) (*model.Balance, error) {
	return nil, nil
}

func (m *mockedSuccessBalanceService) GetPaid(_ context.Context, _ int32, _ string, _, _ time.Time) (decimal.Decimal, error) {
	return decimal.Decimal{}, nil
}

func (m *mockedSuccessBalanceService) View(_ context.Context, _ int64, _ enum.AccountTypeId, _ string, _ bool) (decimal.Decimal, error) {
	return decimal.NewFromInt(1), nil
}
func (m *mockedSuccessBalanceService) Change(_ context.Context, _ []*accountingPb.Transaction) error {
	return nil
}
func (m *mockedSuccessBalanceService) FindOperations(_ context.Context, _ int64, _ string) ([]*accountingPb.OperationSelectionWithBlock, error) {
	return []*accountingPb.OperationSelectionWithBlock{}, nil
}
func (m *mockedSuccessBalanceService) FindTransactions(_ context.Context, _ []int64, _ int64, _ int64, _ []string, _ *timestamppb.Timestamp) ([]*accountingPb.Transaction, error) {
	return []*accountingPb.Transaction{}, nil
}
func (m *mockedSuccessBalanceService) FindTransactionsByCollectorFilter(_ context.Context, _ *model.TransactionCollectorFilter) (*uint64, []*accountingPb.Transaction, error) {
	return nil, []*accountingPb.Transaction{}, nil
}
func (m *mockedSuccessBalanceService) GetTransactionsByActionID(_ context.Context, _ string) ([]*accountingPb.Transaction, error) {
	return []*accountingPb.Transaction{}, nil
}
func (m *mockedSuccessBalanceService) FindTransactionsWithBlocks(_ context.Context, _ *timestamppb.Timestamp) ([]*accountingPb.TransactionSelectionWithBlock, error) {
	return []*accountingPb.TransactionSelectionWithBlock{}, nil
}
func (m *mockedSuccessBalanceService) GetTransactionByID(_ context.Context, _ int64) (*accountingPb.TransactionSelectionWithBlock, error) {
	return &accountingPb.TransactionSelectionWithBlock{}, nil
}
func (m *mockedSuccessBalanceService) FindLastBlockTimeBalances(_ context.Context, _ []int64) ([]*accountingPb.UserBlockTimeBalance, error) {
	return []*accountingPb.UserBlockTimeBalance{}, nil
}
func (m *mockedSuccessBalanceService) FindBalancesDiffMining(_ context.Context, _ []*accountingPb.UserBeforePayoutMining) ([]*accountingPb.UserMiningDiff, error) {
	return []*accountingPb.UserMiningDiff{}, nil
}
func (m *mockedSuccessBalanceService) FindBalancesDiffWallet(_ context.Context, _ []*accountingPb.UserBeforePayoutWallet) ([]*accountingPb.UserWalletDiff, error) {
	return []*accountingPb.UserWalletDiff{}, nil
}
func (m *mockedSuccessBalanceService) FindBatchOperations(_ context.Context, _ []*accountingPb.UserIDCoinID) ([]*accountingPb.BatchOperationSelection, error) {
	return []*accountingPb.BatchOperationSelection{}, nil
}
func (m *mockedSuccessBalanceService) ChangeMultiple(_ context.Context, _ []*accountingPb.Transaction) error {
	return nil
}
func (m *mockedSuccessBalanceService) GetBalances(_ context.Context, _ int32) ([]*model.Balance, error) {
	return []*model.Balance{}, nil
}
func (m *mockedSuccessBalanceService) GetCoinsSummary(_ context.Context, _ int32) ([]*model.CoinSummary, error) {
	return []*model.CoinSummary{}, nil
}
func (m *mockedSuccessBalanceService) GetTransactionIDByAction(_ context.Context, _ string, _ int, _ string) (int64, error) {
	return 1, nil
}

type mockedFailBalanceService struct{}

func (m *mockedFailBalanceService) GetBalanceBeforeTransaction(_ context.Context, _ int64, _ int64) (decimal.Decimal, error) {
	return decimal.Decimal{}, nil
}

func (m *mockedFailBalanceService) FindOperationsAndTransactions(_ context.Context, _ *model.OperationWithTransactionQuery) ([]*model.OperationWithTransaction, int64, error) {
	return nil, 0, errors.New("some error")
}

func (m *mockedFailBalanceService) GetBalanceByCoin(_ context.Context, _ int32, _ string) (*model.Balance, error) {
	return nil, nil
}

func (m *mockedFailBalanceService) GetPaid(_ context.Context, _ int32, _ string, _, _ time.Time) (decimal.Decimal, error) {
	return decimal.Decimal{}, nil
}

func (m *mockedFailBalanceService) View(_ context.Context, _ int64, _ enum.AccountTypeId, _ string, _ bool) (decimal.Decimal, error) {
	return decimal.Decimal{}, errors.New("some error")
}
func (m *mockedFailBalanceService) Change(_ context.Context, _ []*accountingPb.Transaction) error {
	return errors.New("some error")
}
func (m *mockedFailBalanceService) FindOperations(_ context.Context, _ int64, _ string) ([]*accountingPb.OperationSelectionWithBlock, error) {
	return nil, errors.New("some error")
}
func (m *mockedFailBalanceService) FindTransactions(_ context.Context, _ []int64, _ int64, _ int64, _ []string, _ *timestamppb.Timestamp) ([]*accountingPb.Transaction, error) {
	return nil, errors.New("some error")
}
func (m *mockedFailBalanceService) FindTransactionsByCollectorFilter(_ context.Context, _ *model.TransactionCollectorFilter) (*uint64, []*accountingPb.Transaction, error) {
	return nil, nil, errors.New("some error")
}
func (m *mockedFailBalanceService) GetTransactionsByActionID(_ context.Context, _ string) ([]*accountingPb.Transaction, error) {
	return nil, errors.New("some error")
}
func (m *mockedFailBalanceService) FindTransactionsWithBlocks(_ context.Context, _ *timestamppb.Timestamp) ([]*accountingPb.TransactionSelectionWithBlock, error) {
	return nil, errors.New("some error")
}
func (m *mockedFailBalanceService) GetTransactionByID(_ context.Context, _ int64) (*accountingPb.TransactionSelectionWithBlock, error) {
	return nil, errors.New("some error")
}
func (m *mockedFailBalanceService) FindLastBlockTimeBalances(_ context.Context, _ []int64) ([]*accountingPb.UserBlockTimeBalance, error) {
	return nil, errors.New("some error")
}
func (m *mockedFailBalanceService) FindBalancesDiffMining(_ context.Context, _ []*accountingPb.UserBeforePayoutMining) ([]*accountingPb.UserMiningDiff, error) {
	return nil, errors.New("some error")
}
func (m *mockedFailBalanceService) FindBalancesDiffWallet(_ context.Context, _ []*accountingPb.UserBeforePayoutWallet) ([]*accountingPb.UserWalletDiff, error) {
	return nil, errors.New("some error")
}
func (m *mockedFailBalanceService) FindBatchOperations(_ context.Context, _ []*accountingPb.UserIDCoinID) ([]*accountingPb.BatchOperationSelection, error) {
	return nil, errors.New("some error")
}
func (m *mockedFailBalanceService) ChangeMultiple(_ context.Context, _ []*accountingPb.Transaction) error {
	return errors.New("some error")
}
func (m *mockedFailBalanceService) GetBalances(_ context.Context, _ int32) ([]*model.Balance, error) {
	return nil, errors.New("some error")
}
func (m *mockedFailBalanceService) GetCoinsSummary(_ context.Context, _ int32) ([]*model.CoinSummary, error) {
	return nil, errors.New("some error")
}
func (m *mockedFailBalanceService) GetTransactionIDByAction(_ context.Context, _ string, _ int, _ string) (int64, error) {
	return 0, errors.New("some error")
}

// Incomes mocks
type mockedSuccessIncomesService struct{}

func (m *mockedSuccessIncomesService) GetHistory(_ context.Context, _ *accountingPb.GetHistoryRequest) (*accountingPb.GetHistoryResponse, error) {
	return &accountingPb.GetHistoryResponse{}, nil
}

type mockedFailIncomesService struct{}

func (m *mockedFailIncomesService) GetHistory(_ context.Context, _ *accountingPb.GetHistoryRequest) (*accountingPb.GetHistoryResponse, error) {
	return nil, errors.New("some error")
}

// Payouts mocks
type mockedSuccessPayoutsService struct{}

func (m *mockedSuccessPayoutsService) GetHistory(_ context.Context, _ *accountingPb.GetHistoryRequest) (*accountingPb.GetHistoryResponse, error) {
	return &accountingPb.GetHistoryResponse{}, nil
}

type mockedFailPayoutsService struct{}

func (m *mockedFailPayoutsService) GetHistory(_ context.Context, _ *accountingPb.GetHistoryRequest) (*accountingPb.GetHistoryResponse, error) {
	return nil, errors.New("some error")
}

// Wallets mocks
type mockedSuccessWalletsService struct{}

func (m *mockedSuccessWalletsService) GetHistory(_ context.Context, _ *accountingPb.GetHistoryRequest) (*accountingPb.GetHistoryResponse, error) {
	return &accountingPb.GetHistoryResponse{}, nil
}

type mockedFailWalletsService struct{}

func (m *mockedFailWalletsService) GetHistory(_ context.Context, _ *accountingPb.GetHistoryRequest) (*accountingPb.GetHistoryResponse, error) {
	return nil, errors.New("some error")
}
