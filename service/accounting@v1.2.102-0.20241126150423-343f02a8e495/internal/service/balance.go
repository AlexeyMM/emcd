package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	sdkLog "code.emcdtech.com/emcd/sdk/log"

	"code.emcdtech.com/emcd/service/accounting/model/enum"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/types/known/timestamppb"

	"code.emcdtech.com/emcd/service/accounting/internal/config"
	"code.emcdtech.com/emcd/service/accounting/internal/repository"
	"code.emcdtech.com/emcd/service/accounting/model"
	accountingPb "code.emcdtech.com/emcd/service/accounting/protocol/accounting"
)

type Balance interface {
	View(ctx context.Context, userID int64, accountTypeID enum.AccountTypeId, coinID string, totalBalance bool) (decimal.Decimal, error)
	Change(ctx context.Context, transactions []*accountingPb.Transaction) error
	FindOperations(ctx context.Context, userID int64, coinID string) ([]*accountingPb.OperationSelectionWithBlock, error)
	FindBatchOperations(ctx context.Context, usersWithCoins []*accountingPb.UserIDCoinID) ([]*accountingPb.BatchOperationSelection, error)
	FindTransactions(ctx context.Context, types []int64, userID, accountTypeID int64, coinIDs []string, from *timestamppb.Timestamp) ([]*accountingPb.Transaction, error)
	FindTransactionsByCollectorFilter(ctx context.Context, filter *model.TransactionCollectorFilter) (*uint64, []*accountingPb.Transaction, error)
	GetTransactionsByActionID(ctx context.Context, actionID string) ([]*accountingPb.Transaction, error)
	FindTransactionsWithBlocks(ctx context.Context, blockedTill *timestamppb.Timestamp) ([]*accountingPb.TransactionSelectionWithBlock, error)
	GetTransactionByID(ctx context.Context, id int64) (*accountingPb.TransactionSelectionWithBlock, error)
	FindLastBlockTimeBalances(ctx context.Context, userAccountIDs []int64) ([]*accountingPb.UserBlockTimeBalance, error)
	FindBalancesDiffMining(ctx context.Context, data []*accountingPb.UserBeforePayoutMining) ([]*accountingPb.UserMiningDiff, error)
	FindBalancesDiffWallet(ctx context.Context, data []*accountingPb.UserBeforePayoutWallet) ([]*accountingPb.UserWalletDiff, error)
	ChangeMultiple(ctx context.Context, trs []*accountingPb.Transaction) error
	GetBalances(ctx context.Context, userID int32) ([]*model.Balance, error)
	GetBalanceByCoin(ctx context.Context, userID int32, coin string) (*model.Balance, error)
	GetPaid(ctx context.Context, userID int32, coin string, from, to time.Time) (decimal.Decimal, error)
	GetCoinsSummary(ctx context.Context, userID int32) ([]*model.CoinSummary, error)
	GetTransactionIDByAction(ctx context.Context, actionID string, txType int, amount string) (int64, error)
	FindOperationsAndTransactions(ctx context.Context, request *model.OperationWithTransactionQuery) ([]*model.OperationWithTransaction, int64, error)
	GetBalanceBeforeTransaction(ctx context.Context, accountID, transactionID int64) (decimal.Decimal, error)
}

type TxBeginner interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

type changeFunc = func() error

type balance struct {
	walletCoinsIDs    []int
	walletCoinsStrIDs map[int]string
	balanceRepository repository.Balance
	pool              TxBeginner
	serviceData       config.ServiceData
	reward            repository.Reward
	whiteLabel        repository.WhiteLabel
	// referral commission to emcd and wls should be ignored
	ignoreReferralPaymentByUserID map[uuid.UUID]bool
	p2pAdminId                    int32
}

const mining = "mining"
const balanceInitTimeout = time.Second * 10

func NewBalance(
	walletCoinsIDs []int,
	walletCoinsStrIDs map[int]string,
	balanceRepository repository.Balance,
	pool TxBeginner,
	serviceData config.ServiceData,
	reward repository.Reward,
	whiteLabel repository.WhiteLabel,
	ignoreReferralPaymentByUserID map[uuid.UUID]bool,
) (Balance, error) {

	ctx, cancel := context.WithTimeout(context.Background(), balanceInitTimeout)
	defer cancel()

	p2pAdminId, err := balanceRepository.GetP2PAdminId(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting P2PAdminId: %w", err)
	}

	out := &balance{
		walletCoinsIDs:                walletCoinsIDs,
		walletCoinsStrIDs:             walletCoinsStrIDs,
		balanceRepository:             balanceRepository,
		pool:                          pool,
		serviceData:                   serviceData,
		reward:                        reward,
		whiteLabel:                    whiteLabel,
		ignoreReferralPaymentByUserID: ignoreReferralPaymentByUserID,
		p2pAdminId:                    int32(p2pAdminId),
	}

	return out, nil
}

func (s *balance) View(ctx context.Context, userID int64, accountTypeID enum.AccountTypeId, coinID string, totalBalance bool) (decimal.Decimal, error) {
	coinIDInt, err := strconv.Atoi(coinID)
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("can't convert to int: %w", err)
	}

	amount, err := s.balanceRepository.View(ctx, userID, int64(coinIDInt), accountTypeID, totalBalance)
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("repository: %w", err)
	}

	return amount, nil
}

func (s *balance) Change(ctx context.Context, transactions []*accountingPb.Transaction) error {
	_, commissionsCalculated := containsType(transactions, int64(model.UserPaysPoolComsTrTypeID))
	_, refCommissionCalculated := containsType(transactions, int64(model.PoolPaysUsersReferralsTrTypeID))
	idx, isMiningPayments := containsType(transactions, int64(model.MainCoinMiningPayoutTrTypeID))
	if isMiningPayments && !(commissionsCalculated || refCommissionCalculated) &&
		!isLostUserAccountID(transactions[idx].SenderAccountID) {
		miningAccountID := transactions[idx].ReceiverAccountID
		adminAccountID := transactions[idx].SenderAccountID
		amount, err := decimal.NewFromString(transactions[idx].Amount)
		if err != nil {
			return fmt.Errorf("parsing amount: %s. %w", transactions[idx].Amount, err)
		}
		coinID, err := strconv.Atoi(transactions[idx].CoinID)
		if err != nil {
			return fmt.Errorf("parse coin id: %s. %w", transactions[idx].CoinID, err)
		}

		actionID, err := uuid.Parse(transactions[idx].ActionID)
		if err != nil {
			return fmt.Errorf("parse action id: %s. %w", transactions[idx].ActionID, err)
		}
		hashrate := transactions[idx].Hashrate

		rewardTxs, err := s.getRewardTransactionsForMining(ctx, s.walletCoinsStrIDs, int(miningAccountID), int(adminAccountID), coinID, hashrate, amount, actionID, transactions[idx].CreatedAt.AsTime())
		if err != nil {
			return fmt.Errorf("getRewardTransactionsForMining: %w", err)
		}
		transactions = append(transactions, rewardTxs...)
		logUserTxs(ctx, transactions)
	}

	sqlTx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("pool.Begin: %w", err)
	}
	defer func(sqlTx pgx.Tx, ctx context.Context) {
		err := sqlTx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			sdkLog.Error(ctx, err.Error())

		}
	}(sqlTx, ctx)

	for _, transaction := range transactions {
		err = s.changeForTransaction(ctx, sqlTx, transaction)
		if err != nil {
			return fmt.Errorf("change: %w", err)
		}
	}

	err = sqlTx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("sqlTx.Commit: %w", err)
	}
	return nil
}

func isLostUserAccountID(accID int64) bool {
	lostUserAccountIDs := []int64{1788713, 1788715, 1788716, 1788714, 1788717, 1788720, 1788719, 1788718, 6682057}
	for _, lostAccID := range lostUserAccountIDs {
		if accID == lostAccID {
			return true
		}
	}
	return false
}

func containsType(trs []*accountingPb.Transaction, targetType int64) (int, bool) {
	for i := range trs {
		if trs[i].Type == targetType {
			return i, true
		}
	}
	return -1, false
}

func logUserTxs(ctx context.Context, txArr []*accountingPb.Transaction) {
	if len(txArr) == 0 {
		return
	}

	str := "ActionID[" + txArr[0].ActionID + "] :" + " TX: ["

	for _, tx := range txArr {
		str += fmt.Sprintf("{Type: %d, Amount: %s, Sender: %d, Receiver: %d},", tx.Type, tx.Amount, tx.SenderAccountID, tx.ReceiverAccountID)
	}

	str += "]"

	sdkLog.Info(ctx, str)
}

func (s *balance) getRewardTransactionsForMining(
	ctx context.Context,
	walletCoinsStrIDs map[int]string,
	userMiningAccountID, adminAccountID, coinID int,
	hashrate int64,
	amount decimal.Decimal, actionID uuid.UUID, createdAt time.Time,
) ([]*accountingPb.Transaction, error) {
	userID, oldUserID, refID, err := s.balanceRepository.GetUserIDsByAccountID(ctx, userMiningAccountID)
	if err != nil {
		return nil, fmt.Errorf("usersAccounts.GetUserIDByAccountID: %w", err)
	}
	coin := strings.ToUpper(walletCoinsStrIDs[coinID])
	incomes, err := s.reward.GetReward(ctx, userID, mining, coin, amount)
	if err != nil {
		return nil, fmt.Errorf("reward.GetReward: userID: %s. %w", userID.String(), err)
	}
	var (
		wlAccountID       int
		referrerAccountID = -1
	)
	for i := range incomes {
		switch incomes[i].Type {
		case model.UserPaysWlComsTrTypeID:
			wlID := incomes[i].UserID
			wlOldUserID, err := s.whiteLabel.GetUserIDByID(ctx, wlID)
			if err != nil {
				return nil, fmt.Errorf("whiteLabel.GetUserIDByID: %w", err)
			}
			wlAccountID, err = s.balanceRepository.GetUserAccountIDByOldID(ctx, wlOldUserID, enum.MiningAccountTypeID, coinID)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					sdkLog.Error(ctx, "balanceRepository.GetUserAccountIDByOldID: userID: %s. wlID: %s. wlOldUserID: %d. coin_id: %d. %v", userID.String(), wlID.String(), wlOldUserID, coinID, err)
					wlAccountID, err = s.createMiningAccount(ctx, wlOldUserID, coinID)
					if err != nil {
						return nil, fmt.Errorf("createMiningAccount: %w", err)
					}
				} else {
					return nil, fmt.Errorf("balance.GetUserAccountIDByOldID: %w", err)
				}
			}
		case model.PoolPaysUsersReferralsTrTypeID, model.WlPaysUserReferralsTrTypeID:
			// referral commission to emcd and wls should be ignored
			if !s.ignoreReferralPaymentByUserID[incomes[i].UserID] {
				referrerAccountID, err = s.balanceRepository.GetUserAccountIDByNewID(ctx, incomes[i].UserID, enum.ReferralAccountTypeID, coinID)
				if err != nil {
					if errors.Is(err, pgx.ErrNoRows) {
						sdkLog.Error(ctx, "balanceRepository.GetUserAccountIDByNewID: userID: %d, %s. %v", refID, incomes[i].UserID.String(), err)
						referrerAccountID, err = s.createReferralAccount(ctx, refID, coinID)
						if err != nil {
							return nil, fmt.Errorf("createReferralAccount: %w", err)
						}
					} else {
						return nil, fmt.Errorf("balance.GetReferrerAccountID: %w", err)
					}
				}
			}
		}
	}
	rs := s.convertIncomesToTransactions(incomes, oldUserID, userMiningAccountID, adminAccountID, wlAccountID, referrerAccountID, coinID, hashrate, actionID, createdAt)

	return rs, nil
}

func (s *balance) createMiningAccount(ctx context.Context, oldID int32, coinID int) (int, error) {
	miningAccountID, err := s.balanceRepository.CreateUsersAccount(ctx, oldID, enum.MiningAccountTypeID, coinID, model.MinPayDefault[coinID])
	if err != nil {
		return 0, fmt.Errorf("balanceRepository.CreateUsersAccount: %w", err)
	}

	if err = s.balanceRepository.CreateAccountPool(ctx, miningAccountID); err != nil {
		return 0, fmt.Errorf("balanceRepository.CreateAccountPool: %w", err)
	}
	return miningAccountID, nil
}

func (s *balance) createReferralAccount(ctx context.Context, oldUserID int32, coinID int) (int, error) {
	refAccountID, err := s.balanceRepository.CreateUsersAccount(ctx, oldUserID, enum.ReferralAccountTypeID, coinID, 0)
	if err != nil {
		return 0, fmt.Errorf("balanceRepository.CreateUsersAccount: %w", err)
	}

	if err = s.balanceRepository.CreateAccountReferral(ctx, refAccountID, coinID); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("balanceRepository.CreateAccountReferral: %w", err)
	}
	return refAccountID, nil
}

const (
	defaultPrecision = 8
)

func (s *balance) convertIncomesToTransactions(incomes []*model.UserIncome, oldUserID int32, userMiningAccountID, adminAccountID, wlAccountID, referrerAccountID,
	coinID int, hashrate int64, actionID uuid.UUID, createdAt time.Time) []*accountingPb.Transaction {
	receiverToIncome := make(map[int]*model.UserIncome)

	for i := range incomes {
		switch incomes[i].Type {
		case model.UserPaysPoolComsTrTypeID:
			receiverToIncome[adminAccountID] = incomes[i]
		case model.PoolPaysUsersReferralsTrTypeID:
			receiverToIncome[referrerAccountID] = incomes[i]
		case model.UserPaysWlComsTrTypeID:
			receiverToIncome[wlAccountID] = incomes[i]
		case model.WlPaysUserReferralsTrTypeID:
			receiverToIncome[referrerAccountID] = incomes[i]
		}
	}

	for _, income := range receiverToIncome {
		var sender int
		switch income.Type {
		case model.UserPaysPoolComsTrTypeID:
			sender = userMiningAccountID
		case model.PoolPaysUsersReferralsTrTypeID:
			sender = adminAccountID
		case model.UserPaysWlComsTrTypeID:
			sender = userMiningAccountID
		case model.WlPaysUserReferralsTrTypeID:
			sender = wlAccountID
		}
		if _, ok := receiverToIncome[sender]; ok {
			receiverToIncome[sender].Amount = receiverToIncome[sender].Amount.Add(income.Amount)
		}
	}

	res := make([]*accountingPb.Transaction, 0, len(incomes))
	for i := range incomes {
		if s.ignoreReferralPaymentByUserID[incomes[i].UserID] && (incomes[i].Type == model.PoolPaysUsersReferralsTrTypeID || incomes[i].Type == model.WlPaysUserReferralsTrTypeID) {
			continue
		}
		amount := incomes[i].Amount.Truncate(defaultPrecision)
		if amount.Equal(decimal.Zero) {
			continue
		}

		tr := &accountingPb.Transaction{
			Type:              int64(incomes[i].Type),
			CreatedAt:         timestamppb.New(createdAt),
			SenderAccountID:   0,
			ReceiverAccountID: 0,
			CoinID:            strconv.Itoa(coinID),
			Amount:            amount.String(),
			Comment:           "",
			Hash:              "",
			ReceiverAddress:   "",
			TokenID:           0,
			Hashrate:          0,
			FromReferralId:    0,
			ActionID:          actionID.String(),
		}

		switch incomes[i].Type {
		case model.UserPaysPoolComsTrTypeID:
			tr.SenderAccountID = int64(userMiningAccountID)
			tr.ReceiverAccountID = int64(adminAccountID)
		case model.PoolPaysUsersReferralsTrTypeID:
			tr.SenderAccountID = int64(adminAccountID)
			tr.ReceiverAccountID = int64(referrerAccountID)
			tr.FromReferralId = int64(oldUserID)
			tr.Hashrate = hashrate
		case model.UserPaysWlComsTrTypeID:
			tr.SenderAccountID = int64(userMiningAccountID)
			tr.ReceiverAccountID = int64(wlAccountID)
		case model.WlPaysUserReferralsTrTypeID:
			tr.SenderAccountID = int64(wlAccountID)
			tr.ReceiverAccountID = int64(referrerAccountID)
			tr.FromReferralId = int64(oldUserID)
			tr.Hashrate = hashrate
		}
		res = append(res, tr)
	}
	return res
}

func (s *balance) FindOperations(ctx context.Context, userID int64, coinID string) ([]*accountingPb.OperationSelectionWithBlock, error) {
	coinIDInt, err := strconv.Atoi(coinID)
	if err != nil {
		return nil, fmt.Errorf("convert to int: %w", err)
	}

	operations, err := s.balanceRepository.FindOperations(ctx, int(userID), coinIDInt)
	if err != nil {
		return nil, fmt.Errorf("repository: %w", err)
	}

	result := make([]*accountingPb.OperationSelectionWithBlock, len(operations))
	for i, v := range operations {
		result[i] = &accountingPb.OperationSelectionWithBlock{
			Amount:               v.Amount.String(),
			AccountID:            v.AccountID,
			OperationCoinID:      strconv.Itoa(int(v.OperationCoinID)),
			UserAccountCoinID:    strconv.Itoa(int(v.UserAccountCoinID)),
			AccountTypeID:        int64(v.AccountTypeID),
			IsActive:             v.IsActive,
			Type:                 int64(v.Type),
			CreatedAt:            timestamppb.New(v.CreatedAt),
			TransactionID:        v.TransactionID,
			UnblockTransactionID: v.UnblockTransactionID,
			UnblockToAccountID:   v.UnblockToAccountID,
		}
	}

	return result, nil
}

func (s *balance) FindBatchOperations(ctx context.Context, usersWithCoins []*accountingPb.UserIDCoinID) ([]*accountingPb.BatchOperationSelection, error) {
	data := make(map[int]int, len(usersWithCoins))
	for _, v := range usersWithCoins {
		coinIDInt, err := strconv.Atoi(v.CoinID)
		if err != nil {
			return nil, fmt.Errorf("convert to int: %w", err)
		}
		data[int(v.UserID)] = coinIDInt
	}
	operationsByUsers, err := s.balanceRepository.FindBatchOperations(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("repository: %w", err)
	}

	result := make([]*accountingPb.BatchOperationSelection, 0, len(operationsByUsers))
	for k, v := range operationsByUsers {
		operations := make([]*accountingPb.OperationSelection, 0)
		for _, o := range v {
			operations = append(
				operations,
				&accountingPb.OperationSelection{
					Amount:            o.Amount.String(),
					AccountID:         o.AccountID,
					OperationCoinID:   strconv.Itoa(int(o.OperationCoinID)),
					UserAccountCoinID: strconv.Itoa(int(o.UserAccountCoinID)),
					AccountTypeID:     int64(o.AccountTypeID),
					IsActive:          o.IsActive,
					Type:              int64(o.Type),
					CreatedAt:         timestamppb.New(o.CreatedAt),
					TransactionID:     o.TransactionID,
				},
			)
		}
		result = append(
			result,
			&accountingPb.BatchOperationSelection{
				UserID:     int64(k),
				Operations: operations,
			},
		)
	}

	return result, nil
}

func (s *balance) FindTransactions(ctx context.Context, types []int64, userID, accountTypeID int64, coinIDs []string,
	from *timestamppb.Timestamp) ([]*accountingPb.Transaction, error) {
	coinIDsInt := make([]int, 0, len(coinIDs))
	for _, v := range coinIDs {
		coinID, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("convert to int: %w", err)
		}
		coinIDsInt = append(coinIDsInt, coinID)
	}

	preparedTypes := make([]int, 0, len(types))
	for _, v := range types {
		preparedTypes = append(preparedTypes, int(v))
	}

	transactions, err := s.balanceRepository.FindTransactions(
		ctx, preparedTypes, int(userID), int(accountTypeID), coinIDsInt, from.AsTime())
	if err != nil {
		return nil, fmt.Errorf("balanceRepository.FindTransactions: %w", err)
	}

	var result = make([]*accountingPb.Transaction, 0, len(transactions))
	for _, v := range transactions {
		t := accountingPb.Transaction{
			Type:              int64(v.Type),
			CreatedAt:         timestamppb.New(v.CreatedAt),
			SenderAccountID:   v.SenderAccountID,
			ReceiverAccountID: v.ReceiverAccountID,
			CoinID:            strconv.Itoa(int(v.CoinID)),
			Amount:            v.Amount.String(),
			Comment:           v.Comment,
			Hash:              v.Hash,
			ReceiverAddress:   v.ReceiverAddress,
			TokenID:           v.TokenID,
			Hashrate:          v.Hashrate,
			ActionID:          v.ActionID,
			FromReferralId:    0, // TODO: why empty?
		}
		result = append(result, &t)
	}

	return result, nil
}

func (s *balance) FindTransactionsByCollectorFilter(ctx context.Context, filter *model.TransactionCollectorFilter) (*uint64, []*accountingPb.Transaction, error) {
	if totalCount, transactions, err := s.balanceRepository.FindTransactionsByCollectorFilter(ctx, filter); err != nil {

		return nil, nil, fmt.Errorf("balanceRepository.FindTransactionsByCollectorFilter: %w", err)
	} else {
		var result []*accountingPb.Transaction
		for _, v := range transactions {
			t := accountingPb.Transaction{
				Type:              int64(v.Type),
				CreatedAt:         timestamppb.New(v.CreatedAt),
				SenderAccountID:   v.SenderAccountID,
				ReceiverAccountID: v.ReceiverAccountID,
				CoinID:            strconv.Itoa(int(v.CoinID)),
				Amount:            v.Amount.String(),
				Comment:           v.Comment,
				Hash:              v.Hash,
				ReceiverAddress:   v.ReceiverAddress,
				TokenID:           v.TokenID,
				Hashrate:          v.Hashrate,
				ActionID:          v.ActionID,
				FromReferralId:    0, // TODO: why empty?
			}
			result = append(result, &t)
		}

		return totalCount, result, nil
	}
}

func (s *balance) GetTransactionsByActionID(ctx context.Context, actionID string) ([]*accountingPb.Transaction, error) {
	transactions, err := s.balanceRepository.GetTransactionsByActionID(ctx, actionID)
	if err != nil {
		return nil, fmt.Errorf("repository: %w", err)
	}

	var result = make([]*accountingPb.Transaction, 0, len(transactions))
	for _, v := range transactions {
		t := accountingPb.Transaction{
			Type:              int64(v.Type),
			CreatedAt:         timestamppb.New(v.CreatedAt),
			SenderAccountID:   v.SenderAccountID,
			ReceiverAccountID: v.ReceiverAccountID,
			CoinID:            strconv.Itoa(int(v.CoinID)),
			Amount:            v.Amount.String(),
			Comment:           v.Comment,
			Hash:              v.Hash,
			ReceiverAddress:   v.ReceiverAddress,
			TokenID:           v.TokenID,
			Hashrate:          v.Hashrate,
			ActionID:          v.ActionID,
			FromReferralId:    0, // TODO: why empty?
		}
		result = append(result, &t)
	}

	return result, nil
}

func (s *balance) FindTransactionsWithBlocks(ctx context.Context, blockedTill *timestamppb.Timestamp) (
	[]*accountingPb.TransactionSelectionWithBlock, error) {
	transactions, err := s.balanceRepository.FindTransactionsWithBlocks(ctx, blockedTill.AsTime())
	if err != nil {
		return nil, fmt.Errorf("repository: %w", err)
	}

	result := make([]*accountingPb.TransactionSelectionWithBlock, 0, len(transactions))
	for _, v := range transactions {
		result = append(
			result,
			&accountingPb.TransactionSelectionWithBlock{
				ReceiverAccountID:    v.ReceiverAccountID,
				CoinID:               strconv.Itoa(int(v.CoinID)),
				Type:                 int64(v.Type),
				Amount:               v.Amount.String(),
				BlockID:              v.BlockID,
				UnblockToAccountID:   v.UnblockToAccountID,
				SenderAccountID:      0, // TODO: why empty?
				UnblockTransactionID: 0, // TODO: why empty?
				ActionID:             v.ActionID,
			},
		)
	}

	return result, nil
}

func (s *balance) GetTransactionByID(ctx context.Context, id int64) (*accountingPb.TransactionSelectionWithBlock, error) {
	t, err := s.balanceRepository.GetTransactionByID(ctx, int(id))
	if err != nil {
		return nil, fmt.Errorf("repository: %w", err)
	}

	result := accountingPb.TransactionSelectionWithBlock{
		ReceiverAccountID:    t.ReceiverAccountID,
		CoinID:               strconv.Itoa(int(t.CoinID)),
		Type:                 int64(t.Type),
		Amount:               t.Amount.String(),
		BlockID:              0, // TODO: why empty?
		UnblockToAccountID:   t.UnblockToAccountID,
		SenderAccountID:      t.SenderAccountID,
		UnblockTransactionID: t.UnblockTransactionID,
		ActionID:             t.ActionID,
	}

	return &result, nil
}

func (s *balance) FindLastBlockTimeBalances(ctx context.Context, userAccountIDs []int64) (
	[]*accountingPb.UserBlockTimeBalance, error) {
	balances, err := s.balanceRepository.FindLastBlockTimeBalances(ctx, userAccountIDs)
	if err != nil {
		return nil, fmt.Errorf("repository.FindLastBlockTimeBalances: %w", err)
	}

	result := make([]*accountingPb.UserBlockTimeBalance, 0, len(balances))
	for k, v := range balances {
		item := accountingPb.UserBlockTimeBalance{UserID: int64(k), Balance: v.String()}
		result = append(result, &item)
	}

	return result, nil
}

func (s *balance) FindBalancesDiffMining(ctx context.Context, data []*accountingPb.UserBeforePayoutMining) ([]*accountingPb.UserMiningDiff, error) {
	users := make([]*model.UserBeforePayoutMining, 0, len(data))
	for _, u := range data {
		coinID, err := strconv.Atoi(u.CoinID)
		if err != nil {
			return nil, fmt.Errorf("can't convert to int: %w", err)
		}

		users = append(users, &model.UserBeforePayoutMining{
			UserID:        u.UserID,
			CoinID:        int64(coinID),
			BlockID:       u.BlockID,
			AccountTypeID: enum.AccountTypeId(u.AccountTypeID),
			LastPay:       u.LastPay.AsTime(),
		})
	}

	diffs, err := s.balanceRepository.FindBalancesDiffMining(ctx, users)
	if err != nil {
		return nil, fmt.Errorf("repository.FindBalancesDiffMining: %w", err)
	}

	result := make([]*accountingPb.UserMiningDiff, 0, len(diffs))
	for k, v := range diffs {
		item := accountingPb.UserMiningDiff{BlockID: int64(k), Diff: v.String()}
		result = append(result, &item)
	}

	return result, nil
}

func (s *balance) FindBalancesDiffWallet(ctx context.Context, data []*accountingPb.UserBeforePayoutWallet) ([]*accountingPb.UserWalletDiff, error) {
	users := make([]*model.UserBeforePayoutWallet, 0, len(data))
	for _, u := range data {
		coinID, err := strconv.Atoi(u.CoinID)
		if err != nil {
			return nil, fmt.Errorf("can't convert to int: %w", err)
		}

		users = append(users, &model.UserBeforePayoutWallet{
			UserID:         u.UserID,
			CoinID:         int64(coinID),
			AccountTypeID:  enum.AccountTypeId(u.AccountTypeID),
			TransactionIDs: u.TransactionIDs,
		})
	}

	diffs, err := s.balanceRepository.FindBalancesDiffWallet(ctx, users)
	if err != nil {
		return nil, fmt.Errorf("repository.FindBalancesDiffWallet: %w", err)
	}

	result := make([]*accountingPb.UserWalletDiff, 0, len(diffs))
	for i := range diffs {
		item := accountingPb.UserWalletDiff{UserID: diffs[i].UserID, BlockID: diffs[i].BlockID, Diff: diffs[i].Diff.String()}
		result = append(result, &item)
	}

	return result, nil
}

func (s *balance) ChangeMultiple(ctx context.Context, protoTrs []*accountingPb.Transaction) error {
	trs := make([]*model.Transaction, len(protoTrs))
	var err error
	for i := range protoTrs {
		trs[i], err = ParseProtoTransaction(protoTrs[i])
		if err != nil {
			return fmt.Errorf("parse proto transaction: %w", err)
		}
	}
	err = s.balanceRepository.ChangeMultiple(ctx, trs)
	if err != nil {
		return fmt.Errorf("repository: ChangeMultiple: %w", err)
	}
	return nil
}

func (s *balance) GetBalances(ctx context.Context, userID int32) ([]*model.Balance, error) {
	return s.balanceRepository.GetBalances(ctx, userID, s.walletCoinsIDs, s.walletCoinsStrIDs)
}

func (s *balance) GetBalanceByCoin(ctx context.Context, userID int32, coin string) (*model.Balance, error) {
	coinID, err := s.getCoinIDByCoinName(coin)
	if err != nil {
		return nil, err
	}

	return s.balanceRepository.GetBalanceByCoin(ctx, userID, coinID)
}

func (s *balance) GetPaid(ctx context.Context, userID int32, coin string, from, to time.Time) (decimal.Decimal, error) {
	coinID, err := s.getCoinIDByCoinName(coin)
	if err != nil {
		return decimal.Zero, err
	}

	paid, err := s.balanceRepository.GetPaid(ctx, userID, coinID, from, to)
	if err != nil {
		return decimal.Zero, fmt.Errorf("getPaid: %w", err)
	}
	return paid, nil
}

func (s *balance) GetCoinsSummary(ctx context.Context, userID int32) ([]*model.CoinSummary, error) {
	return s.balanceRepository.GetCoinsSummary(ctx, userID, s.walletCoinsIDs, s.walletCoinsStrIDs)
}

func (s *balance) GetTransactionIDByAction(ctx context.Context, actionID string, txType int, amount string) (int64, error) {
	result, err := s.balanceRepository.GetTransactionIDByAction(ctx, actionID, txType, amount)
	if err != nil {
		return 0, fmt.Errorf("repository: GetTransactionIDByAction: %w", err)
	}

	return int64(result), nil
}

func (s *balance) changeForTransaction(ctx context.Context, sqlTx pgx.Tx, pbTransaction *accountingPb.Transaction) error {
	coinID, err := strconv.Atoi(pbTransaction.CoinID)
	if err != nil {
		return fmt.Errorf("can't convert to int: %w", err)
	}

	amount, err := decimal.NewFromString(pbTransaction.Amount)
	if err != nil {
		return fmt.Errorf("invalid amount value: %w", err)
	}

	t := model.Transaction{
		ID:                0,
		Type:              model.TransactionType(pbTransaction.Type),
		CreatedAt:         pbTransaction.CreatedAt.AsTime(),
		SenderAccountID:   pbTransaction.SenderAccountID,
		ReceiverAccountID: pbTransaction.ReceiverAccountID,
		CoinID:            int64(coinID),
		Amount:            amount,
		Comment:           pbTransaction.Comment,
		Hash:              pbTransaction.Hash,
		Hashrate:          pbTransaction.Hashrate,
		FromReferralId:    pbTransaction.FromReferralId,
		ReceiverAddress:   pbTransaction.ReceiverAddress,
		TokenID:           pbTransaction.TokenID,
		ActionID:          pbTransaction.ActionID,
		UnblockAccountId:  0,
		BlockedTill:       time.Time{},
	}

	action, ok := TransactionsActions[t.Type]
	if !ok {
		return fmt.Errorf("unexpected transaction type: %v", t.Type)
	}

	if t.Type > model.LastSupportedTransactionNumber {
		return fmt.Errorf("this transaction type not declared (last declared type %d): %v", model.LastSupportedTransactionNumber, t.Type)
	}

	change := func() error {
		_, err = s.balanceRepository.Change(ctx, sqlTx, &t)
		if err != nil {
			return fmt.Errorf("repo: Change: %w", err)
		}
		return nil
	}

	changeWithBlock := func() error {
		t.UnblockAccountId = t.ReceiverAccountID
		t.ReceiverAccountID, err = s.getBlockAccountID(ctx, sqlTx, &t)
		t.BlockedTill = GetBlockTillByType(t.Type)
		_, err = s.balanceRepository.ChangeWithBlock(ctx, sqlTx, &t)
		if err != nil {
			return fmt.Errorf("repo: ChangeWithBlock: %w", err)
		}
		return nil
	}

	changeWithUnblock := func() error {
		_, err = s.balanceRepository.ChangeWithUnblock(ctx, sqlTx, &t)
		if err != nil {
			return fmt.Errorf("repo: ChangeWithUnblock: %w", err)
		}
		return nil
	}

	return s.changeDispatcher(action, change, changeWithBlock, changeWithUnblock)
}

func (s *balance) getBlockAccountID(ctx context.Context, sqlTx pgx.Tx, t *model.Transaction) (int64, error) {
	coinID := int(t.CoinID)
	switch t.Type {
	case model.PoolPaysUsersBalanceTrTypeID, // 21
		model.WalletMiningTransferTrTypeID: // 31
		result, err := s.balanceRepository.GetBlockAccountIDFor31Type(ctx, sqlTx, coinID, s.serviceData.PayoutsNodeUsername)
		if err != nil {
			return 0, fmt.Errorf("repo: GetBlockAccountIDFor31Type: %w", err)
		}
		return int64(result), nil

	case model.ExchBlockTrTypeID: // 57
		result, err := s.balanceRepository.GetBlockAccountIDFor57Type(ctx, sqlTx, coinID, s.serviceData.ExchangeUsername)
		if err != nil {
			return 0, fmt.Errorf("repo: GetBlockAccountIDFor57Type: %w", err)
		}
		return int64(result), nil

	case
		model.HedgeBuyBlockTrTypeID,            // 45
		model.FiatWithdrawTrTypeID,             // 42
		model.P2PSellTrType,                    // 66
		model.P2PSellCommissionTrType,          // 69
		model.CnhldEarlyCloseTrTypeID,          // 36
		model.BlockTransferFromDepositToWallet: // 83
		result, err := s.balanceRepository.GetBlockAccountIDBySenderAccountID(ctx, sqlTx, int(t.SenderAccountID))
		if err != nil {
			return 0, fmt.Errorf("repo: GetBlockAccountIDBySenderAccountID: %w", err)
		}
		return int64(result), nil

	default:
		return 0, fmt.Errorf("unexpected transaction type: %v", t.Type)
	}
}

func (s *balance) changeDispatcher(
	action string,
	change, changeWithBlock, changeWithUnblock changeFunc,
) error {

	switch action {
	case ChangeBalance:
		return change()
	case ChangeBalanceWithBlock:
		return changeWithBlock()
	case ChangeBalanceWithUnblock:
		return changeWithUnblock()
	case DeprecatedAction:
		return fmt.Errorf("deprecated action: %s", action)

	}

	return fmt.Errorf("unexpected action: %s", action)
}

func (s *balance) getCoinIDByCoinName(coin string) (int, error) {
	var coinID int
	for id, name := range s.walletCoinsStrIDs {
		if name == coin {
			coinID = id
			break
		}
	}
	if coinID == 0 {
		return 0, fmt.Errorf("unexpected coint: %s", coin)
	}
	return coinID, nil
}

func (s *balance) FindOperationsAndTransactions(ctx context.Context, request *model.OperationWithTransactionQuery) ([]*model.OperationWithTransaction, int64, error) {

	operations, totalCount, err := s.balanceRepository.FindOperationsAndTransactions(ctx, request)

	if err != nil {
		return nil, 0, fmt.Errorf("repository: %w", err)
	}

	return operations, totalCount, nil
}

func (s *balance) GetBalanceBeforeTransaction(ctx context.Context, accountID, transactionID int64) (decimal.Decimal, error) {
	amount, err := s.balanceRepository.GetBalanceBeforeTransaction(ctx, accountID, transactionID)
	if err != nil {
		return amount, fmt.Errorf("repository: %w", err)
	}

	return amount, nil
}
