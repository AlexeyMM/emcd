package service

import (
	"context"
	"fmt"
	"time"

	"code.emcdtech.com/emcd/sdk/log"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"code.emcdtech.com/b2b/swap/internal/business_error"
	"code.emcdtech.com/b2b/swap/internal/client"
	"code.emcdtech.com/b2b/swap/internal/repository"
	"code.emcdtech.com/b2b/swap/internal/slack"
	"code.emcdtech.com/b2b/swap/model"
)

type Swap interface {
	// Swap флоу
	SwapEstimate(ctx context.Context, swap *model.Swap) (*model.Estimate, error)
	PrepareSwap(ctx context.Context, swap *model.Swap) (swapID uuid.UUID, depositAddress *model.AddressData, err error)
	StartSwap(ctx context.Context, swapID uuid.UUID, email, language string) error
	WaitDeposit(ctx context.Context, sw *model.Swap) (model.Status, error)
	GetSwapStatus(ctx context.Context, swapID uuid.UUID) (model.Status, error)
	CheckDeposit(ctx context.Context, sw *model.Swap) (model.Status, error)
	TransferToUnified(ctx context.Context, swap *model.Swap) (model.Status, error)
	CreateOrders(ctx context.Context, swap *model.Swap) (model.Orders, model.Status, error)
	PlaceOrder(ctx context.Context, swap *model.Swap, order *model.Order) (model.Status, error)
	CheckOrder(ctx context.Context, swap *model.Swap, orderID uuid.UUID) (model.Status, error)
	TransferFromSubToMaster(ctx context.Context, swap *model.Swap, transfer *model.InternalTransfer) (model.Status, error)
	CheckTransferFromSubToMaster(ctx context.Context, swap *model.Swap, accountID int64) (model.Status, error)
	PrepareWithdraw(ctx context.Context, swap *model.Swap) (model.Status, error)
	Withdraw(ctx context.Context, swap *model.Swap) (model.Status, error)
	CheckWithdraw(ctx context.Context, swap *model.Swap, withdrawID int64) (model.Status, error)

	// Другие методы
	GetAllActiveSwaps(ctx context.Context) ([]*model.Swap, error)
	GetSwapByID(ctx context.Context, id uuid.UUID) (*SwapByID, error)
	GetSwaps(ctx context.Context, filter *model.SwapFilter) ([]*model.Swap, int, error)
	Update(ctx context.Context, filter *model.SwapFilter, partial *model.SwapPartial) error
}

type swap struct {
	market               client.Market
	exchangeAccount      client.ExchangeAccount
	exchangeTransaction  client.ExchangeTransaction
	exchangeSubscriber   client.Subscriber
	emailCli             client.Email
	swapRep              repository.Swap
	activeSwapRep        repository.ActiveSwap
	orderRep             repository.Order
	accountRep           repository.Account
	depositRep           repository.Deposit
	transferRep          repository.Transfer
	withdrawRep          repository.Withdraw
	coinRep              repository.Coin
	feeRep               repository.OrderFee
	symbolRep            repository.Symbol
	orderBookRep         repository.OrderBook
	userRep              repository.User
	swapStatusUpdater    SwapStatusUpdater
	ourFee               decimal.Decimal
	maxLimit             decimal.Decimal
	minLimit             decimal.Decimal
	byBitMasterUid       int
	mainCtx              context.Context
	slack                slack.Slack
	depositWaitingPeriod time.Duration
	swapExecutorCh       chan<- *model.Swap
}

func NewSwap(
	market client.Market,
	exchangeAccount client.ExchangeAccount,
	exchangeTransaction client.ExchangeTransaction,
	exchangeSubscriber client.Subscriber,
	emailCli client.Email,
	swapRep repository.Swap,
	activeSwapRep repository.ActiveSwap,
	orderRep repository.Order,
	accountRep repository.Account,
	depositRep repository.Deposit,
	transferRep repository.Transfer,
	withdrawRep repository.Withdraw,
	coinRep repository.Coin,
	feeRep repository.OrderFee,
	symbolRep repository.Symbol,
	orderBookRep repository.OrderBook,
	userRep repository.User,
	swapStatusUpdater SwapStatusUpdater,
	ourFee decimal.Decimal,
	maxLimit decimal.Decimal,
	minLimit decimal.Decimal,
	byBitMasterUid int,
	mainCtx context.Context,
	slack slack.Slack,
	depositWaitingPeriod time.Duration,
	swapExecutorCh chan<- *model.Swap,
) *swap {
	return &swap{
		market:               market,
		exchangeAccount:      exchangeAccount,
		exchangeTransaction:  exchangeTransaction,
		exchangeSubscriber:   exchangeSubscriber,
		emailCli:             emailCli,
		swapRep:              swapRep,
		activeSwapRep:        activeSwapRep,
		orderRep:             orderRep,
		accountRep:           accountRep,
		depositRep:           depositRep,
		transferRep:          transferRep,
		withdrawRep:          withdrawRep,
		coinRep:              coinRep,
		feeRep:               feeRep,
		symbolRep:            symbolRep,
		orderBookRep:         orderBookRep,
		userRep:              userRep,
		swapStatusUpdater:    swapStatusUpdater,
		ourFee:               ourFee,
		maxLimit:             maxLimit,
		minLimit:             minLimit,
		byBitMasterUid:       byBitMasterUid,
		mainCtx:              mainCtx,
		slack:                slack,
		depositWaitingPeriod: depositWaitingPeriod,
		swapExecutorCh:       swapExecutorCh,
	}
}

func (s *swap) SwapEstimate(ctx context.Context, swap *model.Swap) (*model.Estimate, error) {
	if !swap.AmountFrom.IsZero() {
		// estimate от amountFrom
		return s.swapEstimate(ctx, swap)
	} else {
		// estimate от amountTo
		// Задача сделать, что бы от amountTo получался примерно такой же результат как от amountFrom

		// Запоминаем amountTo, сколько хотим получить
		amountTo := swap.AmountTo

		// Получаем приблизительный результат, нас интересует amountFrom для дальнейших расчётов
		est, err := s.swapEstimate(ctx, swap)
		if err != nil {
			return nil, err
		}
		amountFrom := swap.AmountFrom

		var (
			percentIncrease = decimal.NewFromFloat(1.001)
			percentDecrease = decimal.NewFromFloat(0.999)
		)

		// Если полученный amountTo больше начального - рекурсивно уменьшаем amountFrom
		// Если полученный amountTo меньше начального - рекурсивно увеличиваем amountFrom
		return s.recursiveEstimate(ctx, swap, est, amountTo, amountFrom, percentIncrease, percentDecrease, 1)
	}
}

func (s *swap) recursiveEstimate(ctx context.Context,
	swap *model.Swap,
	est *model.Estimate,
	amountTo,
	amountFrom,
	percentIncrease,
	percentDecrease decimal.Decimal,
	count int,
) (*model.Estimate, error) {
	var err error

	// страховка от бесконечных циклов
	if count > 1000 {
		return nil, fmt.Errorf("count exсeeded %w", businessError.CalculateSwapErr)
	}

	if est.AmountTo.GreaterThan(amountTo) {
		amountFrom = amountFrom.Mul(percentDecrease)
		swap.AmountFrom = amountFrom
		swap.AmountTo = decimal.Zero
		est, err = s.swapEstimate(ctx, swap)
		if err != nil {
			return nil, err
		}
		if est.AmountTo.LessThan(amountTo) {
			est.AmountTo = amountTo
			return est, nil
		}

	} else {
		amountFrom = amountFrom.Mul(percentIncrease)
		swap.AmountFrom = amountFrom
		swap.AmountTo = decimal.Zero
		est, err = s.swapEstimate(ctx, swap)
		if err != nil {
			return nil, err
		}
		if est.AmountTo.GreaterThan(amountTo) {
			est.AmountTo = amountTo
			return est, nil
		}
	}

	count = count + 1

	return s.recursiveEstimate(ctx, swap, est, amountTo, amountFrom, percentIncrease, percentDecrease, count)
}

func (s *swap) swapEstimate(ctx context.Context, swap *model.Swap) (*model.Estimate, error) {
	if swap.CoinFrom == swap.CoinTo {
		return nil, businessError.SameCoinsErr
	}

	err := s.checkOnPossibilityWithdraw(ctx, swap.CoinTo, swap.NetworkTo)
	if err != nil {
		return nil, fmt.Errorf("checkOnPossibilityWithdraw: %w", err)
	}

	maxLimit, err := s.calculateMaxLimit(ctx, swap)
	if err != nil {
		return nil, fmt.Errorf("checkAndReturnMaxLimit: %w", err)
	}

	bestPath, err := s.getSwapPathAndCalculateAmount(ctx, swap)
	if err != nil {
		return nil, fmt.Errorf("getSwapPathAndCalculateAmount: %w", err)
	}

	minimum, err := s.calculateMin(ctx, swap.CoinFrom, swap.CoinTo, swap.NetworkFrom, swap.NetworkTo, bestPath)
	if err != nil {
		return nil, fmt.Errorf("checkAndReturnMinLimit: %w", err)
	}

	netFrom, err := s.coinRep.GetNetwork(ctx, swap.CoinFrom, swap.NetworkFrom)
	if err != nil {
		return nil, fmt.Errorf("getNetwork 1: %w", err)
	}

	netTo, err := s.coinRep.GetNetwork(ctx, swap.CoinTo, swap.NetworkTo)
	if err != nil {
		return nil, fmt.Errorf("getNetwork 2: %w", err)
	}

	// Отрицательный может получиться после расчёта маленьких сумм на дешёвых монетах, после вычета комиссий
	if swap.AmountFrom.LessThan(decimal.Zero) {
		swap.AmountFrom = decimal.Zero
	}
	if swap.AmountTo.LessThan(decimal.Zero) {
		swap.AmountTo = decimal.Zero
	}

	return &model.Estimate{
		AmountFrom: swap.AmountFrom.Truncate(int32(netFrom.AccuracyWithdrawAndDeposit)),
		AmountTo:   swap.AmountTo.Truncate(int32(netTo.AccuracyWithdrawAndDeposit)),
		Rate:       swap.AmountTo.Div(swap.AmountFrom).Truncate(int32(netFrom.AccuracyWithdrawAndDeposit)),
		Limits: &model.Limits{
			Min: minimum.Truncate(int32(netFrom.AccuracyWithdrawAndDeposit)),
			Max: maxLimit.Truncate(int32(netFrom.AccuracyWithdrawAndDeposit)),
		},
	}, nil
}

func (s *swap) PrepareSwap(ctx context.Context, sw *model.Swap) (swapID uuid.UUID, depositAddress *model.AddressData, err error) {
	if !s.activeSwapRep.IsSwapLimitExceeded() {
		return uuid.Nil, nil, fmt.Errorf("isSwapLimitExceeded: %w", businessError.SwapActiveLimitExceededErr)
	}

	err = s.checkOnPossibilityWithdraw(ctx, sw.CoinTo, sw.NetworkTo)
	if err != nil {
		return uuid.Nil, nil, fmt.Errorf("checkOnPossibilityWithdraw: %w", err)
	}

	_, err = s.checkAndReturnMaxLimit(ctx, sw)
	if err != nil {
		return uuid.Nil, nil, fmt.Errorf("checkAndReturnMaxLimit: %w", err)
	}

	bestPath, err := s.getSwapPathAndCalculateAmount(ctx, sw)
	if err != nil {
		return uuid.Nil, nil, fmt.Errorf("getSwapPathAndCalculateAmount: %w", err)
	}

	// Отрицательный может получиться после расчёта маленьких сумм на дешёвых монетах, после вычета комиссий
	if sw.AmountFrom.LessThanOrEqual(decimal.Zero) || sw.AmountTo.LessThanOrEqual(decimal.Zero) {
		return uuid.Nil, nil, fmt.Errorf("invalid amount: amount_from: %s, amount_to: %s, err: %w",
			sw.AmountFrom, sw.AmountTo, businessError.CalculateSwapErr)
	}

	_, err = s.checkAndReturnMinLimit(ctx, sw, bestPath)
	if err != nil {
		return uuid.Nil, nil, fmt.Errorf("checkAndReturnMinLimit: %w", err)
	}

	account, err := s.exchangeAccount.CreateSubAccount(ctx)
	if err != nil {
		return uuid.Nil, nil, fmt.Errorf("createSubAccount: %w", err)
	}
	if !account.IsValid {
		return uuid.Nil, nil, businessError.CreateSubAccountErr
	}

	sw.AccountFrom = account.ID

	account.Keys, err = s.exchangeAccount.CreateSubAPIKey(ctx, int(account.ID))
	if err != nil {
		return uuid.Nil, nil, fmt.Errorf("createSubAPIKey: %w", err)
	}

	log.Debug(ctx, "createSubAcc: %+v", account)
	log.Debug(ctx, "keys: %+v", account.Keys)

	depositAddress, err = s.exchangeAccount.GetSubDepositAddress(ctx, int(account.ID), sw.CoinFrom, sw.NetworkFrom)
	if err != nil {
		return uuid.Nil, nil, fmt.Errorf("getSubDepositAddress: %w", err)
	}

	sw.AddressFrom = depositAddress.Address
	sw.TagFrom = depositAddress.Tag

	err = s.accountRep.Add(ctx, account)
	if err != nil {
		return uuid.Nil, nil, fmt.Errorf("addSubAccount.Add: %w", err)
	}

	sw.ID = uuid.New()
	sw.StartTime = time.Now().UTC()
	sw.Status = model.Unknown

	err = s.swapRep.Add(ctx, sw)
	if err != nil {
		return uuid.Nil, nil, fmt.Errorf("add: %w", err)
	}

	return sw.ID, depositAddress, nil
}

func (s *swap) StartSwap(ctx context.Context, swapID uuid.UUID, email, language string) error {
	if s.activeSwapRep.Exist(swapID) {
		return fmt.Errorf("swap already exist: %w", businessError.SwapAlreadyExistsErr)
	}

	// Создаём юзера
	userID := uuid.New()
	err := s.userRep.Add(ctx, &model.User{
		ID:       userID,
		Email:    email,
		Language: language,
	})
	if err != nil {
		return fmt.Errorf("createUser: %w", err)
	}

	err = s.startSwapExecution(ctx, userID, swapID)
	if err != nil {
		return fmt.Errorf("startSwapExecution: %w", err)
	}

	s.emailCli.SendInitialSwapMessage(ctx, email, language, swapID)

	return nil
}

func (s *swap) startSwapExecution(ctx context.Context, userID, swapID uuid.UUID) error {
	err := s.swapRep.WithinTransaction(ctx, func(ctx context.Context) error {
		sw, err := s.swapRep.FindOne(ctx,
			&model.SwapFilter{
				ID: &swapID,
			})
		if err != nil {
			return fmt.Errorf("findOne: %w", err)
		}

		startTime := time.Now().UTC()
		status := model.WaitDeposit

		err = s.swapRep.Update(ctx, sw,
			&model.SwapFilter{
				ID: &swapID,
			},
			&model.SwapPartial{
				UserID:    &userID,
				Status:    &status,
				StartTime: &startTime,
			})
		if err != nil {
			return fmt.Errorf("update: %w", err)
		}

		select {
		case s.swapExecutorCh <- sw:
		default:
			err = s.swapStatusUpdater.UpdateAndBroadcast(ctx, sw, model.Error)
			if err != nil {
				return fmt.Errorf("updateAndBroadcast: %w", err)
			}
			return fmt.Errorf("swapExecutorCh channel full: %w", businessError.SwapActiveLimitExceededErr)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("withinTransaction: %w", err)
	}
	return nil
}

// WaitDeposit проверяет поступил ли депозит по REST API, не используем websocket из-за
// https://emcdtechltd.slack.com/archives/C076DDWGMLL/p1728896418394689
// При использовании websocket нужно выключить этот шаг, потому что логика дублируется там
// Выключить, означает убрать его из swap_executor и заменить на GetSwapStatus, ожидая когда статус смениться на CheckDeposit,
// wallet_ws поменяет статус когда придёт event.
// Так же нужно вернуть подписку на все wallet websocket при старте приложения, у которых статус = waitDeposit
func (s *swap) WaitDeposit(ctx context.Context, sw *model.Swap) (model.Status, error) {
	log.Debug(ctx, "waitDeposit: swap_id: %s", sw.ID)
	remainTime := sw.StartTime.Add(s.depositWaitingPeriod).Sub(time.Now().UTC())
	if remainTime <= 0 {
		log.Debug(ctx, "waitDeposit: remainTime is less than 0: swapID: %s", sw.ID.String())

		err := s.swapStatusUpdater.UpdateAndBroadcast(ctx, sw, model.Cancel)
		if err != nil {
			return model.Unknown, fmt.Errorf("swapStatusUpdater: %w", err)
		}

		return model.Cancel, nil
	}
	return s.CheckDeposit(ctx, sw)
}

func (s *swap) GetSwapStatus(ctx context.Context, swapID uuid.UUID) (model.Status, error) {
	sw, err := s.swapRep.FindOne(ctx, &model.SwapFilter{
		ID: &swapID,
	})
	if err != nil {
		return 0, fmt.Errorf("swapRep.GetStatus: %w", err)
	}
	return sw.Status, nil
}

func (s *swap) CheckDeposit(ctx context.Context, sw *model.Swap) (model.Status, error) {
	account, err := s.accountRep.FindOne(ctx, &model.AccountFilter{
		ID: &sw.AccountFrom,
	})
	if err != nil {
		return model.Unknown, fmt.Errorf("find: %w", err)
	}

	deposits, err := s.exchangeAccount.GetDepositRecords(ctx, sw.CoinFrom, sw.StartTime, account.Keys.ApiKey, account.Keys.ApiSecret)
	if err != nil {
		return model.Unknown, fmt.Errorf("getDepositRecords: %w", err)
	}
	if len(deposits) == 0 {
		return sw.Status, nil
	}

	// Проверка на AML и failed и pending депозиты
	for _, dep := range deposits {
		if dep.DepositType != model.DepositNormal {
			log.Debug(ctx, "deposit abnormal: %+v", dep)

			err = s.swapStatusUpdater.UpdateAndBroadcast(ctx, sw, model.Error)
			if err != nil {
				return model.Unknown, fmt.Errorf("updateAndBroadcast: %w", err)
			}

			err = s.slack.Send(ctx, fmt.Sprintf("swap: %s, abnormal deposit: возможно не пройдена AML проверка", sw.ID.String()))
			if err != nil {
				return model.Unknown, fmt.Errorf("slack.Send: %w", err)
			}

			return model.Error, nil
		}

		if dep.Status == model.DepositFailed {
			log.Debug(ctx, "deposit failed: %+v", dep)
			err = s.swapStatusUpdater.UpdateAndBroadcast(ctx, sw, model.Error)
			if err != nil {
				return model.Unknown, fmt.Errorf("updateAndBroadcast: %w", err)
			}

			err = s.slack.Send(ctx, fmt.Sprintf("swap: %s, deposit failed", sw.ID.String()))
			if err != nil {
				return model.Unknown, fmt.Errorf("slack.Send: %w", err)
			}

			return model.Error, nil
		}

		if dep.Status == model.DepositPending {
			log.Debug(ctx, "deposit pending: %+v", dep)
			return model.CheckDeposit, nil

		}
	}

	var total decimal.Decimal
	for _, dep := range deposits {
		log.Debug(ctx, "deposit: %+v", dep)
		switch dep.Status {
		case model.DepositSuccess:
			total = total.Add(dep.Amount).Add(dep.Fee)

		default:
			return model.Unknown, fmt.Errorf("unknown deposit status: %+v", dep)
		}
	}

	// Сохраняем в базу
	err = s.depositRep.WithinTransaction(ctx, func(ctx context.Context) error {
		for _, dep := range deposits {
			log.Debug(ctx, "save deposit: %+v", dep)
			dep.SwapID = sw.ID
			err = s.depositRep.Add(ctx, dep)
			if err != nil {
				return fmt.Errorf("add: %w", err)
			}
		}

		// Обновляем amount_from, amount_to если сумма депозита отличается от amount_from
		// Проверяем лимиты, если выходим за рамки - Error
		if !sw.AmountFrom.Equal(total) {
			log.Info(ctx, "swap_id: %s, deposit is different from amount_from: deposit: %s, amount_from: %s", sw.ID, total.String(), sw.AmountFrom.String())
			sw.AmountFrom = total
			err = s.updateAmountsAfterDeposit(ctx, sw)
			if err != nil {
				return fmt.Errorf("updateAmountsAfterDeposit: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("withinTransaction: %w", err)
	}

	err = s.swapStatusUpdater.UpdateAndBroadcast(ctx, sw, model.TransferToUnified)
	if err != nil {
		return model.Unknown, fmt.Errorf("updateAndBroadcast: %w", err)
	}

	return model.TransferToUnified, nil
}

func (s *swap) updateAmountsAfterDeposit(ctx context.Context, sw *model.Swap) error {
	estimate, err := s.swapEstimate(ctx, sw)
	if err != nil {
		return fmt.Errorf("swapEstimate: %w", err)
	}

	log.Debug(ctx, "swap_id: %s, new estimate: %+v", sw.ID, estimate)

	err = s.swapRep.Update(ctx, sw,
		&model.SwapFilter{
			ID: &sw.ID,
		},
		&model.SwapPartial{
			AmountFrom: &estimate.AmountFrom,
			AmountTo:   &estimate.AmountTo,
		})
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}

	if estimate.AmountFrom.LessThan(estimate.Limits.Min) {
		err = s.swapStatusUpdater.UpdateAndBroadcast(ctx, sw, model.Error)
		if err != nil {
			return fmt.Errorf("updateAndBroadcast 1: %w", err)
		}
		return fmt.Errorf("amount_from: %s, min_from: %s err: %w",
			estimate.AmountFrom.String(), estimate.Limits.Min.String(), businessError.BelowMinLimitErr)
	}
	if estimate.AmountTo.GreaterThan(estimate.Limits.Max) {
		err = s.swapStatusUpdater.UpdateAndBroadcast(ctx, sw, model.Error)
		if err != nil {
			return fmt.Errorf("updateAndBroadcast 2: %w", err)
		}
		return fmt.Errorf("amount_to: %s, max_to: %s, err: %w",
			estimate.AmountTo.String(), estimate.Limits.Max.String(), businessError.MaxLimitExceededErr)
	}

	return nil
}

func (s *swap) TransferToUnified(ctx context.Context, sw *model.Swap) (model.Status, error) {
	isLast := true
	trs, err := s.transferRep.Find(ctx, &model.InternalTransferFilter{
		FromAccountID: &sw.AccountFrom,
		IsLast:        &isLast,
	})
	if err != nil {
		return 0, fmt.Errorf("findOne: %w", err)
	}

	var myTransfer *model.InternalTransfer
	if len(trs) == 1 {
		myTransfer = trs[0]
	}

	if len(trs) == 0 {
		// Переводим с FUND в SPOT
		amount, err := s.getBalanceByCoin(ctx, sw.ID, int(sw.AccountFrom), sw.CoinFrom, model.Fund)
		if err != nil {
			return model.Unknown, fmt.Errorf("getBalanceByCoin: %w", err)
		}
		if amount.IsZero() {
			return sw.Status, nil
		}

		if !sw.AmountFrom.Equal(amount) {
			log.Warn(ctx, "transferToUnified: amount_from not equal balance: amount_from: %s, balance: %s",
				sw.AmountFrom.String(), amount.String())
		}

		intTr := model.InternalTransfer{
			ID:              uuid.New(),
			Coin:            sw.CoinFrom,
			Amount:          amount,
			FromAccountType: model.Fund,
			ToAccountType:   model.UNIFIED,
			UpdatedAt:       time.Now(),
		}
		acc, err := s.accountRep.FindOne(ctx, &model.AccountFilter{
			ID: &sw.AccountFrom,
		})
		if err != nil {
			return 0, fmt.Errorf("findOne: %w", err)
		}
		log.Debug(ctx, "try createInternalTransfer: %+v", intTr)
		internalTransfer, err := s.exchangeTransaction.CreateInternalTransfer(ctx, &intTr, acc.Keys.ApiKey, acc.Keys.ApiSecret)
		if err != nil {
			return 0, fmt.Errorf("createInternalTransfer: %w", err)
		}
		err = s.transferRep.Add(ctx, internalTransfer)
		if err != nil {
			return 0, fmt.Errorf("add: %w", err)
		}
		myTransfer = internalTransfer
	} else {
		if myTransfer == nil {
			return model.Unknown, fmt.Errorf("transferRep.Find 2: expected 1 transaction, found %d", len(trs))
		}
		getTransfer, err := s.exchangeTransaction.GetTransfer(ctx, myTransfer.ID)
		if err != nil {
			return 0, fmt.Errorf("getTransfer: %w", err)
		}
		err = s.transferRep.Update(ctx, myTransfer,
			&model.InternalTransferFilter{
				ID: &myTransfer.ID,
			},
			&model.InternalTransferPartial{
				Status: &myTransfer.Status,
			})
		if err != nil {
			return 0, fmt.Errorf("update: %w", err)
		}
		myTransfer = getTransfer
	}

	switch myTransfer.Status {
	case model.ItsSuccess:
		log.Debug(ctx, "transferToUnified successfully: %+v", myTransfer)
		err = s.swapStatusUpdater.UpdateAndBroadcast(ctx, sw, model.CreateOrder)
		if err != nil {
			return 0, fmt.Errorf("updateAndBroadcast: %w", err)
		}
		return model.CreateOrder, nil
	case model.ItsPending:
		return sw.Status, nil
	case model.ItsFailed:
		log.Debug(ctx, "transferToUnifiedFailed: %+v", myTransfer)
		return model.Unknown, fmt.Errorf("transfer failed: swapID: %s, error: %w", sw.ID.String(), err)
	default:
		log.Debug(ctx, "transferToUnifiedFailed default: %+v", myTransfer)
		return model.Unknown, fmt.Errorf("transfer to spot: default: swapID: %s, error: %w", sw.ID.String(), err)
	}
}

func (s *swap) CreateOrders(ctx context.Context, sw *model.Swap) (model.Orders, model.Status, error) {
	direct, indirect, err := s.getSwapOptions(ctx, sw.CoinFrom, sw.CoinTo)
	if err != nil {
		return nil, model.Unknown, fmt.Errorf("getSwapOptions: %w", err)
	}

	_, isDirect, err := s.getAmountToAndTheBestOptionToSwap(ctx, direct, indirect, sw.AmountFrom, false, decimal.Zero, decimal.Zero)
	if err != nil {
		return nil, model.Unknown, fmt.Errorf("getAmountToAndTheBestOptionToSwap: %w", err)
	}

	account, err := s.accountRep.FindOne(ctx, &model.AccountFilter{
		ID: &sw.AccountFrom,
	})
	if err != nil {
		return nil, model.Unknown, fmt.Errorf("getSubAccount: %w", err)
	}

	var orders []*model.Order

	if isDirect {
		sw.AmountFrom = s.truncateAmountBySymbolRules(ctx, sw.AmountFrom, direct.Symbol, direct.Direction, true)

		ord := model.Order{
			ID:         uuid.New(),
			SwapID:     sw.ID,
			AccountID:  account.ID,
			Category:   model.Spot,
			Symbol:     direct.Symbol,
			Direction:  direct.Direction,
			AmountFrom: sw.AmountFrom,
			Status:     model.OrderCreated,
			IsFirst:    true,
		}

		orders = append(orders, &ord)

		err = s.orderRep.Add(ctx, &ord)
		if err != nil {
			return nil, model.Unknown, fmt.Errorf("add 1: %w", err)
		}
	} else {
		for i := range indirect {
			var (
				// amountFrom второго ордера будет устанавливать по факту выполнения первого ордера из amountTo
				amountFrom decimal.Decimal
				isFirst    bool
			)

			if i == 0 {
				isFirst = true
				amountFrom = s.truncateAmountBySymbolRules(ctx, sw.AmountFrom, indirect[i].Symbol, indirect[i].Direction, true)
			}
			ord := model.Order{
				ID:         uuid.New(),
				SwapID:     sw.ID,
				AccountID:  account.ID,
				Category:   model.Spot,
				Symbol:     indirect[i].Symbol,
				Direction:  indirect[i].Direction,
				AmountFrom: amountFrom,
				IsFirst:    isFirst,
				Status:     model.OrderCreated,
			}

			orders = append(orders, &ord)

			err = s.orderRep.Add(ctx, &ord)
			if err != nil {
				return nil, model.Unknown, fmt.Errorf("add 2: %w", err)
			}
		}
	}

	if len(orders) == 0 {
		return nil, model.Unknown, fmt.Errorf("no orders to create")
	}

	err = s.swapStatusUpdater.UpdateAndBroadcast(ctx, sw, model.PlaceOrder)
	if err != nil {
		return nil, model.Unknown, fmt.Errorf("updateStatus: %w", err)
	}

	err = s.exchangeSubscriber.SubscribeOnOrders(ctx, account, orders, false)
	if err != nil {
		return nil, model.Unknown, fmt.Errorf("subscribeOnOrders: %w", err)
	}

	log.Debug(ctx, "create orders successfully: %s", sw.ID.String())

	return orders, model.PlaceOrder, nil
}

func (s *swap) PlaceOrder(ctx context.Context, sw *model.Swap, order *model.Order) (model.Status, error) {
	log.Debug(ctx, "try placeOrder: %+v", order)

	// Обновить AmountFrom второго ордера фактическим amountTo первого ордера
	if sw.Status == model.PlaceAdditionalOrder {
		err := s.updateAmountFromSecondOrder(ctx, sw.AccountFrom, order)
		if err != nil {
			return model.Unknown, fmt.Errorf("updateAmountFromSecondOrder: %w", err)
		}
	}

	// Проверяем accuracy перед размещением ордера, не по символу, а по coin, потому что монеты уйдут со спотового счёта,
	// а там accuracy обычно ниже
	var coinFrom string
	if order.IsFirst {
		coinFrom = sw.CoinFrom
	} else {
		coinFrom = model.CoinUSDT
	}
	amountFrom, err := s.truncateByCoinAccuracy(ctx, coinFrom, order.AmountFrom)
	if err != nil {
		return 0, fmt.Errorf("truncateByCoinAccuracy: %w", err)
	}
	order.AmountFrom = amountFrom

	// TODO Доп проверка, хорошо бы убрать
	balance, err := s.getBalanceByCoin(ctx, sw.ID, int(sw.AccountFrom), coinFrom, model.UNIFIED)
	if err != nil {
		return model.Unknown, fmt.Errorf("getBalanceByCoin: %w", err)
	}
	if balance.IsZero() {
		return sw.Status, nil
	}
	if !balance.Equal(order.AmountFrom) {
		log.Warn(ctx, "placeOrder: balance not equal: real: %s, expected: %s", balance, order.AmountFrom)
		if order.AmountFrom.GreaterThan(balance) {
			order.AmountFrom = balance
		}
		if order.AmountFrom.IsZero() {
			order.AmountFrom = balance
		}
	}

	account, err := s.accountRep.FindOne(ctx, &model.AccountFilter{
		ID: &sw.AccountFrom,
	})
	if err != nil {
		return model.Unknown, fmt.Errorf("findOne: %w", err)
	}

	err = s.market.PlaceOrder(ctx, order, account.Keys)
	if err != nil {
		return model.Unknown, fmt.Errorf("placeOrder: %w", err)
	}

	newOrderStatus := model.OrderPending
	err = s.orderRep.Update(ctx, order,
		&model.OrderFilter{
			ID: &order.ID,
		},
		&model.OrderPartial{
			Status: &newOrderStatus,
		})
	if err != nil {
		return 0, fmt.Errorf("update: %w", err)
	}

	var nextStep model.Status
	if order.IsFirst {
		nextStep = model.CheckOrder
	} else {
		nextStep = model.CheckAdditionalOrder
	}

	err = s.swapStatusUpdater.UpdateAndBroadcast(ctx, sw, nextStep)
	if err != nil {
		return model.Unknown, fmt.Errorf("updateStatus: swap_id: %s, status: %d, err: %s", sw.ID, nextStep, err)
	}

	log.Debug(ctx, "order placed: order: %+v", order)

	return nextStep, nil
}

func (s *swap) CheckOrder(ctx context.Context, sw *model.Swap, orderID uuid.UUID) (model.Status, error) {
	// Находим нужный ордер
	ord, err := s.orderRep.FindOne(ctx, &model.OrderFilter{
		ID: &orderID,
	})
	if err != nil {
		return model.Unknown, fmt.Errorf("orderRep.FindOne: %w", err)
	}

	log.Debug(ctx, "try check order: %+v", ord)

	// Ждём выполнения ордера
	if ord.Status == model.OrderPending {
		account, err := s.accountRep.FindOne(ctx, &model.AccountFilter{
			ID: &sw.AccountFrom,
		})
		if err != nil {
			return model.Unknown, fmt.Errorf("accountRep.FindOne: %w", err)
		}

		status, err := s.market.GetOrderStatus(ctx, ord.ID, account.Keys.ApiKey, account.Keys.ApiSecret)
		if err != nil {
			return model.Unknown, fmt.Errorf("getOrderStatus: %w", err)
		}

		if status == model.OrderPending || status == model.OrderUnknown {
			return sw.Status, nil
		} else {
			err = s.orderRep.Update(ctx, ord,
				&model.OrderFilter{
					ID: &ord.ID,
				},
				&model.OrderPartial{
					Status: &status,
				})
			if err != nil {
				return 0, fmt.Errorf("update: %w", err)
			}
			ord.Status = status
			log.Warn(ctx, "order status updated after check history: %+v", ord)
		}

	} else if ord.Status == model.OrderFailed {
		return model.Unknown, fmt.Errorf("order failed: %+v", ord)
	} else if ord.Status != model.OrderFilled && ord.Status != model.OrderPartiallyFilled {
		return model.Unknown, fmt.Errorf("unexpected order status: %+v", ord)
	}

	log.Debug(ctx, "order is success: %+v", ord)

	//// Обновляем amountTo в базе после выполнения ордера
	//amountTo := s.truncateAmountBySymbolRules(ctx, ord.AmountTo, ord.Symbol, ord.Direction, false)
	//err = s.orderRep.Update(ctx, ord,
	//	&model.OrderFilter{
	//		ID: &ord.ID,
	//	},
	//	&model.OrderPartial{
	//		AmountTo: &amountTo,
	//	})
	//if err != nil {
	//	return 0, fmt.Errorf("update: %w", err)
	//}
	//log.Debug(ctx, "updateOrder: %+v", ord)

	// Определяем следующий шаг, обновляем базу, шлём событие

	orders, err := s.orderRep.Find(ctx, &model.OrderFilter{
		AccountID: &sw.AccountFrom,
	})
	if err != nil {
		return 0, fmt.Errorf("find: %w", err)
	}
	var nextStep model.Status
	switch len(orders) {
	case 1:
		nextStep = model.TransferFromSubToMaster
	case 2:
		firstOrder, err := orders.FindFirst()
		if err != nil {
			return model.Unknown, fmt.Errorf("findFirst: %w", err)
		}
		if firstOrder.ID == orderID {
			nextStep = model.PlaceAdditionalOrder
		} else {
			nextStep = model.TransferFromSubToMaster
		}
	}
	err = s.swapStatusUpdater.UpdateAndBroadcast(ctx, sw, nextStep)
	if err != nil {
		return 0, fmt.Errorf("updateStatus: %w", err)
	}

	// Пыль на master account, ожидаем только при покупке, может остаться котируемая валюта
	if ord.Status == model.OrderPartiallyFilled {
		switch ord.IsFirst {
		case true:
			quoteCoin := sw.CoinFrom

			err = s.transferDustToMaster(ctx, ord.AmountFrom, sw.AmountFrom, quoteCoin, sw.AccountFrom)
			if err != nil {
				// без return, ошибка здесь повлияет только на удаление sub account
				log.Error(ctx, "checkOrder: transferDustToMaster: %s", err.Error())
			}

		case false:
			quoteCoin := model.CoinUSDT
			firstOrder := orders[0]

			err = s.transferDustToMaster(ctx, ord.AmountFrom, firstOrder.AmountTo, quoteCoin, sw.AccountFrom)
			if err != nil {
				// без return, ошибка здесь повлияет только на удаление sub account
				log.Error(ctx, "checkOrder: transferDustToMaster: %s", err.Error())
			}
		}
	}

	log.Debug(ctx, "order is success: %s, next step: %d", orderID, nextStep)

	return nextStep, nil
}

// updateAmountFromSecondOrder обновляет amountFrom второго ордера, значение amountTo первого ордера
func (s *swap) updateAmountFromSecondOrder(ctx context.Context, accountFrom int64, secondOrder *model.Order) error {
	// Ищем первый ордер
	findIsFirst := true
	firstOrder, err := s.orderRep.FindOne(ctx, &model.OrderFilter{
		AccountID: &accountFrom,
		IsFirst:   &findIsFirst,
	})
	if err != nil {
		return fmt.Errorf("findOne: %w", err)
	}

	// Обновляем второй ордер
	amountFrom := s.truncateAmountBySymbolRules(ctx, firstOrder.AmountTo, secondOrder.Symbol, secondOrder.Direction, true)
	err = s.orderRep.Update(ctx, secondOrder,
		&model.OrderFilter{
			ID: &secondOrder.ID,
		}, &model.OrderPartial{
			AmountFrom: &amountFrom,
		})
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}

	log.Debug(ctx, "update second order: %+v", secondOrder)
	return nil
}

func (s *swap) transferDustToMaster(
	ctx context.Context,
	realAmountFrom,
	expectedAmountFrom decimal.Decimal,
	coin string,
	accountID int64) error {
	//if realAmountFrom.GreaterThanOrEqual(expectedAmountFrom) {
	//	return nil
	//}

	myCoin, err := s.coinRep.Get(ctx, coin)
	if err != nil {
		return fmt.Errorf("getCoin: %w", err)
	}
	if myCoin.Accuracy == 0 {
		myCoin.Accuracy = model.DefaultAccuracy
	}

	dust := expectedAmountFrom.Sub(realAmountFrom).Truncate(int32(myCoin.Accuracy))

	log.Debug(ctx, "transferDustToMaster: expAmount: %s, realAmount: %s, dust: %s, coin: %s, accountID: %d",
		expectedAmountFrom.String(), realAmountFrom, dust, coin, accountID)

	err = s.transferFromSubToMaster(ctx, &model.InternalTransfer{
		ID:            uuid.New(),
		Coin:          coin,
		Amount:        dust,
		FromAccountID: accountID,
	})
	if err != nil {
		return fmt.Errorf("transferDustToMaster: %s", err.Error())
	}

	return nil
}

func (s *swap) TransferFromSubToMaster(ctx context.Context, swap *model.Swap, transfer *model.InternalTransfer) (model.Status, error) {
	log.Debug(ctx, "try TransferFromSubToMaster: %+v", transfer)

	// TODO доп проверка, хорошо бы убрать
	balance, err := s.getBalanceByCoin(ctx, swap.ID, int(transfer.FromAccountID), transfer.Coin, model.UNIFIED)
	if err != nil {
		return 0, fmt.Errorf("getBalanceByCoin: %w", err)
	}
	if balance.IsZero() {
		return model.TransferFromSubToMaster, nil
	}
	if !balance.Equal(transfer.Amount) {
		log.Warn(ctx, "transferFromSubToMaster: balance not equal: real: %s, expected: %s",
			balance, transfer.Amount)
		transfer.Amount = balance
	}

	log.Debug(ctx, "swap_id: %s: transferFromSubToMaster: balance: %s", swap.ID.String(), balance.String())

	amount, err := s.truncateByCoinAccuracy(ctx, transfer.Coin, transfer.Amount)
	if err != nil {
		return model.Unknown, fmt.Errorf("truncateByCoinAccuracy: %w", err)
	}
	transfer.Amount = amount

	err = s.transferFromSubToMaster(ctx, transfer)
	if err != nil {
		return model.Unknown, fmt.Errorf("transferFromSubToMaster: %w", err)
	}

	var nextStep model.Status
	if transfer.Status == model.ItsSuccess {
		nextStep = model.PrepareWithdraw
	} else {
		nextStep = model.CheckTransferFromSubToMaster
	}

	err = s.swapStatusUpdater.UpdateAndBroadcast(ctx, swap, nextStep)
	if err != nil {
		return 0, fmt.Errorf("updateStatus: %w", err)
	}

	return nextStep, nil
}

func (s *swap) transferFromSubToMaster(ctx context.Context, transfer *model.InternalTransfer) error {
	status, err := s.exchangeTransaction.TransferFromSubToMaster(ctx, transfer)
	if err != nil {
		return fmt.Errorf("transferFromSubToMaster: %w", err)
	}
	transfer.Status = status
	transfer.UpdatedAt = time.Now().UTC()

	log.Debug(ctx, "transferFromSubToMaster successfully: %+v", transfer)

	err = s.transferRep.Add(ctx, transfer)
	if err != nil {
		return fmt.Errorf("addInternalTransfer: %w", err)
	}

	log.Debug(ctx, "addInternalTransfer successfully: %+v", transfer)

	return nil
}

func (s *swap) CheckTransferFromSubToMaster(ctx context.Context, swap *model.Swap, accountID int64) (model.Status, error) {
	isLast := true
	tr, err := s.transferRep.FindOne(ctx, &model.InternalTransferFilter{
		FromAccountID: &accountID,
		IsLast:        &isLast,
	})
	if err != nil {
		return model.Unknown, fmt.Errorf("findOne: %w", err)
	}

	t, err := s.exchangeTransaction.GetTransfer(ctx, tr.ID)
	if err != nil {
		return model.Unknown, fmt.Errorf("getTransfer: %w", err)
	}

	err = s.transferRep.Update(ctx, tr,
		&model.InternalTransferFilter{
			ID: &tr.ID,
		},
		&model.InternalTransferPartial{
			Status: &t.Status,
		})
	if err != nil {
		return model.Unknown, fmt.Errorf("updateInternalTransferStatus: %w", err)
	}

	var nextStep model.Status

	switch t.Status {
	case model.ItsSuccess:
		nextStep = model.PrepareWithdraw
	case model.ItsPending:
		// Пробуем ещё раз
		nextStep = model.CheckTransferFromSubToMaster
		log.Debug(ctx, "internal transfer is pending: %s", tr.ID.String())
	case model.ItsFailed:
		// Возвращаем на предыдущий шаг
		nextStep = model.TransferFromSubToMaster
		log.Debug(ctx, "internal transfer is failed: %s", tr.ID.String())
	default:
		nextStep = model.Unknown
		log.Debug(ctx, "internal transfer is unknown status: %s, %s", tr.ID.String(), t.Status)
	}

	if nextStep == model.PrepareWithdraw || nextStep == model.TransferFromSubToMaster {
		err = s.swapStatusUpdater.UpdateAndBroadcast(ctx, swap, nextStep)
		if err != nil {
			return 0, fmt.Errorf("updateStatus: %w", err)
		}
	}

	return nextStep, nil
}

func (s *swap) PrepareWithdraw(ctx context.Context, sw *model.Swap) (model.Status, error) {
	// Ищем трансфер на main аккаунт, что бы узнать сумму для вывода
	isLast := true
	tr, err := s.transferRep.FindOne(ctx, &model.InternalTransferFilter{
		FromAccountID: &sw.AccountFrom,
		IsLast:        &isLast,
	})
	if err != nil {
		return model.Unknown, fmt.Errorf("findOne: %w", err)
	}

	wth := &model.Withdraw{
		InternalID: uuid.New(),
		SwapID:     sw.ID,
		Coin:       sw.CoinTo,
		Network:    sw.NetworkTo,
		Address:    sw.AddressTo,
		Tag:        sw.TagTo,
		Amount:     tr.Amount,
	}

	// Проверяем достаточность баланса, ждём возможно монеты в пути
	withdrawalAmount, err := s.exchangeAccount.GetWithdrawalAmount(ctx, wth.Coin)
	if err != nil {
		return model.Unknown, fmt.Errorf("getWithdrawalAmount: %w", err)
	}
	if withdrawalAmount.LessThan(wth.Amount) {
		log.Debug(ctx, "withdraw: available amount is less than withdrawal amount: %s %s", withdrawalAmount.String(), wth.Amount.String())
		return model.PrepareWithdraw, nil
	}

	// Считаем комиссию;
	// Отнимаем нашу комиссию
	ourFee := calculateFee(wth.Amount, s.ourFee)
	wth.Amount = wth.Amount.Sub(ourFee)
	log.Debug(ctx, "swapID: %s, our fee: %s", sw.ID.String(), ourFee.String())

	// Отнимаем комиссию за вывод
	fee, err := s.calculateWithdrawFee(ctx, wth)
	if err != nil {
		return model.Unknown, fmt.Errorf("calculateWithdrawFee: %w", err)
	}
	log.Debug(ctx, "swapID: %s, withdraw fee: %s", sw.ID.String(), fee.String())
	if !fee.IsZero() {
		wth.Amount = wth.Amount.Sub(fee)
		wth.IncludeFeeInAmount = true
	} else {
		wth.IncludeFeeInAmount = false
	}

	wth.Amount, err = s.truncateByWithdrawAccuracy(ctx, wth.Coin, wth.Network, wth.Amount)
	if err != nil {
		return model.Unknown, fmt.Errorf("truncateByWithdrawAccuracy: %w", err)
	}
	log.Debug(ctx, "swapID: %s, withdraw amount after fee: %s", sw.ID.String(), wth.Amount.String())

	// Обновляем базу
	err = s.withdrawRep.WithinTransaction(ctx, func(ctx context.Context) error {
		err = s.withdrawRep.Add(ctx, wth)
		if err != nil {
			return fmt.Errorf("add: %w", err)
		}
		err = s.swapRep.Update(ctx, sw,
			&model.SwapFilter{
				ID: &sw.ID,
			},
			&model.SwapPartial{
				AmountTo: &wth.Amount,
			})
		if err != nil {
			return fmt.Errorf("swapRep.Update: %w", err)
		}
		return nil
	})
	if err != nil {
		return model.Unknown, err
	}

	err = s.swapStatusUpdater.UpdateAndBroadcast(ctx, sw, model.WithdrawSwapStatus)
	if err != nil {
		return model.Unknown, fmt.Errorf("updateStatus: %w", err)
	}

	return model.WithdrawSwapStatus, nil
}

func (s *swap) Withdraw(ctx context.Context, swap *model.Swap) (model.Status, error) {
	wth, err := s.withdrawRep.FindOne(ctx,
		&model.WithdrawFilter{
			SwapID: &swap.ID,
		})
	if err != nil {
		return model.Unknown, fmt.Errorf("findOne: %w", err)
	}

	log.Debug(ctx, "try withdraw: %+v", wth)

	withdrawID, err := s.exchangeTransaction.Withdraw(ctx, wth)
	if err != nil {
		return model.Unknown, fmt.Errorf("withdraw: %w", err)
	}

	err = s.withdrawRep.Update(ctx, wth,
		&model.WithdrawFilter{
			SwapID: &swap.ID,
		},
		&model.WithdrawPartial{
			ID: &withdrawID,
		})
	if err != nil {
		return model.Unknown, fmt.Errorf("update: %w", err)
	}

	log.Debug(ctx, "withdraw successfully: %+v", wth)

	err = s.swapStatusUpdater.UpdateAndBroadcast(ctx, swap, model.WaitWithdraw)
	if err != nil {
		return model.Unknown, fmt.Errorf("updateStatus: %w", err)
	}

	return model.WaitWithdraw, nil
}

func (s *swap) CheckWithdraw(ctx context.Context, swap *model.Swap, withdrawID int64) (model.Status, error) {
	w, err := s.exchangeTransaction.GetWithdraw(ctx, int(withdrawID))
	if err != nil {
		return model.Unknown, fmt.Errorf("getWithdraw: %w", err)
	}

	err = s.withdrawRep.Update(ctx, w,
		&model.WithdrawFilter{
			ID: &withdrawID,
		},
		&model.WithdrawPartial{
			Amount: &w.Amount,
			Status: &w.Status,
			HashID: &w.HashID,
		})
	if err != nil {
		return model.Unknown, fmt.Errorf("updateStatus: %w", err)
	}

	if w.Status == model.WsBlockchainConfirmed {
		return model.WaitWithdraw, nil
	}

	if w.Status == model.WsFailed {
		err = s.swapStatusUpdater.UpdateAndBroadcast(ctx, swap, model.Error)
		if err != nil {
			return model.Unknown, fmt.Errorf("updateAndBroadcast: %w", err)
		}
		return model.Error, nil
	}

	if w.Status != model.WsSuccess {
		return model.WaitWithdraw, nil
	} else if w.Status == model.WsSuccess {
		err = s.swapStatusUpdater.UpdateAndBroadcast(ctx, swap, model.Completed)
		if err != nil {
			return model.Unknown, fmt.Errorf("updateAndBroadcast: %w", err)
		}
		log.Debug(ctx, "checkWithdraw: update status success swapID: %s", swap.ID)

		user, err := s.userRep.FindOne(ctx, &model.UserFilter{
			ID: &swap.UserID,
		})
		if err != nil {
			return model.Unknown, fmt.Errorf("findOne: %w", err)
		}

		s.emailCli.SendSuccessfulSwapMessage(ctx, swap, user)

		// Проверить пыль по USDT перед удалением sub account
		balance, err := s.exchangeAccount.GetBalanceByCoin(ctx, int(swap.AccountFrom), model.CoinUSDT, "SPOT")
		if err != nil {
			return model.Unknown, fmt.Errorf("getBalanceByCoin: %w", err)
		}
		if balance.TransferBalance.IsPositive() {
			log.Debug(ctx, "checkWithdraw: transfer dust to master amount: %s, sub acc: %d",
				balance.TransferBalance.String(),
				swap.AccountFrom,
			)
			err = s.transferFromSubToMaster(ctx, &model.InternalTransfer{
				ID:            uuid.New(),
				Coin:          model.CoinUSDT,
				Amount:        balance.TransferBalance,
				FromAccountID: swap.AccountFrom,
			})
			if err != nil {
				return model.Unknown, fmt.Errorf("transferFromSubToMaster: %w", err)
			}
		}

		// Удаляем sub account
		//account, err := s.accountRep.Find(ctx, accountID)
		//if err != nil {
		//	return model.Unknown, fmt.Errorf("getSubAccount: %w", err)
		//}
		//err = s.exchangeAccount.DeleteSubAPIKey(ctx, account.Keys.ApiKey)
		//if err != nil {
		//	return model.Unknown, fmt.Errorf("deleteSubAPIKey: %w", err)
		//}
		//err = s.exchangeAccount.DeleteSubAccount(ctx, accountID)
		//if err != nil {
		//	return model.Unknown, fmt.Errorf("deleteSubAccount: %w", err)
		//}
	}

	return model.Completed, nil
}

func (s *swap) GetAllActiveSwaps(ctx context.Context) ([]*model.Swap, error) {
	statusComplete := []model.Status{model.Unknown, model.Completed, model.ManualCompleted, model.Cancel, model.DepositError, model.Error}
	swaps, err := s.swapRep.Find(ctx, &model.SwapFilter{
		NotEqStatus: statusComplete,
	})
	if err != nil {
		return nil, fmt.Errorf("find: %w", err)
	}

	return swaps, nil
}

func (s *swap) checkOnPossibilityWithdraw(ctx context.Context, coin, network string) error {
	coinTo, err := s.coinRep.Get(ctx, coin)
	if err != nil {
		return fmt.Errorf("checkOnPossibilityWithdraw: %w", err)
	}
	if coinTo.Title != coin {
		return fmt.Errorf("checkOnPossibilityWithdraw coin: %s", coinTo.Title)
	}
	var isSupportedNetwork bool
	for i := range coinTo.Networks {
		if coinTo.Networks[i].Title == network {
			if !coinTo.Networks[i].WithdrawSupported {
				return businessError.NotSupportedWithdrawByNetwork
			}
			isSupportedNetwork = true
			break
		}
	}
	if isSupportedNetwork {
		return nil
	} else {
		return businessError.NotSupportedWithdrawByNetwork
	}
}

func (s *swap) getSwapPathAndCalculateAmount(ctx context.Context, swap *model.Swap) ([]*step, error) {
	var (
		reversed            bool
		isWithoutCommission bool
	)

	if swap.AmountFrom.IsZero() && !swap.AmountTo.IsZero() {
		reverseSwap(swap)
		reversed = true
		isWithoutCommission = true
	}

	toNetwork, err := s.coinRep.GetNetwork(ctx, swap.CoinTo, swap.NetworkTo)
	if err != nil {
		return nil, fmt.Errorf("get network %s for coin: %s: %w", swap.NetworkTo, swap.CoinTo, err)
	}

	direct, indirect, err := s.getSwapOptions(ctx, swap.CoinFrom, swap.CoinTo)
	if err != nil {
		return nil, fmt.Errorf("getSwapOptions: %w", err)
	}
	amountTo, isDirect, err := s.getAmountToAndTheBestOptionToSwap(ctx, direct, indirect, swap.AmountFrom, isWithoutCommission,
		toNetwork.WithdrawFee.Fee, toNetwork.WithdrawFee.PercentageFee)
	if err != nil {
		return nil, fmt.Errorf("getAmountToAndTheBestOptionToSwap: %w", err)
	}
	swap.AmountTo = amountTo

	if reversed {
		reverseSwap(swap)

		// Если amount был посчитан в обратном направлении, то для swap это не подходит, нужен не reverse путь
		direct, indirect, err = s.getSwapOptions(ctx, swap.CoinFrom, swap.CoinTo)
		if err != nil {
			return nil, fmt.Errorf("getSwapOptions: %w", err)
		}
		_, isDirect, err = s.getAmountToAndTheBestOptionToSwap(ctx, direct, indirect, swap.AmountFrom, isWithoutCommission,
			toNetwork.WithdrawFee.Fee, toNetwork.WithdrawFee.PercentageFee)
		if err != nil {
			return nil, fmt.Errorf("getAmountToAndTheBestOptionToSwap: %w", err)
		}
	}

	var bestPath []*step
	if isDirect {
		bestPath = append(bestPath, direct)
	} else {
		bestPath = indirect
	}

	return bestPath, nil
}

func reverseSwap(swap *model.Swap) {
	swap.AmountFrom, swap.AmountTo = swap.AmountTo, swap.AmountFrom
	swap.CoinFrom, swap.CoinTo = swap.CoinTo, swap.CoinFrom
	swap.NetworkFrom, swap.NetworkTo = swap.NetworkTo, swap.NetworkFrom
}

type step struct {
	Symbol    string
	Direction model.Direction
}

// getSwapOptions calculates possible options of swap
// returns direct swap, indirect swap
// indirect means, for example: ETHBTC buy -> BTCUSDT sell -> ETHUSDT buy
func (s *swap) getSwapOptions(ctx context.Context, from, to string) (*step, []*step, error) {
	// only direct swap
	if from == model.CoinUSDT || to == model.CoinUSDT {
		if from == model.CoinUSDT {
			sym := fmt.Sprintf("%s%s", to, model.CoinUSDT)
			exist := s.orderBookRep.IsExist(sym)
			if !exist {
				return nil, nil, fmt.Errorf("no symbol: %s, err: %w", sym, businessError.NoPathToSwapErr)
			}
			return &step{
				Symbol:    sym,
				Direction: model.Buy,
			}, nil, nil
		}

		sym := fmt.Sprintf("%s%s", from, model.CoinUSDT)
		exist := s.orderBookRep.IsExist(sym)
		if !exist {
			return nil, nil, fmt.Errorf("no symbol: %s, err: %w", sym, businessError.NoPathToSwapErr)
		}
		return &step{
			Symbol:    sym,
			Direction: model.Sell,
		}, nil, nil
	}

	var (
		direct   *step
		indirect []*step
	)

	// from -> to direct swap
	exist := s.orderBookRep.IsExist(fmt.Sprintf("%s%s", from, to))
	if exist {
		direct = &step{
			Symbol:    fmt.Sprintf("%s%s", from, to),
			Direction: model.Sell,
		}
	} else {
		// from -> to direct swap
		exist = s.orderBookRep.IsExist(fmt.Sprintf("%s%s", to, from))
		if exist {
			direct = &step{
				Symbol:    fmt.Sprintf("%s%s", to, from),
				Direction: model.Buy,
			}
		}
	}

	// from -> USDT
	sym := fmt.Sprintf("%s%s", from, model.CoinUSDT)
	exist = s.orderBookRep.IsExist(sym)
	if !exist {
		return nil, nil, fmt.Errorf("no symbol: %s, err: %w", sym, businessError.NoPathToSwapErr)
	}
	indirect = append(indirect, &step{
		Symbol:    sym,
		Direction: model.Sell,
	})

	// USDT -> to
	sym = fmt.Sprintf("%s%s", to, model.CoinUSDT)
	exist = s.orderBookRep.IsExist(sym)
	if !exist {
		return nil, nil, fmt.Errorf("no symbol: %s, err: %w", sym, businessError.NoPathToSwapErr)
	}
	indirect = append(indirect, &step{
		Symbol:    sym,
		Direction: model.Buy,
	})

	log.Debug(ctx, "direct path: %+v", direct)
	for _, i := range indirect {
		log.Debug(ctx, "direct path: %+v", i)
	}

	return direct, indirect, nil
}

func (s *swap) getAmountToAndTheBestOptionToSwap(
	ctx context.Context,
	direct *step,
	indirect []*step,
	amountFrom decimal.Decimal,
	isWithoutCommission bool,
	withdrawFee decimal.Decimal,
	withdrawPercentageFee decimal.Decimal,
) (amountTo decimal.Decimal, isDirect bool, err error) {
	var (
		amountToDirectly decimal.Decimal
	)
	if direct != nil {
		amountToDirectly, err = s.calculateAmountToByAmountFromDirectly(ctx, direct.Symbol, direct.Direction, amountFrom, withdrawFee, withdrawPercentageFee, isWithoutCommission)
		if err != nil {
			return decimal.Zero, false, fmt.Errorf("calculateAmountToByAmountFromDirectly: %w", err)
		}
	}

	var amountToIndirectly decimal.Decimal
	if indirect != nil {
		amountToIndirectly, err = s.calculateAmountToByAmountFromIndirectly(ctx, indirect, amountFrom, withdrawFee, withdrawPercentageFee, isWithoutCommission)
		if err != nil {
			return decimal.Zero, false, fmt.Errorf("calculateAmountToByAmountFromIndirectly: %w", err)
		}
	}

	// Сначала выбираем amount_to полагаясь на доступные пути. Не можем сразу сравнить amountToDirectly.GreaterThan(amountToIndirectly)
	// потому что могут быть отрицательные суммы.
	// Пример: indirect amount_to = -0.1, но direct = nil и amountToDirectly = 0
	if direct != nil && indirect == nil {
		return amountToDirectly, true, nil
	}
	if indirect != nil && direct == nil {
		return amountToIndirectly, false, nil
	}

	// Если доступны оба пути выбираем лучший
	isDirect = amountToDirectly.GreaterThan(amountToIndirectly)
	amountTo = amountToDirectly
	if !isDirect {
		amountTo = amountToIndirectly
	}

	return amountTo, isDirect, nil
}

func (s *swap) checkAndReturnMaxLimit(ctx context.Context, swap *model.Swap) (decimal.Decimal, error) {
	if swap.CoinFrom == model.CoinUSDT {
		if swap.AmountFrom.GreaterThan(s.maxLimit) {
			return decimal.Zero, businessError.MaxLimitExceededErr
		}
		return s.maxLimit, nil
	} else {
		maxLimit, err := s.calculateMaxLimit(ctx, swap)
		if err != nil {
			return decimal.Zero, fmt.Errorf("calculateMaxLimit: %w", err)
		}
		if !maxLimit.IsZero() && swap.AmountFrom.GreaterThanOrEqual(maxLimit) {
			return decimal.Zero, fmt.Errorf("calculated: %s, from: %s: err: %w", maxLimit.String(), swap.AmountFrom, businessError.MaxLimitExceededErr)
		}
		return maxLimit, nil
	}
}

func (s *swap) checkAndReturnMinLimit(ctx context.Context, swap *model.Swap, bestPath []*step) (decimal.Decimal, error) {
	minimum, err := s.calculateMin(ctx, swap.CoinFrom, swap.CoinTo, swap.NetworkFrom, swap.NetworkTo, bestPath)
	if err != nil {
		return decimal.Zero, fmt.Errorf("calculateMin: %w", err)
	}
	if !minimum.IsZero() && swap.AmountFrom.LessThan(minimum) {
		return decimal.Zero, fmt.Errorf("minimum: %s, amountFrom: %s, err: %w", minimum.String(), swap.AmountFrom.String(), businessError.BelowMinLimitErr)
	}
	return minimum, nil
}

// calculateMin ⌈max(min_deposit, (1+slippage)*min_withdraw/score)*(1+emcd_fee) + withdraw_fee/score, s.minLimit, min_order_qty*1.2⌉
// значение в монете from
func (s *swap) calculateMin(ctx context.Context, from, to, networkFrom, networkTo string, steps []*step) (decimal.Decimal, error) {
	// Получаем amountTo без учёта комиссий (иначе score будет посчитан некорректно)
	var (
		formalAmountTo decimal.Decimal
		err            error
	)
	formalAmountFrom := decimal.NewFromInt(1)
	if len(steps) == 1 {
		formalAmountTo, err = s.calculateAmountToByAmountFromDirectly(ctx, steps[0].Symbol, steps[0].Direction, formalAmountFrom,
			decimal.Zero, decimal.Zero, true)
		if err != nil {
			return decimal.Zero, fmt.Errorf("calculateAmountToByAmountFromDirectly: %w", err)
		}
	} else if len(steps) == 2 {
		formalAmountTo, err = s.calculateAmountToByAmountFromIndirectly(ctx, steps, formalAmountFrom, decimal.Zero, decimal.Zero,
			true)
		if err != nil {
			return decimal.Zero, fmt.Errorf("calculateAmountToByAmountFromIndirectly: %w", err)
		}
	}
	score := formalAmountTo.Div(formalAmountFrom)
	if score.IsZero() {
		return decimal.Zero, fmt.Errorf("score is zero")
	}
	log.Debug(ctx, "from: %s, to: %s, score: %s", from, to, score.String())

	network, err := s.coinRep.GetNetwork(ctx, from, networkFrom)
	if err != nil {
		return decimal.Zero, fmt.Errorf("getNetwork 1: %w", err)
	}
	minDeposit := network.DepositMin

	network, err = s.coinRep.GetNetwork(ctx, to, networkTo)
	if err != nil {
		return decimal.Zero, fmt.Errorf("getNetwork 2: %w", err)
	}
	minWithdraw := network.WithdrawMin

	var slippage decimal.Decimal
	if from == model.CoinBTC || to == model.CoinBTC {
		slippage = decimal.NewFromFloat(model.SlippageBTC)
	} else {
		slippage = decimal.NewFromFloat(model.SlippageDefault)
	}

	settingsFee, err := s.coinRep.GetWithdrawFee(ctx, to, networkTo)
	if err != nil {
		return decimal.Zero, fmt.Errorf("getWithdrawFee: %w", err)
	}
	withdrawFee := calculateByBitWithdrawFee(formalAmountTo, settingsFee.Fee, settingsFee.PercentageFee)

	var minimum decimal.Decimal

	// (1+slippage)*min_withdraw/score
	minWithdrawInBaseCoin := decimal.NewFromInt(1).Add(slippage).Mul(minWithdraw).Div(score)

	if minDeposit.GreaterThan(minWithdrawInBaseCoin) {
		minimum = minDeposit
	} else {
		minimum = minWithdrawInBaseCoin
	}

	// Добавляем нашу комиссию
	minimum = minimum.Mul(decimal.NewFromInt(1).Add(s.ourFee))

	// Добавляем комиссию за вывод, конвертируемую в монету депозита
	minimum = minimum.Add(withdrawFee.Div(score))

	accuracy, err := s.coinRep.GetAccuracyForWithdrawAndDeposit(ctx, to, networkTo)
	if err != nil {
		return decimal.Zero, fmt.Errorf("getAccuracy: %w", err)
	}
	if accuracy == 0 {
		accuracy = model.DefaultAccuracy
	}

	minimum = minimum.Mul(decimal.NewFromFloat(model.LimitCoefficient)).Truncate(int32(accuracy))

	// Проверяем, если минимальная сумма from меньше, чем эквивалент USDT, то устанавливаем фиксированный minLimit
	if from == model.CoinUSDT {
		if minimum.LessThan(s.minLimit) {
			minimum = s.minLimit
		}
	} else {
		// Конвертируем minLimit USDT в базовую валюту через score
		direct, _, err := s.getSwapOptions(ctx, from, model.CoinUSDT)
		if err != nil {
			return decimal.Decimal{}, fmt.Errorf("getSwapOptions: %w", err)
		}
		usdtScore, err := s.calculateAmountToByAmountFromDirectly(ctx, direct.Symbol, direct.Direction, decimal.NewFromInt(1), decimal.Zero, decimal.Zero, true)
		if err != nil {
			return decimal.Zero, fmt.Errorf("calculateAmountToByAmountFromDirectly for USDT: %w", err)
		}

		minUsdtInBaseCoin := s.minLimit.Div(usdtScore)

		// Сравниваем минимум с конвертированной суммой
		if minimum.LessThan(minUsdtInBaseCoin) {
			minimum = minUsdtInBaseCoin
		}
	}

	// Проверяем на минимальную сумму размещения ордера
	// Например: TON --> INJ
	// 		order 1: TON --> USDT sell, min по TON
	//      order 2: USDT --> INJ buy, min по USDT
	for i := range steps {
		sym, err := s.symbolRep.Get(ctx, steps[i].Symbol)
		if err != nil {
			return decimal.Decimal{}, fmt.Errorf("getSymbol %s: %w", steps[i].Symbol, err)
		}
		log.Debug(ctx, "calculateMin: steps: %+v", sym)

		// При продаже указываем количество в qty
		// sell - только первый шаг
		if !sym.MinOrderQty.IsZero() && sym.BaseCoin == from && steps[i].Direction == model.Sell {
			if minimum.LessThan(sym.MinOrderQty) {
				minimum = sym.MinOrderQty.Mul(decimal.NewFromFloat(1.2))
			}
		}

		// При покупке указываем количество в amt
		// buy - второй шаг; или первый шаг в прямом свопе
		if !sym.MinOrderAmt.IsZero() && sym.BaseCoin == to && steps[i].Direction == model.Buy {
			var direct *step
			// прямой своп, пример: USDT --> TON, QuoteCoin = from = USDT
			if sym.QuoteCoin == from {
				if minimum.LessThan(sym.MinOrderAmt) {
					minimum = sym.MinOrderAmt.Mul(decimal.NewFromFloat(1.2))
					continue
				}
			} else {
				direct, _, err = s.getSwapOptions(ctx, sym.QuoteCoin, from)
				if err != nil {
					return decimal.Zero, fmt.Errorf("getSwapOptions: %w", err)
				}
			}
			if direct == nil {
				log.Warn(ctx, "calculateMin: no direct path: %s --> %s", sym.QuoteCoin, to)
				continue
			}
			am, err := s.calculateAmountToByAmountFromDirectly(ctx, direct.Symbol, direct.Direction, sym.MinOrderAmt, decimal.Zero, decimal.Zero, true)
			if err != nil {
				return decimal.Zero, fmt.Errorf("calculateAmountToByAmountFromDirectly: %w", err)
			}
			if minimum.LessThan(am) {
				minimum = am.Mul(decimal.NewFromFloat(1.2))
			}
		}
	}

	return minimum, nil
}

func (s *swap) calculateAmountToByAmountFromDirectly(
	ctx context.Context,
	symbol string,
	direction model.Direction,
	amountFrom,
	withdrawFee,
	withdrawPercentageFee decimal.Decimal,
	isWithoutCommission bool,
) (decimal.Decimal, error) {
	return s.calculateAmountToByAmountFrom(ctx, symbol, direction, amountFrom, withdrawFee, withdrawPercentageFee, isWithoutCommission)
}

func (s *swap) calculateAmountToByAmountFromIndirectly(
	ctx context.Context,
	steps []*step,
	amountFrom decimal.Decimal,
	withdrawFee,
	withdrawPercentageFee decimal.Decimal,
	isWithoutCommission bool,
) (decimal.Decimal, error) {
	var (
		amountStableCoin decimal.Decimal
		amountTo         decimal.Decimal
		err              error
	)

	if len(steps) != 2 {
		for i := range steps {
			log.Debug(ctx, "step: %+v", steps[i])
		}
		return decimal.Zero, fmt.Errorf("len(steps) != 2: %w", businessError.CalculateSwapErr)
	}

	for i := range steps {
		if i == 0 {
			// first step
			amountStableCoin, err = s.calculateAmountToByAmountFrom(ctx, steps[i].Symbol, steps[i].Direction, amountFrom, decimal.Zero, decimal.Zero, isWithoutCommission)
			if err != nil {
				return decimal.Zero, fmt.Errorf("calculateAmountToByAmountFrom 1: %w", err)
			}
		} else {
			// second step
			amountTo, err = s.calculateAmountToByAmountFrom(ctx, steps[i].Symbol, steps[i].Direction, amountStableCoin, withdrawFee, withdrawPercentageFee, isWithoutCommission)
			if err != nil {
				return decimal.Zero, fmt.Errorf("calculateAmountToByAmountFrom 2: %w", err)
			}
		}
	}
	return amountTo, nil
}

// calculateAmountToByAmountFrom
// покупка: мы покупаем по Ask цене, так как это цена, по которой продавцы готовы продать.
// продажа: мы продаём по Bid цене, так как это цена, по которой покупатели готовы купить.
func (s *swap) calculateAmountToByAmountFrom(
	ctx context.Context,
	symbol string,
	direction model.Direction,
	amountFrom,
	withdrawFee,
	withdrawPercentageFee decimal.Decimal,
	isWithoutCommissions bool,
) (decimal.Decimal, error) {
	var (
		amountTo decimal.Decimal
		err      error
	)
	if direction == model.Buy {
		// Excluding commission
		amountTo, err = s.calculateAmountToDirectionBuy(symbol, amountFrom, levelsRequestedPerTime)
		if err != nil {
			return decimal.Zero, fmt.Errorf("calculateAmountToDirectionBuy: %w", err)
		}
	} else {
		// Excluding commission
		amountTo, err = s.calculateAmountToDirectionSell(symbol, amountFrom, levelsRequestedPerTime)
		if err != nil {
			return decimal.Zero, fmt.Errorf("calculateAmountToDirectionSell: %w", err)
		}
	}

	log.Debug(ctx, "estimate: amountTo excluding commission: %s, amountFrom: %s", amountTo.String(), amountFrom.String())

	if isWithoutCommissions {
		return amountTo, nil
	}

	fee, err := s.feeRep.GetFee(ctx, symbol)
	if err != nil {
		return decimal.Zero, fmt.Errorf("getFee: %w", err)
	}

	// Including bybit order commission
	calculatedOrderFee := calculateFee(amountTo, fee.TakerFee)
	amountTo = amountTo.Sub(calculatedOrderFee)

	// Including our commission
	amountTo = amountTo.Mul(decimal.NewFromInt(1).Sub(s.ourFee))

	// Including withdraw commission
	calculatedWithdrawFee := calculateByBitWithdrawFee(amountTo, withdrawFee, withdrawPercentageFee)
	amountTo = amountTo.Sub(calculatedWithdrawFee)
	log.Debug(ctx, "estimate: amountTo including commission: %s", amountTo.String())

	return amountTo, nil
}

const levelsRequestedPerTime = 3

// calculateAmountToDirectionBuy в направлении buy amountFrom - котируемая валюта (BTCUSDT, amountFrom - USDT)
func (s *swap) calculateAmountToDirectionBuy(symbol string, amountFrom decimal.Decimal, levelsRequested int) (decimal.Decimal, error) {
	levels, err := s.orderBookRep.GetAskTopLevels(symbol, levelsRequested)
	if err != nil {
		return decimal.Zero, fmt.Errorf("getAskTopLevels: %w", err)
	}
	if len(levels) == 0 {
		return decimal.Zero, fmt.Errorf("getAskTopLevels: levels is nil: %w", businessError.CalculateSwapErr)
	}

	// levels TONUSDT [0] price - 5.29, [1] size - 29.69

	var (
		totalAmountTo        decimal.Decimal
		remainingQuoteAmount = amountFrom
	)

	for _, level := range levels {
		price := decimal.NewFromFloat(level[0])
		size := decimal.NewFromFloat(level[1])

		cost := size.Mul(price)

		if remainingQuoteAmount.LessThanOrEqual(cost) {
			// Если оставшейся суммы USDT достаточно для покупки, рассчитываем объём базовой валюты
			return totalAmountTo.Add(remainingQuoteAmount.Div(price)), nil
		}

		// Если текущего объёма недостаточно, покупаем весь объём на уровне
		totalAmountTo = totalAmountTo.Add(size)
		remainingQuoteAmount = remainingQuoteAmount.Sub(cost)
	}

	// Если не хватило уровней, запрашиваем больше
	if len(levels) < levelsRequested {
		return decimal.Zero, fmt.Errorf("levels is not equal %d: %w", len(levels), businessError.MarketDepthExceededErr)
	} else {
		return s.calculateAmountToDirectionBuy(symbol, amountFrom, levelsRequested+levelsRequestedPerTime)
	}
}

// calculateAmountToDirectionSell в направлении sell amountFrom - базовая валюта (BTCUSDT, amountFrom - BTC)
func (s *swap) calculateAmountToDirectionSell(symbol string, amountFrom decimal.Decimal, levelsRequested int) (decimal.Decimal, error) {
	levels, err := s.orderBookRep.GetBidTopLevels(symbol, levelsRequested)
	if err != nil {
		return decimal.Zero, fmt.Errorf("getbidTopLevels: %w", err)
	}
	if len(levels) == 0 {
		return decimal.Zero, fmt.Errorf("getBidTopLevels: levels is nil: %w", businessError.CalculateSwapErr)
	}

	// levels TONUSDT [0] price - 5.29, [1] size - 29.69

	var (
		totalAmountTo       decimal.Decimal
		remainingBaseAmount = amountFrom
	)

	for _, level := range levels {
		price := decimal.NewFromFloat(level[0])
		size := decimal.NewFromFloat(level[1])

		if remainingBaseAmount.LessThanOrEqual(size) {
			// Если оставшийся объём можно продать на этом уровне
			return totalAmountTo.Add(remainingBaseAmount.Mul(price)), nil
		}

		// Если объёма на уровне недостаточно, продаём весь объём на уровне
		totalAmountTo = totalAmountTo.Add(price.Mul(size))
		remainingBaseAmount = remainingBaseAmount.Sub(size)
	}

	// Если не хватило уровней, запрашиваем больше
	if len(levels) < levelsRequested {
		return decimal.Zero, fmt.Errorf("not enough levels: %d: %w", len(levels), businessError.MarketDepthExceededErr)
	} else {
		return s.calculateAmountToDirectionSell(symbol, amountFrom, levelsRequested+levelsRequestedPerTime)
	}
}

// calculateWithdrawFee считает комиссию за вывод
func (s *swap) calculateWithdrawFee(ctx context.Context, wth *model.Withdraw) (decimal.Decimal, error) {
	withdrawF, err := s.coinRep.GetWithdrawFee(ctx, wth.Coin, wth.Network)
	if err != nil {
		return decimal.Zero, fmt.Errorf("getWithdrawFee: %w", err)
	}
	return calculateByBitWithdrawFee(wth.Amount, withdrawF.Fee, withdrawF.PercentageFee), nil
}

// calculateByBitWithdrawFee
// if withdrawPercentageFee != 0: handlingFee = inputAmount / (1 - withdrawPercentageFee) * withdrawPercentageFee + withdrawFee
// if withdrawPercentageFee = 0: handlingFee = withdrawFee
func calculateByBitWithdrawFee(amount, withdrawFee, withdrawPercentageFee decimal.Decimal) decimal.Decimal {
	if withdrawPercentageFee.IsZero() {
		return withdrawFee
	}

	return amount.Div(decimal.NewFromInt(1).Sub(withdrawPercentageFee)).Mul(withdrawPercentageFee).Add(withdrawFee)
}

func calculateFee(amount, fee decimal.Decimal) decimal.Decimal {
	return amount.Mul(fee)
}

func (s *swap) truncateAmountBySymbolRules(ctx context.Context, amount decimal.Decimal, symbol string, direction model.Direction,
	isAmountFrom bool) decimal.Decimal {
	var baseAccuracy, quoteAccuracy int32
	sym, err := s.symbolRep.Get(ctx, symbol)
	if err != nil {
		log.Error(ctx, "truncateAmountBySymbolRules: symbol %s: %w", symbol, err)
	} else {
		baseAccuracy = sym.Accuracy.BaseAccuracy
		quoteAccuracy = sym.Accuracy.QuoteAccuracy
	}

	if baseAccuracy == 0 {
		baseAccuracy = model.DefaultAccuracy
	}
	if quoteAccuracy == 0 {
		quoteAccuracy = model.DefaultAccuracy
	}

	switch isAmountFrom {
	case true:
		if direction == model.Sell {
			amount = amount.Truncate(baseAccuracy)
		} else {
			amount = amount.Truncate(quoteAccuracy)
		}
	case false:
		if direction == model.Sell {
			amount = amount.Truncate(quoteAccuracy)
		} else {
			amount = amount.Truncate(baseAccuracy)
		}
	}

	log.Debug(ctx, "truncateAmountBySymbolRules: symbol: %s, baseAccuracy: %d, quoteAccuracy: %d, amount: %s",
		symbol, baseAccuracy, quoteAccuracy, amount)

	return amount
}

func (s *swap) truncateByCoinAccuracy(ctx context.Context, coin string, amount decimal.Decimal) (decimal.Decimal, error) {
	myCoin, err := s.coinRep.Get(ctx, coin)
	if err != nil {
		return decimal.Zero, fmt.Errorf("get: %w", err)
	}
	var accuracy int
	if myCoin.Accuracy == 0 {
		accuracy = model.DefaultAccuracy
	} else {
		accuracy = myCoin.Accuracy
	}
	return amount.Truncate(int32(accuracy)), nil
}

func (s *swap) truncateByWithdrawAccuracy(ctx context.Context, coin, network string, amount decimal.Decimal) (decimal.Decimal, error) {
	accuracy, err := s.coinRep.GetAccuracyForWithdrawAndDeposit(ctx, coin, network)
	if err != nil {
		log.Error(ctx, "getAccuracyForWithdrawAndDeposit: %w", err)
	}
	if accuracy == 0 {
		accuracy = model.DefaultWithdrawAccuracy
	}
	return amount.Truncate(int32(accuracy)), nil
}

// calculateMaxLimit вычисляет максимальную сумму в монете, которую собираемся обменять.
// За исходную монету берём USDT и максимальное количество из конфигурации.
// За целевую монету берём ту монету, которую собираемся обменять.
func (s *swap) calculateMaxLimit(ctx context.Context, swap *model.Swap) (decimal.Decimal, error) {
	if swap.CoinFrom == model.CoinUSDT {
		return s.maxLimit, nil
	}

	coinTo := swap.CoinFrom
	networkTo := swap.NetworkFrom
	amountFrom := s.maxLimit

	toNetwork, err := s.coinRep.GetNetwork(ctx, coinTo, networkTo)
	if err != nil {
		return decimal.Zero, fmt.Errorf("get network %s for coin: %s: %w", networkTo, coinTo, err)
	}

	direct, indirect, err := s.getSwapOptions(ctx, model.CoinUSDT, coinTo)
	if err != nil {
		return decimal.Zero, fmt.Errorf("getSwapOptions: %w", err)
	}

	amountTo, _, err := s.getAmountToAndTheBestOptionToSwap(ctx, direct, indirect, amountFrom, false,
		toNetwork.WithdrawFee.Fee, toNetwork.WithdrawFee.PercentageFee)
	if err != nil {
		return decimal.Zero, fmt.Errorf("getAmountToAndTheBestOptionToSwap: %w", err)
	}

	return amountTo, nil
}

type SwapByID struct {
	model.Swap
	Rate         decimal.Decimal `json:"rate"`
	SwapDuration time.Duration   `json:"swap_duration"`
}

// GetSwapByID возвращает Swap по ID
func (s *swap) GetSwapByID(ctx context.Context, id uuid.UUID) (*SwapByID, error) {
	sw, err := s.swapRep.FindOne(ctx, &model.SwapFilter{
		ID: &id,
	})

	if err != nil {
		return nil, fmt.Errorf("swapRep.FindByTxID: %w", err)
	}

	if sw == nil {
		return nil, fmt.Errorf("no symbol: %s, err: %w", id.String(), businessError.TransactionNotFoundErr)
	}

	netFrom, err := s.coinRep.GetNetwork(ctx, sw.CoinFrom, sw.NetworkFrom)
	if err != nil {
		return nil, fmt.Errorf("getNetwork: %w", err)
	}

	sid := &SwapByID{}
	sid.Swap = *sw
	sid.Rate = sw.AmountTo.Div(sw.AmountFrom).Truncate(int32(netFrom.AccuracyWithdrawAndDeposit))
	sid.SwapDuration = s.depositWaitingPeriod

	return sid, nil
}

func (s *swap) Update(ctx context.Context, filter *model.SwapFilter, partial *model.SwapPartial) error {
	sw, err := s.swapRep.FindOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("findOne: %w", err)
	}

	err = s.swapRep.Update(ctx, sw, filter, partial)
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}
	return nil
}

func (s *swap) GetSwaps(ctx context.Context, filter *model.SwapFilter) ([]*model.Swap, int, error) {
	swaps, err := s.swapRep.Find(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("getSwaps: %w", err)
	}

	total, err := s.swapRep.CountTotalWithFilter(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("countTotalByFilter: %w", err)
	}

	return swaps, total, nil
}

func (s *swap) getBalanceByCoin(ctx context.Context, swapID uuid.UUID, account int, coin string, accountType string) (decimal.Decimal, error) {
	balance, err := s.exchangeAccount.GetBalanceByCoin(ctx, account, coin, accountType)
	if err != nil {
		return decimal.Zero, fmt.Errorf("getBalanceByCoin: %w", err)
	}

	if balance.WalletBalance.IsZero() || !balance.WalletBalance.Equal(balance.TransferBalance) {
		log.Warn(ctx, "swap_id: %s: balance wallet and transfer balance not equal: %s %s",
			swapID.String(),
			balance.WalletBalance.String(),
			balance.TransferBalance.String(),
		)
		return decimal.Zero, nil
	}

	return balance.TransferBalance, nil
}
