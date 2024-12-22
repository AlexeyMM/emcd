package worker

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"code.emcdtech.com/emcd/sdk/log"
	"github.com/google/uuid"

	"code.emcdtech.com/b2b/swap/internal/repository"
	"code.emcdtech.com/b2b/swap/internal/service"
	"code.emcdtech.com/b2b/swap/internal/slack"
	"code.emcdtech.com/b2b/swap/model"
)

const (
	ctxPeriod = time.Minute

	checkingDefaultInterval  = 5 * time.Second
	checkingWithdrawInterval = 20 * time.Second

	// Ограничения API - 100req/min
	checkingWaitDepositInterval = 20 * time.Second

	retryAfterErrorCount = 3
	retryInterval        = 10 * time.Second
)

type SwapExecutor struct {
	activeSwapRep                repository.ActiveSwap
	errorCounter                 repository.ErrorCounter
	swap                         service.Swap
	swapStatus                   service.SwapStatusUpdater
	order                        service.Order
	withdraw                     service.Withdraw
	transfer                     service.Transfer
	executeCh                    chan *model.Swap
	wg                           sync.WaitGroup
	slack                        slack.Slack
	depositWaitingPeriod         time.Duration
	workerGroup                  int
	busyWorkers                  int32
	busyWorkersThresholdForAlert int32
}

func NewSwapExecutor(
	activeSwapRep repository.ActiveSwap,
	errorCounter repository.ErrorCounter,
	swap service.Swap,
	swapStatus service.SwapStatusUpdater,
	order service.Order,
	withdraw service.Withdraw,
	transfer service.Transfer,
	executeCh chan *model.Swap,
	slack slack.Slack,
	depositWaitingPeriod time.Duration,
	workerGroup int,
	busyWorkersThresholdForAlert int32,
) *SwapExecutor {
	return &SwapExecutor{
		activeSwapRep:                activeSwapRep,
		errorCounter:                 errorCounter,
		swap:                         swap,
		swapStatus:                   swapStatus,
		order:                        order,
		withdraw:                     withdraw,
		transfer:                     transfer,
		executeCh:                    executeCh,
		slack:                        slack,
		depositWaitingPeriod:         depositWaitingPeriod,
		workerGroup:                  workerGroup,
		busyWorkersThresholdForAlert: busyWorkersThresholdForAlert,
	}
}

// UploadSwaps загружает из базы все незавершённые переводы, при старте приложения.
func (e *SwapExecutor) UploadSwaps(ctx context.Context) error {
	swaps, err := e.swap.GetAllActiveSwaps(ctx)
	if err != nil {
		return fmt.Errorf("uploadSwaps: getAllActiveSwaps: %w", err)
	}

	var activeSwaps []*model.Swap

	for i := range swaps {
		if e.swapExpired(swaps[i]) {
			err = e.swapStatus.UpdateAndBroadcast(ctx, swaps[i], model.Cancel)
			if err != nil {
				return fmt.Errorf("uploadSwaps: updateAndBroadcast: %w", err)
			}
		} else {
			activeSwaps = append(activeSwaps, swaps[i])
		}
	}

	for i := range activeSwaps {
		log.Debug(ctx, "swap start execute: %s", activeSwaps[i].ID)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case e.executeCh <- activeSwaps[i]:
			e.activeSwapRep.Add(activeSwaps[i].ID)
			log.Debug(ctx, "uploadSwaps: activeSwapRep.Add: %s", activeSwaps[i].ID.String())
		default:
			log.Error(ctx, "swapExecutor: uploadSwaps: executeCh is overload: %s", swaps[i].ID)
		}
	}

	return nil
}

func (e *SwapExecutor) Run(ctx context.Context) error {
	// WorkerGroup ограничивает возможное количество одновременно выполняемых свопов
	for i := 0; i < e.workerGroup; i++ {
		e.wg.Add(1)
		go func() {
			defer e.wg.Done()
			for {
				select {
				case <-ctx.Done():
					log.Debug(ctx, "worker stop: %s", ctx.Err())
					return

				case swap := <-e.executeCh:
					atomic.AddInt32(&e.busyWorkers, 1)
					e.activeSwapRep.Add(swap.ID)

					if atomic.LoadInt32(&e.busyWorkers) >= e.busyWorkersThresholdForAlert {
						err := e.slack.Send(ctx, fmt.Sprintf("Одновременно выполняется: %d свопов", atomic.LoadInt32(&e.busyWorkers)))
						if err != nil {
							log.Error(ctx, "slack send failed 1: %s", err)
							// Без return
						}
					}

					e.handle(ctx, swap)
					e.activeSwapRep.Delete(swap.ID)
					atomic.AddInt32(&e.busyWorkers, -1)
				}
			}
		}()
	}

	e.wg.Wait()
	return nil
}

func (e *SwapExecutor) handle(ctx context.Context, swap *model.Swap) {
	// Проверяем не истекло ли время свопа
	if e.swapExpired(swap) {
		err := e.swapStatus.UpdateAndBroadcast(ctx, swap, model.Cancel)
		if err != nil {
			log.Error(ctx, "swap_executor: updateAndBroadcast: %s, err: %s", swap.ID, err)
			return
		}
	}

	newCtx, cancel := context.WithTimeout(ctx, ctxPeriod)
	defer cancel()

	isRepeat, err := e.execute(newCtx, swap)
	if err != nil {
		e.handleError(ctx, swap, err)

		// handle - рекурсивная функция. Defer отработает после завершения последнего рекурсивного вызова
		defer e.errorCounter.Delete(swap.ID)
	}

	if !isRepeat {
		return
	}

	switch swap.Status {
	case model.Unknown:
		log.Error(ctx, "status unknown but repeat is true: %v", swap)
		return
	case model.CheckDeposit:
		<-time.After(checkingWaitDepositInterval)
		// Здесь вызываем handle, вместо того, что бы слать в канал e.executeCh <- swap
		// для того, что бы не удалять swap из кеша активных свопов e.activeSwapRep.Delete(swap.ID)
		// что бы, в момент пиковой нагрузки, место этого свопа не занял другой своп,
		// который может стартовать из service.Swap (новый своп)

		// TODO поменять на e.executeCh <- swap, перейти на очередь

		e.handle(ctx, swap)
	case model.WaitWithdraw:
		<-time.After(checkingWithdrawInterval)
		e.handle(ctx, swap)
	default:
		<-time.After(checkingDefaultInterval)
		e.handle(ctx, swap)
	}
}

func (e *SwapExecutor) execute(ctx context.Context, sw *model.Swap) (bool, error) {
	// Любой блок кода должен начинаться с проверки статуса, например:
	// sw.Status == model.WaitDeposit
	// Все переменные, которые нужны для выполнения блока кода, нужно достать из базы.
	// Нельзя сохранять переменные в общей области видимости метода execute,
	// даже если вы понимаете, что блоки кода будут выполняться последовательно.
	// Такой подход исключит ошибки отсутствия переменных, при перезапуске сервиса, например.

	var err error

	if sw.Status == model.WaitDeposit {
		sw.Status, err = e.swap.WaitDeposit(ctx, sw)
		if err != nil {
			return false, fmt.Errorf("waitDeposit: %w", err)
		}

		if sw.Status == model.WaitDeposit {
			return true, nil
		}
	}

	if sw.Status == model.CheckDeposit {
		sw.Status, err = e.swap.CheckDeposit(ctx, sw)
		if err != nil {
			return false, fmt.Errorf("checkDeposit: %w", err)
		}

		if sw.Status == model.CheckDeposit {
			return true, nil
		}
	}

	if sw.Status == model.TransferToUnified {
		sw.Status, err = e.swap.TransferToUnified(ctx, sw)
		if err != nil {
			return false, fmt.Errorf("transferToUnified: %w", err)
		}

		if sw.Status == model.TransferToUnified {
			return true, nil
		}
	}

	var orders model.Orders

	if sw.Status == model.CreateOrder {
		orders, sw.Status, err = e.swap.CreateOrders(ctx, sw)
		if err != nil {
			return false, fmt.Errorf("createOrder: %w", err)
		}
		if sw.Status == model.CreateOrder {
			return true, nil
		}
	}

	if sw.Status == model.PlaceOrder {
		if len(orders) == 0 {
			orders, err = e.getOrders(ctx, sw)
			if err != nil {
				return false, fmt.Errorf("getOrders: %w", err)
			}
		}
		firstOrder, err := orders.FindFirst()
		if err != nil {
			return false, fmt.Errorf("findFirst: %w", err)
		}

		sw.Status, err = e.swap.PlaceOrder(ctx, sw, firstOrder)
		if err != nil {
			return false, fmt.Errorf("placeOrder 1: %w", err)
		}

		if sw.Status == model.PlaceOrder {
			return true, nil
		}
	}

	if sw.Status == model.CheckOrder {
		if len(orders) == 0 {
			orders, err = e.getOrders(ctx, sw)
			if err != nil {
				return false, fmt.Errorf("getOrders: %w", err)
			}
		}
		firstOrder, err := orders.FindFirst()
		if err != nil {
			return false, fmt.Errorf("findFirstOrder: %w", err)
		}
		sw.Status, err = e.swap.CheckOrder(ctx, sw, firstOrder.ID)
		if err != nil {
			return false, fmt.Errorf("checkOrder 1: %w", err)
		}

		if sw.Status == model.CheckOrder {
			return true, nil
		}
	}

	if sw.Status == model.PlaceAdditionalOrder {
		if len(orders) == 0 {
			orders, err = e.getOrders(ctx, sw)
			if err != nil {
				return false, fmt.Errorf("getOrders: %w", err)
			}
		}

		secondOrder, err := orders.FindSecond()
		if err != nil {
			return false, fmt.Errorf("findSecond: %w", err)
		}

		// Берём из базы, потому что для второго ордера мы не знаем AmountFrom, он устанавливается исходя из фактически
		// полученных монет первого ордера
		additionalOrder, err := e.order.FindOne(ctx, &model.OrderFilter{
			ID: &secondOrder.ID,
		})
		if err != nil {
			return false, fmt.Errorf("findOne: %w", err)
		}
		if additionalOrder == nil {
			return false, fmt.Errorf("additional order doesn't exist: %s", sw.ID)
		}
		sw.Status, err = e.swap.PlaceOrder(ctx, sw, additionalOrder)
		if err != nil {
			return false, fmt.Errorf("placeOrder 2: %w", err)
		}

		if sw.Status == model.PlaceAdditionalOrder {
			return true, nil
		}
	}

	if sw.Status == model.CheckAdditionalOrder {
		if len(orders) == 0 {
			orders, err = e.getOrders(ctx, sw)
			if err != nil {
				return false, fmt.Errorf("getOrders: %w", err)
			}
		}

		secondOrder, err := orders.FindSecond()
		if err != nil {
			return false, fmt.Errorf("findSecond: %w", err)
		}

		sw.Status, err = e.swap.CheckOrder(ctx, sw, secondOrder.ID)
		if err != nil {
			return false, fmt.Errorf("checkOrder 2: %w", err)
		}

		if sw.Status == model.CheckAdditionalOrder {
			return true, nil
		}
	}

	if sw.Status == model.TransferFromSubToMaster {
		if len(orders) == 0 {
			orders, err = e.getOrders(ctx, sw)
			if err != nil {
				return false, fmt.Errorf("getOrders: %w", err)
			}
		}

		// Получаем последний ордер, что бы узнать сколько у нас монет для вывода
		var lastOrder *model.Order
		if len(orders) == 1 {
			firstOrder, err := orders.FindFirst()
			if err != nil {
				return false, fmt.Errorf("findFirst: %w", err)
			}

			lastOrder, err = e.order.FindOne(ctx, &model.OrderFilter{
				ID: &firstOrder.ID,
			})
			if err != nil {
				return false, fmt.Errorf("getOrder 2: %w", err)
			}
		} else if len(orders) == 2 {
			secondOrder, err := orders.FindSecond()
			if err != nil {
				return false, fmt.Errorf("findSecond: %w", err)
			}

			lastOrder, err = e.order.FindOne(ctx, &model.OrderFilter{
				ID: &secondOrder.ID,
			})
			if err != nil {
				return false, fmt.Errorf("getOrder 3: %w", err)
			}
		}
		if lastOrder == nil {
			return false, fmt.Errorf("lastOrder doesn't exist: %s", sw.ID)
		}

		sw.Status, err = e.swap.TransferFromSubToMaster(ctx, sw, &model.InternalTransfer{
			ID:            uuid.New(),
			Coin:          sw.CoinTo,
			Amount:        lastOrder.AmountTo,
			FromAccountID: sw.AccountFrom,
		})
		if err != nil {
			return false, fmt.Errorf("transferFromSubToMaster: %w", err)
		}

		if sw.Status == model.TransferFromSubToMaster {
			return true, nil
		}
	}

	if sw.Status == model.CheckTransferFromSubToMaster {
		sw.Status, err = e.swap.CheckTransferFromSubToMaster(ctx, sw, sw.AccountFrom)
		if err != nil {
			return false, fmt.Errorf("getLastTransferFromSubToMaster: %w", err)
		}

		if sw.Status == model.CheckTransferFromSubToMaster {
			return true, nil
		}
	}

	if sw.Status == model.PrepareWithdraw {
		sw.Status, err = e.swap.PrepareWithdraw(ctx, sw)
		if err != nil {
			return false, fmt.Errorf("prepareWithdraw: %w", err)
		}

		if sw.Status == model.PrepareWithdraw {
			return true, nil
		}
	}

	if sw.Status == model.WithdrawSwapStatus {
		sw.Status, err = e.swap.Withdraw(ctx, sw)
		if err != nil {
			return false, fmt.Errorf("withdraw: %w", err)
		}

		if sw.Status == model.WithdrawSwapStatus {
			return true, nil
		}
	}

	if sw.Status == model.WaitWithdraw {
		w, err := e.withdraw.GetBySwapID(ctx, sw.ID)
		if err != nil {
			return false, fmt.Errorf("getBySwapID: %w", err)
		}
		sw.Status, err = e.swap.CheckWithdraw(ctx, sw, w.ID)
		if err != nil {
			return false, fmt.Errorf("checkWithdraw: %w", err)
		}

		if sw.Status == model.WaitWithdraw {
			return true, nil
		}
	}

	if sw.Status == model.Completed {
		return false, nil
	}

	if sw.Status == model.Cancel {
		return false, nil
	}

	if sw.Status == model.Error {
		// Вернуть ошибку, что бы сделать swapStatus.UpdateAndBroadcast
		return false, fmt.Errorf("swap error: %s", sw.ID)
	}

	if sw.Status == model.ManualCompleted {
		return false, nil
	}

	log.Warn(ctx, "unexpected swap status: %+v", sw)
	return true, nil
}

func (e *SwapExecutor) getOrders(ctx context.Context, sw *model.Swap) (model.Orders, error) {
	orders, err := e.order.Find(ctx, &model.OrderFilter{
		AccountID: &sw.AccountFrom,
	})
	if err != nil {
		return nil, fmt.Errorf("getOrdersByAccountID: %w", err)
	}
	if len(orders) == 0 {
		return nil, fmt.Errorf("len(orders) == 0: %s", sw.ID)
	}
	return orders, nil
}

func (e *SwapExecutor) swapExpired(swap *model.Swap) bool {
	if swap.Status > model.WaitDeposit {
		return false
	}
	return swap.StartTime.Add(e.depositWaitingPeriod).Before(time.Now().UTC())
}

func (e *SwapExecutor) handleError(ctx context.Context, swap *model.Swap, err error) {
	log.Error(ctx, "swap_executor: swap_id: %s: err: %s", swap.ID.String(), err.Error())

	count := e.errorCounter.Inc(swap.ID)

	if count < retryAfterErrorCount {
		log.Debug(ctx, "swap_executor: retry swap: %s", swap.ID.String())
		<-time.After(retryInterval)
		e.handle(ctx, swap)
		return
	}

	// Только отправляем статус на фронт, но не обновляем базу. Что бы исправить проблему и довезти после рестарта
	e.swapStatus.Broadcast(ctx, swap, model.Error)

	err = e.slack.Send(ctx, fmt.Sprintf("swap %s failed: %s", swap.ID.String(), err.Error()))
	if err != nil {
		log.Error(ctx, "slack send failed 2: %s, err: %s", swap.ID.String(), err.Error())
		return
	}
}
