package repository_test

import (
	"code.emcdtech.com/emcd/service/accounting/model"
	accountingPb "code.emcdtech.com/emcd/service/accounting/protocol/accounting"
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strconv"
	"testing"
	"time"

	protoAccountingMock "code.emcdtech.com/emcd/service/accounting/mocks/protocol/accounting"
	externalRepository "code.emcdtech.com/emcd/service/accounting/repository"
)

func TestClientRepository_ViewBalanceSuccess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	protoAccountingClientMock := protoAccountingMock.NewMockAccountingServiceClient(t)
	accountingRepository := externalRepository.NewAccountingRepository(protoAccountingClientMock)

	userId := int64(1)
	coin := "1"
	accountTypeId := int64(1)
	totalBalance := false

	protoAccountingClientMock.EXPECT().ViewBalance(
		ctx,
		&accountingPb.ViewBalanceRequest{
			UserID:        userId,
			CoinID:        coin,
			AccountTypeID: accountTypeId,
			TotalBalance:  totalBalance,
		}).Return(&accountingPb.ViewBalanceResponse{
		Balance: "1",
	}, nil)

	balance, err := accountingRepository.ViewBalance(ctx, userId, coin, accountTypeId, totalBalance)
	require.NoError(t, err)
	require.Equal(t, "1", balance)

}

func TestClientRepository_ViewBalanceFail(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	protoAccountingClientMock := protoAccountingMock.NewMockAccountingServiceClient(t)
	accountingRepository := externalRepository.NewAccountingRepository(protoAccountingClientMock)

	protoAccountingClientMock.EXPECT().ViewBalance(
		ctx,
		&accountingPb.ViewBalanceRequest{
			UserID:        1,
			CoinID:        "1",
			AccountTypeID: 1,
			TotalBalance:  false,
		}).Return(&accountingPb.ViewBalanceResponse{
		Balance: "",
	}, errors.New("error"))

	_, err := accountingRepository.ViewBalance(ctx, int64(1), "1", 1, false)
	require.Error(t, err)

}

func TestClientRepository_ChangeBalanceSuccess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	protoAccountingClientMock := protoAccountingMock.NewMockAccountingServiceClient(t)
	accountingRepository := externalRepository.NewAccountingRepository(protoAccountingClientMock)

	protoAccountingClientMock.EXPECT().ChangeBalance(
		ctx,
		&accountingPb.ChangeBalanceRequest{
			Transactions: []*accountingPb.Transaction{},
		}).Return(nil, nil)

	err := accountingRepository.ChangeBalance(ctx, nil)
	require.NoError(t, err)

}

func TestClientRepository_ChangeBalanceFail(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	protoAccountingClientMock := protoAccountingMock.NewMockAccountingServiceClient(t)
	accountingRepository := externalRepository.NewAccountingRepository(protoAccountingClientMock)

	retError := errors.New("error")

	protoAccountingClientMock.EXPECT().ChangeBalance(
		ctx,
		&accountingPb.ChangeBalanceRequest{
			Transactions: []*accountingPb.Transaction{},
		}).Return(nil, retError)

	err := accountingRepository.ChangeBalance(ctx, nil)
	require.EqualError(t, err, fmt.Errorf("accounting: %w", retError).Error())

}

func TestClientRepository_FindOperationsSuccess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	protoAccountingClientMock := protoAccountingMock.NewMockAccountingServiceClient(t)
	accountingRepository := externalRepository.NewAccountingRepository(protoAccountingClientMock)

	userId := int64(1)
	coinId := int64(1)

	protoAccountingClientMock.EXPECT().FindOperations(
		ctx,
		&accountingPb.FindOperationsRequest{
			UserID: userId,
			CoinID: strconv.Itoa(int(coinId)),
		}).Return(&accountingPb.FindOperationsResponse{Operations: nil}, nil)

	_, err := accountingRepository.FindOperations(ctx, userId, coinId)
	require.NoError(t, err)

}

func TestClientRepository_FindOperationsFail(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	protoAccountingClientMock := protoAccountingMock.NewMockAccountingServiceClient(t)
	accountingRepository := externalRepository.NewAccountingRepository(protoAccountingClientMock)

	userId := int64(1)
	coinId := int64(1)

	retError := errors.New("error")

	protoAccountingClientMock.EXPECT().FindOperations(
		ctx,
		&accountingPb.FindOperationsRequest{
			UserID: userId,
			CoinID: strconv.Itoa(int(coinId)),
		}).Return(nil, retError)

	_, err := accountingRepository.FindOperations(ctx, userId, coinId)
	require.Error(t, err)

}

func TestClientRepository_FindTransactionsSuccess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	protoAccountingClientMock := protoAccountingMock.NewMockAccountingServiceClient(t)
	accountingRepository := externalRepository.NewAccountingRepository(protoAccountingClientMock)

	types := []int64{}
	userId := int(1)
	userAccountId := int(1)
	coinIdIntList := []int{1}
	coinIdStrList := []string{"1"}
	from := time.Time{}

	protoAccountingClientMock.EXPECT().FindTransactions(
		ctx,
		&accountingPb.FindTransactionsRequest{
			Types:         types,
			UserID:        int64(userId),
			AccountTypeID: int64(userAccountId),
			CoinIDs:       coinIdStrList,
			From:          timestamppb.New(from),
		}).Return(&accountingPb.FindTransactionsResponse{
		Transactions: nil,
	}, nil)

	_, err := accountingRepository.FindTransactions(ctx, types, userId, userAccountId, coinIdIntList, from)
	require.NoError(t, err)

}

func TestClientRepository_FindTransactionsFail(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	protoAccountingClientMock := protoAccountingMock.NewMockAccountingServiceClient(t)
	accountingRepository := externalRepository.NewAccountingRepository(protoAccountingClientMock)

	types := []int64{}
	userId := int(1)
	userAccountId := int(1)
	coinIdIntList := []int{1}
	coinIdStrList := []string{"1"}
	from := time.Time{}

	retError := errors.New("error")

	protoAccountingClientMock.EXPECT().FindTransactions(
		ctx,
		&accountingPb.FindTransactionsRequest{
			Types:         types,
			UserID:        int64(userId),
			AccountTypeID: int64(userAccountId),
			CoinIDs:       coinIdStrList,
			From:          timestamppb.New(from),
		}).Return(nil, retError)

	_, err := accountingRepository.FindTransactions(ctx, types, userId, userAccountId, coinIdIntList, from)
	require.Error(t, err)

}

func TestClientRepository_FindTransactionsWithBlocksSuccess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	protoAccountingClientMock := protoAccountingMock.NewMockAccountingServiceClient(t)
	accountingRepository := externalRepository.NewAccountingRepository(protoAccountingClientMock)

	blockedTill := time.Time{}

	protoAccountingClientMock.EXPECT().FindTransactionsWithBlocks(
		ctx,
		&accountingPb.FindTransactionsWithBlocksRequest{
			BlockedTill: timestamppb.New(blockedTill),
		}).Return(&accountingPb.FindTransactionsWithBlocksResponse{
		Transactions: nil,
	}, nil)

	_, err := accountingRepository.FindTransactionsWithBlocks(ctx, blockedTill)
	require.NoError(t, err)

}

func TestClientRepository_FindTransactionsWithBlocksFail(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	protoAccountingClientMock := protoAccountingMock.NewMockAccountingServiceClient(t)
	accountingRepository := externalRepository.NewAccountingRepository(protoAccountingClientMock)

	blockedTill := time.Time{}

	retError := errors.New("error")

	protoAccountingClientMock.EXPECT().FindTransactionsWithBlocks(
		ctx,
		&accountingPb.FindTransactionsWithBlocksRequest{
			BlockedTill: timestamppb.New(blockedTill),
		}).Return(nil, retError)

	_, err := accountingRepository.FindTransactionsWithBlocks(ctx, blockedTill)
	require.Error(t, err)

}

func TestClientRepository_FindLastBlockTimeBalancesSuccess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	protoAccountingClientMock := protoAccountingMock.NewMockAccountingServiceClient(t)
	accountingRepository := externalRepository.NewAccountingRepository(protoAccountingClientMock)

	userAccountIDs := []int64{1}

	protoAccountingClientMock.EXPECT().FindLastBlockTimeBalances(
		ctx,
		&accountingPb.FindLastBlockTimeBalancesRequest{
			UserAccountIDs: userAccountIDs,
		}).Return(&accountingPb.FindLastBlockTimeBalancesResponse{
		Balances: nil,
	}, nil)

	_, err := accountingRepository.FindLastBlockTimeBalances(ctx, userAccountIDs)
	require.NoError(t, err)

}

func TestClientRepository_FindLastBlockTimeBalancesFail(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	protoAccountingClientMock := protoAccountingMock.NewMockAccountingServiceClient(t)
	accountingRepository := externalRepository.NewAccountingRepository(protoAccountingClientMock)

	userAccountIDs := []int64{1}

	retError := errors.New("error")

	protoAccountingClientMock.EXPECT().FindLastBlockTimeBalances(
		ctx,
		&accountingPb.FindLastBlockTimeBalancesRequest{
			UserAccountIDs: userAccountIDs,
		}).Return(nil, retError)

	_, err := accountingRepository.FindLastBlockTimeBalances(ctx, userAccountIDs)
	require.Error(t, err)

}

func TestClientRepository_FindBalancesDiffMiningSuccess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	protoAccountingClientMock := protoAccountingMock.NewMockAccountingServiceClient(t)
	accountingRepository := externalRepository.NewAccountingRepository(protoAccountingClientMock)

	reqUsers := []model.UserBeforePayoutMining{}

	protoAccountingClientMock.EXPECT().FindBalancesDiffMining(
		ctx,
		&accountingPb.FindBalancesDiffMiningRequest{
			Users: []*accountingPb.UserBeforePayoutMining{},
		}).Return(&accountingPb.FindBalancesDiffMiningResponse{
		Diffs: nil,
	}, nil)

	_, err := accountingRepository.FindBalancesDiffMining(ctx, reqUsers)
	require.NoError(t, err)

}

func TestClientRepository_FindBalancesDiffMiningFail(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	protoAccountingClientMock := protoAccountingMock.NewMockAccountingServiceClient(t)
	accountingRepository := externalRepository.NewAccountingRepository(protoAccountingClientMock)

	reqUsers := []model.UserBeforePayoutMining{}

	retError := errors.New("error")

	protoAccountingClientMock.EXPECT().FindBalancesDiffMining(
		ctx,
		&accountingPb.FindBalancesDiffMiningRequest{
			Users: []*accountingPb.UserBeforePayoutMining{},
		}).Return(nil, retError)

	_, err := accountingRepository.FindBalancesDiffMining(ctx, reqUsers)
	require.Error(t, err)

}

func TestClientRepository_FindBalancesDiffWalletSuccess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	protoAccountingClientMock := protoAccountingMock.NewMockAccountingServiceClient(t)
	accountingRepository := externalRepository.NewAccountingRepository(protoAccountingClientMock)

	reqUsers := []model.UserBeforePayoutWallet{}

	protoAccountingClientMock.EXPECT().FindBalancesDiffWallet(
		ctx,
		&accountingPb.FindBalancesDiffWalletRequest{
			Users: []*accountingPb.UserBeforePayoutWallet{},
		}).Return(&accountingPb.FindBalancesDiffWalletResponse{
		Diffs: nil,
	}, nil)

	_, err := accountingRepository.FindBalancesDiffWallet(ctx, reqUsers)
	require.NoError(t, err)

}

func TestClientRepository_FindBalancesDiffWalletFail(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	protoAccountingClientMock := protoAccountingMock.NewMockAccountingServiceClient(t)
	accountingRepository := externalRepository.NewAccountingRepository(protoAccountingClientMock)

	reqUsers := []model.UserBeforePayoutWallet{}

	retError := errors.New("error")

	protoAccountingClientMock.EXPECT().FindBalancesDiffWallet(
		ctx,
		&accountingPb.FindBalancesDiffWalletRequest{
			Users: []*accountingPb.UserBeforePayoutWallet{},
		}).Return(nil, retError)

	_, err := accountingRepository.FindBalancesDiffWallet(ctx, reqUsers)
	require.Error(t, err)

}

func TestClientRepository_FindBatchOperationsSuccess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	protoAccountingClientMock := protoAccountingMock.NewMockAccountingServiceClient(t)
	accountingRepository := externalRepository.NewAccountingRepository(protoAccountingClientMock)

	usersWithCoins := make(map[int]int)

	protoAccountingClientMock.EXPECT().FindBatchOperations(
		ctx,
		&accountingPb.FindBatchOperationsRequest{
			Users: []*accountingPb.UserIDCoinID{},
		}).Return(&accountingPb.FindBatchOperationsResponse{
		OperationsByUsers: nil,
	}, nil)

	_, err := accountingRepository.FindBatchOperations(ctx, usersWithCoins)
	require.NoError(t, err)

}

func TestClientRepository_FindBatchOperationsFail(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	protoAccountingClientMock := protoAccountingMock.NewMockAccountingServiceClient(t)
	accountingRepository := externalRepository.NewAccountingRepository(protoAccountingClientMock)

	usersWithCoins := make(map[int]int)

	retError := errors.New("error")

	protoAccountingClientMock.EXPECT().FindBatchOperations(
		ctx,
		&accountingPb.FindBatchOperationsRequest{
			Users: []*accountingPb.UserIDCoinID{},
		}).Return(nil, retError)

	_, err := accountingRepository.FindBatchOperations(ctx, usersWithCoins)
	require.Error(t, err)

}

func TestClientRepository_GetTransactionByIDSuccess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	protoAccountingClientMock := protoAccountingMock.NewMockAccountingServiceClient(t)
	accountingRepository := externalRepository.NewAccountingRepository(protoAccountingClientMock)

	id := int(1)

	protoAccountingClientMock.EXPECT().GetTransactionByID(
		ctx,
		&accountingPb.GetTransactionByIDRequest{
			Id: int64(id),
		}).Return(&accountingPb.GetTransactionByIDResponse{
		Transaction: &accountingPb.TransactionSelectionWithBlock{
			ReceiverAccountID:    0,
			CoinID:               "1",
			Type:                 0,
			Amount:               "100",
			BlockID:              0,
			UnblockToAccountID:   0,
			SenderAccountID:      0,
			UnblockTransactionID: 0,
			ActionID:             "",
		},
	}, nil)

	_, err := accountingRepository.GetTransactionByID(ctx, id)
	require.NoError(t, err)

}

func TestClientRepository_GetTransactionByIDFail(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	protoAccountingClientMock := protoAccountingMock.NewMockAccountingServiceClient(t)
	accountingRepository := externalRepository.NewAccountingRepository(protoAccountingClientMock)

	id := int(1)

	retError := errors.New("error")

	protoAccountingClientMock.EXPECT().GetTransactionByID(
		ctx,
		&accountingPb.GetTransactionByIDRequest{
			Id: int64(id),
		}).Return(nil, retError)

	_, err := accountingRepository.GetTransactionByID(ctx, id)
	require.Error(t, err)

}
