package handler

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	businessErr "code.emcdtech.com/emcd/sdk/error"
	sdkLog "code.emcdtech.com/emcd/sdk/log"
	coinValidatorRepo "code.emcdtech.com/emcd/service/coin/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"google.golang.org/protobuf/types/known/timestamppb"

	"code.emcdtech.com/emcd/service/accounting/internal/handler/mapping"
	"code.emcdtech.com/emcd/service/accounting/internal/service"
	"code.emcdtech.com/emcd/service/accounting/model"
	"code.emcdtech.com/emcd/service/accounting/model/enum"
	accountingPb "code.emcdtech.com/emcd/service/accounting/protocol/accounting"
)

const PGRaiseException = "P0001"

var (
	// added in list
	undefinedError = businessErr.NewError("acc-1000", "undefined")
	// added in list
	emptyProcessingType = businessErr.NewError("acc-1001", "processing type is empty")
	// added in list
	emptyActionID = businessErr.NewError("acc-1002", "actionId is not defined")
	// added in list
	coinIdCheckFail = businessErr.NewError("acc-1003", "coin check is fail")
	// added in list
	noBalance = businessErr.NewError("acc-1004", "balance is less than the transfer amount")
	// added in list
	doubleBalanceBlock = businessErr.NewError("acc-1005", "balance block is already exists")
	// added in list
	doubleBalanceUnblock = businessErr.NewError("acc-1006", "balance unblock is already exists")
	// added in list
	notFoundBalanceBlock = businessErr.NewError("acc-1007", "balance block is not found")
	// added in list
	passedBlockingTime = businessErr.NewError("acc-1008", "blocking time has passed")
)

type AccountingHandler struct {
	balanceService     service.Balance
	walletsService     service.WalletsHistory
	incomesService     service.IncomesHistory
	payoutsService     service.PayoutsHistory
	limitPayoutService service.LimitPayouts
	userAccountService service.AccountingUserAccount
	transactionService service.Transaction
	payoutsFindService service.Payouts
	checkerService     service.Checker
	coinValidator      coinValidatorRepo.CoinValidatorRepository
	accountingPb.UnimplementedAccountingServiceServer
}

func NewAccountingHandler(
	balanceService service.Balance,
	walletsService service.WalletsHistory,
	incomesService service.IncomesHistory,
	payoutsService service.PayoutsHistory,
	limitPayoutService service.LimitPayouts,
	userAccountService service.AccountingUserAccount,
	transactionService service.Transaction,
	payoutsFindService service.Payouts,
	checkerService service.Checker,
	coinValidator coinValidatorRepo.CoinValidatorRepository,
) *AccountingHandler {
	return &AccountingHandler{
		balanceService:                       balanceService,
		walletsService:                       walletsService,
		incomesService:                       incomesService,
		payoutsService:                       payoutsService,
		limitPayoutService:                   limitPayoutService,
		payoutsFindService:                   payoutsFindService,
		userAccountService:                   userAccountService,
		transactionService:                   transactionService,
		checkerService:                       checkerService,
		coinValidator:                        coinValidator,
		UnimplementedAccountingServiceServer: accountingPb.UnimplementedAccountingServiceServer{},
	}
}

func (h *AccountingHandler) ViewBalance(ctx context.Context, req *accountingPb.ViewBalanceRequest) (*accountingPb.ViewBalanceResponse, error) {
	balance, err := h.balanceService.View(ctx, req.UserID, enum.AccountTypeId(req.AccountTypeID), req.CoinID, req.TotalBalance)
	if err != nil {
		sdkLog.Error(ctx, "service: %v", err)
		return nil, err
	}

	return &accountingPb.ViewBalanceResponse{Balance: balance.String()}, nil
}

func (h *AccountingHandler) ChangeBalance(ctx context.Context, req *accountingPb.ChangeBalanceRequest) (*accountingPb.ChangeBalanceResponse, error) {
	err := h.balanceService.Change(ctx, req.Transactions)
	if err != nil {
		sdkLog.Error(ctx, "service: %v", err)

		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == PGRaiseException {
			sdkLog.Error(ctx, "detected business error: %v", e)
			return nil, parseSqlErrorMessage(e.Message)
		}
		return nil, err
	}

	return &accountingPb.ChangeBalanceResponse{}, nil // TODO: empty or nil return only?
}

func (h *AccountingHandler) FindOperations(ctx context.Context, req *accountingPb.FindOperationsRequest) (*accountingPb.FindOperationsResponse, error) {
	operations, err := h.balanceService.FindOperations(ctx, req.UserID, req.CoinID)
	if err != nil {
		sdkLog.Error(ctx, "service: %v", err)
		return nil, err
	}

	return &accountingPb.FindOperationsResponse{Operations: operations}, nil
}

func (h *AccountingHandler) FindBatchOperations(ctx context.Context, req *accountingPb.FindBatchOperationsRequest) (*accountingPb.FindBatchOperationsResponse, error) {
	operationsByUsers, err := h.balanceService.FindBatchOperations(ctx, req.Users)
	if err != nil {
		sdkLog.Error(ctx, "service: %v", err)
		return nil, err
	}

	return &accountingPb.FindBatchOperationsResponse{OperationsByUsers: operationsByUsers}, nil
}

func (h *AccountingHandler) FindTransactions(ctx context.Context, req *accountingPb.FindTransactionsRequest) (*accountingPb.FindTransactionsResponse, error) {
	transactions, err := h.balanceService.FindTransactions(ctx, req.Types, req.UserID, req.AccountTypeID, req.CoinIDs, req.From)
	if err != nil {
		sdkLog.Error(ctx, "service: %v", err)
		return nil, err
	}

	return &accountingPb.FindTransactionsResponse{Transactions: transactions}, nil
}

func (h *AccountingHandler) FindTransactionsByCollectorFilter(ctx context.Context, req *accountingPb.TransactionCollectorFilter) (*accountingPb.TransactionMultiResponse, error) {
	if filterModel, err := mapping.MapProtoToModelTransactionCollectorFilter(h.coinValidator, req); err != nil {
		sdkLog.Error(ctx, "mapping: %v", err)

		return nil, err
	} else if totalCount, transactions, err := h.balanceService.FindTransactionsByCollectorFilter(ctx, filterModel); err != nil {
		sdkLog.Error(ctx, "service: %v", err)

		return nil, err
	} else {

		return &accountingPb.TransactionMultiResponse{
			Transactions: transactions,
			TotalCount:   totalCount,
		}, nil
	}
}

// FindTransactionsByCollectorFilter(TransactionCollectorFilter) returns (TransactionMultiResponse)

func (h *AccountingHandler) GetTransactionsByActionID(ctx context.Context, req *accountingPb.GetTransactionsByActionIDRequest) (*accountingPb.GetTransactionsByActionIDResponse, error) {
	transactions, err := h.balanceService.GetTransactionsByActionID(ctx, req.ActionID)
	if err != nil {
		sdkLog.Error(ctx, "service: %v", err)
		return nil, err
	}

	return &accountingPb.GetTransactionsByActionIDResponse{Transactions: transactions}, nil
}

func (h *AccountingHandler) FindTransactionsWithBlocks(ctx context.Context, req *accountingPb.FindTransactionsWithBlocksRequest) (*accountingPb.FindTransactionsWithBlocksResponse, error) {
	transactions, err := h.balanceService.FindTransactionsWithBlocks(ctx, req.BlockedTill)
	if err != nil {
		sdkLog.Error(ctx, "service: %v", err)
		return nil, err
	}

	return &accountingPb.FindTransactionsWithBlocksResponse{Transactions: transactions}, nil
}

func (h *AccountingHandler) GetTransactionByID(ctx context.Context, req *accountingPb.GetTransactionByIDRequest) (*accountingPb.GetTransactionByIDResponse, error) {
	transaction, err := h.balanceService.GetTransactionByID(ctx, req.Id)
	if err != nil {
		sdkLog.Error(ctx, "service: %v", err)
		return nil, err
	}

	return &accountingPb.GetTransactionByIDResponse{Transaction: transaction}, nil
}

func (h *AccountingHandler) FindLastBlockTimeBalances(ctx context.Context, req *accountingPb.FindLastBlockTimeBalancesRequest) (*accountingPb.FindLastBlockTimeBalancesResponse, error) {
	balances, err := h.balanceService.FindLastBlockTimeBalances(ctx, req.UserAccountIDs)
	if err != nil {
		sdkLog.Error(ctx, "service: %v", err)
		return nil, err
	}

	return &accountingPb.FindLastBlockTimeBalancesResponse{Balances: balances}, nil
}

func (h *AccountingHandler) FindBalancesDiffMining(ctx context.Context, req *accountingPb.FindBalancesDiffMiningRequest) (*accountingPb.FindBalancesDiffMiningResponse, error) {
	diffs, err := h.balanceService.FindBalancesDiffMining(ctx, req.Users)
	if err != nil {
		sdkLog.Error(ctx, "service: %v", err)
		return nil, err
	}

	return &accountingPb.FindBalancesDiffMiningResponse{Diffs: diffs}, nil
}

func (h *AccountingHandler) FindBalancesDiffWallet(ctx context.Context, req *accountingPb.FindBalancesDiffWalletRequest) (*accountingPb.FindBalancesDiffWalletResponse, error) {
	diffs, err := h.balanceService.FindBalancesDiffWallet(ctx, req.Users)
	if err != nil {
		sdkLog.Error(ctx, "service: %v", err)
		return nil, err
	}

	return &accountingPb.FindBalancesDiffWalletResponse{Diffs: diffs}, nil
}

func (h *AccountingHandler) GetHistory(ctx context.Context, req *accountingPb.GetHistoryRequest) (*accountingPb.GetHistoryResponse, error) {
	var history *accountingPb.GetHistoryResponse
	var err error

	switch model.HistoryType(req.Type) {
	case model.HistoryIncome:
		if history, err = h.incomesService.GetHistory(ctx, req); err != nil {
			sdkLog.Error(ctx, "incomesService: %v", err)
			return nil, fmt.Errorf("incomesService: %w", err)
		}
	case model.HistoryPayout:
		if history, err = h.payoutsService.GetHistory(ctx, req); err != nil {
			sdkLog.Error(ctx, "payoutsService: %v", err)
			return nil, fmt.Errorf("payoutsService: %w", err)
		}
	case model.HistoryWallet, model.HistoryCoinhold:
		if history, err = h.walletsService.GetHistory(ctx, req); err != nil {
			sdkLog.Error(ctx, "walletsService: %v", err)
			return nil, fmt.Errorf("walletsService: %w", err)
		}
	default:
		msg := fmt.Sprintf("unknown history type: %s", model.HistoryType(req.Type))
		sdkLog.Error(ctx, msg)
		return nil, errors.New(msg)
	}

	return history, nil
}

func (h *AccountingHandler) ChangeMultipleBalance(ctx context.Context, req *accountingPb.ChangeMultipleBalanceRequest) (*accountingPb.ChangeMultipleBalanceResponse, error) {
	err := h.balanceService.ChangeMultiple(ctx, req.Transactions)
	if err != nil {
		sdkLog.Error(ctx, "service: ChangeMultiple: %v", err)
	}
	return &accountingPb.ChangeMultipleBalanceResponse{}, nil
}

func (h *AccountingHandler) GetBalances(ctx context.Context, req *accountingPb.UserIDRequest) (*accountingPb.GetBalancesResponse, error) {
	balances, err := h.balanceService.GetBalances(ctx, req.GetUserId())
	if err != nil {
		sdkLog.Error(ctx, "GetBalances: %v", err)
		return nil, err
	}
	res := &accountingPb.GetBalancesResponse{
		CoinBalance: make([]*accountingPb.CoinBalance, len(balances)),
	}
	for i := range res.CoinBalance {
		b := balances[i]
		res.CoinBalance[i] = &accountingPb.CoinBalance{
			CoinId:                      b.CoinID,
			WalletBalance:               b.WalletBalance.String(),
			MiningBalance:               b.MiningBalance.String(),
			CoinholdsBalance:            b.CoinholdsBalance.String(),
			P2PBalance:                  b.P2pBalance.String(),
			BlockedBalanceCoinhold:      b.BlockedBalanceCoinhold.String(),
			BlockedBalanceFreeWithdraw:  b.BlockedBalanceFreeWithdraw.String(),
			BlockedBalanceP2P:           b.BlockedBalanceP2p.String(),
			BlockedBalanceMinimgPayouts: b.BlockedBalanceMiningPayouts.String(),
		}
	}
	return res, nil
}

func (h *AccountingHandler) GetBalanceByCoin(ctx context.Context, req *accountingPb.GetBalanceByCoinRequest) (*accountingPb.GetBalanceByCoinResponse, error) {
	balances, err := h.balanceService.GetBalanceByCoin(ctx, req.GetUserId(), req.GetCoin())
	if err != nil {
		sdkLog.Error(ctx, "GetBalanceByCoin: %v", err)
		return nil, err
	}
	return &accountingPb.GetBalanceByCoinResponse{
		CoinBalance: &accountingPb.CoinBalance{
			CoinId:                      balances.CoinID,
			WalletBalance:               balances.WalletBalance.String(),
			MiningBalance:               balances.MiningBalance.String(),
			CoinholdsBalance:            balances.CoinholdsBalance.String(),
			P2PBalance:                  balances.P2pBalance.String(),
			BlockedBalanceCoinhold:      balances.BlockedBalanceCoinhold.String(),
			BlockedBalanceFreeWithdraw:  balances.BlockedBalanceFreeWithdraw.String(),
			BlockedBalanceP2P:           balances.BlockedBalanceP2p.String(),
			BlockedBalanceMinimgPayouts: "",
		},
	}, nil
}

func (h *AccountingHandler) GetPaid(ctx context.Context, req *accountingPb.GetPaidRequest) (*accountingPb.GetPaidResponse, error) {
	paid, err := h.balanceService.GetPaid(ctx, req.GetUserId(), req.GetCoin(), req.GetFrom().AsTime(), req.GetTo().AsTime())
	if err != nil {
		sdkLog.Error(ctx, "GetPaid: %v", err)
		return nil, err
	}
	return &accountingPb.GetPaidResponse{
		Paid: paid.String(),
	}, nil
}

func (h *AccountingHandler) GetCoinsSummary(ctx context.Context, req *accountingPb.UserIDRequest) (*accountingPb.GetCoinsSummaryResponse, error) {
	coinsSummary, err := h.balanceService.GetCoinsSummary(ctx, req.GetUserId())
	if err != nil {
		sdkLog.Error(ctx, "GetCoinsSummary: %v", err)
		return nil, err
	}
	res := &accountingPb.GetCoinsSummaryResponse{
		CoinSummary: make([]*accountingPb.CoinSummary, len(coinsSummary)),
	}
	for i := range res.CoinSummary {
		b := coinsSummary[i]
		res.CoinSummary[i] = &accountingPb.CoinSummary{
			CoinId:      b.CoinID,
			TotalAmount: b.TotalAmount.String(),
		}
	}

	return res, nil
}

func (h *AccountingHandler) GetTransactionIDByAction(ctx context.Context, req *accountingPb.GetTransactionIDByActionRequest) (*accountingPb.GetTransactionIDByActionResponse, error) {
	txID, err := h.balanceService.GetTransactionIDByAction(ctx, req.GetActionID(), int(req.GetType()), req.GetAmount())
	if err != nil {
		sdkLog.Error(ctx, "GetTransactionIDByAction: %v", err)
		return nil, err
	}

	return &accountingPb.GetTransactionIDByActionResponse{Id: txID}, nil
}

func (h *AccountingHandler) CheckPayoutsLimit(ctx context.Context, req *accountingPb.CheckPayoutsLimitRequest) (*accountingPb.CheckPayoutsLimitResponse, error) {

	err := h.limitPayoutService.CheckLimit(ctx, req.GetUserID(), req.GetCoinID(), float64(req.GetAmount()))
	if err != nil {
		sdkLog.Error(ctx, "check payouts limit: %v", err)
		return nil, err
	}
	return &accountingPb.CheckPayoutsLimitResponse{}, nil
}

func (h *AccountingHandler) GetPayoutsBlockStatus(ctx context.Context, req *accountingPb.GetPayoutsBlockStatusRequest) (*accountingPb.GetPayoutsBlockStatusResponse, error) {

	result, err := h.limitPayoutService.GetBlockStatus(ctx, req.GetUserID())
	if err != nil {
		sdkLog.Error(ctx, "get block status: %v", err)
		return &accountingPb.GetPayoutsBlockStatusResponse{Status: model.StatusUndefined, Message: err.Error()}, nil
	}

	var message string

	switch result {
	case model.StatusBlocked:
		message = "blocked"
	case model.StatusUnblocked:
		message = "unblocked"
	default:
		message = "not blocked"
	}

	return &accountingPb.GetPayoutsBlockStatusResponse{Status: int32(result), Message: message}, nil
}

func (h *AccountingHandler) SetPayoutsBlockStatus(ctx context.Context, req *accountingPb.SetPayoutsBlockStatusRequest) (*accountingPb.SetPayoutsBlockStatusResponse, error) {

	status := req.GetStatus()

	if status != model.StatusBlocked && status != model.StatusUnblocked {
		return &accountingPb.SetPayoutsBlockStatusResponse{Success: false, Message: "status code is not valid"}, nil
	}

	err := h.limitPayoutService.SetBlockStatus(ctx, req.GetUserID(), int(status))
	if err != nil {
		sdkLog.Error(ctx, "set block status: %v", err)
		return &accountingPb.SetPayoutsBlockStatusResponse{Success: false, Message: err.Error()}, nil
	}

	return &accountingPb.SetPayoutsBlockStatusResponse{Success: true, Message: ""}, nil
}

func parseSqlErrorMessage(err string) error {
	switch {
	case strings.Contains(err, "processing type is empty"):
		return emptyProcessingType
	case strings.Contains(err, "action_id is empty"):
		return emptyActionID
	case strings.Contains(err, "coin check is fail"):
		return coinIdCheckFail
	case strings.Contains(err, "balance is less than the transfer amount"):
		return noBalance
	case strings.Contains(err, "balance block is already exists"):
		return doubleBalanceBlock
	case strings.Contains(err, "balance unblock is already exists"):
		return doubleBalanceUnblock
	case strings.Contains(err, "balance block is not found"):
		return notFoundBalanceBlock
	case strings.Contains(err, "blocking time has passed"):
		return passedBlockingTime
	default:
		return undefinedError
	}
}

func (h *AccountingHandler) FindOperationsAndTransactions(ctx context.Context, req *accountingPb.FindOperationsAndTransactionsRequest) (*accountingPb.FindOperationsAndTransactionsResponse, error) {

	var amountFloat float64
	var actionID *uuid.UUID

	var err error
	if req.GetAmount() != "" {
		amountFloat, err = strconv.ParseFloat(req.GetAmount(), 64)
		if err != nil {
			sdkLog.Error(ctx, "Accounting.FindOperationsAndTransactions.ParseAmount: %v", err)
			return nil, err
		}
	}

	if req.GetActionID() != "" {
		actionIDP, err := uuid.Parse(*req.ActionID)
		if err != nil {
			sdkLog.Error(ctx, "Accounting.FindOperationsAndTransactions.ParseActionID: %v", err)
			return nil, err
		}
		actionID = &actionIDP
	}

	request := model.OperationWithTransactionQuery{
		UserID:               req.GetUserID(),
		CoinID:               req.GetCoinID(),
		TokenID:              req.GetTokenID(),
		ActionID:             actionID,
		AccountType:          req.GetAccountType(),
		OperationTypes:       req.GetOperationTypes(),
		DateFrom:             req.GetDateFrom(),
		DateTo:               req.GetDateTo(),
		Amount:               amountFloat,
		Hash:                 req.GetHash(),
		ReceiverAccountID:    req.GetReceiverAccountID(),
		ReceiverAddress:      req.GetReceiverAddress(),
		SenderAccountID:      req.GetSenderAccountID(),
		TransactionBlockID:   req.GetTransactionBlockID(),
		UnblockToAccountId:   req.GetUnblockToAccountId(),
		UnblockTransactionId: req.GetUnblockTransactionId(),
		FromReferralId:       req.GetFromReferralId(),
		Limit:                req.GetLimit(),
		Offset:               req.GetOffset(),
		SortField:            req.GetSort().GetField(),
		Asc:                  req.GetSort().GetAsc(),
	}

	operations, totalCount, err := h.balanceService.FindOperationsAndTransactions(ctx, &request)

	if err != nil {
		sdkLog.Error(ctx, "FindOperationsAndTransactions: %v", err)
		return nil, err
	}

	res := &accountingPb.FindOperationsAndTransactionsResponse{
		TotalCount: totalCount,
		Operations: make([]*accountingPb.OperationWithTransaction, len(operations)),
	}

	for i := range res.Operations {
		o := operations[i]

		res.Operations[i] = &accountingPb.OperationWithTransaction{
			Id:                   o.Id,
			AccountID:            o.AccountID,
			CoinID:               o.CoinID,
			TokenID:              o.TokenID,
			Amount:               o.Amount.String(),
			Type:                 int64(o.Type),
			TransactionID:        o.TransactionID,
			ActionID:             o.ActionID,
			Comment:              o.Comment,
			Fee:                  o.Fee.String(),
			FromReferralId:       o.FromReferralId,
			GasPrice:             o.GasPrice.String(),
			Hash:                 o.Hash,
			Hashrate:             o.Hashrate,
			ReceiverAccountID:    o.ReceiverAccountID,
			ReceiverAddress:      o.ReceiverAddress,
			SenderAccountID:      o.SenderAccountID,
			TransactionBlockID:   o.TransactionBlockID,
			BlockedTill:          timestamppb.New(o.BlockedTill),
			UnblockToAccountId:   o.UnblockToAccountId,
			UnblockTransactionId: o.UnblockTransactionId,
			CreatedAt:            timestamppb.New(o.CreatedAt),
		}
	}

	return res, nil
}

func (h *AccountingHandler) FindPayoutsForBlock(ctx context.Context, req *accountingPb.FindPayoutsForBlockRequest) (*accountingPb.FindPayoutsForBlockResponse, error) {

	payouts, err := h.payoutsFindService.FindPayoutsForBlock(ctx, req.CoinID, req.MinPay, req.Timestamp.AsTime())

	if err != nil {
		sdkLog.Error(ctx, "FindPayoutsForBlock: %v", err)
		return nil, err
	}

	response := &accountingPb.FindPayoutsForBlockResponse{
		Payouts: make([]*accountingPb.PayoutForBlock, len(payouts)),
	}

	for i := range response.Payouts {
		payout := payouts[i]

		response.Payouts[i] = &accountingPb.PayoutForBlock{
			UserID:    payout.UserID,
			AccountID: payout.AccountID,
			Address:   payout.Address,
			Balance:   payout.Balance.String(),
		}
	}

	return response, nil
}

func (h *AccountingHandler) GetCurrentPayoutsBlock(ctx context.Context, req *accountingPb.GetCurrentPayoutsBlockRequest) (*accountingPb.GetCurrentPayoutsBlockResponse, error) {

	transactions, err := h.payoutsFindService.GetCurrentPayoutsBlock(ctx, req.CoinID, req.Username, req.IsService)

	if err != nil {
		sdkLog.Error(ctx, "GetCurrentPayoutsBlock: %v", err)
		return nil, err
	}

	response := &accountingPb.GetCurrentPayoutsBlockResponse{
		Transactions: make([]*accountingPb.PayoutBlockTransaction, len(transactions)),
	}

	for i := range response.Transactions {
		transaction := transactions[i]

		response.Transactions[i] = &accountingPb.PayoutBlockTransaction{
			ID:      transaction.ID,
			Balance: transaction.Balance.String(),
		}
	}

	return response, nil
}

func (h *AccountingHandler) GetFreePayouts(ctx context.Context, req *accountingPb.GetFreePayoutsRequest) (*accountingPb.GetFreePayoutsResponse, error) {
	transactions, err := h.payoutsFindService.GetFreePayouts(ctx, req.CoinID)

	if err != nil {
		sdkLog.Error(ctx, "GetFreePayouts: %v", err)
		return nil, err
	}

	response := &accountingPb.GetFreePayoutsResponse{
		Payouts: make([]*accountingPb.FreePayout, len(transactions)),
	}

	for i := range response.Payouts {
		transaction := transactions[i]

		response.Payouts[i] = &accountingPb.FreePayout{
			AccountId:         transaction.AccountId,
			UserId:            transaction.UserId,
			Username:          transaction.Username,
			ID:                transaction.ID,
			ActionID:          transaction.ActionID,
			Amount:            transaction.Amount.String(),
			CoinID:            transaction.CoinID,
			Comment:           transaction.Comment,
			CreatedAt:         timestamppb.New(transaction.CreatedAt),
			Fee:               transaction.Fee.String(),
			FromReferralID:    transaction.FromReferralID,
			GasPrice:          transaction.GasPrice.String(),
			Hash:              transaction.Hash,
			Hashrate:          transaction.Hashrate,
			IsViewer:          transaction.IsViewer,
			ReceiverAccountID: transaction.ReceiverAccountID,
			ReceiverAddress:   transaction.ReceiverAddress,
			SenderAccountID:   transaction.SenderAccountID,
			TokenID:           transaction.TokenID,
			Type:              transaction.Type,
		}
	}

	return response, nil
}

func (h *AccountingHandler) GetCurrentPayoutsList(ctx context.Context, req *accountingPb.GetCurrentPayoutsListRequest) (*accountingPb.GetCurrentPayoutsListResponse, error) {
	transactions, err := h.payoutsFindService.GetCurrentPayoutsList(ctx, req.CoinId, req.PaymentTransactionType)

	if err != nil {
		sdkLog.Error(ctx, "GetFreePayouts: %v", err)
		return nil, err
	}

	response := &accountingPb.GetCurrentPayoutsListResponse{
		Payouts: make([]*accountingPb.CurrentPayout, len(transactions)),
	}

	for i := range response.Payouts {
		transaction := transactions[i]
		calc := &accountingPb.PayoutCalculationData{
			Coinhold:    transaction.Calc.Coinhold.String(),
			Incomes:     transaction.Calc.Incomes.String(),
			Hashrate:    transaction.Calc.Hashrate.String(),
			FeeAndMore:  transaction.Calc.FeeAndMore.String(),
			Ref:         transaction.Calc.Ref.String(),
			Other:       transaction.Calc.Other.String(),
			Types:       transaction.Calc.Types,
			AccountId:   transaction.Calc.AccountID,
			LastPay:     timestamppb.New(transaction.Calc.LastPay),
			IncomeFirst: timestamppb.New(transaction.Calc.IncomeFirst),
			IncomeLast:  timestamppb.New(transaction.Calc.IncomeLast),
		}

		response.Payouts[i] = &accountingPb.CurrentPayout{
			Id:          transaction.ID,
			AccountID2:  transaction.AccountID2,
			Username:    transaction.Username,
			UserID:      transaction.UserID,
			RefID:       transaction.RefID,
			CoinID:      transaction.CoinID,
			Minpay:      transaction.Minpay.String(),
			BlockCreate: timestamppb.New(transaction.BlockCreate),
			MasterID:    transaction.MasterID,
			Address:     transaction.Address,
			Balance:     transaction.Balance.String(),
			BlockID:     transaction.BlockID,
			Calc:        calc,
		}
	}

	return response, nil
}

func (h *AccountingHandler) GetCurrentReferralsPayoutsList(ctx context.Context, req *accountingPb.GetCurrentReferralsPayoutsListRequest) (*accountingPb.GetCurrentReferralsPayoutsListResponse, error) {
	transactions, err := h.payoutsFindService.GetCurrentReferralsPayoutsList(ctx, req.CoinId, req.PaymentTransactionType, req.ReferralId)

	if err != nil {
		sdkLog.Error(ctx, "GetFreePayouts: %v", err)
		return nil, err
	}

	response := &accountingPb.GetCurrentReferralsPayoutsListResponse{
		Payouts: make([]*accountingPb.CurrentReferralPayout, len(transactions)),
	}

	for i := range response.Payouts {
		transaction := transactions[i]
		calc := &accountingPb.PayoutCalculationData{
			Coinhold:    transaction.Calc.Coinhold.String(),
			Incomes:     transaction.Calc.Incomes.String(),
			Hashrate:    transaction.Calc.Hashrate.String(),
			FeeAndMore:  transaction.Calc.FeeAndMore.String(),
			Ref:         transaction.Calc.Ref.String(),
			Other:       transaction.Calc.Other.String(),
			Types:       transaction.Calc.Types,
			AccountId:   transaction.Calc.AccountID,
			LastPay:     timestamppb.New(transaction.Calc.LastPay),
			IncomeFirst: timestamppb.New(transaction.Calc.IncomeFirst),
			IncomeLast:  timestamppb.New(transaction.Calc.IncomeLast),
		}

		response.Payouts[i] = &accountingPb.CurrentReferralPayout{
			Id:         transaction.ID,
			AccountID2: transaction.AccountID2,
			Username:   transaction.Username,
			UserID:     transaction.UserID,
			RefID:      transaction.RefID,
			CoinID:     transaction.CoinID,
			Minpay:     transaction.Minpay.String(),
			MasterID:   transaction.MasterID,
			Address:    transaction.Address,
			Balance:    transaction.Balance.String(),
			Calc:       calc,
		}
	}

	return response, nil
}

func (h *AccountingHandler) CheckFreePayoutTransaction(ctx context.Context, req *accountingPb.CheckFreePayoutTransactionRequest) (*accountingPb.CheckFreePayoutTransactionResponse, error) {
	balance, err := h.payoutsFindService.CheckFreePayoutTransaction(ctx, req.AccountID, req.TransactionID)
	if err != nil {
		sdkLog.Error(ctx, "service: %v", err)
		return nil, err
	}

	return &accountingPb.CheckFreePayoutTransactionResponse{Sum: balance.String()}, nil
}

func (h *AccountingHandler) CheckPayoutBlockStatus(ctx context.Context, req *accountingPb.CheckPayoutBlockStatusRequest) (*accountingPb.CheckPayoutBlockStatusResponse, error) {
	transactions, err := h.payoutsFindService.CheckPayoutBlockStatus(ctx, req.BlockTransactionIds)

	if err != nil {
		sdkLog.Error(ctx, "CheckPayoutBlockStatus: %v", err)
		return nil, err
	}

	response := &accountingPb.CheckPayoutBlockStatusResponse{
		PayoutBlocks: make([]*accountingPb.PayoutBlock, len(transactions)),
	}

	for i := range response.PayoutBlocks {
		transaction := transactions[i]

		response.PayoutBlocks[i] = &accountingPb.PayoutBlock{
			ToAccountId:     transaction.ToAccountId,
			Type:            transaction.Type,
			ReceiverAddress: transaction.ReceiverAddress,
			UbTrId:          transaction.UnblockTransactionId,
			Amount:          transaction.Amount.String(),
		}
	}

	return response, nil
}

func (h *AccountingHandler) CheckIncomeOperations(ctx context.Context, req *accountingPb.CheckIncomeOperationsRequest) (*accountingPb.CheckIncomeOperationsResponse, error) {

	query := model.CheckIncomeOperationsQuery{
		CreatedAt: req.CreatedAt.AsTime(),
		Coin:      req.Coin,
		UserID:    req.UserID,
		AccountID: req.AccountID,
		LastPayAt: req.LastPayAt.AsTime(),
	}

	transactions, err := h.payoutsFindService.CheckIncomeOperations(ctx, query)

	if err != nil {
		sdkLog.Error(ctx, "CheckIncomeOperations: %v", err)
		return nil, err
	}

	response := &accountingPb.CheckIncomeOperationsResponse{
		Incomes: make([]*accountingPb.IncomeWithFee, len(transactions)),
	}

	for i := range response.Incomes {
		transaction := transactions[i]

		response.Incomes[i] = &accountingPb.IncomeWithFee{
			TransactionId: transaction.TransactionId,
			Hashrate:      transaction.Hashrate,
			CreatedAt:     timestamppb.New(transaction.CreatedAt),
			Fee:           transaction.Fee.String(),
			Amount:        transaction.Amount.String(),
		}
	}

	return response, nil
}

func (h *AccountingHandler) GetAveragePaid(ctx context.Context, req *accountingPb.GetAveragePaidRequest) (*accountingPb.GetAveragePaidResponse, error) {
	query := model.AveragePaidQuery{
		CoinID:            req.CoinID,
		Days:              req.Days,
		TransactionTypeID: req.TransactionTypeID,
		AccountTypeID:     req.AccountTypeID,
		Username:          req.Username,
	}

	averagePaid, err := h.payoutsFindService.GetAveragePaid(ctx, query)
	if err != nil {
		sdkLog.Error(ctx, "service: %v", err)
		return nil, err
	}

	return &accountingPb.GetAveragePaidResponse{Avg: averagePaid.String()}, nil
}

func (h *AccountingHandler) CheckOthers(ctx context.Context, req *accountingPb.CheckOthersRequest) (*accountingPb.CheckOthersResponse, error) {

	query := model.CheckOtherQuery{
		AccountID:      req.AccountID,
		Types:          req.Types,
		LastPayAt:      req.LastPayAt.AsTime(),
		BlockCreatedAt: nil,
	}

	blockCreatedAt := req.BlockCreatedAt.AsTime()

	if req.BlockCreatedAt != nil {
		query.BlockCreatedAt = &blockCreatedAt
	}

	transactions, err := h.payoutsFindService.CheckOthers(ctx, query)

	if err != nil {
		sdkLog.Error(ctx, "service: %v", err)
		return nil, err
	}

	response := &accountingPb.CheckOthersResponse{
		Others: make([]*accountingPb.OtherOperationsWithTransaction, len(transactions)),
	}

	for i := range response.Others {
		transaction := transactions[i]

		response.Others[i] = &accountingPb.OtherOperationsWithTransaction{
			TransactionID:  transaction.TransactionID,
			SenderID:       transaction.SenderID,
			ReceiverID:     transaction.ReceiverID,
			Hash:           transaction.Hash,
			OperationID:    transaction.OperationID,
			SenderUserID:   transaction.SenderUserID,
			ReceiverUserID: transaction.ReceiverUserID,
			Amount:         transaction.Amount.String(),
			Type:           transaction.TransactionTypeID,
			CreatedAt:      timestamppb.New(transaction.CreatedAt),
			Comment:        transaction.Comment,
		}
	}

	return response, nil
}

func (h *AccountingHandler) GetBalanceBeforeTransaction(ctx context.Context, req *accountingPb.GetBalanceBeforeTransactionRequest) (*accountingPb.GetBalanceBeforeTransactionResponse, error) {
	amount, err := h.balanceService.GetBalanceBeforeTransaction(ctx, req.AccountID, req.TransactionID)

	if err != nil {
		sdkLog.Error(ctx, "service: %v", err)
		return nil, err
	}

	return &accountingPb.GetBalanceBeforeTransactionResponse{Sum: amount.String()}, nil
}

func (h *AccountingHandler) GetServiceUserData(ctx context.Context, req *accountingPb.GetServiceUserDataRequest) (*accountingPb.GetServiceUserDataResponse, error) {

	blocks, err := h.payoutsFindService.GetServiceUserData(ctx, req.CoinID, req.Username, req.Limit)

	if err != nil {
		sdkLog.Error(ctx, "service: %v", err)
		return nil, err
	}

	response := &accountingPb.GetServiceUserDataResponse{
		Blocks: make([]*accountingPb.ServiceUserBlock, len(blocks)),
	}

	for i := range response.Blocks {
		transaction := blocks[i]

		response.Blocks[i] = &accountingPb.ServiceUserBlock{
			Amount:      transaction.Amount.String(),
			Address:     transaction.Address,
			SuAccountID: transaction.SuAccountID,
			BlockID:     transaction.BlockID,
			UserID:      transaction.UserID,
			Username:    transaction.Username,
		}
	}

	return response, nil
}

func (h *AccountingHandler) GetIncomesHashrateByDate(ctx context.Context, req *accountingPb.GetIncomesHashrateRequest) (*accountingPb.GetIncomesHashrateResponse, error) {

	data, err := h.checkerService.GetIncomesHashrateByDate(ctx, req.Date.AsTime())

	if err != nil {
		sdkLog.Error(ctx, "GetHashrateIncomesByDate: %v", err)
		return nil, err
	}

	response := &accountingPb.GetIncomesHashrateResponse{
		HashrateByDate: make([]*accountingPb.HashrateByDate, len(data)),
	}

	for i := range response.HashrateByDate {
		res := data[i]

		response.HashrateByDate[i] = &accountingPb.HashrateByDate{
			CoinId:   res.CoinId,
			Hashrate: res.Hashrate,
		}
	}

	return response, nil
}

func (h *AccountingHandler) GetCoinsOperationsSum(ctx context.Context, _ *accountingPb.GetCoinsOperationsSumRequest) (*accountingPb.GetCoinsOperationsSumResponse, error) {

	data, err := h.checkerService.GetCoinsOperationsSum(ctx)

	if err != nil {
		sdkLog.Error(ctx, "GetCoinsOperationsSum: %v", err)
		return nil, err
	}

	response := &accountingPb.GetCoinsOperationsSumResponse{
		Data: make([]*accountingPb.OperationsSumData, len(data)),
	}

	for i := range response.Data {
		res := data[i]

		response.Data[i] = &accountingPb.OperationsSumData{
			CoinId: res.CoinId,
			Sum:    res.Sum.String(),
		}
	}

	return response, nil
}

func (h *AccountingHandler) GetTransactionOperationsIntegrity(ctx context.Context, _ *accountingPb.GetTransactionOperationsIntegrityRequest) (*accountingPb.GetTransactionOperationsIntegrityResponse, error) {

	data, err := h.checkerService.GetTransactionOperationsIntegrity(ctx)

	if err != nil {
		sdkLog.Error(ctx, "GetTransactionOperationsIntegrity: %v", err)
		return nil, err
	}

	response := &accountingPb.GetTransactionOperationsIntegrityResponse{
		Data: make([]*accountingPb.TransactionOperationsIntegrityData, len(data)),
	}

	for i := range response.Data {
		res := data[i]

		response.Data[i] = &accountingPb.TransactionOperationsIntegrityData{
			Count:       res.Count,
			TrId:        res.TrId,
			Op2Id:       res.Op2Id,
			Op1Id:       res.Op1Id,
			OpPairCheck: res.OpPairCheck,
			TrNegChk:    res.TrNegChk,
			OpSumChk:    res.OpSumChk,
			DiffChk:     res.DiffChk,
			TrDateChk:   res.TrDateChk,
			CoinChk:     res.CoinChk,
			AccChk:      res.AccChk,
		}
	}

	return response, nil
}

func (h *AccountingHandler) GetCheckTransactionCoins(ctx context.Context, _ *accountingPb.GetCheckTransactionCoinsRequest) (*accountingPb.GetCheckTransactionCoinsResponse, error) {

	data, err := h.checkerService.GetCheckTransactionCoins(ctx)

	if err != nil {
		sdkLog.Error(ctx, "GetCheckTransactionCoins: %v", err)
		return nil, err
	}

	response := &accountingPb.GetCheckTransactionCoinsResponse{
		TrIds: data.TrIds,
		OpIds: data.TrIds,
	}

	return response, nil
}

func (h *AccountingHandler) GetCheckFreezePayoutsBlocks(ctx context.Context, _ *accountingPb.GetCheckFreezePayoutsBlocksRequest) (*accountingPb.GetCheckFreezePayoutsBlocksResponse, error) {

	data, err := h.checkerService.GetCheckFreezePayoutsBlocks(ctx)

	if err != nil {
		sdkLog.Error(ctx, "GetCheckFreezePayoutsBlocks: %v", err)
		return nil, err
	}

	response := &accountingPb.GetCheckFreezePayoutsBlocksResponse{
		Data: make([]*accountingPb.CheckFreezePayoutsBlocksData, len(data)),
	}

	for i := range response.Data {
		res := data[i]

		response.Data[i] = &accountingPb.CheckFreezePayoutsBlocksData{
			TrId:      res.TrId,
			Type:      res.Type,
			UserId:    res.UserId,
			CreatedAt: timestamppb.New(res.CreatedAt),
		}
	}

	return response, nil
}
