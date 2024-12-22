package service

import (
	"context"
	"fmt"
	"strings"

	"code.emcdtech.com/emcd/sdk/log"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	businessError "code.emcdtech.com/b2b/swap/internal/business_error"
	"code.emcdtech.com/b2b/swap/internal/client"
	"code.emcdtech.com/b2b/swap/internal/repository"
	"code.emcdtech.com/b2b/swap/model"
)

type Admin interface {
	GetBalanceByCoin(ctx context.Context, accountType, coin string) (decimal.Decimal, error)
	TransferBetweenAccountTypes(ctx context.Context, fromAccountType, ToAccountType, coin string, amount decimal.Decimal) error
	PlaceOrderForUSDT(ctx context.Context, coin string, direction model.Direction, amount decimal.Decimal) (uuid.UUID, error)
	CheckOrder(ctx context.Context, id uuid.UUID) (model.OrderStatus, error)
	Withdraw(ctx context.Context, swapID uuid.UUID) (int, error)
	GetWithdrawalLink(ctx context.Context, withdrawalID int) (string, error)
	RequestAQuote(ctx context.Context, from, to, accountType string, amount decimal.Decimal) (*model.Quote, error)
	ConfirmAQuote(ctx context.Context, id string) (string, error)
	GetConvertStatus(ctx context.Context, id, accountType string) (string, error)
	ChangeManualSwapStatus(ctx context.Context, swapID uuid.UUID, status model.Status) error
	GetSwapStatusHistory(ctx context.Context, swapID uuid.UUID) ([]*model.SwapStatusHistoryItem, error)
}

type admin struct {
	market               client.Market
	exchangeAccount      client.ExchangeAccount
	exchangeTransaction  client.ExchangeTransaction
	emailCli             client.Email
	swapRep              repository.Swap
	swapStatusHistoryRep repository.SwapStatusHistory
	transferRep          repository.Transfer
	withdrawRep          repository.Withdraw
	coinRep              repository.Coin
	orderBookRep         repository.OrderBook
	explorerRep          repository.Explorer
	userRep              repository.User
	statusUpdater        SwapStatusUpdater
	byBitMasterUid       int
	masterApiKey         string
	masterApiSecret      string
}

func NewAdmin(
	market client.Market,
	exchangeAccount client.ExchangeAccount,
	exchangeTransaction client.ExchangeTransaction,
	emailCli client.Email,
	swapRep repository.Swap,
	swapStatusHistoryRep repository.SwapStatusHistory,
	transferRep repository.Transfer,
	withdrawRep repository.Withdraw,
	coinRep repository.Coin,
	orderBookRep repository.OrderBook,
	explorerRep repository.Explorer,
	userRep repository.User,
	statusUpdater SwapStatusUpdater,
	byBitMasterUid int,
	masterApiKey string,
	masterApiSecret string,
) Admin {
	return &admin{
		market:               market,
		exchangeAccount:      exchangeAccount,
		exchangeTransaction:  exchangeTransaction,
		swapRep:              swapRep,
		swapStatusHistoryRep: swapStatusHistoryRep,
		emailCli:             emailCli,
		transferRep:          transferRep,
		withdrawRep:          withdrawRep,
		coinRep:              coinRep,
		orderBookRep:         orderBookRep,
		explorerRep:          explorerRep,
		userRep:              userRep,
		statusUpdater:        statusUpdater,
		byBitMasterUid:       byBitMasterUid,
		masterApiKey:         masterApiKey,
		masterApiSecret:      masterApiSecret,
	}
}

func (a *admin) GetSwapStatusHistory(ctx context.Context, swapID uuid.UUID) ([]*model.SwapStatusHistoryItem, error) {
	sw, err := a.swapStatusHistoryRep.Find(ctx, &model.SwapStatusHistoryFilter{SwapID: &swapID})
	if err != nil {
		return nil, fmt.Errorf("getStatusHistory: %w", err)
	}

	return sw, nil
}

func (a *admin) GetBalanceByCoin(ctx context.Context, accountType, coin string) (decimal.Decimal, error) {
	amount, err := a.exchangeAccount.GetBalanceByCoin(ctx, a.byBitMasterUid, coin, accountType)
	if err != nil {
		return decimal.Zero, fmt.Errorf("getBalanceByCoin: %w", err)
	}
	// TODO вернуть всё
	return amount.WalletBalance, nil
}

func (a *admin) TransferBetweenAccountTypes(ctx context.Context, fromAccountType, ToAccountType, coin string, amount decimal.Decimal) error {
	_, err := a.exchangeTransaction.CreateInternalTransfer(ctx, &model.InternalTransfer{
		ID:              uuid.New(),
		Coin:            coin,
		Amount:          amount,
		FromAccountType: fromAccountType,
		ToAccountType:   ToAccountType,
	}, a.masterApiKey, a.masterApiSecret)
	if err != nil {
		return fmt.Errorf("createInternalTransfer: %w", err)
	}
	return nil
}

func (a *admin) PlaceOrderForUSDT(ctx context.Context, coin string, direction model.Direction, amount decimal.Decimal) (uuid.UUID, error) {
	mySymbol := fmt.Sprintf("%s%s", strings.ToUpper(coin), model.CoinUSDT)
	ok := a.orderBookRep.IsExist(mySymbol)
	if !ok {
		return uuid.Nil, fmt.Errorf("iSExist: symbol: %s, err: %w", mySymbol, businessError.SymbolNotFoundAdminErr)
	}

	orderID := uuid.New()

	err := a.market.PlaceOrder(ctx, &model.Order{
		ID:         orderID,
		Category:   model.Spot,
		Symbol:     mySymbol,
		Direction:  direction,
		AmountFrom: amount,
	}, &model.Secrets{
		AccountID: int64(a.byBitMasterUid),
		ApiKey:    a.masterApiKey,
		ApiSecret: a.masterApiSecret,
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("placeOrder: %w", err)
	}

	return orderID, nil
}

func (a *admin) CheckOrder(ctx context.Context, id uuid.UUID) (model.OrderStatus, error) {
	status, err := a.market.GetOrderStatus(ctx, id, a.masterApiKey, a.masterApiSecret)
	if err != nil {
		return model.OrderUnknown, fmt.Errorf("getOrderStatus: %w", err)
	}
	return status, nil
}

func (a *admin) Withdraw(ctx context.Context, swapID uuid.UUID) (int, error) {
	// Получить нужный swap
	mySwap, err := a.swapRep.FindOne(ctx,
		&model.SwapFilter{
			ID: &swapID,
		})
	if err != nil {
		return 0, fmt.Errorf("findOne: %w", err)
	}

	// Проверяем статус, что бы не вывести монеты 2 раза
	// Монеты для ручного вывода доступны, если своп упал с ошибкой
	// TODO убрать CANCEL
	if mySwap.Status != model.Error && mySwap.Status != model.Cancel {
		return 0, businessError.WithdrawImpossibleBecauseOfSwapStatus
	}

	// Проверить, что в базе withdraw ещё нет вывода по swapID
	wtds, err := a.withdrawRep.Find(ctx, &model.WithdrawFilter{
		SwapID: &swapID,
	})
	if err != nil {
		return 0, fmt.Errorf("find: %w", err)
	}
	if len(wtds) > 0 {
		if findAtLeastOneAliveWithdraw(wtds) {
			return 0, businessError.WithdrawImpossibleBecauseWithdrawHasAlreadyBeen
		}
	}

	// Truncate amount используя точность для вывода
	accuracy, err := a.coinRep.GetAccuracyForWithdrawAndDeposit(ctx, mySwap.CoinTo, mySwap.NetworkTo)
	if err != nil {
		log.Error(ctx, "admin: withdraw: getAccuracyForWithdrawAndDeposit: %s", err.Error())
		// Без return, не принципиально
		//return 0, fmt.Errorf("getAccuracyForWithdrawAndDeposit: %w", err)
	}
	if accuracy == 0 {
		accuracy = model.DefaultWithdrawAccuracy
	}
	mySwap.AmountTo = mySwap.AmountTo.Truncate(int32(accuracy))

	wtd := &model.Withdraw{
		InternalID:         uuid.New(),
		ID:                 0,
		SwapID:             mySwap.ID,
		HashID:             "",
		Coin:               mySwap.CoinTo,
		Network:            mySwap.NetworkTo,
		Address:            mySwap.AddressTo,
		Tag:                mySwap.TagTo,
		Amount:             mySwap.AmountTo,
		IncludeFeeInAmount: true, // true - фактическая сумма, которая будет получена
	}

	id, err := a.exchangeTransaction.Withdraw(ctx, wtd)
	if err != nil {
		return 0, fmt.Errorf("withdraw: %w", err)
	}
	wtd.ID = id

	err = a.withdrawRep.Add(ctx, wtd)
	if err != nil {
		log.Error(ctx, "admin: withdraw: withdrawRep.Add: %s", err.Error())
		// Без return, что бы не допустить повторного вывода, т.к., транзакция выполнена exchangeTransaction.Withdraw
		//return 0, fmt.Errorf("add: %w", err)
	}

	err = a.statusUpdater.UpdateAndBroadcast(ctx, mySwap, model.ManualCompleted)
	if err != nil {
		log.Error(ctx, "admin: withdraw: updateAndBroadcast: %s", err.Error())
		// Без return, что бы не допустить повторного вывода, т.к., транзакция выполнена exchangeTransaction.Withdraw
		//return 0, fmt.Errorf("update: %w", err)
	}

	// Всё что ниже, не должно возвращать ошибку, потому что swap уже в статусе ManualCompleted
	a.sendSuccessfullyEmail(ctx, mySwap)

	return int(id), nil
}

func (a *admin) GetWithdrawalLink(ctx context.Context, withdrawalID int) (string, error) {
	wtd, err := a.exchangeTransaction.GetWithdraw(ctx, withdrawalID)
	if err != nil {
		return "", fmt.Errorf("getWithdraw: %w", err)
	}

	link, err := a.explorerRep.GetTransactionLink(ctx, wtd.Coin, wtd.HashID)
	if err != nil {
		return "", fmt.Errorf("getTransactionLink: %w", err)
	}
	if link != "" {
		return link, nil
	} else {
		return wtd.HashID, nil
	}
}

func (a *admin) RequestAQuote(ctx context.Context, from, to, accountType string, amount decimal.Decimal) (*model.Quote, error) {
	quote, err := a.market.RequestAQuote(ctx, from, to, accountType, amount)
	if err != nil {
		return nil, fmt.Errorf("requestAQuote: %w", err)
	}
	return quote, nil
}

func (a *admin) ConfirmAQuote(ctx context.Context, id string) (string, error) {
	quote, err := a.market.ConfirmAQuote(ctx, id)
	if err != nil {
		return "", fmt.Errorf("confirmAQuote: %w", err)
	}
	return quote, nil
}

func (a *admin) GetConvertStatus(ctx context.Context, id, accountType string) (string, error) {
	status, err := a.market.GetConvertStatus(ctx, id, accountType)
	if err != nil {
		return "", fmt.Errorf("getConvertStatus: %w", err)
	}
	return status, nil
}

func (a *admin) ChangeManualSwapStatus(ctx context.Context, swapID uuid.UUID, status model.Status) error {
	sw, err := a.swapRep.FindOne(ctx, &model.SwapFilter{
		ID: &swapID,
	})
	if err != nil {
		return fmt.Errorf("findOne: %w", err)
	}

	err = a.swapRep.WithinTransaction(ctx, func(ctx context.Context) error {
		err := a.swapRep.Update(ctx, sw, &model.SwapFilter{ID: &sw.ID}, &model.SwapPartial{Status: &status})
		if err != nil {
			return fmt.Errorf("update: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("withinTransaction: %w", err)
	}

	return nil
}

func findAtLeastOneAliveWithdraw(withdraws model.Withdraws) bool {
	for _, wt := range withdraws {
		if wt.Status != model.WsFailed && wt.Status != model.WsUnknown {
			return true
		}
	}
	return false
}

func (a *admin) sendSuccessfullyEmail(ctx context.Context, sw *model.Swap) {
	user, err := a.userRep.FindOne(ctx, &model.UserFilter{
		ID: &sw.UserID,
	})
	if err != nil {
		log.Error(ctx, "admin: sendSuccessfullyEmail: userRep.FindOne: %s", err.Error())
		return
	}

	a.emailCli.SendSuccessfulSwapMessage(ctx, sw, user)
}
