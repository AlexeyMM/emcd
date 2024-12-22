package repository

import (
	"code.emcdtech.com/emcd/service/accounting/model"
	"code.emcdtech.com/emcd/service/accounting/model/enum"
	accountingPb "code.emcdtech.com/emcd/service/accounting/protocol/accounting"
	"context"
	"fmt"
	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strconv"
	"time"
)

type AccountingRepository interface {
	ViewBalance(ctx context.Context, userID int64, coinID string, accountTypeID int64, totalBalance bool) (balance string, err error)
	ChangeBalance(ctx context.Context, transactions []*model.Transaction) error
	FindOperations(ctx context.Context, userID, coinID int64) ([]*model.OperationSelectionWithBlock, error)
	FindBatchOperations(ctx context.Context, usersWithCoins map[int]int) (map[int64][]*model.OperationSelection, error)
	FindTransactions(ctx context.Context, types []int64, userID, userAccountID int, coinIDs []int, from time.Time) ([]*model.Transaction, error)
	GetTransactionsByActionID(ctx context.Context, actionID string) ([]*model.Transaction, error)
	FindTransactionsWithBlocks(ctx context.Context, blockedTill time.Time) ([]*model.TransactionSelectionWithBlock, error)
	GetTransactionByID(ctx context.Context, id int) (*model.TransactionSelectionWithBlock, error)
	GetTransactionIDByAction(ctx context.Context, actionID string, amount decimal.Decimal, Type model.TransactionType) (int, error)
	FindLastBlockTimeBalances(ctx context.Context, data []int64) (map[int64]decimal.Decimal, error)
	FindBalancesDiffMining(ctx context.Context, data []model.UserBeforePayoutMining) (map[int64]decimal.Decimal, error)
	FindBalancesDiffWallet(ctx context.Context, data []model.UserBeforePayoutWallet) (map[int64][]model.UserWalletDiff, error)
}

type accountingRepository struct {
	handler accountingPb.AccountingServiceClient
}

func NewAccountingRepository(handler accountingPb.AccountingServiceClient) AccountingRepository {

	return &accountingRepository{
		handler: handler,
	}
}

func (s *accountingRepository) ViewBalance(ctx context.Context, userID int64,
	coinID string, accountTypeID int64, totalBalance bool) (string, error) {
	request := &accountingPb.ViewBalanceRequest{
		UserID:        userID,
		CoinID:        coinID,
		AccountTypeID: accountTypeID,
		TotalBalance:  totalBalance,
	}

	resp, err := s.handler.ViewBalance(ctx, request)
	if err != nil {
		return "0", fmt.Errorf("accounting: %w", err)
	}

	return resp.Balance, nil
}

func (s *accountingRepository) ChangeBalance(ctx context.Context,
	transactions []*model.Transaction) error {
	data := make([]*accountingPb.Transaction, 0, len(transactions))
	for _, v := range transactions {
		data = append(data, &accountingPb.Transaction{
			Type:              int64(v.Type),
			CreatedAt:         timestamppb.New(v.CreatedAt),
			SenderAccountID:   v.SenderAccountID,
			ReceiverAccountID: v.ReceiverAccountID,
			CoinID:            strconv.Itoa(int(v.CoinID)),
			Amount:            v.Amount.String(),
			Comment:           v.Comment,
			Hash:              v.Hash,
			Hashrate:          v.Hashrate,
			FromReferralId:    v.FromReferralId,
			ReceiverAddress:   v.ReceiverAddress,
			TokenID:           v.TokenID,
			ActionID:          v.ActionID,
		})
	}

	request := &accountingPb.ChangeBalanceRequest{Transactions: data}
	_, err := s.handler.ChangeBalance(ctx, request)
	if err != nil {
		return fmt.Errorf("accounting: %w", err)
	}

	return nil
}

func (s *accountingRepository) FindOperations(ctx context.Context, userID, coinID int64) (
	[]*model.OperationSelectionWithBlock, error) {
	request := &accountingPb.FindOperationsRequest{
		UserID: userID, CoinID: strconv.Itoa(int(coinID)),
	}

	resp, err := s.handler.FindOperations(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("accounting: %w", err)
	}

	result := make([]*model.OperationSelectionWithBlock, len(resp.Operations))
	for i, v := range resp.Operations {
		amount, err := decimal.NewFromString(v.Amount)
		if err != nil {
			return nil, fmt.Errorf("invalid decimal: %w", err)
		}
		operationCoinID, err := strconv.Atoi(v.OperationCoinID)
		if err != nil {
			return nil, fmt.Errorf("invalid operationCoinID: %w", err)
		}
		userAccountCoinID, err := strconv.Atoi(v.UserAccountCoinID)
		if err != nil {
			return nil, fmt.Errorf("invalid userAccountCoinID: %w", err)
		}

		result[i] = &model.OperationSelectionWithBlock{
			Amount:               amount,
			AccountID:            v.AccountID,
			OperationCoinID:      int64(operationCoinID),
			UserAccountCoinID:    int64(userAccountCoinID),
			AccountTypeID:        enum.AccountTypeId(v.AccountTypeID),
			IsActive:             v.IsActive,
			Type:                 model.TransactionType(v.Type),
			CreatedAt:            v.CreatedAt.AsTime(),
			TransactionID:        v.TransactionID,
			UnblockTransactionID: v.UnblockTransactionID,
			UnblockToAccountID:   v.UnblockToAccountID,
		}
	}

	return result, nil
}

func (s *accountingRepository) FindBatchOperations(ctx context.Context,
	usersWithCoins map[int]int) (map[int64][]*model.OperationSelection, error) {
	users := make([]*accountingPb.UserIDCoinID, 0, len(usersWithCoins))
	for k, v := range usersWithCoins {
		users = append(
			users,
			&accountingPb.UserIDCoinID{UserID: int64(k), CoinID: strconv.Itoa(v)},
		)
	}
	request := &accountingPb.FindBatchOperationsRequest{
		Users: users,
	}

	resp, err := s.handler.FindBatchOperations(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("accounting: %w", err)
	}

	result := make(map[int64][]*model.OperationSelection)
	for _, v := range resp.OperationsByUsers {
		for _, o := range v.Operations {
			amount, err := decimal.NewFromString(o.Amount)
			if err != nil {
				return nil, fmt.Errorf("invalid decimal: %w", err)
			}
			operationCoinID, err := strconv.Atoi(o.OperationCoinID)
			if err != nil {
				return nil, fmt.Errorf("invalid operationCoinID: %w", err)
			}
			userAccountCoinID, err := strconv.Atoi(o.UserAccountCoinID)
			if err != nil {
				return nil, fmt.Errorf("invalid userAccountCoinID: %w", err)
			}

			result[v.UserID] = append(
				result[v.UserID],
				&model.OperationSelection{
					Amount:            amount,
					AccountID:         o.AccountID,
					OperationCoinID:   int64(operationCoinID),
					UserAccountCoinID: int64(userAccountCoinID),
					AccountTypeID:     enum.AccountTypeId(o.AccountTypeID),
					IsActive:          o.IsActive,
					Type:              model.TransactionType(o.Type),
					CreatedAt:         o.CreatedAt.AsTime(),
					TransactionID:     o.TransactionID,
				},
			)
		}
	}

	return result, nil
}

func (s *accountingRepository) FindTransactions(ctx context.Context,
	types []int64, userID, userAccountID int, coinIDs []int, from time.Time) ([]*model.Transaction, error) {
	coinIDsString := make([]string, 0, len(coinIDs))
	for _, v := range coinIDs {
		coinIDsString = append(coinIDsString, strconv.Itoa(v))
	}

	request := &accountingPb.FindTransactionsRequest{
		Types:         types,
		UserID:        int64(userID),
		AccountTypeID: int64(userAccountID),
		CoinIDs:       coinIDsString,
		From:          timestamppb.New(from),
	}

	resp, err := s.handler.FindTransactions(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("accounting: %w", err)
	}

	transactions := make([]*model.Transaction, 0, len(resp.Transactions))
	for _, v := range resp.Transactions {
		amount, err := decimal.NewFromString(v.Amount)
		if err != nil {
			return nil, fmt.Errorf("invalid decimal: %w", err)
		}
		coinIDInt, err := strconv.Atoi(v.CoinID)
		if err != nil {
			return nil, fmt.Errorf("invalid coinID: %w", err)
		}

		t := model.Transaction{
			ID:                0,
			Type:              model.TransactionType(v.Type),
			CreatedAt:         v.CreatedAt.AsTime(),
			SenderAccountID:   v.SenderAccountID,
			ReceiverAccountID: v.ReceiverAccountID,
			CoinID:            int64(coinIDInt),
			TokenID:           v.TokenID,
			Amount:            amount,
			Comment:           v.Comment,
			FromReferralId:    v.FromReferralId,
			Hash:              v.Hash,
			ReceiverAddress:   v.ReceiverAddress,
			Hashrate:          v.Hashrate,
			ActionID:          v.ActionID,
			UnblockAccountId:  0, // TODO: why empty
			BlockedTill:       time.Time{},
		}
		transactions = append(transactions, &t)
	}

	return transactions, nil
}

func (s *accountingRepository) GetTransactionsByActionID(ctx context.Context, actionID string) ([]*model.Transaction, error) {
	request := &accountingPb.GetTransactionsByActionIDRequest{
		ActionID: actionID,
	}

	resp, err := s.handler.GetTransactionsByActionID(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("accounting: %w", err)
	}

	transactions := make([]*model.Transaction, 0, len(resp.Transactions))
	for _, v := range resp.Transactions {
		amount, err := decimal.NewFromString(v.Amount)
		if err != nil {
			return nil, fmt.Errorf("invalid decimal: %w", err)
		}
		coinIDInt, err := strconv.Atoi(v.CoinID)
		if err != nil {
			return nil, fmt.Errorf("invalid coinID: %w", err)
		}

		t := model.Transaction{
			ID:                0,
			Type:              model.TransactionType(v.Type),
			CreatedAt:         v.CreatedAt.AsTime(),
			SenderAccountID:   v.SenderAccountID,
			ReceiverAccountID: v.ReceiverAccountID,
			CoinID:            int64(coinIDInt),
			TokenID:           v.TokenID,
			Amount:            amount,
			Comment:           v.Comment,
			Hash:              v.Hash,
			ReceiverAddress:   v.ReceiverAddress,
			Hashrate:          v.Hashrate,
			ActionID:          v.ActionID,
			FromReferralId:    0, // TODO: why empty
			UnblockAccountId:  0, // TODO: why empty
			BlockedTill:       time.Time{},
		}
		transactions = append(transactions, &t)
	}

	return transactions, nil
}

func (s *accountingRepository) FindTransactionsWithBlocks(ctx context.Context,
	blockedTill time.Time) ([]*model.TransactionSelectionWithBlock, error) {
	request := &accountingPb.FindTransactionsWithBlocksRequest{BlockedTill: timestamppb.New(blockedTill)}

	resp, err := s.handler.FindTransactionsWithBlocks(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("accounting: %w", err)
	}

	result := make([]*model.TransactionSelectionWithBlock, 0, len(resp.Transactions))
	for _, v := range resp.Transactions {
		amount, err := decimal.NewFromString(v.Amount)
		if err != nil {
			return nil, fmt.Errorf("invalid decimal: %w", err)
		}
		coinID, err := strconv.Atoi(v.CoinID)
		if err != nil {
			return nil, fmt.Errorf("invalid coinID: %w", err)
		}

		result = append(
			result,
			&model.TransactionSelectionWithBlock{
				SenderAccountID:      0, // TODO: why empty?
				ReceiverAccountID:    v.ReceiverAccountID,
				CoinID:               int64(coinID),
				Type:                 model.TransactionType(v.Type),
				Amount:               amount,
				BlockID:              v.BlockID,
				UnblockToAccountID:   v.UnblockToAccountID,
				UnblockTransactionID: 0, // TODO: why empty?
				ActionID:             v.ActionID,
			},
		)
	}

	return result, nil
}

func (s *accountingRepository) GetTransactionIDByAction(ctx context.Context,
	actionID string, amount decimal.Decimal, transactionType model.TransactionType) (int, error) {
	request := &accountingPb.GetTransactionIDByActionRequest{
		ActionID: actionID, Amount: amount.String(), Type: int64(transactionType),
	}

	resp, err := s.handler.GetTransactionIDByAction(ctx, request)
	if err != nil {
		return 0, fmt.Errorf("accounting: %w", err)
	}

	return int(resp.Id), nil
}

func (s *accountingRepository) GetTransactionByID(ctx context.Context, id int) (
	*model.TransactionSelectionWithBlock, error) {
	request := &accountingPb.GetTransactionByIDRequest{
		Id: int64(id),
	}

	resp, err := s.handler.GetTransactionByID(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("accounting: %w", err)
	}

	t := resp.Transaction
	coinID, err := strconv.Atoi(t.CoinID)
	if err != nil {
		return nil, fmt.Errorf("invalid coinID: %w", err)
	}
	amount, err := decimal.NewFromString(t.Amount)
	if err != nil {
		return nil, fmt.Errorf("invalid decimal: %w", err)
	}

	result := model.TransactionSelectionWithBlock{
		SenderAccountID:      t.SenderAccountID,
		ReceiverAccountID:    t.ReceiverAccountID,
		CoinID:               int64(coinID),
		Type:                 model.TransactionType(t.Type),
		Amount:               amount,
		BlockID:              0, // TODO: why empty?
		UnblockToAccountID:   t.UnblockToAccountID,
		UnblockTransactionID: t.UnblockTransactionID,
		ActionID:             t.ActionID,
	}

	return &result, nil
}

func (s *accountingRepository) FindLastBlockTimeBalances(ctx context.Context,
	userAccountIDs []int64) (map[int64]decimal.Decimal, error) {

	request := &accountingPb.FindLastBlockTimeBalancesRequest{
		UserAccountIDs: userAccountIDs,
	}

	resp, err := s.handler.FindLastBlockTimeBalances(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("accounting: %w", err)
	}

	result := make(map[int64]decimal.Decimal, len(resp.Balances))
	for _, v := range resp.Balances {
		balance, err := decimal.NewFromString(v.Balance)
		if err != nil {
			return nil, fmt.Errorf("invalid string: %w", err)
		}
		result[v.UserID] = balance
	}

	return result, nil
}

func (s *accountingRepository) FindBalancesDiffMining(ctx context.Context,
	data []model.UserBeforePayoutMining) (map[int64]decimal.Decimal, error) {
	users := make([]*accountingPb.UserBeforePayoutMining, 0, len(data))
	for i := range data {
		user := accountingPb.UserBeforePayoutMining{
			UserID:        data[i].UserID,
			CoinID:        strconv.Itoa(int(data[i].CoinID)),
			BlockID:       data[i].BlockID,
			AccountTypeID: int64(data[i].AccountTypeID),
			LastPay:       timestamppb.New(data[i].LastPay),
		}
		users = append(users, &user)
	}

	request := &accountingPb.FindBalancesDiffMiningRequest{
		Users: users,
	}
	resp, err := s.handler.FindBalancesDiffMining(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("accounting: %w", err)
	}

	result := make(map[int64]decimal.Decimal, len(resp.Diffs))
	for _, v := range resp.Diffs {
		diff, err := decimal.NewFromString(v.Diff)
		if err != nil {
			return nil, fmt.Errorf("invalid string: %w", err)
		}
		result[v.BlockID] = diff
	}

	return result, nil
}

func (s *accountingRepository) FindBalancesDiffWallet(ctx context.Context,
	data []model.UserBeforePayoutWallet) (map[int64][]model.UserWalletDiff, error) {
	users := make([]*accountingPb.UserBeforePayoutWallet, 0, len(data))
	for i := range data {
		user := accountingPb.UserBeforePayoutWallet{
			UserID:         data[i].UserID,
			CoinID:         strconv.Itoa(int(data[i].CoinID)),
			AccountTypeID:  int64(data[i].AccountTypeID),
			TransactionIDs: data[i].TransactionIDs,
		}
		users = append(users, &user)
	}

	request := &accountingPb.FindBalancesDiffWalletRequest{
		Users: users,
	}
	resp, err := s.handler.FindBalancesDiffWallet(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("accounting: %w", err)
	}

	result := make(map[int64][]model.UserWalletDiff)
	for _, v := range resp.Diffs {
		diff, err := decimal.NewFromString(v.Diff)
		if err != nil {
			return nil, fmt.Errorf("invalid string: %w", err)
		}
		result[v.UserID] = append(result[v.UserID], model.UserWalletDiff{UserID: v.UserID, BlockID: v.BlockID, Diff: diff})
	}

	return result, nil
}
